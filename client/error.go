package client

// 全局错误代码：https://work.weixin.qq.com/api/doc#90000/90139/90313

import "fmt"

type WxMsgInterface interface {
	GetErrCode() int
	GetErrMsg() string
	Error() string
}

type WXMsgError struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (err WXMsgError) Error() string {
	return fmt.Sprintf("errcode=%d, errmsg='%s'", err.ErrCode, err.ErrMsg)
}

func (err WXMsgError) GetErrCode() int {
	return err.ErrCode
}

func (err WXMsgError) GetErrMsg() string {
	return err.ErrMsg
}

func NewWXMsgError(code int, msg string) *WXMsgError {
	return &WXMsgError{ErrCode: code, ErrMsg: msg}
}
