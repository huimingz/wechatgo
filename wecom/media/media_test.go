package media

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/huimingz/wechatgo/testdata"
	"github.com/huimingz/wechatgo/wecom"
)

type MediaTestSuite struct {
	suite.Suite
	mediaId string
	media   *WechatMedia
}

func (s *MediaTestSuite) SetupSuite() {
	client := wecom.NewClient(testdata.TestConf.CorpId, testdata.TestConf.CorpSecret, testdata.TestConf.AgentId)
	s.media = NewWechatMedia(client)
}

func (s *MediaTestSuite) TestUploadMedia() {
	type args struct {
		filename string
		filePath string
		type_    string
	}
	testcases := []struct {
		name      string
		args      args
		expectErr bool
	}{
		{
			name: "case 1",
			args: args{
				filename: "test_file",
				filePath: "test_file.jpg",
				type_:    "image",
			},
			expectErr: false,
		},
	}

	for _, tt := range testcases {
		s.Run(tt.name, func() {
			file, err := os.Open(tt.args.filePath)
			s.NoError(err)

			mediaInfo, err := s.media.UploadMedia(context.Background(), tt.args.filename, tt.args.type_, file)

			if tt.expectErr {
				s.Error(err)
			} else {
				s.NoError(err)
				s.mediaId = mediaInfo.MediaId
			}
		})
	}
}

func (s *MediaTestSuite) TestGetMedia() {
	body, fn, err := s.media.GetMedia(context.Background(), s.mediaId)
	defer body.Close()

	s.NoError(err)
	s.NotEmpty(fn)
}

func TestMediaTestSuite(t *testing.T) {
	suite.Run(t, new(MediaTestSuite))
}

func init() {
	var conf = testdata.TestConf
	var wechatClient = wecom.NewClient(conf.CorpId, conf.CorpSecret, conf.AgentId)
	wechatMedia = NewWechatMedia(wechatClient)
}

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
