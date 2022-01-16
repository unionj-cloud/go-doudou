package memrate

import (
	lru "github.com/hashicorp/golang-lru"
	"github.com/hashicorp/golang-lru/simplelru"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/ratelimit/base"
	"sync"
)

const defaultMaxKeys = 256

type MemoryStore struct {
	keys      *lru.Cache
	maxKeys   int
	onEvicted simplelru.EvictCallback
	limiterFn func(store *MemoryStore, key string) base.Limiter
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

func NewMemoryStore(fn func(store *MemoryStore, key string) base.Limiter, opts ...MemoryStoreOption) *MemoryStore {
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

func (store *MemoryStore) addKey(key string) base.Limiter {
	store.mu.Lock()
	defer store.mu.Unlock()

	limiter, exists := store.keys.Get(key)
	if exists {
		return limiter.(base.Limiter)
	}

	limiter = store.limiterFn(store, key)
	store.keys.Add(key, limiter)

	return limiter.(base.Limiter)
}

// GetLimiter returns the rate limiter for the provided key if it exists,
// otherwise calls addKey to add key to the map
func (store *MemoryStore) GetLimiter(key string) base.Limiter {
	store.mu.RLock()

	limiter, exists := store.keys.Get(key)
	if !exists {
		store.mu.RUnlock()
		return store.addKey(key)
	}

	store.mu.RUnlock()
	return limiter.(base.Limiter)
}

func (store *MemoryStore) DeleteKey(key string) {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.keys.Remove(key)
	logrus.Debugf("[go-doudou] key %s is deleted from store", key)
}
