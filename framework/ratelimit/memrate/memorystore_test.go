package memrate

import (
	"context"
	lru "github.com/hashicorp/golang-lru"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/framework/ratelimit"
	"reflect"
	"testing"
)

func TestMemoryStore_addKey(t *testing.T) {
	type fields struct {
		keys      map[string]ratelimit.Limiter
		limiterFn LimiterFn
	}
	type args struct {
		key string
	}
	limiter := NewLimiter(1, 3)
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ratelimit.Limiter
	}{
		{
			name: "",
			fields: fields{
				keys: make(map[string]ratelimit.Limiter),
				limiterFn: func(ctx context.Context, store *MemoryStore, key string) ratelimit.Limiter {
					return NewLimiter(1, 3)
				},
			},
			args: args{
				key: "192.168.1.6:8080",
			},
			want: limiter,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, _ := lru.New(256)
			store := &MemoryStore{
				keys:      keys,
				limiterFn: tt.fields.limiterFn,
			}
			if got := store.addKeyCtx(context.Background(), tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryStore_GetLimiter(t *testing.T) {
	type fields struct {
		keys      map[string]ratelimit.Limiter
		limiterFn LimiterFn
	}
	type args struct {
		key string
	}
	limiter := NewLimiter(1, 3)
	key := "192.168.1.6:8080"
	keys, _ := lru.New(256)
	store := &MemoryStore{
		keys: keys,
		limiterFn: func(ctx context.Context, store *MemoryStore, key string) ratelimit.Limiter {
			return NewLimiter(1, 3)
		},
	}
	assert.Equal(t, limiter, store.GetLimiter(key))
	assert.Equal(t, limiter, store.GetLimiter(key))
}

func TestMemoryStore_DeleteKey(t *testing.T) {
	key := "192.168.1.6:8080"
	keys, _ := lru.New(256)
	store := &MemoryStore{
		keys: keys,
		limiterFn: func(ctx context.Context, store *MemoryStore, key string) ratelimit.Limiter {
			return NewLimiter(1, 3)
		},
	}
	store.addKeyCtx(context.Background(), key)
	store.addKeyCtx(context.Background(), key)
	if exists := store.keys.Contains(key); !exists {
		t.Error("key should exists")
	}
	store.DeleteKey(key)
	if exists := store.keys.Contains(key); exists {
		t.Error("key should not exists")
	}
}

func TestNewMemoryStore(t *testing.T) {
	type args struct {
		fn   LimiterFn
		opts []MemoryStoreOption
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				fn: func(ctx context.Context, store *MemoryStore, key string) ratelimit.Limiter {
					return NewLimiter(1, 3)
				},
				opts: nil,
			},
		},
		{
			name: "",
			args: args{
				fn: func(ctx context.Context, store *MemoryStore, key string) ratelimit.Limiter {
					return NewLimiter(1, 3)
				},
				opts: []MemoryStoreOption{
					WithMaxKeys(100),
					WithOnEvicted(func(key interface{}, value interface{}) {
						log.Println(key, value)
					}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemoryStore(tt.args.fn, tt.args.opts...); got == nil {
				t.Error("NewMemoryStore() = nil")
			}
		})
	}
}
