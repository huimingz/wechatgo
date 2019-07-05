// Package msg 发送应用消息
package msg

import (
	"strings"

	"github.com/huimingz/wechatgo/ent"
)

const (
	urlSend = "/cgi-bin/message/send"
)

type sendData struct {
	// 成员ID列表（消息接收者，多个接收者用‘|’分隔，最多支持1000个）。
	// 特殊情况：指定为@all，则向该企业应用的全部成员发送
	ToUser string `json:"touser,omitempty"`

	// 部门ID列表，多个接收者用‘|’分隔，最多支持100个。当touser为@all时忽略本参数
	ToParty string `json:"toparty,omitempty"`

	// 标签ID列表，多个接收者用‘|’分隔，最多支持100个。当touser为@all时忽略本参数
	ToTag string `json:"totag,omitempty"`

	// 消息类型
	MsgType string `json:"msgtype"`

	// 企业应用的id，整型。企业内部开发，可在应用的设置页面查看；
	// 第三方服务商，可通过接口 获取企业授权信息 获取该参数值
	AgentId int `json:"agentid"`
}

func (data *sendData) Init(toUser []string, toParty, toTag []int, msgType string, agentid int) {
	data.ToUser = strings.Join(toUser, "|")
	data.ToParty = intSliceToString(toParty, "|")
	data.ToTag = intSliceToString(toTag, "|")
	data.MsgType = msgType
	data.AgentId = agentid
}

type TextMsg struct {
	Content string `json:"content"` // 消息内容，最长不超过2048个字节，超过将截断
}

type ImageMsg struct {
	MediaId string `json:"media_id"` // 图片媒体文件id，可以调用上传临时素材接口获取
}

type VoiceMsg struct {
	MediaId string `json:"media_id"` // 语音文件id，可以调用上传临时素材接口获取
}

type VideoMsg struct {
	MediaId     string `json:"media_id"`              // 视频媒体文件id，可以调用上传临时素材接口获取
	Title       string `json:"title,omitempty"`       // 视频消息的标题，不超过128个字节，超过会自动截断
	Description string `json:"description,omitempty"` // 视频消息的描述，不超过512个字节，超过会自动截断
}

type FileMsg struct {
	MediaId string `json:"media_id"` // 文件id，可以调用上传临时素材接口获取
}

type TextCardMsg struct {
	Title       string `json:"title"`             // 标题，不超过128个字节，超过会自动截断
	Description string `json:"description"`       // 描述，不超过512个字节，超过会自动截断
	Url         string `json:"url"`               // 点击后跳转的链接
	BtnText     string `json:"btntext,omitempty"` // 按钮文字。 默认为“详情”， 不超过4个文字，超过自动截断
}

type Article struct {
	Title       string `json:"title"`                 // 标题，不超过128个字节，超过会自动截断
	Description string `json:"description,omitempty"` // 描述，不超过512个字节，超过会自动截断
	Url         string `json:"url"`                   // 点击后跳转的链接
	PicUrl      string `json:"picurl,omitempty"`      // 图文消息的图片链接，支持JPG、PNG格式，较好的效果为大图 1068*455，小图150*150
}

type NewsMsg struct {
	Articles []Article `json:"articles"` // 图文消息，一个图文消息支持1到8条图文
}

type MPArticle struct {
	// 标题，不超过128个字节，超过会自动截断
	Title string `json:"title"`

	// 图文消息缩略图的media_id, 可以通过素材管理接口获得。
	// 此处thumb_media_id即上传接口返回的media_id
	ThumbMediaId string `json:"thumb_media_id"`

	// 图文消息的作者，不超过64个字节
	Author string `json:"author,omitempty"`

	// 图文消息点击“阅读原文”之后的页面链接
	ContentSourceUrl string `json:"content_source_url,omitempty"`

	// 图文消息的内容，支持html标签，不超过666 K个字节
	Content string `json:"content"`

	// 图文消息的描述，不超过512个字节，超过会自动截断
	Digest string `json:"digest"`
}

type MPNewsMsg struct {
	Articles []MPArticle `json:"articles"` // 图文消息，一个图文消息支持1到8条图文
}

type MarkdownMsg struct {
	Content string `json:"content"` // markdown内容，最长不超过2048个字节，必须是utf8编码
}

type NoticeContentItem struct {
	Key   string `json:"key"`   // 长度10个汉字以内
	Value string `json:"value"` // 长度30个汉字以内
}

