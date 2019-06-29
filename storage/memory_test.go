package storage

import (
	"reflect"
	"testing"
	"time"
)

func TestNewMemoryStorage(t *testing.T) {
	ms := NewMemoryStorage()
	if reflect.ValueOf(ms).Kind() != reflect.Ptr {
		t.Error("NewMemoryStorage() is not pointer")
	}
}

func TestMemoryStorage_Get(t *testing.T) {
	ms := NewMemoryStorage()
	ms.Set("get", "xxx", time.Microsecond*150)

	type test struct {
		name    string
		key     string
		storage *MemoryStorage
		sleep   time.Duration
		want    string
	}

	tests := []test{
		test{
			name:    "case 1",
			key:     "get",
			storage: ms,
			sleep:   0,
			want:    "xxx",
		},
		test{
			name:    "case 2",
			key:     "get",
			storage: ms,
			sleep:   time.Microsecond * 150,
			want:    "",
		},
		test{
			name:    "case 3",
			key:     "getx",
			storage: ms,
			sleep:   0,
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time.Sleep(tt.sleep)
			val := tt.storage.Get(tt.key)
			if val != tt.want {
				t.Errorf("MemoryStorage.Get() error = value [%s] != %s", val, tt.want)
			}
		})
	}
}

func TestMemoryStorage_HasExpired(t *testing.T) {
	ms := NewMemoryStorage()
	ms.Set("key", "val", time.Microsecond*150)

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
			got := tt.storage.HasExpired(tt.key)
			if got != tt.want {
				t.Errorf("MemoryStorage.HasExpired() error = result [%v] != want [%v]", got, tt.want)
			}
		})
	}
}

func TestMemoryStorage_Set(t *testing.T) {
	ms := NewMemoryStorage()

	err := ms.Set("key", "val", time.Nanosecond*10)
	if err != nil {
		t.Errorf("MemoryStorage.Set() error = result [%v] != nil", err)
	}
}
