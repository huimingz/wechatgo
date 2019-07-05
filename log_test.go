package wechatgo

import (
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	var log Logger = DefaultLogger()
	if log == nil {
		t.Error("DefaultLogger() retrun nil type")
	}
}
