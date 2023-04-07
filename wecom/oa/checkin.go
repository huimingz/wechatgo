// Package oa OA数据接口
package oa

import (
	"context"

	"github.com/huimingz/wechatgo/wecom"
)

const (
	urlGetCheckinData   string = "/cgi-bin/checkin/getcheckindata"
	urlGetCheckinOption string = "/cgi-bin/checkin/getcheckinoption"
)

type CheckinData struct {
	UserId         string   `json:"userid"`          // 用户id
	GroupName      string   `json:"groupname"`       // 打卡规则名称
	CheckinType    string   `json:"checkin_type"`    // 打卡类型。字符串，目前有：上班打卡，下班打卡，外出打卡
	ExceptionType  string   `json:"exception_type"`  // 异常类型，字符串，包括：时间异常，地点异常，未打卡，wifi异常，非常用设备。如果有多个异常，以分号间隔
	CheckinTime    int      `json:"checkin_time"`    // 打卡时间。Unix时间戳
	LocationTitle  string   `json:"location_title"`  // 打卡地点title
	LocationDetail string   `json:"location_detail"` // 打卡地点详情
	WifiName       string   `json:"wifiname"`        // 打卡wifi名称
	Notes          string   `json:"notes"`           // 打卡备注
	WifiMac        string   `json:"wifimac"`         // 打卡的MAC地址/bssid
	MediaIds       []string `json:"mediaids"`        // 打卡的附件media_id，可使用media/get获取附件
}

type CheckinOptLocationInfo struct {
	Latitude       int    `json:"lat"`        // 位置打卡地点纬度
	Longitude      int    `json:"lng"`        // 位置打卡地点经度
	LocationTitle  string `json:"loc_title"`  // 位置打卡地点名称
	LocationDetail string `json:"loc_detail"` // 位置打卡地点详情
	Distance       int    `json:"distance"`   // 允许打卡范围（米）
}

type CheckinOptWifiInfo struct {
	WifiName string `json:"wifiname"` // WiFi打卡地点名称
	WifiMac  string `json:"wifimac"`  // WiFi打卡地点MAC地址/bssid
}

type CheckinOptSpeDay struct {
	TimeStamp   int              `json:"timestamp"`   // 特殊日期具体时间
	Notes       string           `json:"notes"`       // 特殊日期备注
	CheckinTime []CheckinOptTime `json:"checkintime"` // 打卡时间
}

type CheckinOptTime struct {
	WorkSec          int `json:"work_sec"`            // 上班时间，表示为距离当天0点的秒数
	OffWorkSec       int `json:"off_work_sec"`        // 下班时间，表示为距离当天0点的秒数
	RemindWorkSec    int `json:"remind_work_sec"`     // 上班提醒时间，表示为距离当天0点的秒数
	RemindOffWorkSec int `json:"remind_off_work_sec"` // 下班提醒时间，表示为距离当天0点的秒数
}

type CheckinOptGroupData struct {
	WorkDays       []int            `json:"workdays"`        // 工作日。若为固定时间上下班或自由上下班，则1到7分别表示星期一到星期日；若为按班次上下班，则表示拉取班次的日期。
	CheckinTime    []CheckinOptTime `json:"checkintime"`     // 打卡时间
	FlexTime       int              `json:"flex_time"`       // 弹性时间（毫秒）
	NoNeedOffWork  bool             `json:"noneed_offwork"`  // 下班不需要打卡
	LimitAheadTime int              `json:"limit_aheadtime"` // 打卡时间限制（毫秒）
}

type CheckinOptGroup struct {
	GroupType              int                      `json:"grouptype"`                // 打卡规则类型。1：固定时间上下班；2：按班次上下班；3：自由上下班
	GroupId                int                      `json:"groupid"`                  // 打卡规则id
	CheckinData            []CheckinOptGroupData    `json:"checkindata"`              // 分组数据
	SpeWorkDays            []CheckinOptSpeDay       `json:"spe_workdays"`             // 特殊上班日
	SpeOffDays             []CheckinOptSpeDay       `json:"spe_offdays"`              // 特殊下班日
	SyncHolidays           bool                     `json:"sync_holidays"`            // 是否同步法定节假日
	GroupName              string                   `json:"groupname"`                // 打卡规则名称
	NeedPhoto              bool                     `json:"need_photo"`               // 是否打卡必须拍照
	WifiMacInfos           []CheckinOptWifiInfo     `json:"wifimac_infos"`            // Wifi打卡信息
	NoteCanUseLocalPic     bool                     `json:"note_can_use_local_pic"`   // 是否备注时允许上传本地图片
	AllowCheckinOffWorkDay bool                     `json:"allow_checkin_offworkday"` // 是否非工作日允许打卡
	AllowApplyOffWorkDay   bool                     `json:"allow_apply_offworkday"`   // 是否允许异常打卡时提交申请
	LocationInfos          []CheckinOptLocationInfo `json:"loc_infos"`                // 位置信息
}

type CheckinOptInfo struct {
	UserId string          `json:"userid"` // 用户id
	Group  CheckinOptGroup `json:"group"`  // 规则组
}

type WechatCheckin struct {
	Client *wecom.Client
}

func NewWechatCheckin(client *wecom.Client) *WechatCheckin {
	return &WechatCheckin{client}
}

// 获取打卡数据
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90262
func (w WechatCheckin) GetData(ctx context.Context, checkinType, startTime, endTime int, userIds []string) ([]CheckinData, error) {
	data := struct {
		OpenCheckinDataType int      `json:"opencheckindatatype"` // 打卡类型。1：上下班打卡；2：外出打卡；3：全部打卡
		StartTime           int      `json:"starttime"`           // 获取打卡记录的开始时间。Unix时间戳
		EndTime             int      `json:"endtime"`             // 获取打卡记录的结束时间。Unix时间戳
		UserIds             []string `json:"useridlist"`          // 需要获取打卡记录的用户列表
	}{
		OpenCheckinDataType: checkinType,
		StartTime:           startTime,
		EndTime:             endTime,
		UserIds:             userIds,
	}

	out := struct {
		CheckinData []CheckinData `json:"checkindata"`
	}{}
	err := w.Client.Post(ctx, urlGetCheckinData, nil, data, nil, &out)
	return out.CheckinData, err
}

// 获取打卡规则
//
// 企业微信文档：https://work.weixin.qq.com/api/doc#90000/90135/90263
func (w WechatCheckin) GetOption(ctx context.Context, dateTime int, userIds []string) ([]CheckinOptInfo, error) {
	data := struct {
		DateTime   int      `json:"datetime"`   // 需要获取规则的日期当天0点的Unix时间戳
		UserIdList []string `json:"useridlist"` // 需要获取打卡规则的用户列表
	}{
		DateTime:   dateTime,
		UserIdList: userIds,
	}

	out := struct {
		Info []CheckinOptInfo `json:"info"`
	}{}

	err := w.Client.Post(ctx, urlGetCheckinOption, nil, data, nil, &out)
	return out.Info, err
}
