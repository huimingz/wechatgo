package user

import (
	"testing"

	"github.com/huimingz/wechatgo"
	"github.com/huimingz/wechatgo/testdata"
	"github.com/huimingz/wechatgo/wecom"
)

var wechatUser *WechatUser

func TestWechatUser_CreateUser(t *testing.T) {
	user := UserForCreate{}
	user.UserId = "example_xxx"
	user.Name = "xxx"
	user.Department = []int{1}
	user.Email = "xxx@example.com"
	user.ToInvite = false

	err := wechatUser.CreateUser(user)
	if err != nil {
		t.Errorf("WechatUser.CreateUser() error = '%s'", err)
	}
}

func TestWechatUser_GetUser(t *testing.T) {
	user, err := wechatUser.GetUser("example_xxx")
	if err != nil {
		t.Errorf("WechatUser.GetUser() error = '%s'", err)
	}
	if user.UserId != "example_xxx" {
		t.Error("WechatUser.GetUser() error = '返回用户数据不匹配'")
	}
}

func TestWechatUser_UpdateUser(t *testing.T) {
	user := UserForUpdate{}
	user.UserId = "example_xxx"
	user.Name = "new_xxx"

	err := wechatUser.UpdateUser(user)
	if err != nil {
		t.Errorf("WechatUser.UpdateUser() error = '%s'", err)
	}

	info, err := wechatUser.GetUser(user.UserId)
	if err != nil {
		t.Errorf("WechatUser.UpdateUser() error = '%s'", err)
	}
	if info.Name != user.Name {
		t.Errorf("WechatUser.UpdateUser() error = '更新后的用户名与期望不符, [%s != %s]'", info.Name, user.Name)
	}
}

func TestWechatUser_Invite(t *testing.T) {
	err := wechatUser.Invite([]string{"example_xxx"}, nil, nil)
	if err != nil {
		t.Errorf("WechatUser.Invite() error = '%s'", err)
	}
}

func TestWechatUser_Verify(t *testing.T) {
	err := wechatUser.Verify("example_xxx")
	if err != nil {
		t.Errorf("WechatUser.Verify() error = '%s'", err)
	}
}

func TestWechatUser_UserId2OpenId(t *testing.T) {
	openId, err := wechatUser.UserId2OpenId(testdata.TestConf.UserId)
	if err != nil {
		t.Errorf("WechatUser.UserId2OpenId() error = '%s'", err)
	}
	if openId != testdata.TestConf.OpenId {
		t.Errorf("WechatUser.UserId2OpenId() error = '返回的openid[%s]与期望不符'", openId)
	}
}

func TestWechatUser_DeleteUser(t *testing.T) {
	err := wechatUser.DeleteUser("example_xxx")
	if err != nil {
		t.Errorf("WechatUser.DeleteUser() error = '%s'", err)
	}

	err = wechatUser.DeleteUser("example_xxx")
	if err != nil {
		if v, ok := err.(*wechatgo.WXMsgError); ok {
			if v.ErrCode != 60111 {
				t.Errorf("WechatUser.DeleteUser() error = '错误代码[%d] != 60111'", v.ErrCode)
			}
		}
	} else {
		t.Error("WechatUser.DeleteUser() error = '未出现错误'")
	}
}

func TestWechatUser_DeleteUsers(t *testing.T) {
	user := UserForCreate{}
	user.UserId = "example_xxx"
	user.Name = "xxx"
	user.Email = "xxx@example.com"
	user.Department = []int{1}
	user.ToInvite = false

	err := wechatUser.CreateUser(user)
	if err != nil {
		t.Errorf("WechatUser.DeleteUsers() error = '%s'", err)
	}

	err = wechatUser.DeleteUsers([]string{user.UserId})
	if err != nil {
		t.Errorf("WechatUser.DeleteUsers() error = '%s'", err)
	}
}

func TestWechatUser_OpenId2UserId(t *testing.T) {
	userId, err := wechatUser.OpenId2UserId(testdata.TestConf.OpenId)
	if err != nil {
		t.Errorf("WechatUser.OpenId2UserId() error = '%s'", err)
	}

	if userId != testdata.TestConf.UserId {
		t.Error("WechatUser.OpenId2UserId() error = '返回的userid与期望不符'")
	}
}

func TestWechatUser_GetJoinQRCode(t *testing.T) {
	qrCode, err := wechatUser.GetJoinQRCode(0)
	if err != nil {
		t.Errorf("WechatUser.GetJoinQRCode() error = '%s'", err)
	}

	if qrCode == "" {
		t.Error("WechatUser.GetJoinQRCode() error = '返回qrcode为空'")
	}
}

func init() {
	var conf = testdata.TestConf
	var wechatClient = wecom.NewWechatClient(conf.CorpId, conf.UserSecret, conf.AgentId)
	wechatUser = NewWechatUser(wechatClient)
}
