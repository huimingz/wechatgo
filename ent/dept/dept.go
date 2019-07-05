// Package dept 部门管理
package dept

import (
	"net/url"
	"strconv"

	"github.com/huimingz/wechatgo/ent"
)

const (
	urlDeptGet           = "/cgi-bin/department/list"
	urlGetUserList       = "/cgi-bin/user/simplelist"
	urlGetUserDetailList = "/cgi-bin/user/list"
	urlDeptCreate        = "/cgi-bin/department/create"
	urlDeptUpdate        = "/cgi-bin/department/update"
	urlDeptDelete        = "/cgi-bin/department/delete"
)

type DeptInfo struct {
	Id       int    `json:"id"`       // 创建的部门id
	Name     string `json:"name"`     // 部门名称
	ParentId int    `json:"parentid"` // 父亲部门id。根部门为1
	Order    int    `json:"order"`    // 在父部门中的次序值。order值大的排序靠前。值范围是[0, 2^32)
}

type UserInfo struct {
	UserId     string `json:"userid"`     // 成员UserID。对应管理端的帐号
	Name       string `json:"name"`       // 成员名称
	Department []int  `json:"department"` // 成员所属部门列表。列表项为部门ID，32位整型
}

type WechatDept struct {
	Client *ent.WechatClient
}

func NewWechatDept(client *ent.WechatClient) *WechatDept {
	return &WechatDept{client}
}

// 获取指定部门及其下的子部门
//
// 权限说明：
// 只能拉取token对应的应用的权限范围内的部门列表
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90208
func (w WechatDept) Get(deptId int) ([]DeptInfo, error) {
	values := url.Values{}
	values.Add("id", strconv.Itoa(deptId))

	out := struct {
		Department []DeptInfo `json:"department"`
	}{}
	err := w.Client.Get(urlDeptGet, values, nil, &out)
	return out.Department, err
}

// 获取部门列表
//
// 权限说明：
// 只能拉取token对应的应用的权限范围内的部门列表
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90208
func (w WechatDept) GetList() ([]DeptInfo, error) {
	out := struct {
		Department []DeptInfo `json:"department"`
	}{}
	err := w.Client.Get(urlDeptGet, nil, nil, &out)
	return out.Department, err
}

// 获取部门成员
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90200
func (w WechatDept) GetUserList(deptId int, fetchChild bool) ([]UserInfo, error) {
	values := url.Values{}
	values.Add("department_id", strconv.Itoa(deptId))
	if fetchChild {
		values.Add("fetch_child", "1")
	} else {
		values.Add("fetch_child", "0")
	}

	out := struct {
		UserList []UserInfo `json:"userlist"`
	}{}

	err := w.Client.Get(urlGetUserList, values, nil, &out)
	return out.UserList, err
}

// 获取部门成员详情
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90201
func (w WechatDept) GetUserDetailList(deptId int, fetchChild bool) ([]UserInfo, error) {
	values := url.Values{}
	values.Add("department_id", strconv.Itoa(deptId))
	if fetchChild {
		values.Add("fetch_child", "1")
	} else {
		values.Add("fetch_child", "0")
	}

	out := struct {
		UserList []UserInfo `json:"userlist"`
	}{}
	err := w.Client.Get(urlGetUserDetailList, values, nil, &out)
	return out.UserList, err
}

// 创建部门
//
// 注意，部门的最大层级为15层；部门总数不能超过3万个；每个部门下的节点不能超过3万个。
// 建议保证创建的部门和对应部门成员是串行化处理。
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90205
func (w WechatDept) Create(name string, parentId, order, deptId int) (id int, err error) {
	data := struct {
		Name     string `json:"name"`
		ParentId int    `json:"parentid,omitempty"`
		Order    int    `json:"order,omitempty"`
		Id       int    `json:"id,omitempty"`
	}{
		Name:     name,
		ParentId: parentId,
		Order:    order,
		Id:       deptId,
	}

	out := struct {
		Id int `json:"id"`
	}{}
	err = w.Client.Post(urlDeptCreate, nil, data, nil, &out)
	return out.Id, err
}

// 更新部门
//
// 注意，部门的最大层级为15层；部门总数不能超过3万个；每个部门下的节点不能超过3万个。
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90206
func (w WechatDept) Update(deptId, parentId, order int, name string) error {
	data := struct {
		Id       int    `json:"id"`
		Name     string `json:"name,omitempty"`
		Parentid int    `json:"parentid,omitempty"`
		Order    int    `json:"order,omitempty"`
	}{
		Id:       deptId,
		Name:     name,
		Parentid: parentId,
		Order:    order,
	}

	return w.Client.Post(urlDeptUpdate, nil, data, nil, nil)
}

// 删除部门
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90207
func (w WechatDept) Delete(deptId int) error {
	values := url.Values{}
	values.Add("id", strconv.Itoa(deptId))

	return w.Client.Get(urlDeptDelete, values, nil, nil)
}
