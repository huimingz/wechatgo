package storage

import (
	"context"
	"sync"
	"time"
)

type storageData struct {
	expireIn time.Time
	value    string
}

type MemoryStorage struct {
	mutex *sync.Mutex
	data  map[string]storageData
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{data: map[string]storageData{}, mutex: &sync.Mutex{}}
}

func (s MemoryStorage) Get(ctx context.Context, key string) string {
	data, ok := s.data[key]
	if ok && s.HasExpired(ctx, key) {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		delete(s.data, key)
		return ""
	}
	return data.value
}

func (s MemoryStorage) Set(ctx context.Context, key string, val string, ttl time.Duration) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[key] = storageData{value: val, expireIn: time.Now().Add(ttl)}
	return nil
}

func (s MemoryStorage) HasExpired(ctx context.Context, key string) bool {
	val, ok := s.data[key]
	if !ok {
		return true
	}
	if val.expireIn.Sub(time.Now()) <= 0 {
		return true
	}
	return false
}
