// Package extcontact 外部联系人管理
package extcontact

import (
	"net/url"

	"github.com/huimingz/wechatgo/ent"
	"github.com/huimingz/wechatgo/ent/user"
)

const (
	urlGetFollowUserList   string = "/cgi-bin/externalcontact/get_follow_user_list"
	urlGetUserList         string = "/cgi-bin/externalcontact/list"
	urlGetUserDetail       string = "/cgi-bin/externalcontact/get"
	urlAddContactWay       string = "/cgi-bin/externalcontact/add_contact_way"
	urlAddMsgTemplate      string = "/cgi-bin/externalcontact/add_msg_template"
	urlGetGroupMsgResult   string = "/cgi-bin/externalcontact/get_group_msg_result"
	urlGetUserBehaviorData string = "/cgi-bin/externalcontact/get_user_behavior_data"
	urlSendWelcomeMsg      string = "/cgi-bin/externalcontact/send_welcome_msg"
	urlGetUnassignedList   string = "/cgi-bin/externalcontact/get_unassigned_list"
	urlTransfer            string = "/cgi-bin/externalcontact/transfer"
)

type ExternalContact struct {
	ExternalUserId  string              `json:"external_userid"`  // 外部联系人的userid
	Name            string              `json:"name"`             // 外部联系人的姓名或别名
	Position        string              `json:"position"`         // 外部联系人的职位，如果外部企业或用户选择隐藏职位，则不返回，仅当联系人类型是企业微信用户时有此字段
	Avatar          string              `json:"avatar"`           // 外部联系人头像，第三方不可获取
	CorpName        string              `json:"corp_name"`        // 外部联系人所在企业的简称，仅当联系人类型是企业微信用户时有此字段
	CorpFullName    string              `json:"corp_full_name"`   // 外部联系人所在企业的主体名称，仅当联系人类型是企业微信用户时有此字段
	Type            int                 `json:"type"`             // 外部联系人的类型，1表示该外部联系人是微信用户，2表示该外部联系人是企业微信用户
	Gender          int                 `json:"gender"`           // 外部联系人性别 0-未知 1-男性 2-女性
	UnionId         string              `json:"unionid"`          // 外部联系人在微信开放平台的唯一身份标识（微信unionid），通过此字段企业可将外部联系人与公众号/小程序用户关联起来。
	ExternalProfile user.UserExtProfile `json:"external_profile"` // 外部联系人的自定义展示信息，可以有多个字段和多种类型，包括文本，网页和小程序，仅当联系人类型是企业微信用户时有此字段，字段详情见对外属性；
}

type FollowUserTag struct {
	GroupName string `json:"group_name"` // 该成员添加此外部联系人所打标签的分组名称
	TagName   string `json:"tag_name"`   // 该成员添加此外部联系人所打标签名称
	Type      int    `json:"type"`       // 该成员添加此外部联系人所打标签类型, 1-企业设置, 2-用户自定义
}

type FollowUser struct {
	UserId        string          `json:"userid"`         // 添加了此外部联系人的企业成员userid
	Remark        string          `json:"remark"`         // 该成员对此外部联系人的备注
	Description   string          `json:"description"`    // 该成员对此外部联系人的描述
	CreateTime    int             `json:"createtime"`     // 该成员添加此外部联系人的时间
	Tags          []FollowUserTag `json:"tags"`           // 该成员添加此外部联系人标签
	RemarkCompany string          `json:"remark_company"` // 该成员对此客户备注的企业名称
	RemarkMobiles []int           `json:"remark_mobiles"` // 该成员对此客户备注的手机号码
	State         string          `json:"state"`          // 该成员添加此客户的渠道，由用户通过创建「联系我」方式指定
}

type ExternalUserDetail struct {
	ExternalContact ExternalContact `json:"external_contact"` // 外部联系人
	FollowUser      []FollowUser    `json:"follow_user"`      // 该外部联系人的成员用户
}

type MiniProgramMsg struct {
	Title      string `json:"title"`        // 小程序消息标题
	PicMediaId string `json:"pic_media_id"` // 小程序消息封面的mediaid，封面图建议尺寸为520*416
	Appid      string `json:"appid"`        // 小程序appid，必须是关联到企业的小程序应用
	Page       string `json:"page"`         // 小程序page路径
}

type LinkMsg struct {
	Title  string `json:"title"`            // 图文消息标题
	PicUrl string `json:"picurl,omitempty"` // 图文消息封面的url（可为空）
	Desc   string `json:"desc,omitempty"`   // 图文消息的描述（可为空）
	Url    string `json:"url"`              // 图文消息的链接
}

