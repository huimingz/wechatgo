// Package storage session存储器
package storage

import (
	"context"
	"time"
)

type Storage interface {
	// Get 获取Value
	Get(ctx context.Context, key string) string

	// Set 设置Value，包含TTL
	Set(ctx context.Context, key, val string, ttl time.Duration) error

	// HasExpired 检测是否过期
	HasExpired(ctx context.Context, key string) bool
}
