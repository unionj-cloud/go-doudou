package memrate

import (
	"context"
	lru "github.com/hashicorp/golang-lru"
	"github.com/hashicorp/golang-lru/simplelru"
	"github.com/unionj-cloud/go-doudou/ratelimit/base"
	"github.com/unionj-cloud/go-doudou/svc/logger"
	"sync"
)

const defaultMaxKeys = 256

type LimiterFn func(ctx context.Context, store *MemoryStore, key string) base.Limiter

type MemoryStore struct {
	keys      *lru.Cache
	maxKeys   int
	onEvicted simplelru.EvictCallback
	limiterFn LimiterFn
	mu        *sync.RWMutex
}

type MemoryStoreOption func(*MemoryStore)

// WithMaxKeys set maxKeys
func WithMaxKeys(maxKeys int) MemoryStoreOption {
	return func(ls *MemoryStore) {
		ls.maxKeys = maxKeys
	}
}

// WithOnEvicted set onEvicted
func WithOnEvicted(onEvicted func(key interface{}, value interface{})) MemoryStoreOption {
	return func(ls *MemoryStore) {
		ls.onEvicted = onEvicted
	}
}

func NewMemoryStore(fn LimiterFn, opts ...MemoryStoreOption) *MemoryStore {
	store := &MemoryStore{
		maxKeys:   defaultMaxKeys,
		limiterFn: fn,
		mu:        &sync.RWMutex{},
	}

	for _, opt := range opts {
		opt(store)
	}

	if store.onEvicted != nil {
		store.keys, _ = lru.NewWithEvict(store.maxKeys, store.onEvicted)
	} else {
		store.keys, _ = lru.New(store.maxKeys)
	}

	return store
}

// GetLimiter returns the rate limiter for the provided key if it exists,
// otherwise calls addKey to add key to the map
func (store *MemoryStore) GetLimiter(key string) base.Limiter {
	return store.GetLimiterCtx(context.Background(), key)
}

func (store *MemoryStore) addKeyCtx(ctx context.Context, key string) base.Limiter {
	store.mu.Lock()
	defer store.mu.Unlock()

	limiter, exists := store.keys.Get(key)
	if exists {
		return limiter.(base.Limiter)
	}

	limiter = store.limiterFn(ctx, store, key)
	store.keys.Add(key, limiter)

	return limiter.(base.Limiter)
}

// GetLimiterCtx returns the rate limiter for the provided key if it exists,
// otherwise calls addKey to add key to the map
func (store *MemoryStore) GetLimiterCtx(ctx context.Context, key string) base.Limiter {
	store.mu.RLock()

	limiter, exists := store.keys.Get(key)
	if !exists {
		store.mu.RUnlock()
		return store.addKeyCtx(ctx, key)
	}

	store.mu.RUnlock()
	return limiter.(base.Limiter)
}

func (store *MemoryStore) DeleteKey(key string) {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.keys.Remove(key)
	logger.Debugf("[go-doudou] key %s is deleted from store", key)
}
