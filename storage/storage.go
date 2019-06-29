// Package storage session存储器
package storage

import (
	"time"
)

type Storage interface {
	Get(key string) string
	Set(key, val string, ttl time.Duration) error
	HasExpired(key string) bool
}