type MiniProgramNoticeMsg struct {
	// 小程序appid，必须是与当前小程序应用关联的小程序
	AppId string `json:"appid"`

	// 点击消息卡片后的小程序页面，仅限本小程序内的页面。
	// 该字段不填则消息点击后不跳转。
	Page string `json:"page,omitempty"`

	// 消息标题，长度限制4-12个汉字
	Title string `json:"title"`

	// 消息描述，长度限制4-12个汉字
	Description string `json:"description,omitempty"`

	// 是否放大第一个content_item
	EmphasisFirstItem bool `json:"emphasis_first_item,omitempty"`

	// 消息内容键值对，最多允许10个item
	ContentItem []NoticeContentItem `json:"content_item,omitempty"`
}

type TaskCardBtn struct {
	// 按钮key值，用户点击后，会产生任务卡片回调事件，回调事件会带上该key值，
	// 只能由数字、字母和“_-@.”组成，最长支持128字节
	Key string `json:"key"`

	// 按钮名称
	Name string `json:"name"`

	// 点击按钮后显示的名称，默认为“已处理”
	ReplaceName string `json:"replace_name,omitempty"`

	// 按钮字体颜色，可选“red”或者“blue”,默认为“blue”
	Color string `json:"color,omitempty"`

	// 按钮字体是否加粗，默认false
	IsBold bool `json:"is_bold,omitempty"`
}

type TaskCardMsg struct {
	// 标题，不超过128个字节，超过会自动截断
	Title string `json:"title"`

	// 描述，不超过512个字节，超过会自动截断
	Description string `json:"description"`

	// 点击后跳转的链接。最长2048字节，请确保包含了协议头(http/https)
	Url string `json:"url,omitempty"`

	// 任务id，同一个应用发送的任务卡片消息的任务id不能重复，
	// 只能由数字、字母和“_-@.”组成，最长支持128字节
	TaskId string `json:"task_id"`

	// 按钮列表，按钮个数为为1~2个。
	Btn []TaskCardBtn `json:"btn"`
}

type WechatMsg struct {
	Client *ent.WechatClient
}

func NewWechatMsg(client *ent.WechatClient) *WechatMsg {
	return &WechatMsg{client}
}

// 发送消息
func (w WechatMsg) send(data interface{}) error {
	errmsg := MsgError{}
	err := w.Client.Post(urlSend, nil, data, &errmsg, nil)

	if errmsg.InvalidUser != "" || errmsg.InvalidTag != "" || errmsg.InvalidParty != "" {
		return &errmsg
	} else {
		return err
	}
}

// 文本消息
//
// toUser、toParty、toTag不能同时为空
//
// 其中text参数的content字段可以支持换行、以及A标签，即可打开自定义的网页
// （可参考以上示例代码）(注意：换行符请用转义过的\n)
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90236/%E6%96%87%E6%9C%AC%E6%B6%88%E6%81%AF
func (w WechatMsg) SendText(toUser []string, toParty, toTag []int, text TextMsg, safe bool) error {
	data := struct {
		sendData
		Text TextMsg `json:"text"`
		Safe int     `json:"safe,omitempty"`
	}{}
	data.Init(toUser, toParty, toTag, "text", w.Client.AgentId)
	data.Text = text
	if safe {
		data.Safe = 1
	} else {
		data.Safe = 0
	}
	return w.send(data)
}

// 图片消息
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90236/%E5%9B%BE%E7%89%87%E6%B6%88%E6%81%AF
func (w WechatMsg) SendImage(toUser []string, toParty, toTag []int, image ImageMsg, safe bool) error {
	data := struct {
		sendData
		Image ImageMsg `json:"image"`
		Safe  int      `json:"safe,omitempty"`
	}{}
	data.Init(toUser, toParty, toTag, "image", w.Client.AgentId)
	data.Image = image
	if safe {
		data.Safe = 1
	} else {
		data.Safe = 0
	}
	return w.send(data)
}

// 语音消息
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90236/%E8%AF%AD%E9%9F%B3%E6%B6%88%E6%81%AF
func (w WechatMsg) SendVoice(toUser []string, toParty, toTag []int, voice VoiceMsg) error {
	data := struct {
		sendData
		Voice VoiceMsg `json:"voice"`
	}{}
	data.Init(toUser, toParty, toTag, "voice", w.Client.AgentId)
	data.Voice = voice
	return w.send(data)
}

// 视频消息
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90236/%E8%A7%86%E9%A2%91%E6%B6%88%E6%81%AF
func (w WechatMsg) SendVideo(toUser []string, toParty, toTag []int, video VideoMsg, safe bool) error {
	data := struct {
		sendData
		Video VideoMsg `json:"video"`
		Safe  int      `json:"safe,omitempty"`
	}{}
	data.Init(toUser, toParty, toTag, "video", w.Client.AgentId)
	data.Video = video
	if safe {
		data.Safe = 1
	} else {
		data.Safe = 0
	}
	return w.send(data)
}

