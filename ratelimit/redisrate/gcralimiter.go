package redisrate

import (
	"context"
	"github.com/unionj-cloud/go-doudou/ratelimit/base"
	"github.com/unionj-cloud/go-doudou/svc/logger"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const redisPrefix = "rate:"

type Rediser interface {
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
	EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd
	ScriptExists(ctx context.Context, hashes ...string) *redis.BoolSliceCmd
	ScriptLoad(ctx context.Context, script string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type GcraLimiter struct {
	rdb   Rediser
	key   string
	limit base.Limit
	timer *time.Timer
	mu    sync.RWMutex
}

func (rll *GcraLimiter) AllowCtx(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		if ctx.Err() != nil {
			logger.Error(ctx.Err())
		}
		return false
	default:
		return rll.Allow()
	}
}

func (rll *GcraLimiter) ReserveCtx(ctx context.Context) (time.Duration, bool) {
	select {
	case <-ctx.Done():
		if ctx.Err() != nil {
			logger.Error(ctx.Err())
		}
		return 0, false
	default:
		return rll.Reserve()
	}
}

func (rll *GcraLimiter) Allow() bool {
	allow, err := rll.AllowE()
	if err != nil {
		logger.Error(err)
		return false
	}
	return allow
}

func (rll *GcraLimiter) Reserve() (time.Duration, bool) {
	allow, err := rll.AllowN(context.Background(), rll.key, rll.limit, 1)
	if err != nil {
		logger.Error(err)
	}
	return allow.RetryAfter, allow.Allowed > 0
}

// TODO
func (rll *GcraLimiter) Wait(ctx context.Context) error {
	select {}
}

func (rll *GcraLimiter) AllowE() (bool, error) {
	allow, err := rll.AllowN(context.Background(), rll.key, rll.limit, 1)
	return allow.Allowed > 0, err
}

func (rll *GcraLimiter) ReserveE() (time.Duration, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (rll *GcraLimiter) AllowECtx(ctx context.Context) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (rll *GcraLimiter) ReserveECtx(ctx context.Context) (time.Duration, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (rll *GcraLimiter) resetTimer(resetAfter time.Duration) {
	if rll.timer != nil && rll.timer.Stop() {
		rll.timer.Reset(resetAfter)
	}
}

// NewGcraLimiter returns a new Limiter.
func NewGcraLimiter(rdb Rediser, key string, r float64, period time.Duration, b int) base.Limiter {
	return &GcraLimiter{
		rdb: rdb,
		key: key,
		limit: base.Limit{
			Rate:   r,
			Burst:  b,
			Period: period,
		},
	}
}

// NewGcraLimiterLimit returns a new Limiter.
func NewGcraLimiterLimit(rdb Rediser, key string, l base.Limit) base.Limiter {
	return &GcraLimiter{
		rdb:   rdb,
		key:   key,
		limit: l,
	}
}

// AllowN reports whether n events may happen at time now.
func (rll *GcraLimiter) AllowN(
	ctx context.Context,
	key string,
	limit base.Limit,
	n int,
) (res *Result, err error) {
	defer func() {
		if res != nil {
			rll.resetTimer(res.ResetAfter)
		}
	}()
	values := []interface{}{limit.Burst, limit.Rate, limit.Period.Seconds(), n}
	v, err := allowN.Run(ctx, rll.rdb, []string{redisPrefix + key}, values...).Result()
	if err != nil {
		return nil, err
	}

	values = v.([]interface{})

	retryAfter, err := strconv.ParseFloat(values[2].(string), 64)
	if err != nil {
		return nil, err
	}

	resetAfter, err := strconv.ParseFloat(values[3].(string), 64)
	if err != nil {
		return nil, err
	}

	res = &Result{
		Limit:      limit,
		Allowed:    int(values[0].(int64)),
		Remaining:  int(values[1].(int64)),
		RetryAfter: dur(retryAfter),
		ResetAfter: dur(resetAfter),
	}
	return res, nil
}

func dur(f float64) time.Duration {
	if f == -1 {
		return -1
	}
	return time.Duration(f * float64(time.Second))
}

type Result struct {
	// Limit is the limit that was used to obtain this result.
	Limit base.Limit

	// Allowed is the number of events that may happen at time now.
	Allowed int

	// Remaining is the maximum number of requests that could be
	// permitted instantaneously for this key given the current
	// state. For example, if a rate limiter allows 10 requests per
	// second and has already received 6 requests for this key this
	// second, Remaining would be 4.
	Remaining int

	// RetryAfter is the time until the next request will be permitted.
	// It should be -1 unless the rate limit has been exceeded.
	RetryAfter time.Duration

	// ResetAfter is the time until the RateLimiter returns to its
	// initial state for a given key. For example, if a rate limiter
	// manages requests per second and received one request 200ms ago,
	// Reset would return 800ms. You can also think of this as the time
	// until Limit and Remaining will be equal.
	ResetAfter time.Duration
}
