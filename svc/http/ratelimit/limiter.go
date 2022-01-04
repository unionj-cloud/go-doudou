package ratelimit

import (
	"context"
	"net/http"
	"time"
)

type Limiter interface {
	Allow() bool
	Reserve() (time.Duration, bool)
	Wait(ctx context.Context) error
	SetKeyFunc(keyFunc func(req *http.Request) string)
}
