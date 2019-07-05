package msg

import (
	"github.com/huimingz/wechatgo"
)

type MsgError struct {
	wechatgo.WXMsgError
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
}
