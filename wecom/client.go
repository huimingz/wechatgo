// Package wecom 客服端
//
// 特性：
// 线程安全；access token过期自动
package wecom

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/huimingz/wechatgo"
	"github.com/huimingz/wechatgo/storage"
)

const EntBaseUrl = "https://qyapi.weixin.qq.com"

type WechatClient struct {
	CorpId      string          // 企业ID
	CorpSecret  string          // 应用的凭证密钥
	AgentId     int             // agentid
	HttpClient  *http.Client    // http.Client
	accessToken string          // 凭证
	ExpiresIn   time.Duration   // 凭证的有效时间
	lastFresh   time.Time       // 最后一次token刷新时间
	Mutex       *sync.Mutex     // 互斥锁
	BaseUrl     string          // 微信服务器Url
	Storage     storage.Storage // token存储器
	storageKey  string          // token的key值
	Log         wechatgo.Logger // 日志
}

type WechatClientOption func(client *WechatClient)

func WechatClientWithLogger(logger wechatgo.Logger) WechatClientOption {
	return func(client *WechatClient) {
		client.Log = logger
	}
}

func WechatClientWithStorage(storage storage.Storage) WechatClientOption {
	return func(client *WechatClient) {
		client.Storage = storage
	}
}

func WechatClientWithHTTPClient(httpClient *http.Client) WechatClientOption {
	return func(client *WechatClient) {
		client.HttpClient = httpClient
	}
}

func WechatClientWithExpiresIn(sec time.Duration) WechatClientOption {
	if sec <= 0 || sec > 7200 {
		sec = time.Second * 7200
	}

	return func(client *WechatClient) {
		client.ExpiresIn = sec
	}
}

func WechatClientWithMutex(lock *sync.Mutex) WechatClientOption {
	return func(client *WechatClient) {
		client.Mutex = lock
	}
}

func NewWechatClient(corpid, corpSecret string, agentId int, options ...WechatClientOption) *WechatClient {
	client := WechatClient{}
	client.CorpId = corpid
	client.CorpSecret = corpSecret
	client.AgentId = agentId
	client.BaseUrl = EntBaseUrl

	for _, opt := range options {
		opt(&client)
	}

	if client.HttpClient == nil {
		client.HttpClient = &http.Client{}
	}
	if client.ExpiresIn == 0 {
		client.ExpiresIn = time.Second * 7200
	}
	if client.Mutex == nil {
		client.Mutex = &sync.Mutex{}
	}
	if client.Storage == nil {
		client.Storage = storage.NewMemoryStorage()
	}

	if client.Log == nil {
		client.Log = wechatgo.DefaultLogger()
	}

	return &client
}

// GetAccessTokenStorageKey 获取认证令牌缓存Key
func (client WechatClient) GetAccessTokenStorageKey() string {
	if client.storageKey == "" {
		client.storageKey = "accesstoken_" + client.CorpSecret
	}
	return client.storageKey
}

// IsExpired 检查access token是否过期
func (client *WechatClient) IsExpired(ctx context.Context) bool {
	storageKey := client.GetAccessTokenStorageKey()
	if client.Storage.HasExpired(ctx, storageKey) {
		return true
	}
	return false
}

// GetAccessToken 提供access_token的获取接口
//
// 当access_token过期或者为空字符串时，会重新获取一次access_token
func (client *WechatClient) GetAccessToken(ctx context.Context) (string, error) {
	storageKey := client.GetAccessTokenStorageKey()
	val := client.Storage.Get(ctx, storageKey)

	if val == "" {
		client.Log.Info("The access token is expired, try to get a new access token")

		err := client.FetchAccessToken(ctx)
		if err != nil {
			client.Log.Error(
				fmt.Sprintf("An error has occurred during getting access token，Error: %s\n", err.Error()))
			return "", err
		}
		val = client.Storage.Get(ctx, storageKey)
	}
	return val, nil
}

// FetchAccessToken 重新获取access token
//
// 当设置的过期时间为无效（<0 || >7200）时，将自动重置过期时间为远程指定时间
func (client *WechatClient) FetchAccessToken(ctx context.Context) error {
	values := url.Values{}
	values.Add("corpid", client.CorpId)
	values.Add("corpsecret", client.CorpSecret)

	url_ := client.UrlCompletion("/cgi-bin/gettoken")
	url_ += "?" + values.Encode()

	request, err := http.NewRequest("GET", url_, nil)
	if err != nil {
		return err
	}

	request = request.WithContext(ctx)
	resp, err := client.HttpClient.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Json解码
	accessToken := struct {
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExporesIn   uint64 `json:"expires_in"`
	}{}
	// accessToken := RespAccessToken{}
	err = json.Unmarshal(content, &accessToken)
	if err != nil {
		return err
	}

	// 检查过期时间是否有效，无效则重新设置
	if accessToken.ErrCode == 0 && accessToken.AccessToken != "" {
		if client.ExpiresIn == 0 || client.ExpiresIn > time.Second*7200 {
			client.Mutex.Lock()
			defer client.Mutex.Unlock()
			client.Log.Info(fmt.Sprintf("Old expiresIn is't valid, set new expiresIn = %ds", accessToken.ExporesIn))
			client.ExpiresIn = time.Duration(accessToken.ExporesIn) * 1000 * 1000 * 1000
		}

		key := client.GetAccessTokenStorageKey()
		client.Log.Info("Set new access token to storage.")
		err = client.Storage.Set(ctx, key, accessToken.AccessToken, client.ExpiresIn)
		if err != nil {
			return err
		}
	} else if accessToken.ErrCode != 0 {
		return wechatgo.NewWXMsgError(accessToken.ErrCode, accessToken.ErrMsg)
	}
	return nil
}

