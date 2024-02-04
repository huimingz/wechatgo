// 使用通讯录管理接口，原则上需要使用 通讯录管理secret，也可以使用 应用secret。
// 但是使用应用secret只能进行“查询”、“邀请”等非写操作，而且只能操作应用可见范围内的通讯录。
// 文档网址：https://work.weixin.qq.com/api/doc#90000/90135/90193

// Package user 成员管理
package wecom

import (
	"context"
	"net/url"
	"strconv"
	"time"
)

const (
	urlUserGet         = "/cgi-bin/user/get"
	urlUserCreate      = "/cgi-bin/user/create"
	urlUserUpdate      = "/cgi-bin/user/update"
	urlUserDelete      = "/cgi-bin/user/delete"
	urlBatchUserDelete = "/cgi-bin/user/batchdelete"
	urlUserId2OpenId   = "/cgi-bin/user/convert_to_openid"
	urlOpenId2UserId   = "/cgi-bin/user/convert_to_userid"
	urlUserVerify      = "/cgi-bin/user/authsucc"
	urlUerInvate       = "/cgi-bin/batch/invite"
	urlGetJoinQRCode   = "/cgi-bin/corp/get_join_qrcode"
	urlGetActiveStat   = "/cgi-bin/user/get_active_stat" // 获取企业活跃成员数
)

// 创建用户时使用
type UserForCreate struct {
	// 成员UserID。对应管理端的帐号，企业内必须唯一。不区分大小写，长度为1~64个字节。
	// 只能由数字、字母和“_-@.”四种字符组成，且第一个字符必须是数字或字母。
	UserId string `json:"userid"`

	// 成员名称。长度为1~64个utf8字符
	Name string `json:"name"`

	// 成员别名。长度1~32个utf8字符
	Alias string `json:"alias,omitempty"`

	// 手机号码。企业内必须唯一，mobile/email二者不能同时为空
	Mobile string `json:"mobile,omitempty"`

	// 成员所属部门id列表,不超过20个
	Department []int `json:"department"`

	// 部门内的排序值，默认为0，成员次序以创建时间从小到大排列。
	// 数量必须和department一致，数值越大排序越前面。有效的值范围是[0, 2^32)
	Order []int `json:"order,omitempty"`

	// 职务信息。长度为0~128个字符
	Position string `json:"position,omitempty"`

	// 性别。1表示男性，2表示女性
	Gender string `json:"gender,omitempty"`

	// 邮箱。长度6~64个字节，且为有效的email格式。企业内必须唯一，
	// mobile/email二者不能同时为空
	Email string `json:"email,omitempty"`

	// 个数必须和department一致，表示在所在的部门内是否为上级。
	// 1表示为上级，0表示非上级。在审批等应用里可以用来标识上级审批人
	IsLeaderInDept []int `json:"is_leader_in_dept,omitempty"`

	// 启用/禁用成员。1表示启用成员，0表示禁用成员
	Enable int `json:"enable,omitempty"`

	// 成员头像的mediaid，通过素材管理接口上传图片获得的mediaid
	AvatarMediaId string `json:"avatar_mediaid,omitempty"`

	// 座机。32字节以内，由纯数字或’-‘号组成
	Telephone string `json:"telephone,omitempty"`

	// 地址。长度最大128个字符
	Address string `json:"address,omitempty"`

	// 自定义字段。自定义字段需要先在WEB管理端添加，见扩展属性添加方法，否则忽略未知属性
	// 的赋值。与对外属性一致，不过只支持type=0的文本和type=1的网页类型，详细描述查看对
	// 外属性
	ExtAttr *UserExtAttr `json:"extattr,omitempty"`

	// 是否邀请该成员使用企业微信（将通过微信服务通知或短信或邮件下发邀请，
	// 每天自动下发一次，最多持续3个工作日），默认值为true
	ToInvite bool `json:"to_invite,omitempty"`

	// 对外职务，如果设置了该值，则以此作为对外展示的职务，否则以position来展示。
	// 长度12个汉字内
	ExternalPosition string `json:"external_position,omitempty"`

	// 成员对外属性，字段详情见对外属性
	ExternalProfile *UserExtProfile `json:"external_profile,omitempty"`
}