type ImageMsg struct {
	MediaId string `json:"media_id"` // 图片的media_id
}

type TextMsg struct {
	Content string `json:"content,omitempty"` // 消息文本内容（可为空）
}

type MsgTemplate struct {
	ExternalUserId []string       `json:"external_userid,omitempty"` // 客户的外部联系人id列表，不可与sender同时为空，最多可传入1万个客户
	Sender         string         `json:"sender,omitempty"`          // 发送企业群发消息的成员userid，不可与external_userid同时为空
	Text           TextMsg        `json:"text"`                      // 消息文本
	Image          ImageMsg       `json:"image"`                     // 图片
	Link           LinkMsg        `json:"link"`                      // 图文消息
	MiniProgram    MiniProgramMsg `json:"miniprogram"`               // 小程序
}

type GroupMsgResultDetail struct {
	ExternalUserId string `json:"external_userid"` // 外部联系人userid
	UserId         string `json:"userid"`          // 企业服务人员的userid
	Status         int    `json:"status"`          // 发送状态 0-未发送 1-已发送 2-因客户不是好友导致发送失败 3-因客户已经收到其他群发消息导致发送失败
	SendTime       int    `json:"send_time"`       // 发送时间，发送状态为1时返回
}

type GroupMsgResult struct {
	CheckStatus int                    `json:"check_status"` // 模板消息的审核状态 0-审核中 1-审核成功 2-审核失败
	DetailList  []GroupMsgResultDetail `json:"detail_list"`  // 详细列表
}

type UserBehavior struct {
	StatTime        int `json:"stat_time"`        // 数据日期，为当日0点的时间戳
	ChatCnt         int `json:"chat_cnt"`         // 成员有主动发送过消息的聊天数，包括单聊和群聊
	MessageCnt      int `json:"message_cnt"`      // 成员在单聊和群聊中发送的消息总数
	ReplyPercentage int `json:"reply_percentage"` // 已回复聊天占比
	AvgReplyTime    int `json:"avg_reply_time"`   // 平均首次回复时长，单位为分钟
}

type WelcomeMsg struct {
	WelcomeCode    string         `json:"welcome_code"` // 通过添加外部联系人事件推送给企业的发送欢迎语的凭证，有效期为20秒
	Text           TextMsg        `json:"text"`         // 消息文本
	Image          ImageMsg       `json:"image"`        // 图片
	Link           LinkMsg        `json:"link"`         // 图文消息
	MiniProgramMsg MiniProgramMsg `json:"miniprogram"`  // 小程序
}

type UnassignedUser struct {
	HandoverUserId string `json:"handover_userid"` // 离职成员的userid
	ExternalUserId string `json:"external_userid"` // 外部联系人userid
	DimissionTime  int    `json:"dimission_time"`  // 成员离职时间
}

type WechatContact struct {
	Client *ent.WechatClient
}

func NewWechatContact(client *ent.WechatClient) *WechatContact {
	return &WechatContact{client}
}

// 获取配置了客户联系功能的成员列表
//
// 企业和第三方服务商可通过此接口获取配置了客户联系功能的成员列表。
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/91554
func (w WechatContact) GetFollowUserList() (followUser []string, err error) {
	out := struct {
		FollowUser []string `json:"follow_user"` // 配置了客户联系功能的成员userid列表
	}{}

	err = w.Client.Get(urlGetFollowUserList, nil, nil, &out)
	followUser = out.FollowUser
	return
}

// 获取外部联系人列表
//
// 企业可通过此接口获取指定成员添加的客户列表。客户是指配置了客户联系功能的成员所添加
// 的外部联系人。没有配置客户联系功能的成员，所添加的外部联系人将不会作为客户返回。
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/91555
func (w WechatContact) GetUserList() (externalUserId []string, err error) {
	out := struct {
		ExternalUserId []string `json:"external_userid"` // 外部联系人的userid列表
	}{}

	err = w.Client.Get(urlGetUserList, nil, nil, &out)
	externalUserId = out.ExternalUserId
	return
}

// 获取外部联系人详情
//
// 企业可通过此接口，根据外部联系人的userid，拉取外部联系人详情。
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/91556
func (w WechatContact) GetUserDetail(userId string) (*ExternalUserDetail, error) {
	values := url.Values{}
	values.Add("external_userid", userId)

	out := ExternalUserDetail{}
	err := w.Client.Get(urlGetUserDetail, values, nil, &out)
	return &out, err
}

