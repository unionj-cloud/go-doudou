package database

import (
	"context"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/caches"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
)

var _ caches.Cacher = (*CacherAdapter)(nil)

type CacherAdapter struct {
	marshaler *Marshaler
}

func (c *CacherAdapter) Delete(tag string, tags ...string) error {
	invalidateTags := []string{tag}
	invalidateTags = append(invalidateTags, tags...)
	if err := c.marshaler.Invalidate(context.Background(), store.WithInvalidateTags(invalidateTags)); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *CacherAdapter) Get(key string) *caches.Query {
	result := new(caches.Query)
	_, err := c.marshaler.Get(context.Background(), key, result)
	if err != nil {
		zlogger.Err(err).Msg(err.Error())
		return nil
	}
	return result
}

func (c *CacherAdapter) Store(key string, val *caches.Query) error {
	if err := c.marshaler.Set(context.Background(), key, val, store.WithTags(val.Tags)); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func NewCacherAdapter(cacheManager cache.CacheInterface[any]) *CacherAdapter {
	return &CacherAdapter{
		marshaler: NewMarshaler(cacheManager),
	}
}
