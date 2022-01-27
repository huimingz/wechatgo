// 外部联系人管理
package extcontact

import (
	"testing"

	"github.com/huimingz/wechatgo/testdata"
	"github.com/huimingz/wechatgo/wecom"
)

var wechatContact *WechatContact

func TestWechatContact_GetUserList(t *testing.T) {
	_, err := wechatContact.GetUserList()
	if err == nil {
		t.Error("WechatContact.GetUserList() error = '未出现期望的错误'")
	}
}

func init() {
	var conf = testdata.TestConf
	var wechatClient = wecom.NewWechatClient(conf.CorpId, conf.UserSecret, conf.AgentId)
	wechatContact = NewWechatContact(wechatClient)
}