// 配置客户联系「联系我」方式
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/91559
func (w WechatContact) AddContactWay(type_, scene, style int, remark, state string, skipVerify bool, user []string, party []int) (configId string, err error) {
	data := struct {
		Type       int      `json:"type"`             // 联系方式类型,1-单人, 2-多人
		Scene      int      `json:"scene"`            // 场景，1-在小程序中联系，2-通过二维码联系
		Style      int      `json:"style,omitempty"`  // 在小程序中联系时使用的控件样式，详见附表
		Remark     string   `json:"remark,omitempty"` // 联系方式的备注信息，用于助记，不超过30个字符
		SkipVerify bool     `json:"skip_verify"`      // 外部客户添加时是否无需验证
		State      string   `json:"state,omitempty"`  // 企业自定义的state参数，用于区分不同的添加渠道，在调用“获取外部联系人详情”时会返回该参数值
		User       []string `json:"user,omitempty"`   // 使用该联系方式的用户userID列表，在type为1时为必填，且只能有一个
		Party      []int    `json:"party,omitempty"`  // 使用该联系方式的部门id列表，只在type为2时有效
	}{
		Type:   type_,
		Scene:  scene,
		Style:  style,
		Remark: remark,
		State:  state,
		User:   user,
		Party:  party,
	}

	out := struct {
		ConfigId string `json:"config_id"`
	}{}
	err = w.Client.Post(urlAddContactWay, nil, data, nil, &out)
	configId = out.ConfigId
	return
}

// 添加企业群发消息模板
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/91560
func (w WechatContact) AddMsgTemplate(msgTemplate MsgTemplate) (failList []string, msgId string, err error) {
	out := struct {
		FailList []string `json:"fail_list"`
		MsgId    string   `json:"msgid"`
	}{}

	err = w.Client.Post(urlAddMsgTemplate, nil, msgTemplate, nil, &out)
	failList = out.FailList
	msgId = out.MsgId
	return
}

// 获取企业群发消息发送结果
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/91561
func (w WechatContact) GetGroupMsgResult(msgId string) (*GroupMsgResult, error) {
	data := struct {
		MsgId string `json:"msgid"` // 群发消息的id，通过添加企业群发消息模板接口返回
	}{}

	out := GroupMsgResult{}
	err := w.Client.Post(urlGetGroupMsgResult, nil, data, nil, &out)
	return &out, err
}

// 获取员工行为数据
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/91580
func (w WechatContact) GetUserBehaviorData(userIds []string, startTime, endTime int) (behavior []UserBehavior, err error) {
	data := struct {
		UserId    []string `json:"userid"`     // userid列表
		StartTime int      `json:"start_time"` // 数据起始时间
		EndTime   int      `json:"end_time"`   // 数据结束时间
	}{}
	out := struct {
		BehaviorData []UserBehavior `json:"behaviro_data"`
	}{}

	err = w.Client.Post(urlGetUserBehaviorData, nil, data, nil, &out)
	behavior = out.BehaviorData
	return
}

// 发送新客户欢迎语
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/91688
func (w WechatContact) SendWelcomeMsg(welcomeMsg WelcomeMsg) error {
	return w.Client.Post(urlSendWelcomeMsg, nil, welcomeMsg, nil, nil)
}

// 获取离职成员的客户列表
//
// 企业和第三方可通过此接口，获取所有离职成员的客户列表，并可进一步调用离职成员的外部
// 联系人再分配接口将这些客户重新分配给其他企业成员。
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/91563
func (w WechatContact) GetUnassignedList(pageId, pageSize int) (userlist []UnassignedUser, isLast bool, err error) {
	data := struct {
		PageId   int `json:"page_id,omitempty"`   // 分页查询，要查询页号，从0开始
		PageSize int `json:"page_size,omitempty"` // 每次返回的最大记录数，默认为1000，最大值为1000
	}{}
	out := struct {
		Info   []UnassignedUser `json:"info"`
		IsLast bool             `json:"is_last"`
	}{}

	err = w.Client.Post(urlGetUnassignedList, nil, data, nil, &out)
	userlist = out.Info
	isLast = out.IsLast
	return
}

// 离职成员的外部联系人再分配
//
// 企业可通过此接口，将已离职成员的外部联系人分配给另一个成员接替联系。
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/91564
func (w WechatContact) Transfer(externalUserId, handoverUserId, takeoverUserId string) error {
	data := struct {
		ExternalUserId string `json:"external_userid"` // 外部联系人的userid，注意不是企业成员的帐号
		HandoverUserId string `json:"handover_userid"` // 离职成员的userid
		TakeoverUserId string `json:"takeover_userid"` // 接替成员的userid
	}{}

	return w.Client.Post(urlTransfer, nil, data, nil, nil)
}
