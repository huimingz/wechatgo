package wecom

type DeptInfo struct {
	Id       int    `json:"id"`       // 创建的部门id
	Name     string `json:"name"`     // 部门名称
	ParentId int    `json:"parentid"` // 父亲部门id。根部门为1
	Order    int    `json:"order"`    // 在父部门中的次序值。order值大的排序靠前。值范围是[0, 2^32)
}

type DepartmentUserInfo struct {
	UserId     string `json:"userid"`     // 成员UserID。对应管理端的帐号
	Name       string `json:"name"`       // 成员名称
	Department []int  `json:"department"` // 成员所属部门列表。列表项为部门ID，32位整型
}

type DepartmentDetail struct {
	Id               int      `json:"id"`                // 部门id
	Name             string   `json:"name"`              // 部门名称
	NameEnglish      string   `json:"name_en"`           // 部门英文名称
	DepartmentLeader []string `json:"department_leader"` // 部门领导人userid列表
	ParentId         int      `json:"parentid"`          // 父部门id。根部门为1
	Order            int      `json:"order"`             // 在父部门中的次序值。order值大的排序靠前。值范围是[0, 2^32)
}
