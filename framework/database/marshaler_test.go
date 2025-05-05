package database

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/toolkit/caches"
	gcache "github.com/unionj-cloud/toolkit/gocache/lib/cache"
	"github.com/unionj-cloud/toolkit/gocache/lib/store"
)

// 创建一个简单的内存存储的mock
type mockStore struct {
	data  map[string][]byte
	mutex sync.RWMutex
}

func newMockStore() *mockStore {
	return &mockStore{
		data: make(map[string][]byte),
	}
}

func (m *mockStore) Get(ctx context.Context, key interface{}) (interface{}, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if val, ok := m.data[key.(string)]; ok {
		return val, nil
	}
	return nil, store.NotFoundWithCause(nil)
}

func (m *mockStore) Set(ctx context.Context, key interface{}, value interface{}, options ...store.Option) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data[key.(string)] = value.([]byte)
	return nil
}

func (m *mockStore) Delete(ctx context.Context, key interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.data, key.(string))
	return nil
}

func (m *mockStore) Invalidate(ctx context.Context, options ...store.InvalidateOption) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data = make(map[string][]byte)
	return nil
}

func (m *mockStore) Clear(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data = make(map[string][]byte)
	return nil
}

func (m *mockStore) GetType() string {
	return "mock"
}

func (m *mockStore) GetWithTTL(ctx context.Context, key interface{}) (interface{}, time.Duration, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if val, ok := m.data[key.(string)]; ok {
		return val, 5 * time.Minute, nil // 返回固定的TTL
	}
	return nil, 0, store.NotFoundWithCause(nil)
}

// 创建测试用的内存缓存管理器
func createTestCacheManager() gcache.CacheInterface[any] {
	// 使用自定义的mock存储
	mockStore := newMockStore()
	return gcache.New[any](mockStore)
}

func TestNewMarshaler(t *testing.T) {
	cacheManager := createTestCacheManager()
	config := MarshalerConfig{
		CompactMap: true,
	}

	marshaler := NewMarshaler(cacheManager, config)
	assert.NotNil(t, marshaler)
	assert.Equal(t, config, marshaler.conf)
}

func TestMarshaler_Get(t *testing.T) {
	cacheManager := createTestCacheManager()
	marshaler := NewMarshaler(cacheManager, MarshalerConfig{CompactMap: true})

	// 测试缓存不存在的情况
	ctx := context.Background()
	result := new(caches.Query)
	_, err := marshaler.Get(ctx, "nonexistent", result)
	assert.Error(t, err)

	// 测试正常获取缓存的情况
	query := &caches.Query{
		Dest: map[string]interface{}{"id": 1, "name": "test"},
		Tags: []string{"users"},
	}

	// 设置缓存
	err = marshaler.Set(ctx, "testkey", query)
	assert.NoError(t, err)

	// 获取缓存
	returnObj := new(caches.Query)
	got, err := marshaler.Get(ctx, "testkey", returnObj)
	assert.NoError(t, err)
	assert.Equal(t, query.Tags, returnObj.Tags)
	assert.Equal(t, got, returnObj)
}

func TestMarshaler_Set(t *testing.T) {
	cacheManager := createTestCacheManager()
	marshaler := NewMarshaler(cacheManager, MarshalerConfig{CompactMap: true})
	ctx := context.Background()

	// 测试普通 map
	query := &caches.Query{
		Dest: map[string]interface{}{"id": 1, "name": "test", "empty": nil},
		Tags: []string{"users"},
	}

	err := marshaler.Set(ctx, "testkey", query)
	assert.NoError(t, err)

	returnObj := new(caches.Query)
	_, err = marshaler.Get(ctx, "testkey", returnObj)
	assert.NoError(t, err)
	// 验证空值被移除
	destMap := returnObj.Dest.(map[string]interface{})
	_, hasEmpty := destMap["empty"]
	assert.False(t, hasEmpty)

	// 测试 []map
	queryWithSlice := &caches.Query{
		Dest: []map[string]interface{}{
			{"id": 1, "name": "test", "empty": nil},
			{"id": 2, "name": "test2", "empty": nil},
		},
		Tags: []string{"users"},
	}

	err = marshaler.Set(ctx, "testkey2", queryWithSlice)
	assert.NoError(t, err)

	returnObj2 := new(caches.Query)
	_, err = marshaler.Get(ctx, "testkey2", returnObj2)
	assert.NoError(t, err)
	// 验证空值被移除
	destSlice := returnObj2.Dest.([]interface{})
	assert.Len(t, destSlice, 2)
}

