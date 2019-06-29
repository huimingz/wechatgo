package user

import (
	"github.com/huimingz/wechatgo/client"
)

type InvitError struct {
	client.WXMsgError
	InvalidUser  []string `json:"invaliduser,omitempty"`  // 非法成员列表
	InvalidParty []int    `json:"invalidparty,omitempty"` // 非法部门列表
	InvalidTag   []int    `json:"invalidtag,omitempty"`   // 非法标签列表
}
