// Package log 日志
package log

import (
	"os"

	"github.com/Sirupsen/logrus"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
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

func DefaultLogger() *logrus.Logger {
	textFormatter := logrus.TextFormatter{}
	textFormatter.FullTimestamp = true
	textFormatter.ForceColors = true
	textFormatter.TimestampFormat = "2006/01/02 15:04:05"

	log := logrus.New()
	log.SetFormatter(&textFormatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)

	return log
}
