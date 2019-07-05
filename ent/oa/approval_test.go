package oa

import (
	"testing"

	"github.com/huimingz/wechatgo/ent"
	"github.com/huimingz/wechatgo/testdata"
)

var wechatApproval *WechatApproval

func TestWechatApproval_GetApprovalData(t *testing.T) {
	data, err := wechatApproval.GetApprovalData(1559048367, 1561640367, 0)
	if err != nil {
		t.Errorf("WechatApproval.GetApprovalData() error = '%s'", err)
	}
	if data.Count == 0 {
		t.Error("WechatApproval.GetApprovalData() error = '审批数据为空'")
	}

}

func init() {
	var conf = testdata.TestConf
	var wechatClient = ent.NewWechatClient(conf.CorpId, conf.ApprovalSecret, conf.ApprovalAgentId, nil, 0, nil, nil, nil)
	wechatApproval = NewWechatApproval(wechatClient)
}
