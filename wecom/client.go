// Package client 客服端
//
// 特性：
// 线程安全；access token过期自动
package wecom

import (
	"bytes"
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

func NewWechatClient(corpid, corpSecret string, agentId int, httpClient *http.Client,
	expiresIn time.Duration, mutex *sync.Mutex, storage_ storage.Storage, log_ wechatgo.Logger) *WechatClient {
	client := WechatClient{}
	client.CorpId = corpid
	client.CorpSecret = corpSecret
	client.AgentId = agentId
	client.BaseUrl = EntBaseUrl

	if httpClient == nil {
		client.HttpClient = &http.Client{}
	} else {
		client.HttpClient = httpClient
	}

	if expiresIn <= 0 || expiresIn > 7200 {
		client.ExpiresIn = time.Second * 7200
	} else {
		client.ExpiresIn = expiresIn
	}

	if mutex == nil {
		client.Mutex = &sync.Mutex{}
	} else {
		client.Mutex = mutex
	}

	if storage_ == nil {
		client.Storage = storage.NewMemoryStorage()
	} else {
		client.Storage = storage_
	}

	if log_ == nil {
		client.Log = wechatgo.DefaultLogger()
	} else {
		client.Log = log_
	}

	return &client
}

func (client WechatClient) GetAccessTokenStorageKey() string {
	if client.storageKey == "" {
		client.storageKey = "accesstoken_" + client.CorpSecret
	}
	return client.storageKey
}

// 检查access token是否过期
func (client *WechatClient) IsExpired() bool {
	storageKey := client.GetAccessTokenStorageKey()
	if client.Storage.HasExpired(storageKey) {
		return true
	}
	return false
}

// 提供access_token的获取接口
//
// 当access_token过期或者为空字符串时，会重新获取一次access_token
func (client *WechatClient) GetAccessToken() (string, error) {
	storageKey := client.GetAccessTokenStorageKey()
	val := client.Storage.Get(storageKey)

	if val == "" {
		client.Log.Info("The access token is expired, try to get a new access token")

		err := client.FetchAccessToken()
		if err != nil {
			client.Log.Error(
				fmt.Sprintf("An error has occurred during getting access token，Error: %s\n", err.Error()))
			return "", err
		}
		val = client.Storage.Get(storageKey)
	}
	return val, nil
}

// 重新获取access token
//
// 当设置的过期时间为无效（<0 || >7200）时，将自动重置过期时间为远程指定时间
func (client *WechatClient) FetchAccessToken() error {
	values := url.Values{}
	values.Add("corpid", client.CorpId)
	values.Add("corpsecret", client.CorpSecret)

	url_ := client.UrlCompletion("/cgi-bin/gettoken")
	url_ += "?" + values.Encode()

	request, err := http.NewRequest("GET", url_, nil)
	if err != nil {
		return err
	}

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
		err = client.Storage.Set(key, accessToken.AccessToken, client.ExpiresIn)
		if err != nil {
			return err
		}
	} else if accessToken.ErrCode != 0 {
		return wechatgo.NewWXMsgError(accessToken.ErrCode, accessToken.ErrMsg)
	}
	return nil
}

// func (client *WechatClient) RespJsonUnmarshal(response *http.Response, v interface{}) error {
// 	if response.StatusCode != 200 {
// 		return errors.New("Response.StatusCode != 200")
// 	}
//
// 	// defer response.Body.Close()
// 	content, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = json.Unmarshal(content, v)
// 	if err != nil {
// 		return err
// 	}
// 	// TODO: 检查errcode
//
// 	return nil
// }

// func (client WechatClient) WXMsgChecker(msg common.WXMsgInterface) error {
// 	if msg.GetErrCode() != 0 {
// 		if v, ok := msg.(error); ok {
// 			if reflect.ValueOf(v).Kind() == reflect.Ptr {
// 				return msg
// 			} else {
// 				return error(msg)
// 			}
// 		} else {
// 			return fmt.Errorf("errcode=%d", msg.GetErrCode())
// 		}
// 	}
// 	return nil
// }

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

func (client WechatClient) valuesTokenCompletion(values url.Values) (url.Values, error) {
	if values == nil {
		values = url.Values{}
	}

	if values.Get("access_token") == "" {
		token, err := client.GetAccessToken()
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

func (client WechatClient) respHandler(resp *http.Response, errmsg wechatgo.WxMsgInterface, out interface{}) error {
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

func (client *WechatClient) Get(url_ string, values url.Values, errmsg wechatgo.WxMsgInterface, out interface{}) error {
	values, err := client.valuesTokenCompletion(values)
	if err != nil {
		return err
	}

	url_ = client.UrlCompletion(url_)
	url_ += "?" + values.Encode()

	request, err := http.NewRequest("GET", url_, nil)
	if err != nil {
		return err
	}

	resp, err := client.HttpClient.Do(request)
	if err != nil {
		return err
	}
	err = client.respHandler(resp, errmsg, out)
	return err
}

func (client *WechatClient) RawGet(url_ string, values url.Values) (resp *http.Response, err error) {
	values, err = client.valuesTokenCompletion(values)
	if err != nil {
		return nil, err
	}

	url_ = client.UrlCompletion(url_)
	url_ += "?" + values.Encode()

	request, err := http.NewRequest("GET", url_, nil)
	if err != nil {
		return nil, err
	}

	resp, err = client.HttpClient.Do(request)
	return resp, err
}

func (client *WechatClient) AdvPost(url_, contentType string, values url.Values, data interface{}, errmsg wechatgo.WxMsgInterface, out interface{}) error {
	values, err := client.valuesTokenCompletion(values)
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
	request.Header.Add("Content-Type", contentType)
	resp, err := client.HttpClient.Do(request)
	if err != nil {
		return err
	}

	return client.respHandler(resp, errmsg, out)
}

func (client WechatClient) Post(url_ string, values url.Values, data interface{}, errmsg wechatgo.WxMsgInterface, out interface{}) error {
	return client.AdvPost(url_, "application/json", values, data, errmsg, out)
}
