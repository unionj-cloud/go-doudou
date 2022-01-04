package ratelimit

import (
	"context"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

type MemoryLimiter struct {
	limiter *rate.Limiter
	keyFunc func(req *http.Request) string
}

func (m *MemoryLimiter) SetKeyFunc(keyFunc func(req *http.Request) string) {
	m.keyFunc = keyFunc
}

func (m *MemoryLimiter) Allow() bool {
	return m.limiter.Allow()
}

func (m *MemoryLimiter) Reserve() (time.Duration, bool) {
	r := m.limiter.Reserve()
	return r.Delay(), r.OK()
}

func (m *MemoryLimiter) Wait(ctx context.Context) error {
	return m.limiter.Wait(ctx)
}