type UserForUpdate struct {
	// 成员UserID。对应管理端的帐号，企业内必须唯一。不区分大小写，长度为1~64个字节
	UserId string `json:"userid"`

	// 成员名称。长度为1~64个utf8字符
	Name string `json:"name,omitempty"`

	// 别名。长度为1-32个utf8字符
	Alias string `json:"alias,omitempty"`

	// 手机号码。企业内必须唯一。若成员已激活企业微信，则需成员自行修改
	// （此情况下该参数被忽略，但不会报错）
	Mobile string `json:"mobile,omitempty"`

	// 成员所属部门id列表，不超过20个
	Department []int `json:"department,omitempty"`

	// 部门内的排序值，默认为0。数量必须和department一致，数值越大排序越前面。
	// 有效的值范围是[0, 2^32)
	Order []int `json:"order,omitempty"`

	// 职务信息。长度为0~128个字符
	Position string `json:"position,omitempty"`

	// 性别。1表示男性，2表示女性
	Gender string `json:"gender,omitempty"`

	// 邮箱。长度不超过64个字节，且为有效的email格式。企业内必须唯一
	Email string `json:"email,omitempty"`

	// 上级字段，个数必须和department一致，表示在所在的部门内是否为上级
	IsLeaderInDept []int `json:"is_leader_in_dept,omitempty"`

	// 启用/禁用成员。1表示启用成员，0表示禁用成员
	Enable int `json:"enable,omitempty"`

	// 成员头像的mediaid，通过素材管理接口上传图片获得的mediaid
	AvatarMediaId string `json:"avatar_mediaid,omitempty"`

	// 座机。由1-32位的纯数字或’-‘号组成
	Telephone string `json:"telephone,omitempty"`

	// 地址。长度最大128个字符
	Address string `json:"address,omitempty"`

	// 自定义字段。自定义字段需要先在WEB管理端添加，见扩展属性添加方法，
	// 否则忽略未知属性的赋值。与对外属性一致，不过只支持type=0的文本和type=1的网页类型，
	// 详细描述查看对外属性
	ExtAttr *UserExtAttr `json:"extattr,omitempty"`

	// 对外职务，如果设置了该值，则以此作为对外展示的职务，否则以position来展示。
	// 不超过12个汉字
	ExternalPosition string `json:"external_position,omitempty"`

	// 成员对外属性，字段详情见对外属性
	ExternalProfile *UserExtProfile `json:"external_profile,omitempty"`
}

type UserManager struct {
	Client *Client
}

func newUserManager(client *Client) *UserManager {
	return &UserManager{client}
}

// 读取成员
//
// 在通讯录同步助手中此接口可以读取企业通讯录的所有成员信息，而自建应用可以读取该应用设置
// 的可见范围内的成员信息。
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90196
func (w *UserManager) GetUser(ctx context.Context, userId string) (*UserInfo, error) {
	values := url.Values{}
	values.Add("userid", userId)

	out := UserInfo{}
	err := w.Client.Get(ctx, urlUserGet, values, nil, &out)
	return &out, err
}

// 创建成员
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90195
func (w *UserManager) CreateUser(ctx context.Context, data UserForCreate) error {
	return w.Client.Post(ctx, urlUserCreate, nil, data, nil, nil)
}

// 更新成员
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90197
func (w *UserManager) UpdateUser(ctx context.Context, data UserForUpdate) error {
	return w.Client.Post(ctx, urlUserUpdate, nil, data, nil, nil)
}

// 删除成员
//
// https://work.weixin.qq.com/api/doc#90000/90135/90198
func (w *UserManager) DeleteUser(ctx context.Context, userId string) error {
	values := url.Values{}
	values.Add("userid", userId)

	return w.Client.Get(ctx, urlUserDelete, values, nil, nil)
}

