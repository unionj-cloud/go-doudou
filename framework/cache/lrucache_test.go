package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLruCache(t *testing.T) {
	// 测试正常创建
	cache := NewLruCache(100, 1*time.Minute)
	assert.NotNil(t, cache)
	assert.NotNil(t, cache.base)
	assert.Equal(t, 1*time.Minute, cache.base.ttl)

	// 测试尺寸为0时应该会panic
	assert.Panics(t, func() {
		NewLruCache(0, 1*time.Minute)
	})

	// 负数尺寸也应该会panic
	assert.Panics(t, func() {
		NewLruCache(-1, 1*time.Minute)
	})
}

func TestLruCacheAdapter(t *testing.T) {
	// 创建LruCache
	cache := NewLruCache(10, 1*time.Minute)

	// 使用底层的Adapter
	adapter, ok := cache.base.store.(*LruCacheAdapter)
	assert.True(t, ok)

	// 测试Add和Get操作
	key := "test-key"
	value := "test-value"

	// 添加键值
	adapter.Add(key, value)

	// 获取并验证
	got, ok := adapter.Get(key)
	assert.True(t, ok)
	assert.Equal(t, value, got)

	// 移除键
	adapter.Remove(key)

	// 确认已删除
	_, ok = adapter.Get(key)
	assert.False(t, ok)
}

func TestLruCache_Operations(t *testing.T) {
	// 创建容量有限的LruCache
	capacity := 2
	cache := NewLruCache(capacity, 1*time.Minute)

	// 测试Set和Get
	key1 := "key1"
	data1 := []byte("data1")

	cache.Set(key1, data1)
	result1, ok := cache.Get(key1)
	assert.True(t, ok)
	assert.Equal(t, data1, result1)

	// 添加第二个元素
	key2 := "key2"
	data2 := []byte("data2")

	cache.Set(key2, data2)
	result2, ok := cache.Get(key2)
	assert.True(t, ok)
	assert.Equal(t, data2, result2)

	// 添加第三个元素，这应该会淘汰最早的元素key1
	key3 := "key3"
	data3 := []byte("data3")

	cache.Set(key3, data3)
	result3, ok := cache.Get(key3)
	assert.True(t, ok)
	assert.Equal(t, data3, result3)

	// 验证key1已被淘汰
	_, ok = cache.Get(key1)
	assert.False(t, ok)

	// 验证key2仍然存在
	result2, ok = cache.Get(key2)
	assert.True(t, ok)
	assert.Equal(t, data2, result2)

	// 测试删除
	cache.Del(key2)
	_, ok = cache.Get(key2)
	assert.False(t, ok)

	// 测试过期
	shortTTLCache := NewLruCache(10, 50*time.Millisecond)

	expKey := "exp-key"
	expData := []byte("exp-data")

	shortTTLCache.Set(expKey, expData)

	// 刚设置后应该能获取到
	expResult, ok := shortTTLCache.Get(expKey)
	assert.True(t, ok)
	assert.Equal(t, expData, expResult)

	// 等待过期
	time.Sleep(100 * time.Millisecond)

	// 过期后应该获取不到
	_, ok = shortTTLCache.Get(expKey)
	assert.False(t, ok)
}
