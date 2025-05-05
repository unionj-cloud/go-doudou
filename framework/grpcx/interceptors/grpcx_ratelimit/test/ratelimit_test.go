package test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/unionj-cloud/go-doudou/v2/framework/grpcx/interceptors/grpcx_ratelimit"
	"github.com/unionj-cloud/go-doudou/v2/framework/ratelimit"
	"github.com/unionj-cloud/go-doudou/v2/framework/ratelimit/memrate"
)

// 模拟 ServerStream 接口
type mockServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (m *mockServerStream) Context() context.Context {
	return m.ctx
}

func (m *mockServerStream) SendMsg(a interface{}) error {
	return nil
}

func (m *mockServerStream) RecvMsg(a interface{}) error {
	return nil
}

// 自定义 KeyGetter 实现
type testKeyGetter struct{}

func (g *testKeyGetter) GetKey(ctx context.Context, fullMethod string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "default-key"
	}
	userIDs := md.Get("user-id")
	if len(userIDs) == 0 {
		return "default-key"
	}
	// 为每个请求创建唯一的键，这样测试时每个请求都有独立的限流器
	return userIDs[0] + ":" + fullMethod
}

// 创建限流器函数 - 每次请求创建新的限流器，让请求总是成功
func createLimiter(_ context.Context, _ *memrate.MemoryStore, key string) ratelimit.Limiter {
	// 设置每秒5个请求，但初始有10个令牌，让测试中的请求都能成功
	return memrate.NewLimiter(5, 10)
}

// 创建限流器函数 - 用于测试并发
func createLimiter10(_ context.Context, _ *memrate.MemoryStore, key string) ratelimit.Limiter {
	// 关键：共享相同用户ID的请求应该共享限流器，这里我们限制每秒10个请求
	return memrate.NewLimiter(10, 10)
}

// 创建限流器函数 - 用于第6个请求必定失败的情况
func createLimiterLimit5(_ context.Context, _ *memrate.MemoryStore, key string) ratelimit.Limiter {
	// 设置限制只有5个请求，没有初始令牌
	return memrate.NewLimiter(1, 5)
}

func TestUnaryServerInterceptor(t *testing.T) {
	// 创建内存存储的限流器，限制只有5个请求通过
	store := memrate.NewMemoryStore(createLimiterLimit5)

	// 创建带有标识符的上下文
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("user-id", "test-user"))

	// 创建模拟的unary处理函数
	mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}

	// 创建拦截器
	interceptor := grpcx_ratelimit.NewRateLimitInterceptor(
		grpcx_ratelimit.WithMemoryStore(store),
	)

	// 创建键获取器
	keyGetter := &testKeyGetter{}

	// 获取拦截器函数
	interceptorFunc := interceptor.UnaryServerInterceptor(keyGetter)

	// 测试限流效果，前5个请求应该成功，后面的应该失败
	for i := 0; i < 5; i++ {
		resp, err := interceptorFunc(ctx, "request", &grpc.UnaryServerInfo{FullMethod: "/test.Service/TestMethod"}, mockHandler)
		assert.NoError(t, err, "第%d个请求应该成功", i+1)
		assert.Equal(t, "success", resp)
	}

	// 第6个请求应该被限流
	resp, err := interceptorFunc(ctx, "request", &grpc.UnaryServerInfo{FullMethod: "/test.Service/TestMethod"}, mockHandler)
	assert.Error(t, err, "第6个请求应该被限流")
	assert.Nil(t, resp)
	assert.Equal(t, codes.ResourceExhausted, status.Code(err))

	// 等待令牌桶重新填充
	time.Sleep(1 * time.Second)

	// 应该可以再次请求
	resp, err = interceptorFunc(ctx, "request", &grpc.UnaryServerInfo{FullMethod: "/test.Service/TestMethod"}, mockHandler)
	assert.NoError(t, err, "等待后的请求应该成功")
	assert.Equal(t, "success", resp)
}

func TestStreamServerInterceptor(t *testing.T) {
	// 创建内存存储的限流器，限制只有5个请求通过
	store := memrate.NewMemoryStore(createLimiterLimit5)

	// 创建带有标识符的上下文
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("user-id", "test-user"))
	stream := &mockServerStream{ctx: ctx}

	// 创建模拟的stream处理函数
	mockHandler := func(srv interface{}, stream grpc.ServerStream) error {
		return nil
	}

	// 创建拦截器
	interceptor := grpcx_ratelimit.NewRateLimitInterceptor(
		grpcx_ratelimit.WithMemoryStore(store),
	)

	// 创建键获取器
	keyGetter := &testKeyGetter{}

	// 获取拦截器函数
	interceptorFunc := interceptor.StreamServerInterceptor(keyGetter)

	// 测试限流效果，前5个请求应该成功，后面的应该失败
	for i := 0; i < 5; i++ {
		err := interceptorFunc(nil, stream, &grpc.StreamServerInfo{FullMethod: "/test.Service/TestStream"}, mockHandler)
		assert.NoError(t, err, "第%d个请求应该成功", i+1)
	}

	// 第6个请求应该被限流
	err := interceptorFunc(nil, stream, &grpc.StreamServerInfo{FullMethod: "/test.Service/TestStream"}, mockHandler)
	assert.Error(t, err, "第6个请求应该被限流")
	assert.Equal(t, codes.ResourceExhausted, status.Code(err))

	// 等待令牌桶重新填充
	time.Sleep(1 * time.Second)

	// 应该可以再次请求
	err = interceptorFunc(nil, stream, &grpc.StreamServerInfo{FullMethod: "/test.Service/TestStream"}, mockHandler)
	assert.NoError(t, err, "等待后的请求应该成功")
}

func TestConcurrentRequests(t *testing.T) {
	// 创建内存存储的限流器，每秒允许10个请求
	store := memrate.NewMemoryStore(createLimiter10)

	// 创建拦截器
	interceptor := grpcx_ratelimit.NewRateLimitInterceptor(
		grpcx_ratelimit.WithMemoryStore(store),
	)

	// 创建键获取器
	keyGetter := &testKeyGetter{}

	// 获取拦截器函数
	interceptorFunc := interceptor.UnaryServerInterceptor(keyGetter)

	// 创建模拟的unary处理函数
	mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}

	// 并发请求
	var wg sync.WaitGroup
	successCount := int32(0)
	failCount := int32(0)

	// 使用互斥锁保护计数器
	var mu sync.Mutex

	// 启动20个并发请求
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 使用相同的用户ID来共享限流器
			ctx := metadata.NewIncomingContext(context.Background(),
				metadata.Pairs("user-id", "test-user"))

			resp, err := interceptorFunc(ctx, "request", &grpc.UnaryServerInfo{FullMethod: "/test.Service/TestMethod"}, mockHandler)

			mu.Lock()
			defer mu.Unlock()

			if err == nil {
				successCount++
				assert.Equal(t, "success", resp)
			} else {
				failCount++
				assert.Equal(t, codes.ResourceExhausted, status.Code(err))
			}
		}(i)
	}

	wg.Wait()

	// 验证大约有10个请求成功（由于并发性，可能会有少量偏差）
	t.Logf("成功请求数: %d, 失败请求数: %d", successCount, failCount)
	assert.True(t, successCount <= 11, "成功请求不应超过限流器的容量")
	assert.True(t, failCount >= 9, "应该有请求被限流")
}
