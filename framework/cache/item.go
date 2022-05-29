package cache

import "time"

type Item struct {
	Key      string
	Value    interface{}
	ExpireAt time.Time
}

func (item Item) expired() bool {
	return !item.ExpireAt.IsZero() && time.Now().After(item.ExpireAt)
}
