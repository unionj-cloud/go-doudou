package mock

import (
	"context"
	"github.com/slok/goresilience"
)

//go:generate mockgen -destination ./mock_goresilience_runner_interface.go -package mock -source=./goresilience_runner_interface.go

type Runner interface {
	// Run will run the unit of execution passed on f.
	Run(ctx context.Context, f goresilience.Func) error
}
