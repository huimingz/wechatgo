package wecom

import (
	"context"
)

type CorpGroup struct {
	client *Client
}

func newCorpGroup(client *Client) *CorpGroup {
	return &CorpGroup{client: client}
}

func (g *CorpGroup) AppShareInfos(ctx context.Context, agentId int, option *corpGroupListOption) ([]CorpGroupData, error) {
	if option == nil {
		option = NewCorpGroupListOption()
	}
	payload := struct {
		AgentId      int                    `json:"agentid"`
		BusinessType *corpGroupBusinessType `json:"business_type,omitempty"`
		Corpid       *string                `json:"corpid,omitempty"`
		Limit        *uint64                `json:"limit,omitempty"`
		Cursor       *string                `json:"cursor,omitempty"`
	}{
		AgentId:      agentId,
		BusinessType: option.BusinessType,
		Corpid:       option.CorpId,
		Limit:        option.Limit,
		Cursor:       option.Cursor,
	}
	uri := "/cgi-bin/corpgroup/corp/list_app_share_info"
	var out struct {
		CorpList []CorpGroupData `json:"corp_list"`
	}
	if err := g.client.Post(ctx, uri, nil, payload, nil, &out); err != nil {
		return nil, err
	}
	return out.CorpList, nil
}

type CorpGroupData struct {
	CorpId   string `json:"corpid"`    // 企业ID
	AgentId  int    `json:"agentid"`   // 应用ID
	CorpName string `json:"corp_name"` // 企业名称
}

const (
	CorpGroupBusinessTypeInter    corpGroupBusinessType = 0 // 企业互联/局校互联
	CorpGroupBusinessTypeUpstream corpGroupBusinessType = 1 // 上下游企业
)

type corpGroupBusinessType int

type corpGroupListOption struct {
	BusinessType *corpGroupBusinessType `json:"business_type"` // 企业类型
	CorpId       *string                `json:"corpid"`        // 企业ID
	Limit        *uint64                `json:"limit"`         // 分页大小
	Cursor       *string                `json:"cursor"`        // 分页游标
}

func NewCorpGroupListOption() *corpGroupListOption {
	return &corpGroupListOption{}
}

func (o *corpGroupListOption) SetBusinessType(businessType corpGroupBusinessType) *corpGroupListOption {
	o.BusinessType = &businessType
	return o
}

func (o *corpGroupListOption) SetCorpId(corpId string) *corpGroupListOption {
	o.CorpId = &corpId
	return o
}

func (o *corpGroupListOption) SetLimit(limit uint64) *corpGroupListOption {
	o.Limit = &limit
	return o
}

func (o *corpGroupListOption) SetCursor(cursor string) *corpGroupListOption {
	o.Cursor = &cursor
	return o
}
