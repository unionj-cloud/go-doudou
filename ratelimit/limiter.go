package ratelimit

import (
	"context"
	"time"
)

type Limiter interface {
	Allow() bool
	Reserve() (time.Duration, bool)
	Wait(ctx context.Context) error
}
