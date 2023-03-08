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

type defaultLogger struct{}

func (l defaultLogger) Debug(ctx context.Context, args ...interface{}) {
	log.Print(args...)
}

func (l defaultLogger) Info(ctx context.Context, args ...interface{}) {
	log.Print(args...)
}

func (l defaultLogger) Warn(ctx context.Context, args ...interface{}) {
	log.Print(args...)
}

func (l defaultLogger) Error(cxt context.Context, args ...interface{}) {
	log.Print(args...)
}

func DefaultLogger() Logger {
	return defaultLogger{}
}
