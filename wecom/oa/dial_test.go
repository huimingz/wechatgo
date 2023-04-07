// 企业微信给出的文档存在错误，返回结果中存在列表和字典混用的情况，加之返回结果中出现了
// 字段说明中未说明的字段。所以不能保证返回结果的正确性。

package oa

import (
	"context"
	"testing"

	"github.com/huimingz/wechatgo/testdata"
	"github.com/huimingz/wechatgo/wecom"
)

var wechatDial *WechatDial

func TestWechatDial_GetRecord(t *testing.T) {
	record, err := wechatDial.GetRecord(context.Background(), 1559048367, 1561640367, 0, 0)
	if err != nil {
		t.Errorf("WechatDial.GetRecord() error = '%s'", err)
	}
	if len(record) != 0 {
		t.Error("WechatDial.GetRecord() error = '返回记录非空'")
	}
}

func init() {
	var conf = testdata.TestConf
	var wechatClient = wecom.NewClient(conf.CorpId, conf.DialSecret, conf.DialAgentId)
	wechatDial = NewWechatDial(wechatClient)
}
