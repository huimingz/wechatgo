package wecom

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"

	"github.com/huimingz/wechatgo"
	"github.com/huimingz/wechatgo/testdata"
)

var httpClient = &http.Client{}

type TestSuite struct {
	suite.Suite
}

func (s *TestSuite) SetupSuite() {
	httpmock.ActivateNonDefault(httpClient)
}

func (s *TestSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func (s *TestSuite) BeforeTest(suiteName, testName string) {
	httpmock.Reset()

	responder := httpmock.NewStringResponder(http.StatusOK, s.readFixture("response_cgi-bin_gettoken.json"))
	httpmock.RegisterResponder("GET", BASE_URL+"/cgi-bin/gettoken", responder)
}

func (s *TestSuite) readFixture(filename string) string {
	content, err := os.ReadFile("fixtures/" + filename)
	s.NoError(err)

	return string(content)
}

type AppTestSuite struct {
	TestSuite
	wechatAppManage *WechatAppManage
	httpClient      *http.Client
}

func (s *AppTestSuite) SetupSuite() {
	s.TestSuite.SetupSuite()

	var conf = testdata.TestConf
	wechatClient := NewWechatClient(conf.CorpId, conf.CorpSecret, conf.AgentId, WechatClientWithHTTPClient(httpClient))
	s.wechatAppManage = NewWechatAppManage(wechatClient)
}

func (s *AppTestSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func (s *AppTestSuite) TestCreateApp() {
	app := AppInfo{}
	app.AgentId = 1000321
	app.Description = "some comment for app"
	err := s.wechatAppManage.CreateApp(context.Background(), app)
	s.NoError(err)

	v, ok := err.(*wechatgo.WechatMessageError)
	s.True(ok)
	s.NotEqual(301002, v.ErrCode, "error code != 301002")
}

func (s *AppTestSuite) TestShouldGetAllApp() {
	appIntro, err := s.wechatAppManage.GetAllApp(context.Background())

	s.NoError(err)
	s.NotEmpty(appIntro)
}

func (s *AppTestSuite) TestShouldGetApp() {
	responder := httpmock.NewStringResponder(http.StatusOK, s.readFixture("response_cgi-bin_agent_get.json"))
	httpmock.RegisterResponder("GET", BASE_URL+urlGetApp, responder)

	appDetail, err := s.wechatAppManage.GetApp(context.Background(), testdata.TestConf.AgentId)

	s.NoError(err)
	s.NotEmpty(appDetail.Name)
}

func (s *AppTestSuite) TestShouldCreateMenu() {
	menu := Menu{}
	button := Button{
		Type:      "view",
		Name:      "golang",
		Url:       "https://www.golang.org",
		Key:       "",
		SubButton: nil,
	}
	menu.Button = append(menu.Button, button)

	err := s.wechatAppManage.CreateMenu(context.Background(), menu, 0)

	s.NoError(err)
}

func (s *AppTestSuite) TestShouldGetMenu() {
	menu, err := s.wechatAppManage.GetMenu(context.Background(), testdata.TestConf.AgentId)

	s.NoError(err)
	s.NotEmpty(menu.Button)
	s.Equal("golang", menu.Button[0].Name)
}

func (s *AppTestSuite) TestShouldDeleteMenu() {
	err := s.wechatAppManage.DeleteMenu(context.Background(), 0)

	s.NoError(err)
}

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}
