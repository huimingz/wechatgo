package storage

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type MemoryStorageTestSuite struct {
	suite.Suite
	storage *MemoryStorage
}

func (s *MemoryStorageTestSuite) SetupTest() {
	s.storage = NewMemoryStorage()
}

func (s *MemoryStorageTestSuite) TestGet() {
	ctx := context.Background()
	s.storage.Set(ctx, "get", "xxx", time.Microsecond*1150)
	s.Equal("xxx", s.storage.Get(ctx, "get"))

	testcases := []struct {
		name  string
		key   string
		sleep time.Duration
		want  string
	}{
		{
			name:  "case 1",
			key:   "get",
			sleep: 0,
			want:  "xxx",
		},
		{
			name:  "case 2",
			key:   "get",
			sleep: time.Microsecond * 1150,
			want:  "",
		},
	}
	for _, tc := range testcases {
		s.Run(tc.name, func() {
			time.Sleep(tc.sleep)
			s.Equal(tc.want, s.storage.Get(ctx, tc.key))
		})
	}

}

func (s *MemoryStorageTestSuite) TestHasExpired() {
	ms := NewMemoryStorage()
	ms.Set(context.Background(), "key", "val", time.Microsecond*150)

	type test struct {
		name    string
		key     string
		storage *MemoryStorage
		sleep   time.Duration
		want    bool
	}

	tests := []test{
		test{
			name:    "case 1",
			key:     "key",
			storage: ms,
			sleep:   0,
			want:    false,
		},
		test{
			name:    "case 2",
			key:     "key",
			storage: ms,
			sleep:   time.Microsecond * 150,
			want:    true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			time.Sleep(tt.sleep)
			got := tt.storage.HasExpired(context.Background(), tt.key)
			s.Equal(tt.want, got, "MemoryStorage.HasExpired() error = result [%v] != want [%v]", got, tt.want)
		})
	}

}

func TestMemoryStorageTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryStorageTestSuite))
}

func TestNewMemoryStorage(t *testing.T) {
	ms := NewMemoryStorage()
	if reflect.ValueOf(ms).Kind() != reflect.Ptr {
		t.Error("NewMemoryStorage() is not pointer")
	}
}

func TestMemoryStorage_HasExpired(t *testing.T) {
	ms := NewMemoryStorage()
	ms.Set(context.Background(), "key", "val", time.Microsecond*150)

	type test struct {
		name    string
		key     string
		storage *MemoryStorage
		sleep   time.Duration
		want    bool
	}

	tests := []test{
		test{
			name:    "case 1",
			key:     "key",
			storage: ms,
			sleep:   0,
			want:    false,
		},
		test{
			name:    "case 2",
			key:     "key",
			storage: ms,
			sleep:   time.Microsecond * 150,
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time.Sleep(tt.sleep)
			got := tt.storage.HasExpired(context.Background(), tt.key)
			if got != tt.want {
				t.Errorf("MemoryStorage.HasExpired() error = result [%v] != want [%v]", got, tt.want)
			}
		})
	}
}

func TestMemoryStorage_Set(t *testing.T) {
	ms := NewMemoryStorage()

	err := ms.Set(context.Background(), "key", "val", time.Nanosecond*10)
	if err != nil {
		t.Errorf("MemoryStorage.Set() error = result [%v] != nil", err)
	}
}
