package oa

import (
	"context"
	"testing"

	"github.com/huimingz/wechatgo/testdata"
	"github.com/huimingz/wechatgo/wecom"
)

var wechatApproval *WechatApproval

func TestWechatApproval_GetApprovalData(t *testing.T) {
	data, err := wechatApproval.GetApprovalData(context.Background(), 1559048367, 1561640367, 0)
	if err != nil {
		t.Errorf("WechatApproval.GetApprovalData() error = '%s'", err)
	}
	if data.Count == 0 {
		t.Error("WechatApproval.GetApprovalData() error = '审批数据为空'")
	}

}

func init() {
	var conf = testdata.TestConf
	var wechatClient = wecom.NewWechatClient(conf.CorpId, conf.ApprovalSecret, conf.ApprovalAgentId)
	wechatApproval = NewWechatApproval(wechatClient)
}
