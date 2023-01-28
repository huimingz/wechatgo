package wechatgo

// 全局错误代码：https://work.weixin.qq.com/api/doc#90000/90139/90313

import "fmt"

type WechatMsgInterface interface {
	GetErrCode() int
	GetErrMsg() string
	Error() string
}

type WechatMessageError struct {
	ErrCode int    `json:"errcode"` // 响应错误状态码
	ErrMsg  string `json:"errmsg"`  // 响应错误消息内容
}

func (err WechatMessageError) Error() string {
	return fmt.Sprintf("errcode=%d, errmsg='%s'", err.ErrCode, err.ErrMsg)
}

func (err WechatMessageError) GetErrCode() int {
	return err.ErrCode
}

func (err WechatMessageError) GetErrMsg() string {
	return err.ErrMsg
}

func NewWXMsgError(code int, msg string) *WechatMessageError {
	return &WechatMessageError{ErrCode: code, ErrMsg: msg}
}
