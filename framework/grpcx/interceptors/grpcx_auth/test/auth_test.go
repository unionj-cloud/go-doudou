package test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/unionj-cloud/go-doudou/v2/framework/grpcx/interceptors/grpcx_auth"
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

// 实现 Authorizer 接口
type mockAuthorizer struct{}

func (m *mockAuthorizer) Authorize(ctx context.Context, fullMethod string) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, errors.New("无法获取metadata")
	}

	if tokens := md.Get("authorization"); len(tokens) > 0 && tokens[0] == "valid-token" {
		return ctx, nil
	}
	return ctx, errors.New("认证失败")
}

func TestUnaryServerInterceptor(t *testing.T) {
	// 创建带有认证token的上下文
	validCtx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "valid-token"))
	invalidCtx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "invalid-token"))
	noTokenCtx := context.Background()

	// 创建模拟的unary处理函数
	mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}

	// 创建认证器
	authorizer := &mockAuthorizer{}

	// 测试有效token
	interceptor := grpcx_auth.UnaryServerInterceptor(authorizer)
	resp, err := interceptor(validCtx, "request", &grpc.UnaryServerInfo{}, mockHandler)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp)

	// 测试无效token
	resp, err = interceptor(invalidCtx, "request", &grpc.UnaryServerInfo{}, mockHandler)
	assert.Error(t, err)
	assert.Equal(t, "认证失败", err.Error())
	assert.Nil(t, resp)

	// 测试没有token
	resp, err = interceptor(noTokenCtx, "request", &grpc.UnaryServerInfo{}, mockHandler)
	assert.Error(t, err)
	assert.Equal(t, "无法获取metadata", err.Error())
	assert.Nil(t, resp)
}

func TestStreamServerInterceptor(t *testing.T) {
	// 创建带有认证token的上下文
	validCtx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "valid-token"))
	invalidCtx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "invalid-token"))
	noTokenCtx := context.Background()

	// 创建模拟的stream处理函数
	mockHandler := func(srv interface{}, stream grpc.ServerStream) error {
		return nil
	}

	// 创建认证器
	authorizer := &mockAuthorizer{}

	// 测试有效token
	validStream := &mockServerStream{ctx: validCtx}
	interceptor := grpcx_auth.StreamServerInterceptor(authorizer)
	err := interceptor(nil, validStream, &grpc.StreamServerInfo{}, mockHandler)
	assert.NoError(t, err)

	// 测试无效token
	invalidStream := &mockServerStream{ctx: invalidCtx}
	err = interceptor(nil, invalidStream, &grpc.StreamServerInfo{}, mockHandler)
	assert.Error(t, err)
	assert.Equal(t, "认证失败", err.Error())

	// 测试没有token
	noTokenStream := &mockServerStream{ctx: noTokenCtx}
	err = interceptor(nil, noTokenStream, &grpc.StreamServerInfo{}, mockHandler)
	assert.Error(t, err)
	assert.Equal(t, "无法获取metadata", err.Error())
}
