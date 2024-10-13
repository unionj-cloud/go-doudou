package cache

import (
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/toolkit/caller"
)

type ARCCache struct {
	*base
}

func NewARCCache(size int, ttl time.Duration) *ARCCache {
	store, err := lru.NewARC(size)
	if err != nil {
		panic(errors.Wrap(err, caller.NewCaller().String()))
	}
	return &ARCCache{
		newBase(store, ttl),
	}
}
