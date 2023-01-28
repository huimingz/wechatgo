package app

import (
	"context"
	"testing"

	"github.com/huimingz/wechatgo"
	"github.com/huimingz/wechatgo/testdata"
	"github.com/huimingz/wechatgo/wecom"
)

var wechatAppManage *WechatAppManage

func TestWechatAppManage_CreateApp(t *testing.T) {
	app := AppInfo{}
	app.AgentId = 1000321
	app.Description = "some comment for app"
	err := wechatAppManage.CreateApp(context.Background(), app)
	if err != nil {
		if v, ok := err.(*wechatgo.WechatMessageError); ok {
			if v.ErrCode != 301002 {
				t.Error("WechatAppManage.CreateApp() error = 'error code != 301002'")
			}
		} else {
			t.Error("WechatAppManage.CreateApp() error = 'error isn't `*enterprise.WechatMessageError` type'")
		}
	}
}

func TestWechatAppManage_GetAllApp(t *testing.T) {
	appIntro, err := wechatAppManage.GetAllApp(context.Background())
	if err != nil {
		t.Errorf("WechatAppManage.GetAllApp() error = '%s'", err)
	}
	if len(appIntro) == 0 {
		t.Error("WechatAppManage.GetAllApp() error = '返回了空内容'")
	}
}

func TestWechatAppManage_GetApp(t *testing.T) {
	appDetail, err := wechatAppManage.GetApp(context.Background(), testdata.TestConf.AgentId)
	if err != nil {
		t.Errorf("WechatAppManage.GetApp(%d) error = '%s'", testdata.TestConf.AgentId, err)
	}
	if appDetail.Name == "" {
		t.Error("WechatAppManage.GetAllApp() error = '返回的应用详情为空内容'")
	}
}

func TestWechatAppManage_CreateMenu(t *testing.T) {
	menu := Menu{}
	button1 := Button{
		Type:      "view",
		Name:      "golang",
		Url:       "https://www.golang.org",
		Key:       "",
		SubButton: nil,
	}
	menu.Button = append(menu.Button, button1)
	err := wechatAppManage.CreateMenu(context.Background(), menu, 0)
	if err != nil {
		t.Errorf("WechatAppManage.CreateMenu() error = '%s'", err)
	}
}

func TestWechatAppManage_GetMenu(t *testing.T) {
	menu, err := wechatAppManage.GetMenu(context.Background(), testdata.TestConf.AgentId)
	if err != nil {
		t.Errorf("WechatAppManage.GetMenu(%d) error = '%s'", testdata.TestConf.AgentId, err)
	}
	if len(menu.Button) == 0 {
		t.Errorf("WechatAppManage.GetMenu(%d) error = '返回的菜单按钮为空'", testdata.TestConf.AgentId)
	}
	if menu.Button[0].Name != "golang" {
		t.Errorf("WechatAppManage.GetMenu(%d) error = 'Button.Name != golang'", testdata.TestConf.AgentId)
	}
}

func TestWechatAppManage_DeleteMenu(t *testing.T) {
	err := wechatAppManage.DeleteMenu(context.Background(), 0)
	if err != nil {
		t.Errorf("WechatAppManage.DeleteMenu() error = '%s'", err)
	}
}

func init() {
	var conf = testdata.TestConf
	var wechatClient = wecom.NewWechatClient(conf.CorpId, conf.CorpSecret, conf.AgentId)
	wechatAppManage = NewWechatAppManage(wechatClient)
}
