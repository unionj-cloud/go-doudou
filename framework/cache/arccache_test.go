package cache

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewARCCache(t *testing.T) {
	// 测试正常创建
	cache := NewARCCache(100, 1*time.Minute)
	assert.NotNil(t, cache)
	assert.NotNil(t, cache.base)
	assert.Equal(t, 1*time.Minute, cache.base.ttl)

	// 测试尺寸为0时应该会panic
	assert.Panics(t, func() {
		NewARCCache(0, 1*time.Minute)
	})

	// 负数尺寸也应该会panic
	assert.Panics(t, func() {
		NewARCCache(-1, 1*time.Minute)
	})
}

func TestARCCache_Operations(t *testing.T) {
	// 创建ARC缓存
	cache := NewARCCache(10, 1*time.Minute)

	// 测试Set和Get
	key1 := "key1"
	data1 := []byte("data1")

	cache.Set(key1, data1)
	result1, ok := cache.Get(key1)
	assert.True(t, ok)
	assert.Equal(t, data1, result1)

	// 添加多个元素
	for i := 0; i < 10; i++ {
		key := "key" + strconv.Itoa(i)
		data := []byte("data" + strconv.Itoa(i))
		cache.Set(key, data)
	}

	// 再次获取第一个元素，应该仍然存在（ARC应该能保留常用的元素）
	result1, ok = cache.Get(key1)
	assert.True(t, ok)
	assert.Equal(t, data1, result1)

	// 测试删除
	cache.Del(key1)
	_, ok = cache.Get(key1)
	assert.False(t, ok)

	// 测试过期
	shortTTLCache := NewARCCache(10, 50*time.Millisecond)

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

func TestARCCache_CapacityLimit(t *testing.T) {
	// 创建一个小容量的ARC缓存，用于测试超出容量时的行为
	capacity := 5
	cache := NewARCCache(capacity, 1*time.Minute)

	// 添加容量+1个元素
	for i := 0; i < capacity+1; i++ {
		key := "key" + strconv.Itoa(i)
		data := []byte("data" + strconv.Itoa(i))
		cache.Set(key, data)
	}

	// 验证仍然可以获取到最近添加的元素
	latestKey := "key" + strconv.Itoa(capacity)
	latestData := []byte("data" + strconv.Itoa(capacity))
	result, ok := cache.Get(latestKey)
	assert.True(t, ok)
	assert.Equal(t, latestData, result)

	// ARC算法很复杂，不容易简单地断言哪些元素被淘汰，
	// 我们只能验证缓存的总体行为正确
}
