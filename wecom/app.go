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

type applicationManager struct {
	client *Client
}

func newWechatAppManage(client *Client) *applicationManager {
	return &applicationManager{client}
}

// 获取access_token对应的应用列表
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90227
func (w *applicationManager) GetAllApp(ctx context.Context) ([]AppIntro, error) {
	agentList := struct {
		AgentList []AppIntro `json:"agentlist"`
	}{}
	err := w.client.Get(ctx, urlGetAllApp, nil, nil, &agentList)
	return agentList.AgentList, err
}

// 获取指定的应用详情
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90227
func (w *applicationManager) GetApp(ctx context.Context, agentId int) (*AppDetail, error) {
	values := url.Values{}
	values.Add("agentid", w.toStringAgentId(w.orClientAgentId(agentId)))

	appDetail := AppDetail{}
	err := w.client.Get(ctx, urlGetApp, values, nil, &appDetail)
	return &appDetail, err
}

// 设置应用
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90228
func (w *applicationManager) CreateApp(ctx context.Context, appInfo AppInfo) error {
	return w.client.Post(ctx, urlCreateApp, nil, appInfo, nil, nil)
}

// 创建菜单
//
// 参考资料：https://work.weixin.qq.com/api/doc#90000/90135/90230
func (w *applicationManager) CreateMenu(ctx context.Context, menu Menu, agentId int) error {
	values := url.Values{}
	values.Add("agentid", w.toStringAgentId(w.orClientAgentId(agentId)))

	return w.client.Post(ctx, urlCreateMenu, values, menu, nil, nil)
}

// 获取菜单
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90232
func (w *applicationManager) GetMenu(ctx context.Context, agentId int) (*Menu, error) {
	values := url.Values{}
	values.Add("agentid", w.toStringAgentId(w.orClientAgentId(agentId)))

	menu := Menu{}
	err := w.client.Get(ctx, urlGetMenu, values, nil, &menu)
	return &menu, err
}

// 删除菜单
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90233
func (w *applicationManager) DeleteMenu(ctx context.Context, agentId int) error {
	values := url.Values{}
	values.Add("agentid", w.toStringAgentId(w.orClientAgentId(agentId)))

	return w.client.Get(ctx, urlDelMenu, values, nil, nil)
}

func (w *applicationManager) orClientAgentId(agentId int) int {
	if agentId == 0 {
		return w.client.AgentId
	}
	return agentId
}

func (w *applicationManager) toStringAgentId(agentId int) string {
	return strconv.Itoa(w.orClientAgentId(agentId))
}
