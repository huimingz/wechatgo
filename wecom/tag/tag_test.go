package tag

import (
	"testing"

	"github.com/huimingz/wechatgo/testdata"
	"github.com/huimingz/wechatgo/wecom"
)

var wechatTag *WechatTag

func TestWechatTag_Create(t *testing.T) {
	tagId, err := wechatTag.Create(99, "example")
	if err != nil {
		t.Errorf("WechatTag.Create() error = '%s'", err)
	}
	if tagId != 99 {
		t.Errorf("WechatTag.Create() error = 'tagId[%d] != 99'", tagId)
	}
}

func TestWechatTag_Update(t *testing.T) {
	err := wechatTag.Update(99, "example_2")
	if err != nil {
		t.Errorf("WechatTag.Update() error = '%s'", err)
	}
}

func TestWechatTag_GetTagList(t *testing.T) {
	tagList, err := wechatTag.GetTagList()
	if err != nil {
		t.Errorf("WechatTag.GetTagList() error = '%s'", err)
	}

	hasExpectTag := false
	for _, tag := range tagList {
		if tag.TagId == 99 && tag.TagName == "example_2" {
			hasExpectTag = true
		}
	}
	if !hasExpectTag {
		t.Error("WechatTag.GetTagList() error = '无指定的标签'")
	}
}

func TestWechatTag_AddUser(t *testing.T) {
	err := wechatTag.AddUser(99, []string{testdata.TestConf.UserId}, nil)
	if err != nil {
		t.Errorf("WechatTag.AddUser() error = '%s'", err)
	}

	err = wechatTag.AddUser(99, []string{"huimingzxxx"}, nil)
	if err == nil {
		t.Error("WechatTag.AddUser() error = '未出现期望的错误'")
	} else {
		if v, ok := err.(*TagError); ok {
			if v.ErrCode != 40070 {
				t.Errorf("WechatTag.AddUser() error = 'err.ErrCode[%s] != 40070'", v.InvalidList)
			}
		}
	}
}

func TestWechatTag_GetUserList(t *testing.T) {
	tagName, userList, partyList, err := wechatTag.GetUserList(99)
	if err != nil {
		t.Errorf("WechatTag.GetUserList() error = '%s'", err)
	}
	if tagName != "example_2" {
		t.Errorf("WechatTag.GetUserList() error = 'tagName[%s] != example_2'", tagName)
	}
	if len(userList) != 1 {
		t.Error("WechatTag.GetUserList() error = '用户数 != 1'")
	}
	if len(partyList) != 0 {
		t.Error("WechatTag.GetUserList() error = '部门数 != 0'")
	}
}

func TestWechatTag_DeleteUser(t *testing.T) {
	err := wechatTag.DeleteUser(99, []string{testdata.TestConf.UserId}, nil)
	if err != nil {
		t.Errorf("WechatTag.DeleteUser() error = '%s'", err)
	}
}

func TestWechatTag_Delete(t *testing.T) {
	err := wechatTag.Delete(99)
	if err != nil {
		t.Errorf("WechatTag.Delete() error = '%s'", err)
	}
}

func init() {
	var conf = testdata.TestConf
	var wechatClient = wecom.NewWechatClient(conf.CorpId, conf.UserSecret, conf.AgentId)
	wechatTag = NewWechatTag(wechatClient)
}
