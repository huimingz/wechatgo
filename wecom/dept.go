// Package dept 部门管理
package wecom

import (
	"context"
	"net/url"
	"strconv"
)

const (
	urlGetDepartments    = "/cgi-bin/department/list"
	urlGetDepartment     = "/cgi-bin/department/get"
	urlGetUserList       = "/cgi-bin/user/simplelist"
	urlGetUserDetailList = "/cgi-bin/user/list"
	urlCreateDepartment  = "/cgi-bin/department/create"
	urlUpdateDepartment  = "/cgi-bin/department/update"
	urlRemoveDepartment  = "/cgi-bin/department/delete"
	urlDeptSimpleList    = "/cgi-bin/department/simplelist"
)

type WechatDept struct {
	Client *WechatClient
}

func NewWechatDept(client *WechatClient) *WechatDept {
	return &WechatDept{client}
}

// Get 获取指定部门及其下的子部门
//
// 权限说明：
// 只能拉取token对应的应用的权限范围内的部门列表
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90208
func (w WechatDept) Get(ctx context.Context, deptId int) ([]DeptInfo, error) {
	values := url.Values{}
	values.Add("id", strconv.Itoa(deptId))

	out := struct {
		Department []DeptInfo `json:"department"`
	}{}
	err := w.Client.Get(ctx, urlGetDepartments, values, nil, &out)
	return out.Department, err
}

// GetDetail 获取部门详情
func (w WechatDept) GetDetail(ctx context.Context, deptId int) (*DepartmentDetail, error) {
	values := url.Values{}
	values.Add("id", strconv.Itoa(deptId))

	out := struct {
		Department DepartmentDetail `json:"department"`
	}{}
	err := w.Client.Get(ctx, urlGetDepartment, values, nil, &out)
	return &out.Department, err
}

// GetList 获取部门列表
//
// 权限说明：
// 只能拉取token对应的应用的权限范围内的部门列表
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90208
func (w WechatDept) GetList(ctx context.Context) ([]DeptInfo, error) {
	out := struct {
		Department []DeptInfo `json:"department"`
	}{}
	err := w.Client.Get(ctx, urlGetDepartments, nil, nil, &out)
	return out.Department, err
}

func (w WechatDept) GetSubList(ctx context.Context, id int) ([]DeptInfo, error) {
	values := url.Values{}
	values.Add("id", strconv.Itoa(id))

	out := struct {
		Department []DeptInfo `json:"department"`
	}{}
	err := w.Client.Get(ctx, urlDeptSimpleList, values, nil, &out)
	return out.Department, err
}

// GetUserList 获取部门成员
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90200
func (w WechatDept) GetUserList(ctx context.Context, deptId int, fetchChild bool) ([]AppUserInfo, error) {
	values := url.Values{}
	values.Add("department_id", strconv.Itoa(deptId))
	if fetchChild {
		values.Add("fetch_child", "1")
	} else {
		values.Add("fetch_child", "0")
	}

	out := struct {
		UserList []AppUserInfo `json:"userlist"`
	}{}

	err := w.Client.Get(ctx, urlGetUserList, values, nil, &out)
	return out.UserList, err
}

// 获取部门成员详情
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90201
func (w WechatDept) GetUserDetailList(ctx context.Context, deptId int, fetchChild bool) ([]UserInfo, error) {
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
	err := w.Client.Get(ctx, urlGetUserDetailList, values, nil, &out)
	return out.UserList, err
}

// 创建部门
//
// 注意，部门的最大层级为15层；部门总数不能超过3万个；每个部门下的节点不能超过3万个。
// 建议保证创建的部门和对应部门成员是串行化处理。
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90205
func (w WechatDept) Create(ctx context.Context, name string, parentId, order, deptId int) (id int, err error) {
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
	err = w.Client.Post(ctx, urlCreateDepartment, nil, data, nil, &out)
	return out.Id, err
}

// Update 更新部门
//
// 注意，部门的最大层级为15层；部门总数不能超过3万个；每个部门下的节点不能超过3万个。
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90206
func (w WechatDept) Update(ctx context.Context, deptId, parentId, order int, name string) error {
	data := struct {
		Id       int    `json:"id"`
		Name     string `json:"name,omitempty"`
		ParentId int    `json:"parentid,omitempty"`
		Order    int    `json:"order,omitempty"`
	}{
		Id:       deptId,
		Name:     name,
		ParentId: parentId,
		Order:    order,
	}

	return w.Client.Post(ctx, urlUpdateDepartment, nil, data, nil, nil)
}

// Delete 删除部门
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90207
func (w WechatDept) Delete(ctx context.Context, deptId int) error {
	values := url.Values{}
	values.Add("id", strconv.Itoa(deptId))

	return w.Client.Get(ctx, urlRemoveDepartment, values, nil, nil)
}
