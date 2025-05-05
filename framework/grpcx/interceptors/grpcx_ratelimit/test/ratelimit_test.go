package test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/v2/framework/grpcx/interceptors/grpcx_ratelimit"
	"github.com/unionj-cloud/go-doudou/v2/framework/ratelimit/memrate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

func TestUnaryServerInterceptor(t *testing.T) {
	// 创建内存存储的限流器，每秒允许5个请求
	limiter := memrate.NewLimiter(5, 1)

	// 创建带有标识符的上下文
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("user-id", "test-user"))

	// 创建模拟的unary处理函数
	mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}

	// 创建拦截器
	opts := []grpcx_ratelimit.Option{
		grpcx_ratelimit.WithMemoryStore(limiter),
	}
	interceptor := grpcx_ratelimit.UnaryServerInterceptor(opts...)

	// 测试限流效果，前5个请求应该成功，后面的应该失败
	for i := 0; i < 5; i++ {
		resp, err := interceptor(ctx, "request", &grpc.UnaryServerInfo{}, mockHandler)
		assert.NoError(t, err, "第%d个请求应该成功", i+1)
		assert.Equal(t, "success", resp)
	}

	// 第6个请求应该被限流
	resp, err := interceptor(ctx, "request", &grpc.UnaryServerInfo{}, mockHandler)
	assert.Error(t, err, "第6个请求应该被限流")
	assert.Nil(t, resp)
	assert.Equal(t, codes.ResourceExhausted, status.Code(err))

	// 等待令牌桶重新填充
	time.Sleep(1 * time.Second)

	// 应该可以再次请求
	resp, err = interceptor(ctx, "request", &grpc.UnaryServerInfo{}, mockHandler)
	assert.NoError(t, err, "等待后的请求应该成功")
	assert.Equal(t, "success", resp)
}

func TestStreamServerInterceptor(t *testing.T) {
	// 创建内存存储的限流器，每秒允许5个请求
	limiter := memrate.NewLimiter(5, 1)

	// 创建带有标识符的上下文
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("user-id", "test-user"))
	stream := &mockServerStream{ctx: ctx}

	// 创建模拟的stream处理函数
	mockHandler := func(srv interface{}, stream grpc.ServerStream) error {
		return nil
	}

	// 创建拦截器
	opts := []grpcx_ratelimit.Option{
		grpcx_ratelimit.WithMemoryStore(limiter),
	}
	interceptor := grpcx_ratelimit.StreamServerInterceptor(opts...)

	// 测试限流效果，前5个请求应该成功，后面的应该失败
	for i := 0; i < 5; i++ {
		err := interceptor(nil, stream, &grpc.StreamServerInfo{}, mockHandler)
		assert.NoError(t, err, "第%d个请求应该成功", i+1)
	}

	// 第6个请求应该被限流
	err := interceptor(nil, stream, &grpc.StreamServerInfo{}, mockHandler)
	assert.Error(t, err, "第6个请求应该被限流")
	assert.Equal(t, codes.ResourceExhausted, status.Code(err))

	// 等待令牌桶重新填充
	time.Sleep(1 * time.Second)

	// 应该可以再次请求
	err = interceptor(nil, stream, &grpc.StreamServerInfo{}, mockHandler)
	assert.NoError(t, err, "等待后的请求应该成功")
}

func TestConcurrentRequests(t *testing.T) {
	// 创建内存存储的限流器，每秒允许10个请求
	limiter := memrate.NewLimiter(10, 1)

	// 创建拦截器选项
	opts := []grpcx_ratelimit.Option{
		grpcx_ratelimit.WithMemoryStore(limiter),
	}

	// 创建模拟的unary处理函数
	mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}

	// 创建拦截器
	interceptor := grpcx_ratelimit.UnaryServerInterceptor(opts...)

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

			// 使用不同的用户ID避免共享相同的限流器
			ctx := metadata.NewIncomingContext(context.Background(),
				metadata.Pairs("user-id", "test-user-"+string(id)))

			resp, err := interceptor(ctx, "request", &grpc.UnaryServerInfo{}, mockHandler)

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
	assert.True(t, successCount <= 10, "成功请求不应超过限流器的容量")
	assert.True(t, failCount >= 10, "应该有请求被限流")
}
