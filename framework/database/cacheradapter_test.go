package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/toolkit/caches"
	"github.com/unionj-cloud/toolkit/gocache/lib/store"
)

// createTestCacheManager 函数已在 marshaler_test.go 中定义

func TestNewCacherAdapter(t *testing.T) {
	cacheManager := createTestCacheManager()
	config := CacherAdapterConfig{
		MarshalerConfig: MarshalerConfig{
			CompactMap: true,
		},
	}

	adapter := NewCacherAdapter(cacheManager, config)
	assert.NotNil(t, adapter)
	assert.NotNil(t, adapter.marshaler)
}

func TestCacherAdapter_Get(t *testing.T) {
	cacheManager := createTestCacheManager()
	adapter := NewCacherAdapter(cacheManager, CacherAdapterConfig{})
	ctx := context.Background()

	// 测试获取不存在的缓存
	result := adapter.Get("nonexistent")
	assert.Nil(t, result)

	// 测试获取存在的缓存
	query := &caches.Query{
		Dest: map[string]interface{}{"id": 1, "name": "test"},
		Tags: []string{"users"},
	}

	// 通过marshaler直接设置缓存
	err := adapter.marshaler.Set(ctx, "testkey", query, store.WithTags(query.Tags))
	assert.NoError(t, err)

	// 获取缓存
	result = adapter.Get("testkey")
	assert.NotNil(t, result)
	assert.Equal(t, query.Tags, result.Tags)
}

func TestCacherAdapter_Store(t *testing.T) {
	cacheManager := createTestCacheManager()
	adapter := NewCacherAdapter(cacheManager, CacherAdapterConfig{})

	// 测试存储缓存
	query := &caches.Query{
		Dest: map[string]interface{}{"id": 1, "name": "test"},
		Tags: []string{"users"},
	}

	err := adapter.Store("testkey", query)
	assert.NoError(t, err)

	// 验证缓存已存储
	result := adapter.Get("testkey")
	assert.NotNil(t, result)
	assert.Equal(t, query.Tags, result.Tags)
}

func TestCacherAdapter_Delete(t *testing.T) {
	cacheManager := createTestCacheManager()
	adapter := NewCacherAdapter(cacheManager, CacherAdapterConfig{})

	// 创建缓存
	query1 := &caches.Query{
		Dest: map[string]interface{}{"id": 1, "name": "user1"},
		Tags: []string{"users", "user1"},
	}
	query2 := &caches.Query{
		Dest: map[string]interface{}{"id": 2, "name": "user2"},
		Tags: []string{"users", "user2"},
	}

	err := adapter.Store("user1", query1)
	assert.NoError(t, err)
	err = adapter.Store("user2", query2)
	assert.NoError(t, err)

	// 使用标签删除缓存
	err = adapter.Delete("users")
	assert.NoError(t, err)

	// 验证所有带users标签的缓存都已被删除
	result1 := adapter.Get("user1")
	assert.Nil(t, result1)
	result2 := adapter.Get("user2")
	assert.Nil(t, result2)

	// 测试多标签删除
	query3 := &caches.Query{
		Dest: map[string]interface{}{"id": 3, "name": "user3"},
		Tags: []string{"admins", "user3"},
	}
	query4 := &caches.Query{
		Dest: map[string]interface{}{"id": 4, "name": "user4"},
		Tags: []string{"members", "user4"},
	}

	err = adapter.Store("user3", query3)
	assert.NoError(t, err)
	err = adapter.Store("user4", query4)
	assert.NoError(t, err)

	err = adapter.Delete("admins", "members")
	assert.NoError(t, err)

	result3 := adapter.Get("user3")
	assert.Nil(t, result3)
	result4 := adapter.Get("user4")
	assert.Nil(t, result4)
}
