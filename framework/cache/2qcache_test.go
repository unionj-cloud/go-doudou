package cache

import (
	"strconv"
	"testing"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/stretchr/testify/assert"
)

func TestNewTwoQueueCache(t *testing.T) {
	// 测试默认参数创建
	cache := NewTwoQueueCache(100, 1*time.Minute)
	assert.NotNil(t, cache)
	assert.NotNil(t, cache.base)
	assert.Equal(t, 1*time.Minute, cache.base.ttl)
	assert.Equal(t, lru.Default2QRecentRatio, cache.recentRatio)
	assert.Equal(t, lru.Default2QGhostEntries, cache.ghostRatio)

	// 测试自定义参数创建
	customRecentRatio := 0.3
	customGhostRatio := 0.4
	cache = NewTwoQueueCache(
		100,
		1*time.Minute,
		WithRecentRatio(customRecentRatio),
		WithGhostRatio(customGhostRatio),
	)
	assert.Equal(t, customRecentRatio, cache.recentRatio)
	assert.Equal(t, customGhostRatio, cache.ghostRatio)

	// 测试尺寸为0时应该会panic
	assert.Panics(t, func() {
		NewTwoQueueCache(0, 1*time.Minute)
	})

	// 负数尺寸也应该会panic
	assert.Panics(t, func() {
		NewTwoQueueCache(-1, 1*time.Minute)
	})
}

func TestWithRecentRatio(t *testing.T) {
	// 测试WithRecentRatio选项函数
	tqc := &TwoQueueCache{recentRatio: 0.1}
	opt := WithRecentRatio(0.3)
	opt(tqc)
	assert.Equal(t, 0.3, tqc.recentRatio)
}

func TestWithGhostRatio(t *testing.T) {
	// 测试WithGhostRatio选项函数
	tqc := &TwoQueueCache{ghostRatio: 0.1}
	opt := WithGhostRatio(0.3)
	opt(tqc)
	assert.Equal(t, 0.3, tqc.ghostRatio)
}

func TestTwoQueueCache_Operations(t *testing.T) {
	// 创建2Q缓存
	cache := NewTwoQueueCache(10, 1*time.Minute)

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

	// 再次获取第一个元素，应该仍然存在（2Q缓存会根据访问模式保留元素）
	result1, ok = cache.Get(key1)
	assert.True(t, ok)
	assert.Equal(t, data1, result1)

	// 测试删除
	cache.Del(key1)
	_, ok = cache.Get(key1)
	assert.False(t, ok)

	// 测试过期
	shortTTLCache := NewTwoQueueCache(10, 50*time.Millisecond)

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

func TestTwoQueueCache_CapacityLimit(t *testing.T) {
	// 创建小容量的2Q缓存
	capacity := 5
	cache := NewTwoQueueCache(capacity, 1*time.Minute)

	// 添加超出容量的元素
	for i := 0; i < capacity*2; i++ {
		key := "key" + strconv.Itoa(i)
		data := []byte("data" + strconv.Itoa(i))
		cache.Set(key, data)
	}

	// 验证最近添加的元素存在
	latestKey := "key" + strconv.Itoa(capacity*2-1)
	latestData := []byte("data" + strconv.Itoa(capacity*2-1))
	result, ok := cache.Get(latestKey)
	assert.True(t, ok)
	assert.Equal(t, latestData, result)

	// 2Q算法会根据访问频率保留元素，不容易精确断言哪些元素被淘汰，
	// 这里只验证最基本的行为正确
}
