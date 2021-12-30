// Package storage session存储器
package storage

import (
	"time"
)

type Storage interface {
	// Get 获取Value
	Get(key string) string

	// Set 设置Value，包含TTL
	Set(key, val string, ttl time.Duration) error

	// HasExpired 检测是否过期
	HasExpired(key string) bool
}
