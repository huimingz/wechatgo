// Package oa OA数据接口
package oa

import (
	"github.com/huimingz/wechatgo/wecom"
)

const urlGetApprovalData string = "/cgi-bin/corp/getapprovaldata"

type ApprovalEntryExpenseItem struct {
	ExpenseItemType int    `json:"expenseitem_type"` // 费用类型：1飞机票；2火车票；3的士费；4住宿费；5餐饮费；6礼品费；7活动费；8通讯费；9补助；10其他
	Time            int    `json:"time"`             // 发生时间，unix时间 (历史单据字段，新申请单据不再提供)
	Sums            int    `json:"sums"`             // 费用金额，单位元 (历史单据字段，新申请单据不再提供)
	Reason          string `json:"reason"`           // 明细事由 (历史单据字段，新申请单据不再提供)
}

type ApprovalEntryExpense struct {
	ExpenseType int                        `json:"expense_type"` // 报销类型：1差旅费；2交通费；3招待费；4其他报销
	Reason      string                     `json:"reason"`       // 报销事由
	Item        []ApprovalEntryExpenseItem `json:"item"`         // 报销明细 (历史单据字段，新申请单据不再提供)
}

type ApprovalEntryComm struct {
	ApplyData string `json:"apply_data"` // 审批申请的单据数据
}

type ApprovalEntryLeave struct {
	TimeUnit  int    `json:"timeunit"`   // 请假时间单位：0半天；1小时
	LeaveType int    `json:"leave_type"` // 请假类型：1年假；2事假；3病假；4调休假；5婚假；6产假；7陪产假；8其他
	StartTime int    `json:"start_time"` // 请假开始时间，unix时间
	EndTime   int    `json:"end_time"`   // 请假结束时间，unix时间
	Duration  int    `json:"duration"`   // 请假时长，单位小时
	Reason    string `json:"reason"`     // 请假事由
}

type ApprovalEntry struct {
	Spname       string                `json:"spname"`        // 审批名称(请假，报销，自定义审批名称)
	ApplyName    string                `json:"apply_name"`    // 申请人姓名
	ApplyOrg     string                `json:"apply_org"`     // 申请人部门
	ApprovalName []string              `json:"approval_name"` // 审批人姓名
	NotifyName   []string              `json:"notify_name"`   // 抄送人姓名
	SpStatus     int                   `json:"sp_status"`     // 审批状态
	SpNum        int                   `json:"sp_num"`        // 审批单号
	MediaIds     []string              `json:"mediaids"`      // 审批的附件media_id，可使用media/get获取附件
	ApplyTime    int                   `json:"apply_time"`    // 审批单提交时间
	ApplyUserId  string                `json:"apply_user_id"` // 审批单提交者的userid
	Expense      *ApprovalEntryExpense `json:"expense"`       // 报销类型
	Comm         *ApprovalEntryComm    `json:"comm"`          // 审批模板信息
	Leave        *ApprovalEntryLeave   `json:"leave"`         // 请假类型
}

type ApprovalData struct {
	Count     int             `json:"count"`      // 拉取的审批单个数，最大值为100，当total参数大于100时，可运用next_spnum参数进行多次拉取
	Total     int             `json:"total"`      // 时间段内的总审批单个数
	NextSpnum int             `json:"next_spnum"` // 拉取列表的最后一个审批单号
	Data      []ApprovalEntry `json:"data"`
}

type WechatApproval struct {
	Client *wecom.WechatClient
}

func NewWechatApproval(client *wecom.WechatClient) *WechatApproval {
	return &WechatApproval{client}
}

// 获取审批数据
//
// 通过本接口来获取公司一段时间内的审批记录。一次拉取调用最多拉取100个审批记录，可以通
// 过多次拉取的方式来满足需求，但调用频率不可超过600次/分。
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/91530
func (w WechatApproval) GetApprovalData(startTime, endTime, nextSpNum int) (*ApprovalData, error) {
	data := struct {
		StartTime int `json:"starttime"`            // 获取审批记录的开始时间。Unix时间戳
		EndTime   int `json:"endtime"`              // 获取审批记录的结束时间。Unix时间戳
		NextSpNum int `json:"next_spnum,omitempty"` // 第一个拉取的审批单号，不填从该时间段的第一个审批单拉取
	}{
		StartTime: startTime,
		EndTime:   endTime,
		NextSpNum: nextSpNum,
	}

	out := ApprovalData{}
	err := w.Client.Post(urlGetApprovalData, nil, data, nil, &out)
	return &out, err
}
