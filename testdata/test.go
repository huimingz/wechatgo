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

var TestConf = Conf{}
