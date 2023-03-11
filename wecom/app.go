// Package app 应用管理
package wecom

import (
	"context"
	"net/url"
	"strconv"
)

const (
	urlGetApp     = "/cgi-bin/agent/get"
	urlGetAllApp  = "/cgi-bin/agent/list"
	urlCreateApp  = "/cgi-bin/agent/set"
	urlCreateMenu = "/cgi-bin/menu/create"
	urlGetMenu    = "/cgi-bin/menu/get"
	urlDelMenu    = "/cgi-bin/menu/delete"
)

type WechatAppManage struct {
	Client *WechatClient
}

func NewWechatAppManage(client *WechatClient) *WechatAppManage {
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
