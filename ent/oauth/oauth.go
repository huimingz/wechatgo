// Package oauth 身份认证
package oauth

import (
	"net/url"
	"strconv"

	"github.com/huimingz/wechatgo/ent"
)

const (
	baseWebOAuthUrl = "https://open.weixin.qq.com/connect/oauth2/authorize"
	baseQRCodeUrl   = "https://open.work.weixin.qq.com/wwopen/sso/qrConnect"
	userInfoUrl     = "/cgi-bin/user/getuserinfo"
)

type UserInfo struct {
	UserId   string `json:"UserId"`   // 成员UserID
	DeviceId string `json:"DeviceId"` // 手机设备号
}

type WechatOAuth struct {
	Client *ent.WechatClient
}

func NewWechatOAuth(client *ent.WechatClient) *WechatOAuth {
	return &WechatOAuth{Client: client}
}

func (w WechatOAuth) WebOAuthUrl(redirectUrl, responseType, state string) string {
	if redirectUrl == "" {
		return ""
	}
	if state == "" {
		state = "state"
	}
	if responseType == "" {
		responseType = "code"
	}

	var oauthUrl string
	values := url.Values{}
	values.Add("appid", w.Client.CorpId)
	values.Add("redirect_url", redirectUrl)
	values.Add("response_type", "code")
	values.Add("scope", "snsapi_base")
	values.Add("state", state)

	oauthUrl = baseWebOAuthUrl + "?" + values.Encode() + "#wechat_redirect"
	return oauthUrl
}

func (w WechatOAuth) QRCodeUrl(redirectUrl, state string) string {
	if redirectUrl == "" {
		return ""
	}
	if state == "" {
		state = "state"
	}

	var qrCodeUrl string
	values := url.Values{}
	values.Add("appid", w.Client.CorpId)
	values.Add("agentid", strconv.Itoa(w.Client.AgentId))
	values.Add("state", state)
	values.Add("redirect_uri", redirectUrl)

	qrCodeUrl = baseQRCodeUrl + "?" + values.Encode()
	return qrCodeUrl
}

func (w WechatOAuth) GetUserInfo(code string) (*UserInfo, error) {
	values := url.Values{}
	values.Add("code", code)

	userInfo := UserInfo{}
	err := w.Client.Get(userInfoUrl, values, nil, &userInfo)
	return &userInfo, err
}
