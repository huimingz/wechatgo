package wecom

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/huimingz/wechatgo/testdata"
)

var wechatClient = NewClient(testdata.TestConf.CorpId, testdata.TestConf.CorpSecret, testdata.TestConf.AgentId)

type WechatClientSuite struct {
	suite.Suite
	client *Client
}

func (s *WechatClientSuite) SetupSuite() {
	s.client = NewClient(
		testdata.TestConf.CorpId,
		testdata.TestConf.CorpSecret,
		testdata.TestConf.AgentId,
	)
}

func (s *WechatClientSuite) TestShouldBeExpiredIfNotRefreshAccessToken() {
	client := NewClient("", "", 0)

	s.True(client.IsExpired(context.Background()))
}

func (s *WechatClientSuite) TestShouldBeNotExpiredIfRefreshAccessToken() {
	s.client.GetAccessToken(context.Background())

	s.False(s.client.IsExpired(context.Background()))
}

func (s *WechatClientSuite) TestShouldGetAccessTokenSuccessfully() {
	accessToken, err := s.client.GetAccessToken(context.Background())

	s.NoError(err)
	s.NotEmpty(accessToken)
}

func (s *WechatClientSuite) TestShouldRaiseErrorIfInvalidClientWhenGetAccessToken() {
	client := NewClient("xxx", testdata.TestConf.CorpSecret, testdata.TestConf.AgentId)
	accessToken, err := client.GetAccessToken(context.Background())

	s.Error(err)
	s.Empty(accessToken)
}

func (s *WechatClientSuite) TestShoulFetchAccessTokenSuccessfully() {
	err := s.client.FetchAccessToken(context.Background())

	s.NoError(err)
}

func (s *WechatClientSuite) TestShouldRaiseErrorIfInvalidClientWhenFetchAccessToken() {
	client := NewClient("xxx", testdata.TestConf.CorpSecret, testdata.TestConf.AgentId)
	err := client.FetchAccessToken(context.Background())

	s.Error(err)
}

func (s *WechatClientSuite) TestShouldGetAccessTokenStorageKeySuccessfully() {
	key := s.client.GetAccessTokenStorageKey()
	expect := "accesstoken_" + s.client.CorpSecret

	s.Equal(key, expect)
}

func TestWechatClientSuite(t *testing.T) {
	suite.Run(t, new(WechatClientSuite))
}
