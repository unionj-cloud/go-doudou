package base

import (
	"context"
	"time"
)

type Limiter interface {
	Allow() bool
	AllowE() (bool, error)
	AllowCtx(ctx context.Context) bool
	AllowECtx(ctx context.Context) (bool, error)

	Reserve() (time.Duration, bool)
	ReserveE() (time.Duration, bool, error)
	ReserveCtx(ctx context.Context) (time.Duration, bool)
	ReserveECtx(ctx context.Context) (time.Duration, bool, error)

	Wait(ctx context.Context) error
}
