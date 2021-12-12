// Package log 日志
package wechatgo

import (
	"log"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type LoggerExample struct{}

func (l LoggerExample) Debug(args ...interface{}) {
	log.Print(args...)
}

func (l LoggerExample) Info(args ...interface{}) {
	log.Print(args...)
}

func (l LoggerExample) Warn(args ...interface{}) {
	log.Print(args...)
}

func (l LoggerExample) Error(args ...interface{}) {
	log.Print(args...)
}

func DefaultLogger() Logger {
	return LoggerExample{}
}
