package storage

import (
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
	ms := MemoryStorage{}
	ms.data = map[string]storageData{}
	ms.mutex = &sync.Mutex{}
	return &ms
}

func (s MemoryStorage) Get(key string) string {
	data, ok := s.data[key]
	if ok && s.HasExpired(key) {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		delete(s.data, key)
		return ""
	}
	return data.value
}

func (s MemoryStorage) Set(key, val string, ttl time.Duration) error {
	now := time.Now()
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[key] = storageData{value: val, expireIn: now.Add(ttl)}
	return nil
}

func (s MemoryStorage) HasExpired(key string) bool {
	now := time.Now()
	val, ok := s.data[key]
	if !ok {
		return true
	}
	if val.expireIn.Sub(now) <= 0 {
		return true
	}
	return false
}
