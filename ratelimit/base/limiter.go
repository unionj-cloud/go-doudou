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

	ReserveE() (time.Duration, bool, error)
	ReserveECtx(ctx context.Context) (time.Duration, bool, error)

	Wait(ctx context.Context) error
}
