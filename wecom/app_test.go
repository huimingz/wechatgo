package wecom

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"

	"github.com/huimingz/wechatgo/testdata"
)

type TestSuite struct {
	suite.Suite
	httpClient *http.Client
}

func (s *TestSuite) SetupSuite() {
	s.httpClient = &http.Client{}
	httpmock.ActivateNonDefault(s.httpClient)
}

func (s *TestSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func (s *TestSuite) BeforeTest(suiteName, testName string) {
	httpmock.Reset()

	s.registerResponder(http.MethodGet, "/cgi-bin/gettoken")
}

func (s *TestSuite) registerResponder(method, path string) {
	filename := fmt.Sprintf("response%s.json", strings.ReplaceAll(path, "/", "_"))
	responder := httpmock.NewStringResponder(http.StatusOK, s.readFixture(filename))
	httpmock.RegisterResponder(method, _BASE_URL+path, responder)
}

func (s *TestSuite) registerSuccessResponder(method, path string) {
	responder := httpmock.NewStringResponder(http.StatusOK, `{"errcode":0,"errmsg":"ok"}`)
	httpmock.RegisterResponder(method, _BASE_URL+path, responder)
}

func (s *TestSuite) readFixture(filename string) string {
	content, err := os.ReadFile("fixtures/" + filename)
	s.NoError(err)

	return string(content)
}

type AppTestSuite struct {
	TestSuite
	wecom *Wecom
}

func (s *AppTestSuite) SetupSuite() {
	s.TestSuite.SetupSuite()

	var conf = testdata.TestConf
	s.wecom = NewWecom(conf.CorpId, conf.CorpSecret, conf.AgentId, WechatClientWithHTTPClient(s.httpClient))
}

func (s *AppTestSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func (s *AppTestSuite) TestCreateApp() {
	s.registerSuccessResponder(http.MethodPost, urlCreateApp)

	err := s.wecom.App.CreateApp(context.Background(), AppInfo{
		AgentId:     1000321,
		Description: "some comment for app",
	})

	s.NoError(err)
}

func (s *AppTestSuite) TestShouldGetAllApp() {
	s.registerResponder(http.MethodGet, urlGetAllApp)

	appIntro, err := s.wecom.App.GetAllApp(context.Background())

	s.NoError(err)
	s.NotEmpty(appIntro)
}

func (s *AppTestSuite) TestShouldGetApp() {
	s.registerResponder(http.MethodGet, urlGetApp)

	appDetail, err := s.wecom.App.GetApp(context.Background(), testdata.TestConf.AgentId)

	s.NoError(err)
	s.NotEmpty(appDetail.Name)
}

func (s *AppTestSuite) TestShouldCreateMenu() {
	responder := httpmock.NewStringResponder(http.StatusOK, s.readFixture("response_success.json"))
	httpmock.RegisterResponder(http.MethodPost, _BASE_URL+urlCreateMenu, responder)

	menu := Menu{
		Button: []Button{
			{
				Type:      "view",
				Name:      "golang",
				Url:       "https://www.golang.org",
				Key:       "",
				SubButton: nil,
			},
		},
	}
	err := s.wecom.App.CreateMenu(context.Background(), menu, 0)

	s.NoError(err)
}

func (s *AppTestSuite) TestShouldGetMenu() {
	s.registerResponder(http.MethodGet, urlGetMenu)

	menu, err := s.wecom.App.GetMenu(context.Background(), testdata.TestConf.AgentId)

	s.NoError(err)
	s.NotEmpty(menu.Button)
	s.NotEmpty(menu.Button[0].Name)
}

func (s *AppTestSuite) TestShouldDeleteMenu() {
	responder := httpmock.NewStringResponder(http.StatusOK, s.readFixture("response_success.json"))
	httpmock.RegisterResponder(http.MethodGet, _BASE_URL+urlDelMenu, responder)

	err := s.wecom.App.DeleteMenu(context.Background(), 0)

	s.NoError(err)
}

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}
