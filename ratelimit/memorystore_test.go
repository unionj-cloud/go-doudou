package ratelimit

import (
	"reflect"
	"sync"
	"testing"
)

func TestMemoryStore_addKey(t *testing.T) {
	type fields struct {
		keys      map[string]Limiter
		limiterFn func(store *MemoryStore, key string) Limiter
		mu        *sync.RWMutex
	}
	type args struct {
		key string
	}
	limiter := NewTokenLimiter(1, 3)
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Limiter
	}{
		{
			name: "",
			fields: fields{
				keys: make(map[string]Limiter),
				limiterFn: func(store *MemoryStore, key string) Limiter {
					return limiter
				},
				mu: &sync.RWMutex{},
			},
			args: args{
				key: "192.168.1.6:8080",
			},
			want: limiter,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &MemoryStore{
				keys:      tt.fields.keys,
				limiterFn: tt.fields.limiterFn,
				mu:        tt.fields.mu,
			}
			if got := store.addKey(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryStore_GetLimiter(t *testing.T) {
	type fields struct {
		keys      map[string]Limiter
		limiterFn func(store *MemoryStore, key string) Limiter
		mu        *sync.RWMutex
	}
	type args struct {
		key string
	}
	limiter := NewTokenLimiter(1, 3)
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Limiter
	}{
		{
			name: "",
			fields: fields{
				keys: make(map[string]Limiter),
				limiterFn: func(store *MemoryStore, key string) Limiter {
					return limiter
				},
				mu: &sync.RWMutex{},
			},
			args: args{
				key: "192.168.1.6:8080",
			},
			want: limiter,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &MemoryStore{
				keys:      tt.fields.keys,
				limiterFn: tt.fields.limiterFn,
				mu:        tt.fields.mu,
			}
			if got := store.GetLimiter(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLimiter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryStore_DeleteKey(t *testing.T) {
	type fields struct {
		keys      map[string]Limiter
		limiterFn func(store *MemoryStore, key string) Limiter
		mu        *sync.RWMutex
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "",
			fields: fields{
				keys: make(map[string]Limiter),
				limiterFn: func(store *MemoryStore, key string) Limiter {
					return NewTokenLimiter(1, 3)
				},
				mu: &sync.RWMutex{},
			},
			args: args{
				key: "192.168.1.6:8080",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &MemoryStore{
				keys:      tt.fields.keys,
				limiterFn: tt.fields.limiterFn,
				mu:        tt.fields.mu,
			}
			store.addKey(tt.args.key)
			if _, exists := store.keys[tt.args.key]; !exists {
				t.Error("key should exists")
			}
			store.DeleteKey(tt.args.key)
			if _, exists := store.keys[tt.args.key]; exists {
				t.Error("key should not exists")
			}
		})
	}
}

func TestNewMemoryStore(t *testing.T) {
	type args struct {
		opts []MemoryStoreOption
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				opts: []MemoryStoreOption{
					WithLimiterFn(func(store *MemoryStore, key string) Limiter {
						return NewTokenLimiter(1, 3)
					}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemoryStore(tt.args.opts...); got == nil {
				t.Error("NewMemoryStore() = nil")
			}
		})
	}
}
