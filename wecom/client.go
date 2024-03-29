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

const _BASE_URL = "https://qyapi.weixin.qq.com"

type Client struct {
	CorpId           string          // 企业ID
	CorpSecret       string          // 应用的凭证密钥
	AgentId          int             // agent id
	httpClient       *http.Client    // http.Client 对象
	accessToken      string          // 凭证
	expiresIn        time.Duration   // 凭证的有效时间
	lastFresh        time.Time       // 最后一次token刷新时间
	tokenWriterMutex *sync.Mutex     // 互斥锁
	baseUrl          string          // 微信服务器Url
	storage          storage.Storage // token存储器
	storageKey       string          // token的key值
	log              wechatgo.Logger // 日志
}

type ClientOptionFn func(client *Client)

func ClientWithLogger(logger wechatgo.Logger) ClientOptionFn {
	return func(client *Client) {
		client.log = logger
	}
}

func ClientWithStorage(storage storage.Storage) ClientOptionFn {
	return func(client *Client) {
		client.storage = storage
	}
}

func ClientWithHTTPClient(httpClient *http.Client) ClientOptionFn {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}

func ClientWithExpiresIn(sec time.Duration) ClientOptionFn {
	if sec <= 0 || sec > 7200 {
		sec = time.Second * 7200
	}

	return func(client *Client) {
		client.expiresIn = sec
	}
}

func ClientWithMutex(lock *sync.Mutex) ClientOptionFn {
	return func(client *Client) {
		client.tokenWriterMutex = lock
	}
}

func NewClient(corpId, corpSecret string, agentId int, options ...ClientOptionFn) *Client {
	client := Client{}
	client.CorpId = corpId
	client.CorpSecret = corpSecret
	client.AgentId = agentId
	client.baseUrl = _BASE_URL

	for _, opt := range options {
		opt(&client)
	}

	if client.httpClient == nil {
		client.httpClient = &http.Client{}
	}
	if client.expiresIn == 0 {
		client.expiresIn = time.Second * 7200
	}
	if client.tokenWriterMutex == nil {
		client.tokenWriterMutex = &sync.Mutex{}
	}
	if client.storage == nil {
		client.storage = storage.NewMemoryStorage()
	}

	if client.log == nil {
		client.log = wechatgo.DefaultLogger()
	}

	return &client
}

// GetAccessTokenStorageKey 获取认证令牌缓存Key
func (client *Client) GetAccessTokenStorageKey() string {
	if client.storageKey == "" {
		client.storageKey = "accesstoken_" + client.CorpSecret
	}
	return client.storageKey
}

// IsExpired 检查access token是否过期
func (client *Client) IsExpired(ctx context.Context) bool {
	storageKey := client.GetAccessTokenStorageKey()
	if client.storage.HasExpired(ctx, storageKey) {
		return true
	}
	return false
}

// GetAccessToken 提供access_token的获取接口
//
// 当access_token过期或者为空字符串时，会重新获取一次access_token
func (client *Client) GetAccessToken(ctx context.Context) (string, error) {
	storageKey := client.GetAccessTokenStorageKey()
	val := client.storage.Get(ctx, storageKey)

	if val == "" {
		client.log.Info(ctx, "The access token is expired, try to get a new access token")

		err := client.FetchAccessToken(ctx)
		if err != nil {
			client.log.Error(ctx, fmt.Sprintf("An error has occurred during getting access token，Error: %s\n", err.Error()))
			return "", err
		}
		val = client.storage.Get(ctx, storageKey)
	}
	return val, nil
}

// FetchAccessToken 重新获取access token
//
// 当设置的过期时间为无效（<0 || >7200）时，将自动重置过期时间为远程指定时间
func (client *Client) FetchAccessToken(ctx context.Context) error {
	values := url.Values{}
	values.Add("corpid", client.CorpId)
	values.Add("corpsecret", client.CorpSecret)

	request, err := http.NewRequest("GET", client.resourceURL("/cgi-bin/gettoken", values), nil)
	if err != nil {
		return err
	}

	request = request.WithContext(ctx)
	resp, err := client.httpClient.Do(request)
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
		ExpiresIn   uint64 `json:"expires_in"`
	}{}
	// accessToken := RespAccessToken{}
	err = json.Unmarshal(content, &accessToken)
	if err != nil {
		return err
	}

	// 检查过期时间是否有效，无效则重新设置
	if accessToken.ErrCode == 0 && accessToken.AccessToken != "" {
		if client.expiresIn == 0 || client.expiresIn > time.Second*7200 {
			client.tokenWriterMutex.Lock()
			defer client.tokenWriterMutex.Unlock()
			client.log.Info(ctx, fmt.Sprintf("Old expiresIn is't valid, set new expiresIn = %ds", accessToken.ExpiresIn))
			client.expiresIn = time.Duration(accessToken.ExpiresIn) * 1000 * 1000 * 1000
		}

		key := client.GetAccessTokenStorageKey()
		client.log.Info(ctx, "Set new access token to storage.")
		err = client.storage.Set(ctx, key, accessToken.AccessToken, client.expiresIn)
		if err != nil {
			return err
		}
	} else if accessToken.ErrCode != 0 {
		return wechatgo.NewWXMsgError(accessToken.ErrCode, accessToken.ErrMsg)
	}
	return nil
}

