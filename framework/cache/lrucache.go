package cache

import (
	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/caller"
	"time"
)

type LruCache struct {
	*base
}

func NewLruCache(size int, ttl time.Duration) *LruCache {
	store, err := lru.New(size)
	if err != nil {
		panic(errors.Wrap(err, caller.NewCaller().String()))
	}
	return &LruCache{
		newBase(&LruCacheAdapter{store}, ttl),
	}
}

type LruCacheAdapter struct {
	store *lru.Cache
}

func (l *LruCacheAdapter) Get(key interface{}) (value interface{}, ok bool) {
	return l.store.Get(key)
}

func (l *LruCacheAdapter) Add(key, value interface{}) {
	l.store.Add(key, value)
}

func (l *LruCacheAdapter) Remove(key interface{}) {
	l.store.Remove(key)
}
