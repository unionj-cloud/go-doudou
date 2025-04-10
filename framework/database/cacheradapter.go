package database

import (
	"context"
	"github.com/unionj-cloud/toolkit/zlogger"

	"github.com/pkg/errors"
	"github.com/unionj-cloud/toolkit/caches"
	"github.com/unionj-cloud/toolkit/gocache/lib/cache"
	"github.com/unionj-cloud/toolkit/gocache/lib/store"
)

var _ caches.Cacher = (*CacherAdapter)(nil)

type CacherAdapterConfig struct {
	MarshalerConfig
}

type CacherAdapter struct {
	marshaler *Marshaler
	conf      CacherAdapterConfig
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
		zlogger.Info().Msgf("Cache missing: key %s %s\n", key, err.Error())
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

func NewCacherAdapter(cacheManager cache.CacheInterface[any], config CacherAdapterConfig) *CacherAdapter {
	return &CacherAdapter{
		marshaler: NewMarshaler(cacheManager, config.MarshalerConfig),
	}
}
