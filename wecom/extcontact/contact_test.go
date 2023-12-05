// 外部联系人管理
package extcontact

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/huimingz/wechatgo/testdata"
	"github.com/huimingz/wechatgo/wecom"
)

type wechatContactTestSuite struct {
	suite.Suite
	contact *WechatContact
}

func (s *wechatContactTestSuite) SetupSuite() {
	conf := testdata.TestConf
	s.contact = NewWechatContact(wecom.NewClient(conf.CorpId, conf.UserSecret, conf.AgentId))
}

func (s *wechatContactTestSuite) TestGetUserList() {
	userList, err := s.contact.GetUserList(context.Background())
	s.Error(err)
	s.Empty(userList)
}

func TestWechatContact(t *testing.T) {
	suite.Run(t, new(wechatContactTestSuite))
}
