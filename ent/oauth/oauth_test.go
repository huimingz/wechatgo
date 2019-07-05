package oauth

import (
	"fmt"
	"testing"

	"github.com/huimingz/wechatgo"
	"github.com/huimingz/wechatgo/ent"
	"github.com/huimingz/wechatgo/testdata"
)

var wechatOauth *WechatOAuth

func TestWechatOAuth_WebOAuthUrl(t *testing.T) {
	uri := wechatOauth.WebOAuthUrl("http://example.com", "", "")
	if uri == "" {
		t.Error("WechatOauth.WebOauthUrl() error = '生成的授权连接为空'")
	}
	fmt.Println(uri)
}

func TestWechatOAuth_QRCodeUrl(t *testing.T) {
	uri := wechatOauth.QRCodeUrl("http://example.com", "")
	if uri == "" {
		t.Error("WechatOauth.WebOauthUrl() error = '生成的二维码授权连接为空'")
	}
	fmt.Println(uri)
}

func TestWechatOAuth_GetUserInfo(t *testing.T) {
	code := "KyUxwLmpx5coUKU_BiPA2ICBdlDYZQtcPmeyocC_QUY"
	_, err := wechatOauth.GetUserInfo(code)
	if err != nil {
		if v, ok := err.(*wechatgo.WXMsgError); ok {
			if v.ErrCode != 40029 {
				t.Errorf("WechatOauth.GetUserInfo() error = 'error code != 40029, current code = %d'", v.ErrCode)
			}
		}

	}
}

func init() {
	var conf = testdata.TestConf
	var wechatClient = ent.NewWechatClient(conf.CorpId, conf.CorpSecret, conf.AgentId, nil, 0, nil, nil, nil)
	wechatOauth = NewWechatOAuth(wechatClient)
}
