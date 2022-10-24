package grpcx_ratelimit

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/unionj-cloud/go-doudou/v2/framework/ratelimit"
	"github.com/unionj-cloud/go-doudou/v2/framework/ratelimit/memrate"
	"github.com/unionj-cloud/go-doudou/v2/framework/ratelimit/redisrate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type GetKey func(ctx context.Context, fullMethod string) string

type RateLimitInterceptor struct {
	mstore *memrate.MemoryStore
	rdb    redisrate.Rediser
	fn     redisrate.LimitFn
}

const (
	limit = 1000
	burst = 10000
	ttl   = 60 * time.Second
)

func (r *RateLimitInterceptor) limitKey(key string) bool {
	if r == nil {
		r = &RateLimitInterceptor{
			mstore: memrate.NewMemoryStore(func(_ context.Context, store *memrate.MemoryStore, key string) ratelimit.Limiter {
				return memrate.NewLimiter(limit, burst, memrate.WithTimer(10*time.Second, func() {
					store.DeleteKey(key)
				}))
			}),
		}
	}
	if r.rdb != nil && r.fn != nil {
		return redisrate.NewGcraLimiterLimitFn(r.rdb, key, r.fn).Allow()
	} else if r.mstore != nil {
		return r.mstore.GetLimiter(key).Allow()
	}
	return false
}

type RateLimitInterceptorOption func(*RateLimitInterceptor)

func WithMemoryStore(mstore *memrate.MemoryStore) RateLimitInterceptorOption {
	return func(interceptor *RateLimitInterceptor) {
		interceptor.mstore = mstore
	}
}

func WithRedisStore(rdb redisrate.Rediser, fn redisrate.LimitFn) RateLimitInterceptorOption {
	return func(interceptor *RateLimitInterceptor) {
		interceptor.rdb = rdb
		interceptor.fn = fn
	}
}

func NewRateLimitInterceptor(options ...RateLimitInterceptorOption) *RateLimitInterceptor {
	interceptor := RateLimitInterceptor{
		mstore: memrate.NewMemoryStore(func(_ context.Context, store *memrate.MemoryStore, key string) ratelimit.Limiter {
			return memrate.NewLimiter(limit, burst, memrate.WithTimer(ttl, func() {
				store.DeleteKey(key)
			}))
		}),
	}
	for i := range options {
		options[i](&interceptor)
	}
	return &interceptor
}

func (r *RateLimitInterceptor) UnaryServerInterceptor(getKey GetKey) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		key := getKey(ctx, info.FullMethod)
		if !r.limitKey(key) {
			return nil, status.Errorf(codes.ResourceExhausted, "%s is rejected by grpcx_ratelimit middleware, please retry later.", key)
		}
		return handler(ctx, req)
	}
}

func (r *RateLimitInterceptor) StreamServerInterceptor(getKey GetKey) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		if stream1, ok := stream.(*grpc_middleware.WrappedServerStream); ok {
			ctx = stream1.WrappedContext
		}
		key := getKey(ctx, info.FullMethod)
		if !r.limitKey(key) {
			return status.Errorf(codes.ResourceExhausted, "%s is rejected by grpcx_ratelimit middleware, please retry later.", key)
		}
		return handler(srv, stream)
	}
}
