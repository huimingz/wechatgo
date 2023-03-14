package wecom

import "context"

type Wecom struct {
	client *WechatClient
	App    *applicationManager
	User   *UserManager
}

func NewWecom(corpId, corpSecret string, agentId int, optionFns ...WechatClientOption) *Wecom {
	client := NewWechatClient(corpId, corpSecret, agentId, optionFns...)
	return &Wecom{
		client: client,
		App:    newWechatAppManage(client),
		User:   newManager(client),
	}
}

func (w *Wecom) GetAccessToken(ctx context.Context) (token string, err error) {
	return w.client.GetAccessToken(ctx)
}
