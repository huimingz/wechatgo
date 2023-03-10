package app

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"

	"github.com/huimingz/wechatgo"
	"github.com/huimingz/wechatgo/testdata"
	"github.com/huimingz/wechatgo/wecom"
)

type AppTestSuite struct {
	suite.Suite
	wechatAppManage *WechatAppManage
	httpClient      *http.Client
}

func (s *AppTestSuite) SetupSuite() {
	s.httpClient = &http.Client{}
	httpmock.ActivateNonDefault(s.httpClient)

	content, err := os.ReadFile("fixtures/response_cgi-bin_gettoken.json")
	s.NoError(err)
	responder := httpmock.NewStringResponder(http.StatusOK, string(content))
	httpmock.RegisterResponder("GET", wecom.BASE_URL+"/cgi-bin/gettoken", responder)

	var conf = testdata.TestConf
	wechatClient := wecom.NewWechatClient(conf.CorpId, conf.CorpSecret, conf.AgentId, wecom.WechatClientWithHTTPClient(s.httpClient))
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
	wd, err := os.Getwd()
	s.NoError(err)
	content, err := os.ReadFile(wd + "/fixtures/response_cgi-bin_agent_get.json")
	s.NoError(err)

	responder := httpmock.NewStringResponder(http.StatusOK, string(content))
	httpmock.RegisterResponder("GET", wecom.BASE_URL+urlGetApp, responder)

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
