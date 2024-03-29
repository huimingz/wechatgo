package oa

import (
	"context"
	"testing"

	"github.com/huimingz/wechatgo/testdata"
	"github.com/huimingz/wechatgo/wecom"
)

var wechatCheckin *WechatCheckin

func TestWechatCheckin_GetData(t *testing.T) {
	data, err := wechatCheckin.GetData(context.Background(), 3, 1559048367, 1561640367, []string{testdata.TestConf.UserId})
	if err != nil {
		t.Errorf("WechatCheckin.GetData() error = '%s'", err)
	}
	if len(data) != 0 {
		t.Error("WechatCheckin.GetData() error = '返回打卡数据非空'")
	}
}

func TestWechatCheckin_GetOption(t *testing.T) {
	_, err := wechatCheckin.GetOption(context.Background(), 1511971200, []string{testdata.TestConf.UserId})
	if err == nil {
		t.Error("WechatCheckin.GetOption() error = '未出现错误'")
	}
}

func init() {
	var conf = testdata.TestConf
	var wechatClient = wecom.NewClient(conf.CorpId, conf.CheckinSecret, conf.CheckinAgentId)
	wechatCheckin = NewWechatCheckin(wechatClient)
}
