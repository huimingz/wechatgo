// Package app 应用管理
package app

import (
	"context"
	"net/url"
	"strconv"

	"github.com/huimingz/wechatgo/wecom"
)

const (
	urlGetApp     = "/cgi-bin/agent/get"
	urlGetAllApp  = "/cgi-bin/agent/list"
	urlCreateApp  = "/cgi-bin/agent/set"
	urlCreateMenu = "/cgi-bin/menu/create"
	urlGetMenu    = "/cgi-bin/menu/get"
	urlDelMenu    = "/cgi-bin/menu/delete"
)

type UserInfo struct {
	UserId string `json:"userid"` // 用户id
}

type UserInfos struct {
	User []UserInfo `json:"user"` // userid列表
}

type Partys struct {
	PartyId []int `json:"partyid"` // 部门id列表
}

type Tags struct {
	TagId []int `json:"tagid"` // 标签id列表
}

type AppDetail struct {
	AppIntro
	Description        string    `json:"description"`          // 企业应用详情
	AllowUserInfos     UserInfos `json:"allow_userinfos"`      // 企业应用可见范围（人员），其中包括userid
	AllowPartys        Partys    `json:"allow_partys"`         // 企业应用可见范围（部门）
	AllowTags          Tags      `json:"allow_tags"`           // 企业应用可见范围（标签）
	Close              int       `json:"close"`                // 企业应用是否被停用
	RedirectDomain     string    `json:"redirect_domain"`      // 企业应用可信域名
	ReportLocationFlag int       `json:"report_location_flag"` // 企业应用是否打开地理位置上报 0：不上报；1：进入会话上报；
	IsResportEnter     int       `json:"isreportenter"`        // 是否上报用户进入应用事件。0：不接收；1：接收
	HomeUrl            string    `json:"home_url"`             // 应用主页url
}

type AppIntro struct {
	AgentId       int    `json:"agentid"`         // 企业应用id
	Name          string `json:"name"`            // 企业应用名称
	SquareLogoUrl string `json:"square_logo_url"` // 企业应用方形头像
}

type AppInfo struct {
	AgentId            int    `json:"agentid"`                        // 企业应用的id（必须）
	ReportLocationFlag int    `json:"report_location_flag,omitempty"` // 企业应用是否打开地理位置上报 0：不上报；1：进入会话上报；
	LogoMediaId        string `json:"logo_mediaid,omitempty"`         // 企业应用头像的mediaid，通过素材管理接口上传图片获得mediaid，上传后会自动裁剪成方形和圆形两个头像
	Name               string `json:"name,omitempty"`                 // 企业应用名称，长度不超过32个utf8字符
	Description        string `json:"description,omitempty"`          // 企业应用详情，长度为4至120个utf8字符
	RedirectDomain     string `json:"redirect_domain,omitempty"`      // 企业应用可信域名。注意：域名需通过所有权校验，否则jssdk功能将受限，此时返回错误码85005
	IsReportEnter      int    `json:"isreportenter,omitempty"`        // 是否上报用户进入应用事件。0：不接收；1：接收
	HomeUrl            string `json:"home_url,omitempty"`             // 应用主页url。url必须以http或者https开头（为了提高安全性，建议使用https）
}

type Button struct {
	Type      string   `json:"type"`                 // 菜单的响应动作类型（必须）
	Name      string   `json:"name"`                 // 菜单的名字。不能为空，主菜单不能超过16字节，子菜单不能超过40字节（必须）
	Url       string   `json:"url,omitempty"`        // 网页链接，成员点击菜单可打开链接，不超过1024字节。为了提高安全性，建议使用https的url（view类型必须）
	Key       string   `json:"key,omitempty"`        // 菜单KEY值，用于消息接口推送，不超过128字节（click等点击类型必须）
	SubButton []Button `json:"sub_button,omitempty"` // 二级菜单数组，个数应为1~5个
}

type Menu struct {
	Button []Button `json:"button"` // 一级菜单数组，个数应为1~3个（必须）
}

type WechatAppManage struct {
	Client *wecom.WechatClient
}

func NewWechatAppManage(client *wecom.WechatClient) *WechatAppManage {
	return &WechatAppManage{client}
}

// 获取access_token对应的应用列表
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90227
func (w WechatAppManage) GetAllApp(ctx context.Context) ([]AppIntro, error) {
	agentList := struct {
		AgentList []AppIntro `json:"agentlist"`
	}{}
	err := w.Client.Get(ctx, urlGetAllApp, nil, nil, &agentList)
	return agentList.AgentList, err
}

// 获取指定的应用详情
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90227
func (w WechatAppManage) GetApp(ctx context.Context, agentId int) (*AppDetail, error) {
	values := url.Values{}
	if agentId == 0 {
		agentId = w.Client.AgentId
	}
	values.Add("agentid", strconv.Itoa(agentId))

	appDetail := AppDetail{}
	err := w.Client.Get(ctx, urlGetApp, values, nil, &appDetail)
	return &appDetail, err
}

// 设置应用
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90228
func (w WechatAppManage) CreateApp(ctx context.Context, appInfo AppInfo) error {
	return w.Client.Post(ctx, urlCreateApp, nil, appInfo, nil, nil)
}

// 创建菜单
//
// 参考资料：https://work.weixin.qq.com/api/doc#90000/90135/90230
func (w WechatAppManage) CreateMenu(ctx context.Context, menu Menu, agentId int) error {
	values := url.Values{}
	if agentId == 0 {
		agentId = w.Client.AgentId
	}
	if agentId == 0 {
		agentId = w.Client.AgentId
	}
	values.Add("agentid", strconv.Itoa(agentId))

	return w.Client.Post(ctx, urlCreateMenu, values, menu, nil, nil)
}

// 获取菜单
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90232
func (w WechatAppManage) GetMenu(ctx context.Context, agentId int) (*Menu, error) {
	values := url.Values{}
	if agentId == 0 {
		agentId = w.Client.AgentId
	}
	values.Add("agentid", strconv.Itoa(agentId))

	menu := Menu{}
	err := w.Client.Get(ctx, urlGetMenu, values, nil, &menu)
	return &menu, err
}

// 删除菜单
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90233
func (w WechatAppManage) DeleteMenu(ctx context.Context, agentId int) error {
	values := url.Values{}
	if agentId == 0 {
		agentId = w.Client.AgentId
	}
	values.Add("agentid", strconv.Itoa(agentId))

	return w.Client.Get(ctx, urlDelMenu, values, nil, nil)
}
