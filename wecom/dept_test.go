package wecom

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/huimingz/wechatgo/testdata"
)

type departmentTestSuite struct {
	TestSuite
	dept *WechatDept
}

func (s *departmentTestSuite) SetupSuite() {
	s.TestSuite.SetupSuite()

	conf := testdata.TestConf
	s.dept = NewWechatDept(NewClient(conf.CorpId, conf.CorpSecret, conf.AgentId, ClientWithHTTPClient(s.httpClient)))
}

func (s *departmentTestSuite) TestShouldGetDepartment() {
	s.registerResponder(http.MethodGet, urlGetDepartment)

	dept, err := s.dept.GetDetail(context.Background(), 2)

	s.NoError(err)
	s.Equal(2, dept.Id)
}

func (s *departmentTestSuite) TestShouldGetDepartments() {
	s.registerResponder(http.MethodGet, urlGetDepartments)

	departments, err := s.dept.GetList(context.Background())

	s.NoError(err)
	s.NotEmpty(departments)
}

func (s *departmentTestSuite) TestShouldRemoveDepartment() {
	s.registerSuccessResponder(http.MethodGet, urlRemoveDepartment)

	err := s.dept.Delete(context.Background(), 1)

	s.NoError(err)
}

func (s *departmentTestSuite) TestShouldCreateDepartment() {
	s.registerResponder(http.MethodPost, urlCreateDepartment)

	deptId, err := s.dept.Create(context.Background(), "RDGZ", 1, 1, 1)

	s.NoError(err)
	s.Equal(2, deptId)
}

func (s *departmentTestSuite) TestShouldUpdateDepartment() {
	s.registerResponder(http.MethodPost, urlUpdateDepartment)

	err := s.dept.Update(context.Background(), 2, 1, 1, "RDGZ")

	s.NoError(err)
}

var wechatDept *WechatDept

func TestWechatDept_GetUserList(t *testing.T) {
	users, err := wechatDept.GetUserList(context.Background(), 1, false)
	if err != nil {
		t.Errorf("WechatDept.GetUserList() error = '%s'", err)
	}
	if len(users) == 0 {
		t.Error("WechatDept.GetUserList() error = '返回部门成员为空'")
	}

	cusers, err := wechatDept.GetUserList(context.Background(), 1, true)
	if err != nil {
		t.Errorf("WechatDept.GetUserList() error = '%s'", err)
	}
	if len(cusers) <= len(users) {
		t.Error("WechatDept.GetUserList() error = '返回部门和子部门成员数少于指定部门成员数'")
	}
}

func TestWechatDept_GetUserDetailList(t *testing.T) {
	users, err := wechatDept.GetUserDetailList(context.Background(), 1, false)
	if err != nil {
		t.Errorf("WechatDept.GetUserDetailList() error = '%s'", err)
	}
	if len(users) == 0 {
		t.Error("WechatDept.GetUserDetailList() error = '返回部门成员详情为空'")
	}

	cusers, err := wechatDept.GetUserDetailList(context.Background(), 1, true)
	if err != nil {
		t.Errorf("WechatDept.GetUserDetailList() error = '%s'", err)
	}
	if len(cusers) <= len(users) {
		t.Error("WechatDept.GetUserList() error = '返回部门和子部门成员数少于指定部门成员数'")
	}
}

func TestWechatDept_Update(t *testing.T) {
	err := wechatDept.Update(context.Background(), 18, 0, 0, "test_test")
	if err != nil {
		t.Errorf("WechatDept.Update() error = '%s'", err)
	}

	infos, err := wechatDept.Get(context.Background(), 18)
	if err != nil {
		t.Errorf("WechatDept.Update() error = '%s'", err)
	}
	if infos[0].Name != "test_test" {
		t.Errorf("WechatDept.Update() error = '更新部门名称失败'")
	}
}

func init() {
	var conf = testdata.TestConf
	var wechatClient = NewClient(conf.CorpId, conf.UserSecret, conf.AgentId)
	wechatDept = NewWechatDept(wechatClient)
}

func TestDepartment(t *testing.T) {
	suite.Run(t, new(departmentTestSuite))
}
