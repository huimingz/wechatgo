package wecom

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/huimingz/wechatgo/testdata"
)

type CorpGroupTestSuite struct {
	TestSuite
	service *CorpGroup
}

func (s *CorpGroupTestSuite) SetupSuite() {
	conf := testdata.TestConf
	wechatClient := NewClient(conf.CorpId, conf.UserSecret, conf.AgentId)
	s.service = newCorpGroup(wechatClient)
}

func (s *CorpGroupTestSuite) TestAppShareInfos() {
	infos, err := s.service.AppShareInfos(context.Background(), 123, nil)

	s.NoError(err)
	s.NotEmpty(infos)
}

func TestCorpGroupTestSuite(t *testing.T) {
	suite.Run(t, new(CorpGroupTestSuite))
}