func (client *WechatClient) UrlCompletion(reqUrl string) string {
	re, _ := regexp.Compile("^https?://.*")
	if !re.MatchString(reqUrl) {
		if strings.HasSuffix(client.BaseUrl, "/") {
			client.Mutex.Lock()
			defer client.Mutex.Unlock()
			client.BaseUrl = strings.TrimRight(client.BaseUrl, "/")
		}
		reqUrl = strings.TrimLeft(reqUrl, "/")
		return client.BaseUrl + "/" + reqUrl
	} else {
		return reqUrl
	}
}

func (client WechatClient) valuesTokenCompletion(ctx context.Context, values url.Values) (url.Values, error) {
	if values == nil {
		values = url.Values{}
	}

	if values.Get("access_token") == "" {
		token, err := client.GetAccessToken(ctx)
		if err != nil {
			return values, err
		}
		if token == "" {
			return values, errors.New("fresh access token  failed")
		}
		values.Add("access_token", token)
	}
	return values, nil
}

func (client WechatClient) respHandler(ctx context.Context, resp *http.Response, errmsg wechatgo.WxMsgInterface, out interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("http response status code[%d] != 200", resp.StatusCode)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if errmsg == nil {
		errmsg = &wechatgo.WXMsgError{}
	}
	err = json.Unmarshal(content, errmsg)
	if err != nil {
		return err
	}
	client.Log.Debug(fmt.Sprintf("ResponseHandler response message: %s", errmsg))
	if errmsg.GetErrCode() != 0 {
		return errmsg
	}

	if out == nil {
		return nil
	}

	err = json.Unmarshal(content, out)
	if err != nil {
		client.Log.Debug("Get response content: failed!")
		return err
	}
	client.Log.Debug("Get response content: successful!")
	return nil
}

func (client *WechatClient) Get(ctx context.Context, url_ string, values url.Values, errmsg wechatgo.WxMsgInterface, out interface{}) error {
	values, err := client.valuesTokenCompletion(ctx, values)
	if err != nil {
		return err
	}

	url_ = client.UrlCompletion(url_)
	url_ += "?" + values.Encode()

	request, err := http.NewRequest("GET", url_, nil)
	if err != nil {
		return err
	}

	request = request.WithContext(ctx)
	resp, err := client.HttpClient.Do(request)
	if err != nil {
		return err
	}
	err = client.respHandler(ctx, resp, errmsg, out)
	return err
}

func (client *WechatClient) RawGet(ctx context.Context, url_ string, values url.Values) (resp *http.Response, err error) {
	values, err = client.valuesTokenCompletion(ctx, values)
	if err != nil {
		return nil, err
	}

	url_ = client.UrlCompletion(url_)
	url_ += "?" + values.Encode()

	request, err := http.NewRequest("GET", url_, nil)
	if err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)
	resp, err = client.HttpClient.Do(request)
	return resp, err
}

func (client *WechatClient) AdvPost(ctx context.Context, url_, contentType string, values url.Values, data interface{}, errmsg wechatgo.WxMsgInterface, out interface{}) error {
	values, err := client.valuesTokenCompletion(ctx, values)
	if err != nil {
		return err
	}

	url_ = client.UrlCompletion(url_)
	url_ += "?" + values.Encode()

	var body io.Reader
	if v, ok := data.(io.Reader); ok {
		body = v
	} else {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		body = bytes.NewReader(jsonData)
	}

	request, err := http.NewRequest("POST", url_, body)
	if err != nil {
		return err
	}

	request = request.WithContext(ctx)
	request.Header.Add("Content-Type", contentType)
	resp, err := client.HttpClient.Do(request)
	if err != nil {
		return err
	}

	return client.respHandler(ctx, resp, errmsg, out)
}

func (client WechatClient) Post(ctx context.Context, url_ string, values url.Values, data interface{}, errmsg wechatgo.WxMsgInterface, out interface{}) error {
	return client.AdvPost(ctx, url_, "application/json", values, data, errmsg, out)
}
