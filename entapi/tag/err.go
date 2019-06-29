package tag

import (
	"github.com/huimingz/wechatgo/client"
)

type TagError struct {
	client.WXMsgError
	InvalidList  string `json:"invalidlist"`  // 非法的成员帐号列表
	InvalidParty []int  `json:"invalidparty"` // 非法的部门id列表
}
