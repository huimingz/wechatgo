package wecom

import (
	"context"
	"testing"

	"github.com/huimingz/wechatgo/testdata"
)

var corpId string = testdata.TestConf.CorpId
var corpSecret string = testdata.TestConf.CorpSecret
var agentId int = testdata.TestConf.AgentId

var wechatClient = NewWechatClient(corpId, corpSecret, agentId)

func TestWechatClient_IsExpired(t *testing.T) {
	client := NewWechatClient("", "", 0)
	if client.IsExpired(context.Background()) != true {
		t.Error("WechatClient.IsExpired() error = 'isExpired != true'")
	}

	_, _ = wechatClient.GetAccessToken(context.Background())
	if wechatClient.IsExpired(context.Background()) != false {
		t.Error("WechatClient.IsExpired() error = 'isExpired != false'")
	}
}

func TestWechatClient_GetAccessToken(t *testing.T) {
	accessToken, _ := wechatClient.GetAccessToken(context.Background())
	if accessToken == "" {
		t.Error("WechatClient.GetAccessToken() 返回结果access_token为空字符串")
	}

	// invalid client
	client := NewWechatClient("xxx", corpSecret, agentId)
	accessToken, _ = client.GetAccessToken(context.Background())
	if accessToken != "" {
		t.Error("WechatClient.GetAccessToken() 返回结果access_token值为非空字符串")
	}
}

func TestWechatClient_FetchAccessToken(t *testing.T) {
	err := wechatClient.FetchAccessToken(context.Background())
	if err != nil {
		t.Errorf("WechatClient.FetchAccessToken() error = %s", err)
	}
	if token, _ := wechatClient.GetAccessToken(context.Background()); token == "" {
		t.Error("WechatClient.FetchAccessToken() accessToken值为空字符串")
	}
}

func TestWechatClient_GetAccessTokenStorageKey(t *testing.T) {
	key := wechatClient.GetAccessTokenStorageKey()
	expect := "accesstoken_" + wechatClient.CorpSecret
	if key != expect {
		t.Errorf("WechatClient.GetAccessTokenStorageKey() result != expect, %s != %s", key, expect)
	}

}
