// Package testdata 测试数据
package testdata

type Conf struct {
	CorpId          string `yaml:"corp_id"`
	CorpSecret      string `yaml:"corp_secret"`
	UserSecret      string `yaml:"user_secret"`
	ApprovalSecret  string `yaml:"approval_secret"`
	ApprovalAgentId int    `yaml:"approval_agentid"`
	CheckinSecret   string `yaml:"checkin_secret"`
	CheckinAgentId  int    `yaml:"checkin_agentid"`
	DialSecret      string `yaml:"dial_secret"`
	DialAgentId     int    `yaml:"dial_agentid"`
	AgentId         int    `yaml:"agentid"`
	BaseUrl         string `yaml:"base_url"`
	UserId          string `yaml:"userid"`
	OpenId          string `yaml:"openid"`
}

var TestConf = Conf{
	CorpId:          "wwa6bba55a03adfc19",
	CorpSecret:      "b2Y8glnO6-jgzxiehnW_oB1DBp3uJYa3OiAvST0ni80",
	UserSecret:      "_RVFTP9FIfyHfxAxaENioLNp2RQmxyiI5B3fKvQI9N0",
	ApprovalSecret:  "DOMizAY3139oqkSeJaDmKDwKFSBgTNebrx9HI4B0AVA",
	ApprovalAgentId: 3010040,
	CheckinSecret:   "bNu_sTSJRN3UXUZE11cXFcq0BDMlW5cLRLqRJh0g870",
	CheckinAgentId:  3010011,
	DialSecret:      "",
	DialAgentId:     0,
	AgentId:         1000004,
	BaseUrl:         "",
	UserId:          "ZhouHuiMing",
	OpenId:          "odG42w2ZECpIVr8ii1RUvhwCRnZc",
}
