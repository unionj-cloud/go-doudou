package ratelimit

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type Limiter interface {
	Allow() bool
	Reserve() (time.Duration, bool)
	Wait(ctx context.Context) error
	Expire() bool
}

type LimiterStore struct {
	keys          map[string]Limiter
	keyFn         func(req *http.Request) string
	limiterFn     func() Limiter
	mu            *sync.RWMutex
	clearInterval time.Duration
}

type LimiterStoreOption func(*LimiterStore)

// WithKeyFunc set function for getting key
func WithKeyFunc(fn func(req *http.Request) string) LimiterStoreOption {
	return func(ls *LimiterStore) {
		ls.keyFn = fn
	}
}

// WithLimiterFunc set function for getting Limiter instance
func WithLimiterFunc(fn func() Limiter) LimiterStoreOption {
	return func(ls *LimiterStore) {
		ls.limiterFn = fn
	}
}

// WithClearInterval set clear interval for clearing keys that didn't request tokens for a long time
// default is zero, means not clear
func WithClearInterval(ci time.Duration) LimiterStoreOption {
	return func(ls *LimiterStore) {
		ls.clearInterval = ci
	}
}

func NewLimiterStore(opts ...LimiterStoreOption) *LimiterStore {
	store := &LimiterStore{
		keys: make(map[string]Limiter),
		mu:   &sync.RWMutex{},
	}

	for _, opt := range opts {
		opt(store)
	}

	return store
}

// TODO
func (store *LimiterStore) clear() {
	store.mu.Lock()
	defer store.mu.Unlock()
	for key, limiter := range store.keys {
		if limiter.Expire() {
			delete(store.keys, key)
		}
	}
}

// AddKey adds key and limiter pair
func (store *LimiterStore) AddKey(key string) Limiter {
	store.mu.Lock()
	defer store.mu.Unlock()

	limiter := store.limiterFn()
	store.keys[key] = limiter

	return limiter
}

// GetLimiter returns the rate limiter for the provided key if it exists,
// otherwise calls AddKey to add key to the map
func (store *LimiterStore) GetLimiter(key string) Limiter {
	store.mu.Lock()
	limiter, exists := store.keys[key]

	if !exists {
		store.mu.Unlock()
		return store.AddKey(key)
	}

	store.mu.Unlock()
	return limiter
}
