package wecom

type UserInfo struct {
	// 成员UserID。对应管理端的帐号，企业内必须唯一。不区分大小写，长度为1~64个字节
	UserId string `json:"userid"`

	// 成员名称
	Name string `json:"name"`

	// 别名；第三方仅通讯录应用可获取
	Alias string `json:"alias,omitempty"`

	// 手机号码，第三方仅通讯录应用可获取
	Mobile string `json:"mobile,omitempty"`

	// 成员所属部门id列表，仅返回该应用有查看权限的部门id
	Department []int `json:"department"`

	// 主部门ID
	MainDepartMent int `json:"main_department"`

	// 部门内的排序值，默认为0。数量必须和department一致，数值越大排序越前面
	// 值范围是[0, 2^32)
	Order []int `json:"order,omitempty"`

	// 职务信息；第三方仅通讯录应用可获取
	Position string `json:"position,omitempty"`

	// 性别。0表示未定义，1表示男性，2表示女性
	Gender string `json:"gender,omitempty"`

	// 邮箱，第三方仅通讯录应用可获取
	Email string `json:"email,omitempty"`

	// 表示在所在的部门内是否为上级。；第三方仅通讯录应用可获取
	IsLeaderInDept []int `json:"is_leader_in_dept,omitempty"`

	// 成员启用状态。1表示启用的成员，0表示被禁用。注意，服务商调用接口不会返回此字段
	Enable int `json:"enable,omitempty"`

	// 座机。第三方仅通讯录应用可获取
	Telephone string `json:"telephone,omitempty"`

	// 地址
	Address string `json:"address,omitempty"`

	// 扩展属性，第三方仅通讯录应用可获取
	ExtAttr *UserExtAttr `json:"extattr,omitempty"`

	// 头像url。注：如果要获取小图将url最后的”/0”改成”/100”即可。
	// 第三方仅通讯录应用可获取
	Avatar string `json:"avatar,omitempty"`

	// 头像缩略图url。第三方仅通讯录应用可获取；对于非第三方创建的成员，第三方通讯录应用也不可获取
	ThumbAvatar string `json:"thumb_avatar,omitempty"`

	// 激活状态: 1=已激活，2=已禁用，4=未激活。
	// 已激活代表已激活企业微信或已关注微工作台（原企业号）。
	// 未激活代表既未激活企业微信又未关注微工作台（原企业号）。
	Status int `json:"status,omitempty"`

	// 员工个人二维码，扫描可添加为外部联系人；第三方仅通讯录应用可获取
	QrCode string `json:"qr_code,omitempty"`

	// 对外职务，如果设置了该值，则以此作为对外展示的职务，否则以position来展示
	ExternalPosition string `json:"external_position,omitempty"`

	// 成员对外属性，字段详情见对外属性；第三方仅通讯录应用可获取
	ExternalProfile *UserExtProfile `json:"external_profile,omitempty"`
}

type UserExtAttr struct {
	// 属性列表
	Attrs []UserAttr `json:"attrs,omitempty"`
}

type UserExtProfile struct {
	// 企业对外简称，需从已认证的企业简称中选填。可在“我的企业”页中查看企业简称认证状态
	ExternalCorpName string `json:"external_corp_name,omitempty"`

	// 成员对外信息
	ExternalAttr []UserAttr `json:"external_attr,omitempty"`
}

type UserAttr struct {
	// 属性类型: 0-本文 1-网页 2-小程序
	Type int `json:"type"`

	// 属性名称： 需要先确保在管理端有创建该属性，否则会忽略
	Name string `json:"name"`

	// 文本类型的属性(type为0时必填)
	Text *UserAttrText `json:"text,omitempty"`

	// 网页类型的属性，url和title字段要么同时为空表示清除该属性，要么同时不为空
	// type为1时必填
	Web *UserAttrWeb `json:"web,omitempty"`

	// 小程序类型的属性，appid和title字段要么同时为空表示清除改属性，要么同时不为空
	// type为2时必填
	MiniProgram *UserAttrMiniProgram `json:"miniprogram,omitempty"`
}

type UserAttrMiniProgram struct {
	// 小程序appid，必须是有在本企业安装授权的小程序，否则会被忽略
	Appid string `json:"appid,omitempty"`

	// 小程序的展示标题,长度限制12个UTF8字符
	PagePath string `json:"pagepath,omitempty"`

	// 小程序的页面路径
	Title string `json:"title,omitempty"`
}

type UserAttrWeb struct {
	// 网页的url,必须包含http或者https头
	Url string `json:"url,omitempty"`

	// 网页的展示标题,长度限制12个UTF8字符
	Title string `json:"title,omitempty"`
}

type UserAttrText struct {
	// 文本属性内容,长度限制12个UTF8字符
	Value string `json:"value,omitempty"`
}