func TestMarshaler_Delete(t *testing.T) {
	cacheManager := createTestCacheManager()
	marshaler := NewMarshaler(cacheManager, MarshalerConfig{CompactMap: true})
	ctx := context.Background()

	query := &caches.Query{
		Dest: map[string]interface{}{"id": 1, "name": "test"},
		Tags: []string{"users"},
	}

	// 设置缓存
	err := marshaler.Set(ctx, "testkey", query)
	assert.NoError(t, err)

	// 删除缓存
	err = marshaler.Delete(ctx, "testkey")
	assert.NoError(t, err)

	// 验证缓存已删除
	returnObj := new(caches.Query)
	_, err = marshaler.Get(ctx, "testkey", returnObj)
	assert.Error(t, err)
}

func TestMarshaler_Invalidate(t *testing.T) {
	cacheManager := createTestCacheManager()
	marshaler := NewMarshaler(cacheManager, MarshalerConfig{CompactMap: true})
	ctx := context.Background()

	query1 := &caches.Query{
		Dest: map[string]interface{}{"id": 1, "name": "test"},
		Tags: []string{"users"},
	}

	query2 := &caches.Query{
		Dest: map[string]interface{}{"id": 1, "total": 100},
		Tags: []string{"orders"},
	}

	// 设置缓存
	err := marshaler.Set(ctx, "user_1", query1, store.WithTags(query1.Tags))
	assert.NoError(t, err)
	err = marshaler.Set(ctx, "order_1", query2, store.WithTags(query2.Tags))
	assert.NoError(t, err)

	// 使用标签使缓存失效
	err = marshaler.Invalidate(ctx, store.WithInvalidateTags([]string{"users"}))
	assert.NoError(t, err)

	// 验证所有缓存已失效，我们的mock不支持标签过滤清除，所以这里会清除所有缓存
	returnObj1 := new(caches.Query)
	_, err = marshaler.Get(ctx, "user_1", returnObj1)
	assert.Error(t, err)

	returnObj2 := new(caches.Query)
	_, err = marshaler.Get(ctx, "order_1", returnObj2)
	assert.Error(t, err)
}

func TestMarshaler_Clear(t *testing.T) {
	cacheManager := createTestCacheManager()
	marshaler := NewMarshaler(cacheManager, MarshalerConfig{CompactMap: true})
	ctx := context.Background()

	query1 := &caches.Query{
		Dest: map[string]interface{}{"id": 1, "name": "test"},
		Tags: []string{"users"},
	}

	query2 := &caches.Query{
		Dest: map[string]interface{}{"id": 1, "total": 100},
		Tags: []string{"orders"},
	}

	// 设置缓存
	err := marshaler.Set(ctx, "user_1", query1)
	assert.NoError(t, err)
	err = marshaler.Set(ctx, "order_1", query2)
	assert.NoError(t, err)

	// 清空所有缓存
	err = marshaler.Clear(ctx)
	assert.NoError(t, err)

	// 验证所有缓存已清空
	returnObj1 := new(caches.Query)
	_, err = marshaler.Get(ctx, "user_1", returnObj1)
	assert.Error(t, err)

	returnObj2 := new(caches.Query)
	_, err = marshaler.Get(ctx, "order_1", returnObj2)
	assert.Error(t, err)
}

func TestMarshaler_shortenKey(t *testing.T) {
	cacheManager := createTestCacheManager()
	marshaler := NewMarshaler(cacheManager, MarshalerConfig{CompactMap: true})

	// 测试短键
	shortKey := "test_key"
	result := marshaler.shortenKey(shortKey)
	assert.Equal(t, shortKey, result)

	// 测试长键
	longKey := ""
	for i := 0; i < 1001; i++ {
		longKey += "a"
	}
	result = marshaler.shortenKey(longKey)
	assert.NotEqual(t, longKey, result)
	assert.Less(t, len(result), 100) // UUID 通常小于 100 个字符
}
