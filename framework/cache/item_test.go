package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestItem_expired(t *testing.T) {
	// 创建一个已过期的Item
	pastItem := Item{
		Key:      "past",
		Value:    "past-value",
		ExpireAt: time.Now().Add(-1 * time.Hour),
	}
	assert.True(t, pastItem.expired(), "过期时间在过去，应该返回true")

	// 创建一个未来到期的Item
	futureItem := Item{
		Key:      "future",
		Value:    "future-value",
		ExpireAt: time.Now().Add(1 * time.Hour),
	}
	assert.False(t, futureItem.expired(), "过期时间在未来，应该返回false")

	// 创建一个零时间的Item
	zeroTimeItem := Item{
		Key:      "zero",
		Value:    "zero-value",
		ExpireAt: time.Time{}, // 零值
	}
	assert.False(t, zeroTimeItem.expired(), "零时间应该意味着不过期")

	// 创建一个刚好现在过期的Item（这个测试可能不稳定，因为执行非常接近当前时间）
	nowItem := Item{
		Key:      "now",
		Value:    "now-value",
		ExpireAt: time.Now(),
	}
	// 因为现在的时间可能已经过了那么一点点，所以比较难断言结果，
	// 但是理论上如果完全精确地是现在，应该返回false
	// 这里不作具体断言，只是展示逻辑
	_ = nowItem.expired()
}