// 文件消息
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90236/%E6%96%87%E4%BB%B6%E6%B6%88%E6%81%AF
func (w WechatMsg) SendFile(toUser []string, toParty, toTag []int, file FileMsg, safe bool) error {
	data := struct {
		sendData
		File FileMsg `json:"file"`
		Safe int     `json:"safe,omitempty"`
	}{}
	data.Init(toUser, toParty, toTag, "file", w.Client.AgentId)
	data.File = file
	if safe {
		data.Safe = 1
	} else {
		data.Safe = 0
	}
	return w.send(data)
}

// 文本卡片消息
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90236/%E6%96%87%E6%9C%AC%E5%8D%A1%E7%89%87%E6%B6%88%E6%81%AF
func (w WechatMsg) SendTextCard(toUser []string, toParty, toTag []int, textCard TextCardMsg) error {
	data := struct {
		sendData
		TextCard TextCardMsg `json:"textcard"`
	}{}
	data.Init(toUser, toParty, toTag, "textcard", w.Client.AgentId)
	data.TextCard = textCard
	return w.send(data)
}

// 图文消息
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90236/%E5%9B%BE%E6%96%87%E6%B6%88%E6%81%AF
func (w WechatMsg) SendNews(toUser []string, toParty, toTag []int, news NewsMsg) error {
	data := struct {
		sendData
		News NewsMsg `json:"news"`
	}{}
	data.Init(toUser, toParty, toTag, "news", w.Client.AgentId)
	data.News = news
	return w.send(data)
}

// 图文消息（mpnews）
//
// mpnews类型的图文消息，跟普通的图文消息一致，唯一的差异是图文内容存储在企业微信。
// 多次发送mpnews，会被认为是不同的图文，阅读、点赞的统计会被分开计算。
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90236/%E5%9B%BE%E6%96%87%E6%B6%88%E6%81%AF%EF%BC%88mpnews%EF%BC%89
func (w WechatMsg) SendMPNews(toUser []string, toParty, toTag []int, mpNews MPNewsMsg, safe bool) error {
	data := struct {
		sendData
		MPNews MPNewsMsg `json:"mpnews"`

		// 表示是否是保密消息，0表示可对外分享，1表示不能分享且内容显示水印，
		// 2表示仅限在企业内分享，默认为0；
		//
		// 注意仅mpnews类型的消息支持safe值为2，其他消息类型不支持
		Safe int `json:"safe,omitempty"`
	}{}
	data.Init(toUser, toParty, toTag, "mpnews", w.Client.AgentId)
	data.MPNews = mpNews
	if safe {
		data.Safe = 1
	} else {
		data.Safe = 0
	}
	return w.send(data)
}

// markdown消息
//
// 目前仅支持markdown语法的子集。
// 微工作台（原企业号）不支持展示markdown消息
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90236/markdown%E6%B6%88%E6%81%AF
func (w WechatMsg) SendMarkdown(toUser []string, toParty, toTag []int, md MarkdownMsg) error {
	data := struct {
		sendData
		Markdown MarkdownMsg `json:"markdown"`
	}{}
	data.Init(toUser, toParty, toTag, "markdown", w.Client.AgentId)
	data.Markdown = md
	return w.send(data)
}

// 小程序通知消息
//
// 小程序通知消息只允许小程序应用发送，消息会通过【小程序通知】发送给用户。
// 小程序应用仅支持发送小程序通知消息，暂不支持文本、图片、语音、视频、图文等其他类型的消息。
// 不支持@all全员发送
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90236/%E5%B0%8F%E7%A8%8B%E5%BA%8F%E9%80%9A%E7%9F%A5%E6%B6%88%E6%81%AF
func (w WechatMsg) SendMiniProgramNotice(toUser []string, toParty, toTag []int, mpn MiniProgramNoticeMsg) error {
	data := struct {
		sendData
		MiniProgramNotice MiniProgramNoticeMsg `json:"miniprogram_notice"`
	}{}
	data.Init(toUser, toParty, toTag, "miniprogram_notice", w.Client.AgentId)
	data.MiniProgramNotice = mpn
	return w.send(data)
}

// 任务卡片消息
//
// 仅企业微信2.8.2及以上版本支持
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90236/%E4%BB%BB%E5%8A%A1%E5%8D%A1%E7%89%87%E6%B6%88%E6%81%AF
func (w WechatMsg) SendTaskCard(toUser []string, toParty, toTag []int, taskCard TaskCardMsg) error {
	data := struct {
		sendData
		TaskCard TaskCardMsg `json:"taskcard"`
	}{}
	data.Init(toUser, toParty, toTag, "taskcard", w.Client.AgentId)
	data.TaskCard = taskCard
	return w.send(data)
}