// GetDomainIpList 获取微信服务器IP地址
func (client *Client) GetDomainIpList(ctx context.Context) ([]string, error) {
	uri := "/cgi-bin/get_api_domain_ip"
	ipList := struct {
		IpList []string `json:"ip_list"`
	}{}

	err := client.Get(ctx, uri, nil, nil, &ipList)
	return ipList.IpList, err
}

func (client *Client) resourceURL(path string, query url.Values) string {
	var uri string
	re, _ := regexp.Compile("^https?://.*")

	if re.MatchString(path) {
		uri = path
	} else {
		uri = client.baseUrl + "/" + strings.TrimLeft(path, "/")
	}

	if query != nil && len(query) > 0 {
		uri += "?" + query.Encode()
	}
	return uri
}

func (client *Client) valuesTokenCompletion(ctx context.Context, values url.Values) (url.Values, error) {
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

func (client *Client) respHandler(ctx context.Context, resp *http.Response, errmsg wechatgo.WechatMsgInterface, out any) error {
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("http response status code[%d] != 200", resp.StatusCode)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return client.handleResult(ctx, errmsg, content, out)
}

func (client *Client) handleResult(ctx context.Context, errmsg wechatgo.WechatMsgInterface, content []byte, out any) error {
	if err := client.verifyResult(ctx, errmsg, content); err != nil {
		return err
	}

	if out == nil {
		return nil
	}

	if err := json.Unmarshal(content, out); err != nil {
		client.log.Debug(ctx, "Get response content: failed!")
		return err
	}
	client.log.Debug(ctx, "Get response content: successful!")
	return nil
}

func (client *Client) verifyResult(ctx context.Context, errmsg wechatgo.WechatMsgInterface, content []byte) error {
	if errmsg == nil {
		errmsg = &wechatgo.WechatMessageError{}
	}

	if err := json.Unmarshal(content, errmsg); err != nil {
		return err
	}

	client.log.Debug(ctx, fmt.Sprintf("ResponseHandler response message: %s", errmsg))
	if errmsg.GetErrCode() != 0 {
		return errmsg
	}
	return nil
}

func (client *Client) Get(ctx context.Context, path string, values url.Values, errmsg wechatgo.WechatMsgInterface, out any) error {
	values, err := client.valuesTokenCompletion(ctx, values)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("GET", client.resourceURL(path, values), nil)
	if err != nil {
		return err
	}

	request = request.WithContext(ctx)
	resp, err := client.httpClient.Do(request)
	if err != nil {
		return err
	}
	return client.respHandler(ctx, resp, errmsg, out)
}

func (client *Client) RawGet(ctx context.Context, path string, values url.Values) (resp *http.Response, err error) {
	values, err = client.valuesTokenCompletion(ctx, values)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("GET", client.resourceURL(path, values), nil)
	if err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)
	resp, err = client.httpClient.Do(request)
	return resp, err
}

func (client *Client) AdvPost(ctx context.Context, path, contentType string, urlValues url.Values, data interface{}, errmsg wechatgo.WechatMsgInterface, out interface{}) error {
	values, err := client.valuesTokenCompletion(ctx, urlValues)
	if err != nil {
		return err
	}

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

	request, err := http.NewRequest("POST", client.resourceURL(path, values), body)
	if err != nil {
		return err
	}

	request = request.WithContext(ctx)
	request.Header.Add("Content-Type", contentType)
	resp, err := client.httpClient.Do(request)
	if err != nil {
		return err
	}

	return client.respHandler(ctx, resp, errmsg, out)
}

func (client *Client) Post(ctx context.Context, url_ string, values url.Values, data interface{}, errmsg wechatgo.WechatMsgInterface, out interface{}) error {
	return client.AdvPost(ctx, url_, "application/json", values, data, errmsg, out)
}
