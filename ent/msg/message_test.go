package msg

import (
	"testing"

	"github.com/huimingz/wechatgo/ent"
	"github.com/huimingz/wechatgo/testdata"
)

var wechatMsg *WechatMsg

func TestWechatMsg_SendText(t *testing.T) {
	type args struct {
		toUser  []string
		toParty []int
		toTag   []int
		text    TextMsg
		safe    bool
	}
	type test struct {
		name    string
		args    args
		wantErr bool
	}

	tests := []test{
		test{
			name: "case 1",
			args: args{
				toUser:  []string{testdata.TestConf.UserId},
				toParty: nil,
				toTag:   nil,
				text:    TextMsg{"hello world"},
				safe:    false,
			},
			wantErr: false,
		},
		test{
			name: "case 2",
			args: args{
				toUser:  []string{"xxxxxxxxxx"},
				toParty: nil,
				toTag:   nil,
				text:    TextMsg{"hello world"},
				safe:    true,
			},
			wantErr: true,
		},
		test{
			name: "case 3",
			args: args{
				toUser:  nil,
				toParty: nil,
				toTag:   nil,
				text:    TextMsg{"hello world"},
				safe:    false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wechatMsg.SendText(tt.args.toUser, tt.args.toParty, tt.args.toTag, tt.args.text, tt.args.safe)
			if (err != nil) != tt.wantErr {
				t.Errorf("WechatMsg.SendText() error = %v wantError = %v", err, tt.wantErr)
			}
		})
	}
}

func TestWechatMsg_SendMarkdown(t *testing.T) {
	text := `您的会议室已经预定，稍后会同步到邮箱
            >**事项详情** 
            >事　项：<font color=\"info\">开会</font> 
            >组织者：@miglioguan 
            >参与者：@miglioguan、@kunliu、@jamdeezhou、@kanexiong、@kisonwang 
            > 
            >会议室：<font color=\"info\">广州TIT 1楼 301</font> 
            >日　期：<font color=\"warning\">2018年5月18日</font> 
            >时　间：<font color=\"comment\">上午9:00-11:00</font> 
            > 
            >请准时参加会议。 
            > 
            >如需修改会议信息，请点击：[修改会议信息](https://work.weixin.qq.com)`

	err := wechatMsg.SendMarkdown([]string{testdata.TestConf.UserId}, nil, nil, MarkdownMsg{Content: text})
	if err != nil {
		t.Errorf("WechatMsg.SendMarkdown() error = '%s'", err)
	}
}

func init() {
	var conf = testdata.TestConf
	var wechatClient = ent.NewWechatClient(conf.CorpId, conf.CorpSecret, conf.AgentId, nil, 0, nil, nil, nil)
	wechatMsg = NewWechatMsg(wechatClient)
}
