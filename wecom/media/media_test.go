package media

import (
	"context"
	"os"
	"testing"

	"github.com/huimingz/wechatgo/testdata"
	"github.com/huimingz/wechatgo/wecom"
)

var mediaId string
var wechatMedia *WechatMedia

func TestWechatMedia_UploadMedia(t *testing.T) {
	type test struct {
		name     string
		filename string
		filePath string
		type_    string
		wantErr  bool
	}

	tests := []test{
		test{
			name:     "case 1",
			filename: "test_file",
			filePath: "test_file.jpg",
			type_:    "image",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("WechatMedia_UploadMedia() error = %v, wantErr = %v", err, tt.wantErr)
			}

			mediaInfo, err := wechatMedia.UploadMedia(context.Background(), tt.filename, tt.type_, file)
			if (err != nil) != tt.wantErr {
				t.Errorf("WechatMedia_UploadMedia() error = %v, wantErr = %v", err, tt.wantErr)
			}
			mediaId = mediaInfo.MediaId
		})
	}
}

func TestWechatMedia_GetMedia(t *testing.T) {
	body, fn, err := wechatMedia.GetMedia(context.Background(), mediaId)
	if err != nil {
		t.Errorf("WechatMedia.GetMedia() error = '%s'", err)
	}
	defer body.Close()

	if fn == "" {
		t.Error("WechatMedia.GetMedia() error = '文件名为空'")
	}
}

func init() {
	var conf = testdata.TestConf
	var wechatClient = wecom.NewWechatClient(conf.CorpId, conf.CorpSecret, conf.AgentId)
	wechatMedia = NewWechatMedia(wechatClient)
}
