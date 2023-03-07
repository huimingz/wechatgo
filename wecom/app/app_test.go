package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/huimingz/wechatgo"
	"github.com/huimingz/wechatgo/testdata"
	"github.com/huimingz/wechatgo/wecom"
)

var wechatAppManage *WechatAppManage

type AppTestSuite struct {
	suite.Suite
}

func (s *AppTestSuite) SetupSuite() {
	var conf = testdata.TestConf
	var wechatClient = wecom.NewWechatClient(conf.CorpId, conf.CorpSecret, conf.AgentId)
	wechatAppManage = NewWechatAppManage(wechatClient)
}

func (s *AppTestSuite) TestCreateApp() {
	app := AppInfo{}
	app.AgentId = 1000321
	app.Description = "some comment for app"
	err := wechatAppManage.CreateApp(context.Background(), app)
	s.NoError(err)

	v, ok := err.(*wechatgo.WechatMessageError)
	s.True(ok)
	s.NotEqual(301002, v.ErrCode, "error code != 301002")
}

func (s *AppTestSuite) TestShouldGetAllApp() {
	appIntro, err := wechatAppManage.GetAllApp(context.Background())

	s.NoError(err)
	s.NotEmpty(appIntro)
}

func (s *AppTestSuite) TestShouldGetApp() {
	appDetail, err := wechatAppManage.GetApp(context.Background(), testdata.TestConf.AgentId)

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

	err := wechatAppManage.CreateMenu(context.Background(), menu, 0)

	s.NoError(err)
}

func (s *AppTestSuite) TestShouldGetMenu() {
	menu, err := wechatAppManage.GetMenu(context.Background(), testdata.TestConf.AgentId)

	s.NoError(err)
	s.NotEmpty(menu.Button)
	s.Equal("golang", menu.Button[0].Name)
}

func (s *AppTestSuite) TestShouldDeleteMenu() {
	err := wechatAppManage.DeleteMenu(context.Background(), 0)

	s.NoError(err)
}

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}
