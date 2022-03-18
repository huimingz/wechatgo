// 企业微信给出的文档存在错误，返回结果中存在列表和字典混用的情况，加之返回结果中出现了
// 字段说明中未说明的字段。所以不能保证返回结果的正确性。

// Package oa OA数据接口
package oa

import (
	"context"

	"github.com/huimingz/wechatgo/wecom"
)

const urlGetDialRecord = "/cgi-bin/dial/get_dial_record"

type DialCaller struct {
	UserId   string `json:"userid"`   // 主叫用户的userid
	Duration int    `json:"duration"` // 主叫用户的通话时长
}

type DialCallee struct {
	Phone    string `json:"phone"`    // 被叫用户的号码，当被叫用户为外部用户时返回
	UserId   string `json:"userid"`   // 被叫用户的userid，当被叫用户为企业内用户时返回
	Duration int    `json:"duration"` // 被叫用户的通话时长
}

type DialRecord struct {
	CallTime      int          `json:"call_time"`      // 拨出时间
	TotalDuration int          `json:"total_duration"` // 总通话时长，单位为分钟
	CallType      int          `json:"call_type"`      // 通话类型，1-单人通话 2-多人通话
	Caller        DialCaller   `json:"caller"`         // 主叫
	Callee        []DialCallee `json:"callee"`         // 被叫
}

type WechatDial struct {
	Client *wecom.WechatClient
}

func NewWechatDial(client *wecom.WechatClient) *WechatDial {
	return &WechatDial{client}
}

// 获取公费电话拨打记录
//
// 企业可通过此接口，按时间范围拉取成功接通的公费电话拨打记录。
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90267
func (w WechatDial) GetRecord(ctx context.Context, startTime, endTime, offset, limit int) ([]DialRecord, error) {
	data := struct {
		StartTime int `json:"start_time,omitempty"` // 查询的起始时间戳
		EndTime   int `json:"end_time,omitempty"`   // 查询的结束时间戳
		Offset    int `json:"offset,omitempty"`     // 分页查询的偏移量
		Limit     int `json:"limit,omitempty"`      // 分页查询的每页大小,默认为100条，如该参数大于100则按100处理
	}{
		StartTime: startTime,
		EndTime:   endTime,
		Offset:    offset,
		Limit:     limit,
	}

	out := struct {
		Record []DialRecord `json:"record"`
	}{}

	err := w.Client.Post(ctx, urlGetDialRecord, nil, data, nil, &out)
	return out.Record, err
}