// 批量删除成员
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90199
func (w *UserManager) DeleteUsers(ctx context.Context, userIdlist []string) error {
	users := struct {
		UserIdList []string `json:"useridlist"`
	}{
		UserIdList: userIdlist,
	}

	return w.Client.Post(ctx, urlBatchUserDelete, nil, users, nil, nil)
}

// userid转openid
//
// 该接口使用场景为企业支付，在使用企业红包和向员工付款时，需要自行将企业微信的userid转成openid。
//
// 注：需要成员使用微信登录企业微信或者关注微工作台（原企业号）才能转成openid
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90202
func (w *UserManager) UserId2OpenId(ctx context.Context, userId string) (openId string, err error) {
	data := struct {
		UserId string `json:"userid"` // 企业内的成员id
	}{
		UserId: userId,
	}

	out := struct {
		OpenId string `json:"openid"` // 企业微信成员userid对应的openid
	}{}
	err = w.Client.Post(ctx, urlUserId2OpenId, nil, data, nil, &out)
	openId = out.OpenId
	return
}

// openid转userid
//
// 该接口主要应用于使用企业支付之后的结果查询。
//
// 开发者需要知道某个结果事件的openid对应企业微信内成员的信息时，
// 可以通过调用该接口进行转换查询。
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90202
func (w *UserManager) OpenId2UserId(ctx context.Context, openid string) (userId string, err error) {
	data := struct {
		OpenId string `json:"openid"` // 在使用企业支付之后，返回结果的openid
	}{
		OpenId: openid,
	}

	out := struct {
		UserId string `json:"userid"` // 该openid在企业微信对应的成员userid
	}{}
	err = w.Client.Post(ctx, urlOpenId2UserId, nil, data, nil, &out)
	userId = out.UserId
	return
}

// Phone2UserId 手机号换userid
func (w *UserManager) Phone2UserId(ctx context.Context, phone string) (userId string, err error) {
	if phone == "" {
		return "", nil
	}

	return "", nil
}

// 二次验证
//
// 此接口可以满足安全性要求高的企业进行成员加入验证。
// 开启二次验证后，用户加入企业时需要跳转企业自定义的页面进行验证。
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90203
func (w *UserManager) Verify(ctx context.Context, userId string) error {
	values := url.Values{}
	values.Add("userid", userId)

	return w.Client.Get(ctx, urlUserVerify, values, nil, nil)
}

// 邀请成员
//
// 企业可通过接口批量邀请成员使用企业微信，邀请后将通过短信或邮件下发通知。
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90975
func (w *UserManager) Invite(ctx context.Context, user []string, party, tag []int) error {
	data := struct {
		User  []string `json:"user"`  // 成员ID列表, 最多支持1000个
		Party []int    `json:"party"` // 部门ID列表，最多支持100个
		Tag   []int    `json:"tag"`   // 标签ID列表，最多支持100个
	}{
		User:  user,
		Party: party,
		Tag:   tag,
	}

	return w.Client.Post(ctx, urlUerInvate, nil, data, &InvalidError{}, nil)
}

// 获取加入企业二维码
//
// 支持企业用户获取实时成员加入二维码。
// 二维码链接，有效期7天
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/91714
func (w *UserManager) GetJoinQRCode(ctx context.Context, sizeType int) (joinQRCode string, err error) {
	values := url.Values{}
	if sizeType != 0 {
		values.Add("size_type", strconv.Itoa(sizeType))
	}

	out := struct {
		JoinQRCode string `json:"join_qrcode"`
	}{}

	err = w.Client.Get(ctx, urlGetJoinQRCode, values, nil, &out)
	joinQRCode = out.JoinQRCode
	return
}

// GetActiveStat 获取企业活跃成员数
func (w *UserManager) GetActiveStat(ctx context.Context, date time.Time) (activeCount int, err error) {
	payload := struct {
		Date string `json:"date"`
	}{
		Date: date.Format("2006-01-02"),
	}

	out := struct {
		ActiveCount int `json:"active_cnt"` // 活跃成员数
	}{}

	err = w.Client.Post(ctx, urlGetActiveStat, nil, payload, nil, &out)
	return out.ActiveCount, err
}
