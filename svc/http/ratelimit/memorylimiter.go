package ratelimit

import (
	"context"
	"golang.org/x/time/rate"
	"time"
)

// MemoryLimiter wraps rate.Limiter to implement Limiter interface
type MemoryLimiter struct {
	limiter *rate.Limiter
	timeout time.Duration
	last    time.Time
}

func NewMemoryLimiter(r rate.Limit, b int) Limiter {
	return &MemoryLimiter{
		limiter: rate.NewLimiter(r, b),
	}
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

func (m *MemoryLimiter) Expire() bool {
	return time.Now().Sub(m.last) > m.timeout
}
