package dept

import (
	"testing"

	"github.com/huimingz/wechatgo"
	"github.com/huimingz/wechatgo/testdata"
	"github.com/huimingz/wechatgo/wecom"
)

var wechatDept *WechatDept

func TestWechatDept_Get(t *testing.T) {
	info, err := wechatDept.Get(1)
	if err != nil {
		t.Errorf("WechatDept.Get() error = '%s'", err)
	}
	if len(info) == 0 {
		t.Error("WechatDept.Get() error = '返回指定部门内容为空'")
	}
}

func TestWechatDept_GetList(t *testing.T) {
	info, err := wechatDept.GetList()
	if err != nil {
		t.Errorf("WechatDept.GetList() error = '%s'", err)
	}
	if len(info) == 0 {
		t.Error("WechatDept.GetList() error = '返回指定部门内容为空'")
	}
}

func TestWechatDept_GetUserList(t *testing.T) {
	users, err := wechatDept.GetUserList(1, false)
	if err != nil {
		t.Errorf("WechatDept.GetUserList() error = '%s'", err)
	}
	if len(users) == 0 {
		t.Error("WechatDept.GetUserList() error = '返回部门成员为空'")
	}

	cusers, err := wechatDept.GetUserList(1, true)
	if err != nil {
		t.Errorf("WechatDept.GetUserList() error = '%s'", err)
	}
	if len(cusers) <= len(users) {
		t.Error("WechatDept.GetUserList() error = '返回部门和子部门成员数少于指定部门成员数'")
	}
}

func TestWechatDept_GetUserDetailList(t *testing.T) {
	users, err := wechatDept.GetUserDetailList(1, false)
	if err != nil {
		t.Errorf("WechatDept.GetUserDetailList() error = '%s'", err)
	}
	if len(users) == 0 {
		t.Error("WechatDept.GetUserDetailList() error = '返回部门成员详情为空'")
	}

	cusers, err := wechatDept.GetUserDetailList(1, true)
	if err != nil {
		t.Errorf("WechatDept.GetUserDetailList() error = '%s'", err)
	}
	if len(cusers) <= len(users) {
		t.Error("WechatDept.GetUserList() error = '返回部门和子部门成员数少于指定部门成员数'")
	}
}

func TestWechatDept_Create(t *testing.T) {
	id, err := wechatDept.Create("testdata", 1, 0, 18)
	if err != nil {
		if v, ok := err.(*wechatgo.WXMsgError); ok {
			if v.ErrCode != 60008 {
				t.Errorf("WechatDept.Create() error = '%s'", err)
			}
		}
	}
	if id != 18 {
		t.Errorf("WechatDept.Create() error = 'id[%d] != 18'", id)
	}
}

func TestWechatDept_Update(t *testing.T) {
	err := wechatDept.Update(18, 0, 0, "test_test")
	if err != nil {
		t.Errorf("WechatDept.Update() error = '%s'", err)
	}

	infos, err := wechatDept.Get(18)
	if err != nil {
		t.Errorf("WechatDept.Update() error = '%s'", err)
	}
	if infos[0].Name != "test_test" {
		t.Errorf("WechatDept.Update() error = '更新部门名称失败'")
	}
}

func TestWechatDept_Delete(t *testing.T) {
	err := wechatDept.Delete(18)
	if err != nil {
		t.Errorf("WechatDept.Delete() error = '%s'", err)
	}
}

func init() {
	var conf = testdata.TestConf
	var wechatClient = wecom.NewWechatClient(conf.CorpId, conf.UserSecret, conf.AgentId)
	wechatDept = NewWechatDept(wechatClient)
}
