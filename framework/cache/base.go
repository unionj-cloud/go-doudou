package cache

import (
	"time"

	"golang.org/x/exp/rand"
)

type IStore interface {
	Get(key interface{}) (value interface{}, ok bool)
	Add(key, value interface{})
	Remove(key interface{})
}

type base struct {
	store  IStore
	ttl    time.Duration
	offset time.Duration
	rand   *rand.Rand
}

func (l *base) Set(key string, data []byte) {
	ttl := l.ttl
	if l.offset > 0 {
		ttl += time.Duration(l.rand.Int63n(int64(l.offset)))
	}
	l.store.Add(key, &Item{
		Key:      key,
		Value:    data,
		ExpireAt: time.Now().Add(ttl),
	})
}

func (l *base) Get(key string) ([]byte, bool) {
	value, ok := l.store.Get(key)
	if !ok {
		return nil, false
	}
	item := value.(*Item)
	if item.expired() {
		l.store.Remove(key)
		return nil, false
	}
	return item.Value.([]byte), true
}

func (l *base) Del(key string) {
	l.store.Remove(key)
}

const maxOffset = 10 * time.Second

func newBase(store IStore, ttl time.Duration) *base {
	offset := ttl / 10
	if offset > maxOffset {
		offset = maxOffset
	}
	return &base{
		store:  store,
		ttl:    ttl,
		offset: offset,
		rand:   rand.New(rand.NewSource(uint64(time.Now().UnixNano()))),
	}
}
