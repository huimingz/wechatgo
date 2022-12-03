// Package log 日志
package wechatgo

import (
	"context"
	"log"
)

type Logger interface {
	Debug(ctx context.Context, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Error(cxt context.Context, args ...interface{})
}

type LoggerExample struct{}

func (l LoggerExample) Debug(ctx context.Context, args ...interface{}) {
	log.Print(args...)
}

func (l LoggerExample) Info(ctx context.Context, args ...interface{}) {
	log.Print(args...)
}

func (l LoggerExample) Warn(ctx context.Context, args ...interface{}) {
	log.Print(args...)
}

func (l LoggerExample) Error(cxt context.Context, args ...interface{}) {
	log.Print(args...)
}

func DefaultLogger() Logger {
	return LoggerExample{}
}
