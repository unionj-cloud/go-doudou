package caches

import (
	"time"
)

type mockTask struct {
	delay  time.Duration
	actRes string
	expRes string
	id     string
}

func (q *mockTask) GetId() string {
	return q.id
}

func (q *mockTask) Run() {
	time.Sleep(q.delay)
	q.actRes = q.expRes
}
