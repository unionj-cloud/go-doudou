package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 实现一个简单的内存存储来测试base
type mockStore struct {
	data map[interface{}]interface{}
}

func newMockStore() *mockStore {
	return &mockStore{
		data: make(map[interface{}]interface{}),
	}
}

func (m *mockStore) Get(key interface{}) (value interface{}, ok bool) {
	value, ok = m.data[key]
	return
}

func (m *mockStore) Add(key, value interface{}) {
	m.data[key] = value
}

func (m *mockStore) Remove(key interface{}) {
	delete(m.data, key)
}

func TestNewBase(t *testing.T) {
	// 测试正常ttl
	store := newMockStore()
	ttl := 10 * time.Minute
	b := newBase(store, ttl)

	assert.Equal(t, store, b.store)
	assert.Equal(t, ttl, b.ttl)
	assert.Equal(t, ttl/60, b.offset)
	assert.NotNil(t, b.rand)

	// 测试零值ttl
	zeroTTL := time.Duration(0)
	b = newBase(store, zeroTTL)
	assert.Equal(t, store, b.store)
	assert.Equal(t, zeroTTL, b.ttl)
	assert.Equal(t, zeroTTL, b.offset)
	assert.NotNil(t, b.rand)
}

func TestBase_Get(t *testing.T) {
	store := newMockStore()
	ttl := 5 * time.Minute
	b := newBase(store, ttl)

	// 测试获取不存在的键
	val, ok := b.Get("not-exist")
	assert.False(t, ok)
	assert.Nil(t, val)

	// 测试获取存在但已过期的键
	expiredItem := &Item{
		Key:      "expired",
		Value:    []byte("expired-value"),
		ExpireAt: time.Now().Add(-1 * time.Hour),
	}
	store.Add("expired", expiredItem)
	val, ok = b.Get("expired")
	assert.False(t, ok, "过期项应该返回false")
	assert.Nil(t, val, "过期项不应该返回值")
	// 验证过期项已被删除
	_, exists := store.Get("expired")
	assert.False(t, exists, "过期项应该被删除")

	// 测试获取有效项
	validItem := &Item{
		Key:      "valid",
		Value:    []byte("valid-value"),
		ExpireAt: time.Now().Add(1 * time.Hour),
	}
	store.Add("valid", validItem)
	val, ok = b.Get("valid")
	assert.True(t, ok)
	assert.Equal(t, []byte("valid-value"), val)
}

func TestBase_Set(t *testing.T) {
	store := newMockStore()
	ttl := 5 * time.Minute
	b := newBase(store, ttl)

	// 测试设置值
	b.Set("test-key", []byte("test-value"))

	// 验证值已经被设置
	item, ok := store.Get("test-key")
	assert.True(t, ok)

	// 检查Item内容
	itemObj, ok := item.(*Item)
	assert.True(t, ok)
	assert.Equal(t, "test-key", itemObj.Key)
	assert.Equal(t, []byte("test-value"), itemObj.Value)

	// 检查过期时间（不精确匹配，只检查是否在合理范围内）
	expectedExpiry := time.Now().Add(ttl)
	diff := expectedExpiry.Sub(itemObj.ExpireAt)
	assert.LessOrEqual(t, diff.Abs(), ttl/60+time.Second, "过期时间应该在预期附近")
}

func TestBase_Del(t *testing.T) {
	store := newMockStore()
	ttl := 5 * time.Minute
	b := newBase(store, ttl)

	// 添加一个测试键
	testItem := &Item{
		Key:      "test-key",
		Value:    []byte("test-value"),
		ExpireAt: time.Now().Add(ttl),
	}
	store.Add("test-key", testItem)

	// 验证键存在
	_, exists := store.Get("test-key")
	assert.True(t, exists)

	// 删除键
	b.Del("test-key")

	// 验证键已被删除
	_, exists = store.Get("test-key")
	assert.False(t, exists)
}
