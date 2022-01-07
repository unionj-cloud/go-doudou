package ratelimit

import (
	"context"
	"golang.org/x/time/rate"
	"testing"
	"time"
)

func TestNewTokenLimiter(t *testing.T) {
	type args struct {
		r    rate.Limit
		b    int
		opts []TokenLimiterOption
	}
	store := NewMemoryStore()
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				r: 1,
				b: 3,
				opts: []TokenLimiterOption{
					WithTimer(10*time.Second, func() {
						store.DeleteKey("any")
					}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTokenLimiter(tt.args.r, tt.args.b, tt.args.opts...); got == nil {
				t.Error("got should not be nil")
			}
		})
	}
}

func TestTokenLimiter_Allow(t *testing.T) {
	tl := NewTokenLimiter(1, 3)
	if got := tl.Allow(); got != true {
		t.Errorf("Allow() should return true")
	}
	if got := tl.Allow(); got != true {
		t.Errorf("Allow() should return true")
	}
	if got := tl.Allow(); got != true {
		t.Errorf("Allow() should return true")
	}
	if got := tl.Allow(); got != false {
		t.Errorf("Allow() should return false")
	}
}

func TestTokenLimiter_Reserve(t *testing.T) {
	tl := NewTokenLimiter(1, 3)
	if d, ok := tl.Reserve(); ok != true && d != 0 {
		t.Errorf("Reserve() should return true and d should equal to 0")
	}
	if d, ok := tl.Reserve(); ok != true && d != 0 {
		t.Errorf("Reserve() should return true and d should equal to 0")
	}
	if d, ok := tl.Reserve(); ok != true && d != 0 {
		t.Errorf("Reserve() should return true and d should equal to 0")
	}
	if d, ok := tl.Reserve(); ok != true && d <= 0 {
		t.Errorf("Reserve() should return true and d should greater than 0")
	}
	tl = NewTokenLimiter(1, 0)
	if _, ok := tl.Reserve(); ok != false {
		t.Errorf("Reserve() should return false")
	}
}

func TestTokenLimiter_Wait(t *testing.T) {
	tl := NewTokenLimiter(1, 3)
	ctx := context.Background()
	if err := tl.Wait(ctx); err != nil {
		t.Errorf("Wait() shouldn't return error")
	}
	if err := tl.Wait(ctx); err != nil {
		t.Errorf("Wait() shouldn't return error")
	}
	if err := tl.Wait(ctx); err != nil {
		t.Errorf("Wait() shouldn't return error")
	}
	ctx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
	defer cancel()
	if err := tl.Wait(ctx); err.Error() != "rate: Wait(n=1) would exceed context deadline" {
		t.Errorf("Wait() should return error: rate: Wait(n=1) would exceed context deadline, but actual error: %s", err.Error())
	}
}
