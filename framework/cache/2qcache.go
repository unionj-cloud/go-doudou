package cache

import (
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/caller"
)

type TwoQueueCache struct {
	*base
	recentRatio, ghostRatio float64
}

type TwoQueueCacheOption func(*TwoQueueCache)

func WithRecentRatio(recentRatio float64) TwoQueueCacheOption {
	return func(tqc *TwoQueueCache) {
		tqc.recentRatio = recentRatio
	}
}

func WithGhostRatio(ghostRatio float64) TwoQueueCacheOption {
	return func(tqc *TwoQueueCache) {
		tqc.ghostRatio = ghostRatio
	}
}

func NewTwoQueueCache(size int, ttl time.Duration, options ...TwoQueueCacheOption) *TwoQueueCache {
	tqc := &TwoQueueCache{
		recentRatio: lru.Default2QRecentRatio,
		ghostRatio:  lru.Default2QGhostEntries,
	}
	for _, opt := range options {
		opt(tqc)
	}
	store, err := lru.New2QParams(size, tqc.recentRatio, tqc.ghostRatio)
	if err != nil {
		panic(errors.Wrap(err, caller.NewCaller().String()))
	}
	tqc.base = newBase(store, ttl)
	return tqc
}
