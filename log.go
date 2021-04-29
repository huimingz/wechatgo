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
	// Fatal(args ...interface{})
	//
	// Debugln(args ...interface{})
	// Infoln(args ...interface{})
	// Warnln(args ...interface{})
	// Errorln(args ...interface{})
	// Fatalln(args ...interface{})
	//
	// Debugf(format string, args ...interface{})
	// Infof(format string, args ...interface{})
	// Warnf(format string, args ...interface{})
	// Errorf(format string, args ...interface{})
	// Fatalf(format string, args ...interface{})
}


type LoggerExample struct {}

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
