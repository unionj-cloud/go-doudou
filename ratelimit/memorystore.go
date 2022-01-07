package ratelimit

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type MemoryStore struct {
	keys      map[string]Limiter
	limiterFn func(store *MemoryStore, key string) Limiter
	mu        *sync.RWMutex
}

type MemoryStoreOption func(*MemoryStore)

// WithLimiterFn set function for getting Limiter instance
func WithLimiterFn(fn func(store *MemoryStore, key string) Limiter) MemoryStoreOption {
	return func(ls *MemoryStore) {
		ls.limiterFn = fn
	}
}

func NewMemoryStore(opts ...MemoryStoreOption) *MemoryStore {
	store := &MemoryStore{
		keys: make(map[string]Limiter),
		mu:   &sync.RWMutex{},
	}

	for _, opt := range opts {
		opt(store)
	}

	return store
}

func (store *MemoryStore) addKey(key string) Limiter {
	store.mu.Lock()
	defer store.mu.Unlock()

	// check again if key exists because it's read lock in GetLimiter
	limiter, exists := store.keys[key]
	if exists {
		return limiter
	}

	limiter = store.limiterFn(store, key)
	store.keys[key] = limiter

	return limiter
}

// GetLimiter returns the rate limiter for the provided key if it exists,
// otherwise calls addKey to add key to the map
func (store *MemoryStore) GetLimiter(key string) Limiter {
	store.mu.RLock()
	limiter, exists := store.keys[key]

	if !exists {
		store.mu.RUnlock()
		return store.addKey(key)
	}

	store.mu.RUnlock()
	return limiter
}

func (store *MemoryStore) DeleteKey(key string) {
	store.mu.Lock()
	defer store.mu.Unlock()

	delete(store.keys, key)
	logrus.Debugf("[go-doudou] key %s is deleted from store", key)
}
