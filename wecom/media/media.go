// Package media 素材管理
package media

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/huimingz/wechatgo"
	"github.com/huimingz/wechatgo/wecom"
)

const (
	urlUploadMedia        string = "/cgi-bin/media/upload"
	urlUploadPerpetualImg string = "/cgi-bin/media/uploadimg"
	urlGetMedia           string = "/cgi-bin/media/get"
	urlGetVoice           string = "/cgi-bin/media/get/jssdk"
)

type MediaInfo struct {
	Type      string `json:"type"`       // 媒体文件类型，分别有图片（image）、语音（voice）、视频（video），普通文件(file)
	MediaId   string `json:"media_id"`   // 媒体文件上传后获取的唯一标识，3天内有效
	CreatedAt string `json:"created_at"` // 媒体文件上传时间戳
}

type WechatMedia struct {
	Client *wecom.WechatClient
}

func NewWechatMedia(client *wecom.WechatClient) *WechatMedia {
	return &WechatMedia{Client: client}
}

// 上传临时素材
//
// 素材上传得到media_id，该media_id仅三天内有效，media_id在同一企业内应用之间可以共享
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90253
func (w WechatMedia) UploadMedia(ctx context.Context, filename, type_ string, r io.Reader) (*MediaInfo, error) {
	mediaInfo := MediaInfo{}
	bodyBuf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("media", filename)
	if err != nil {
		return &mediaInfo, err
	}
	_, err = io.Copy(fileWriter, r)
	if err != nil {
		return &mediaInfo, err
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	values := url.Values{}
	values.Add("type", type_)

	err = w.Client.AdvPost(ctx, urlUploadMedia, contentType, values, bodyBuf, nil, &mediaInfo)
	return &mediaInfo, err
}

// 检查错误信息
func (w WechatMedia) checkMediaHeader(header http.Header) error {
	if v, ok := header["Error-Code"]; ok {
		code, err := strconv.Atoi(v[0])
		if err != nil {
			return err
		}
		var errMsg string
		if v, ok := header["Error-Msg"]; ok {
			errMsg = v[0]
		}
		return wechatgo.NewWXMsgError(code, errMsg)
	}
	return nil
}

// 获取素材文件名
func (w WechatMedia) getMediaFilename(header http.Header) (fn string, err error) {
	var disposition string
	if v, ok := header["Content-Disposition"]; ok {
		disposition = v[0]
	} else {
		err = errors.New("can't find 'Content-Disposition' from response header")
		return
	}

	compile, err := regexp.Compile(".*filename=\"(?P<filename>[^\"]*)\"")
	if err != nil {
		return
	}

	if v := compile.FindStringSubmatch(disposition); len(v) == 2 {
		fn = v[1]
		return
	} else {
		err = errors.New("can't find filename from 'Content-Disposition'")
		return
	}
}

// 获取临时素材
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90254
func (w WechatMedia) GetMedia(ctx context.Context, mediaId string) (body io.ReadCloser, fn string, err error) {
	values := url.Values{}
	values.Add("media_id", mediaId)

	resp, err := w.Client.RawGet(ctx, urlGetMedia, values)
	if err != nil {
		return
	}

	// 检查错误信息
	err = w.checkMediaHeader(resp.Header)
	if err != nil {
		return
	}

	// 获取文件名
	fn, err = w.getMediaFilename(resp.Header)
	if err != nil {
		return
	}

	body = resp.Body
	return
}

// 上传永久图片
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90256
func (w WechatMedia) UploadPerpetualImg(ctx context.Context, filename string, r io.Reader) (url string, err error) {
	bodyBuf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWritter, err := bodyWriter.CreateFormFile("media", filename)
	if err != nil {
		return
	}

	_, err = io.Copy(fileWritter, r)
	if err != nil {
		return
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	type imgInfo struct {
		Url string `json:"url"`
	}

	imgInfo_ := imgInfo{}
	err = w.Client.AdvPost(ctx, urlUploadPerpetualImg, contentType, nil, bodyBuf, nil, &imgInfo_)
	if err != nil {
		return
	}
	url = imgInfo_.Url
	return
}

// 获取高清语音素材
//
// 可以使用本接口获取从JSSDK的uploadVoice接口上传的临时语音素材，格式为speex，16K
// 采样率。该音频比上文的临时素材获取接口（格式为amr，8K采样率）更加清晰，适合用作语音
// 识别等对音质要求较高的业务。
//
// 参考文档：https://work.weixin.qq.com/api/doc#90000/90135/90255
func (w WechatMedia) GetVoice(ctx context.Context, mediaId string) (body io.ReadCloser, fn string, err error) {
	values := url.Values{}
	values.Add("media_id", mediaId)

	resp, err := w.Client.RawGet(ctx, urlGetVoice, values)
	if err != nil {
		return nil, "", err
	}

	err = w.checkMediaHeader(resp.Header)
	if err != nil {
		return
	}

	fn, err = w.getMediaFilename(resp.Header)
	if err != nil {
		return
	}

	body = resp.Body
	return
}
