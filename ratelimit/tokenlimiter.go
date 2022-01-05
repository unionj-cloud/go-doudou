package ratelimit

import (
	"context"
	"golang.org/x/time/rate"
	"time"
)

// TokenLimiter wraps rate.Limiter to implement Limiter interface
type TokenLimiter struct {
	limiter *rate.Limiter
	timeout time.Duration
	timer   *time.Timer
}

type TokenLimiterOption func(*TokenLimiter)

func WithTimer(timeout time.Duration, fn func()) TokenLimiterOption {
	return func(tl *TokenLimiter) {
		tl.timeout = timeout
		tl.timer = time.AfterFunc(timeout, fn)
	}
}

func NewTokenLimiter(r rate.Limit, b int, opts ...TokenLimiterOption) Limiter {
	tl := &TokenLimiter{
		limiter: rate.NewLimiter(r, b),
	}
	for _, opt := range opts {
		opt(tl)
	}
	return tl
}

func (tl *TokenLimiter) after() {
	if tl.timer != nil && tl.timer.Stop() {
		tl.timer.Reset(tl.timeout)
	}
}

func (tl *TokenLimiter) Allow() bool {
	ok := tl.limiter.Allow()
	tl.after()
	return ok
}

func (tl *TokenLimiter) Reserve() (time.Duration, bool) {
	r := tl.limiter.Reserve()
	tl.after()
	return r.Delay(), r.OK()
}

func (tl *TokenLimiter) Wait(ctx context.Context) error {
	err := tl.limiter.Wait(ctx)
	tl.after()
	return err
}
