// Package tag 标签管理
package tag

import (
	"net/url"
	"strconv"

	"github.com/huimingz/wechatgo/ent"
)

const (
	urlGetTagList    = "/cgi-bin/tag/list"
	urlCreateTag     = "/cgi-bin/tag/create"
	urlUpdateTag     = "/cgi-bin/tag/update"
	urlDeleteTag     = "/cgi-bin/tag/delete"
	urlGetTagUser    = "/cgi-bin/tag/get"
	urlAddTagUser    = "/cgi-bin/tag/addtagusers"
	urlDeleteTagUser = "/cgi-bin/tag/deltagusers"
)

type TagInfo struct {
	TagId   int    `json:"tagid"`   // 标签id
	TagName string `json:"tagname"` // 标签名
}

type UserInfo struct {
	UserId string `json:"userid"` // 成员帐号
	Name   string `json:"name"`   // 成员名
}

type WechatTag struct {
	Client *ent.WechatClient
}

func NewWechatTag(client *ent.WechatClient) *WechatTag {
	return &WechatTag{client}
}

// 创建标签
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90210
func (w WechatTag) Create(tagId int, tagName string) (id int, err error) {
	data := struct {
		// 标签名称，长度限制为32个字以内（汉字或英文字母），标签名不可与其他标签重名
		TagId int `json:"tagid,omitempty"`

		// 标签id，非负整型，指定此参数时新增的标签会生成对应的标签id，
		// 不指定时则以目前最大的id自增
		TagName string `json:"tagname"`
	}{
		TagId:   tagId,
		TagName: tagName,
	}

	out := struct {
		TagId int `json:"tagid"`
	}{}
	err = w.Client.Post(urlCreateTag, nil, data, nil, &out)
	id = out.TagId
	return
}

// 删除标签
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90212
func (w WechatTag) Delete(tagId int) error {
	values := url.Values{}
	values.Add("tagid", strconv.Itoa(tagId))

	return w.Client.Get(urlDeleteTag, values, nil, nil)
}

// 更新标签名字
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90211
func (w WechatTag) Update(tagId int, tagName string) error {
	data := struct {
		// 标签ID
		TagId int `json:"tagid"`
		// 标签名称，长度限制为32个字（汉字或英文字母），标签不可与其他标签重名
		TagName string `json:"tagname"`
	}{
		TagId:   tagId,
		TagName: tagName,
	}

	return w.Client.Post(urlUpdateTag, nil, data, nil, nil)
}

// 获取标签列表
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90216
func (w WechatTag) GetTagList() (tagList []TagInfo, err error) {
	out := struct {
		TagList []TagInfo `json:"taglist"` // 标签列表
	}{}

	err = w.Client.Get(urlGetTagList, nil, nil, &out)
	tagList = out.TagList
	return
}

// 获取标签成员
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90213
func (w WechatTag) GetUserList(tagId int) (tagName string, userList []UserInfo, partyList []int, err error) {
	values := url.Values{}
	values.Add("tagid", strconv.Itoa(tagId))

	out := struct {
		TagName   string     `json:"tagname"`   // 标签名
		UserList  []UserInfo `json:"userlist"`  // 标签中包含的成员列表
		PartyList []int      `json:"partylist"` // 标签中包含的部门id列表
	}{}
	err = w.Client.Get(urlGetTagUser, values, nil, &out)
	tagName = out.TagName
	userList = out.UserList
	partyList = out.PartyList
	return
}

// 增加标签成员
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90214
func (w WechatTag) AddUser(tagId int, userList []string, partyList []int) error {
	data := struct {
		// 标签ID
		TagId int `json:"tagid"`

		// 企业成员ID列表，
		// 注意：userlist、partylist不能同时为空，单次请求长度不超过1000
		UserList []string `json:"userlist,omitempty"`

		// 企业部门ID列表，
		// 注意：userlist、partylist不能同时为空，单次请求长度不超过100
		PartyList []int `json:"partylist,omitempty"`
	}{
		TagId:     tagId,
		UserList:  userList,
		PartyList: partyList,
	}

	return w.Client.Post(urlAddTagUser, nil, data, &TagError{}, nil)
}

// 删除标签成员
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90215
func (w WechatTag) DeleteUser(tagId int, userList []string, partyList []int) error {
	data := struct {
		TagId     int      `json:"tagid"`     // 标签ID
		UserList  []string `json:"userlist"`  // 企业成员ID列表，注意：userlist、partylist不能同时为空
		PartyList []int    `json:"partylist"` // 企业部门ID列表，注意：userlist、partylist不能同时为空
	}{
		TagId:     tagId,
		UserList:  userList,
		PartyList: partyList,
	}

	return w.Client.Post(urlDeleteTagUser, nil, data, &TagError{}, nil)
}
