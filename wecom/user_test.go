package wecom

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/huimingz/wechatgo/testdata"
)

type userTestSuite struct {
	TestSuite
	user *UserManager
}

func (s *userTestSuite) SetupSuite() {
	conf := testdata.TestConf
	wechatClient := NewClient(conf.CorpId, conf.UserSecret, conf.AgentId)
	s.user = newUserManager(wechatClient)
}

func (s *userTestSuite) TestCreateUser() {
	user := UserForCreate{}
	user.UserId = "example_xxx"
	user.Name = "xxx"
	user.Department = []int{1}
	user.Email = "xxx@example.com"
	user.ToInvite = false

	err := s.user.CreateUser(context.Background(), user)

	s.NoError(err)
}

func (s *userTestSuite) TestGetUser() {
	user, err := s.user.GetUser(context.Background(), "example_xxx")

	s.NoError(err)
	s.Equal("example_xxx", user.UserId)
}

func (s *userTestSuite) TestUpdateUser() {
	user := UserForUpdate{
		UserId: "example_xxx",
		Name:   "new_xxx",
	}

	err := s.user.UpdateUser(context.Background(), user)
	s.NoError(err)

	info, err := s.user.GetUser(context.Background(), user.UserId)

	s.NoError(err)
	s.Equal(user.Name, info.Name)
}

func (s *userTestSuite) TestInvite() {
	err := s.user.Invite(context.Background(), []string{"example_xxx"}, nil, nil)

	s.NoError(err)
}

func (s *userTestSuite) TestVerify() {
	err := s.user.Verify(context.Background(), "example_xxx")

	s.NoError(err)
}

func (s *userTestSuite) TestUserId2OpenId() {
	openId, err := s.user.UserId2OpenId(context.Background(), testdata.TestConf.UserId)

	s.NoError(err)
	s.Equal(testdata.TestConf.OpenId, openId)
}

func (s *userTestSuite) TestDeleteUser() {
	err := s.user.DeleteUser(context.Background(), "example_xxx")
	s.NoError(err)

	err = s.user.DeleteUser(context.Background(), "example_xxx")
	s.NoError(err)
}

func (s *userTestSuite) TestDeleteUsers() {
	user := UserForCreate{}
	user.UserId = "example_xxx"
	user.Name = "xxx"
	user.Email = "xxx@example.com"
	user.Department = []int{1}
	user.ToInvite = false

	err := s.user.CreateUser(context.Background(), user)
	s.NoError(err)

	err = s.user.DeleteUsers(context.Background(), []string{user.UserId})
	s.NoError(err)
}

func (s *userTestSuite) TestOpenId2UserId() {
	userId, err := s.user.OpenId2UserId(context.Background(), testdata.TestConf.OpenId)

	s.NoError(err)
	s.Equal(testdata.TestConf.UserId, userId)
}

func (s *userTestSuite) TestGetJoinQRCode() {
	qrCode, err := s.user.GetJoinQRCode(context.Background(), 0)

	s.NoError(err)
	s.NotEmpty(qrCode)
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(userTestSuite))
}
