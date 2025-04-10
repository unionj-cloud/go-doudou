package database

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bytedance/sonic/decoder"
	"github.com/pkg/errors"
	"reflect"

	"github.com/bytedance/sonic"
	"github.com/lithammer/shortuuid/v4"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/unionj-cloud/toolkit/caches"
	"github.com/unionj-cloud/toolkit/gocache/lib/cache"
	"github.com/unionj-cloud/toolkit/gocache/lib/store"
	"github.com/unionj-cloud/toolkit/reflectutils"
)

var json = sonic.ConfigDefault

type MarshalerConfig struct {
	CompactMap bool
}

// Marshaler is the struct that marshal and unmarshal cache values
type Marshaler struct {
	cache cache.CacheInterface[any]
	conf  MarshalerConfig
}

// NewMarshaler creates a new marshaler that marshals/unmarshals cache values
func NewMarshaler(cache cache.CacheInterface[any], config MarshalerConfig) *Marshaler {
	return &Marshaler{
		cache: cache,
		conf:  config,
	}
}

type idecoder interface {
	UseInt64()
	Decode(val interface{}) (err error)
}

// Get obtains a value from cache and unmarshal value with given object
func (c *Marshaler) Get(ctx context.Context, key any, returnObj any) (any, error) {
	key = c.shortenKey(key)
	result, err := c.cache.Get(ctx, key)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var dec idecoder
	switch v := result.(type) {
	case []byte:
		dec = decoder.NewStreamDecoder(bytes.NewBuffer(v))
	case string:
		dec = decoder.NewDecoder(v)
	}

	dec.UseInt64()
	if err = dec.Decode(returnObj); err != nil {
		return nil, errors.WithStack(err)
	}

	return returnObj, nil
}

// Set sets a value in cache by marshaling value
func (c *Marshaler) Set(ctx context.Context, key, object any, options ...store.Option) error {
	key = c.shortenKey(key)
	query := object.(*caches.Query)

	if c.conf.CompactMap {
		source := reflectutils.ValueOf(query.Dest).Interface()
		t := fmt.Sprintf("%T", source)
		if t == "map[string]interface {}" {
			compactMap := lo.OmitBy[string, interface{}](source.(map[string]interface{}), func(key string, value interface{}) bool {
				return value == nil || reflect.ValueOf(value).IsZero()
			})
			query.Dest = compactMap
		} else if t == "[]map[string]interface {}" {
			rows := source.([]map[string]interface{})
			_rows := make([]map[string]interface{}, len(rows))
			lo.ForEach[map[string]interface{}](rows, func(item map[string]interface{}, index int) {
				compactMap := lo.OmitBy[string, interface{}](item, func(key string, value interface{}) bool {
					return value == nil || reflect.ValueOf(value).IsZero()
				})
				_rows[index] = compactMap
			})
			query.Dest = _rows
		}
	}

	bytes, err := json.Marshal(query)
	if err != nil {
		return err
	}

	return c.cache.Set(ctx, key, bytes, options...)
}

// Delete removes a value from the cache
func (c *Marshaler) Delete(ctx context.Context, key any) error {
	key = c.shortenKey(key)
	return c.cache.Delete(ctx, key)
}

// Invalidate invalidate cache values using given options
func (c *Marshaler) Invalidate(ctx context.Context, options ...store.InvalidateOption) error {
	return c.cache.Invalidate(ctx, options...)
}

// Clear reset all cache data
func (c *Marshaler) Clear(ctx context.Context) error {
	return c.cache.Clear(ctx)
}

func (c *Marshaler) shortenKey(key any) string {
	_key := cast.ToString(key)
	if len(_key) > 1000 {
		return shortuuid.NewWithNamespace(_key)
	}
	return _key
}
