package redisrate

import (
	"context"
	logger "github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	"github.com/unionj-cloud/go-doudou/v2/framework/ratelimit"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const redisPrefix = "go-doudou:rate:"

type Rediser interface {
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
	EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd
	ScriptExists(ctx context.Context, hashes ...string) *redis.BoolSliceCmd
	ScriptLoad(ctx context.Context, script string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type LimitFn func(ctx context.Context) ratelimit.Limit

type GcraLimiter struct {
	rdb     Rediser
	key     string
	limit   ratelimit.Limit
	limitFn LimitFn
}

func (gl *GcraLimiter) AllowCtx(ctx context.Context) bool {
	allow, err := gl.AllowECtx(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("")
		return false
	}
	return allow
}

func (gl *GcraLimiter) Allow() bool {
	allow, err := gl.AllowE()
	if err != nil {
		logger.Error().Err(err).Msg("")
		return false
	}
	return allow
}

// Wait you'd better pass a timeout or cancelable context
func (gl *GcraLimiter) Wait(ctx context.Context) error {
	for {
		retryAfter, allow, err := gl.ReserveECtx(ctx)
		if err != nil {
			return err
		}
		if allow {
			return nil
		}
		time.Sleep(retryAfter)
	}
}

func (gl *GcraLimiter) AllowE() (bool, error) {
	allow, err := gl.AllowN(context.Background(), 1)
	return allow.Allowed > 0, err
}

func (gl *GcraLimiter) ReserveE() (time.Duration, bool, error) {
	allow, err := gl.AllowN(context.Background(), 1)
	if err != nil {
		return 0, false, err
	}
	return allow.RetryAfter, allow.Allowed > 0, nil
}

func (gl *GcraLimiter) AllowECtx(ctx context.Context) (bool, error) {
	allow, err := gl.AllowN(ctx, 1)
	return allow.Allowed > 0, err
}

func (gl *GcraLimiter) ReserveECtx(ctx context.Context) (time.Duration, bool, error) {
	allow, err := gl.AllowN(ctx, 1)
	if err != nil {
		return 0, false, err
	}
	return allow.RetryAfter, allow.Allowed > 0, nil
}

// NewGcraLimiter returns a new Limiter.
func NewGcraLimiter(rdb Rediser, key string, r float64, period time.Duration, b int) ratelimit.Limiter {
	return &GcraLimiter{
		rdb: rdb,
		key: key,
		limit: ratelimit.Limit{
			Rate:   r,
			Burst:  b,
			Period: period,
		},
	}
}

// NewGcraLimiterLimit returns a new Limiter.
func NewGcraLimiterLimit(rdb Rediser, key string, l ratelimit.Limit) ratelimit.Limiter {
	return &GcraLimiter{
		rdb:   rdb,
		key:   key,
		limit: l,
	}
}

// NewGcraLimiterLimitFn returns a new Limiter.
func NewGcraLimiterLimitFn(rdb Rediser, key string, fn LimitFn) ratelimit.Limiter {
	return &GcraLimiter{
		rdb:     rdb,
		key:     key,
		limitFn: fn,
	}
}

// AllowN reports whether n events may happen at time now.
func (gl *GcraLimiter) AllowN(ctx context.Context, n int) (res *Result, err error) {
	limit := gl.limit
	if gl.limitFn != nil {
		limit = gl.limitFn(ctx)
	}
	values := []interface{}{limit.Burst, limit.Rate, limit.Period.Seconds(), n}
	v, err := allowN.Run(ctx, gl.rdb, []string{redisPrefix + gl.key}, values...).Result()
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

// Reset gets a key and reset all limitations and previous usages
func (gl *GcraLimiter) Reset(ctx context.Context) error {
	return gl.rdb.Del(ctx, redisPrefix+gl.key).Err()
}

func dur(f float64) time.Duration {
	if f == -1 {
		return -1
	}
	return time.Duration(f * float64(time.Second))
}

type Result struct {
	// Limit is the limit that was used to obtain this result.
	Limit ratelimit.Limit

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
