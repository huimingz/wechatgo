package msg

import (
	"github.com/huimingz/wechatgo/client"
)

type MsgError struct {
	client.WXMsgError
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
}
