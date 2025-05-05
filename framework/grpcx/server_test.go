package grpcx

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestNewGrpcServer(t *testing.T) {
	// 测试基本创建
	server := NewGrpcServer()
	assert.NotNil(t, server)
	assert.NotNil(t, server.Server)

	// 测试带选项创建
	server = NewGrpcServer(grpc.MaxRecvMsgSize(1024 * 1024))
	assert.NotNil(t, server)
	assert.NotNil(t, server.Server)
}

func TestNewEmptyGrpcServer(t *testing.T) {
	server := NewEmptyGrpcServer()
	assert.NotNil(t, server)
	assert.Nil(t, server.Server)
}

func TestNewGrpcServerWithData(t *testing.T) {
	// 准备测试数据
	data := map[string]interface{}{
		"key": "value",
	}

	// 测试创建
	server := NewGrpcServerWithData(data)
	assert.NotNil(t, server)
	assert.NotNil(t, server.Server)
	assert.Equal(t, data, server.data)

	// 测试带选项创建
	server = NewGrpcServerWithData(data, grpc.MaxRecvMsgSize(1024*1024))
	assert.NotNil(t, server)
	assert.NotNil(t, server.Server)
	assert.Equal(t, data, server.data)
}

func TestPrintServices(t *testing.T) {
	// 创建服务器并注册健康检查服务
	server := NewGrpcServer()
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(server.Server, healthServer)

	// 测试打印服务方法，这里只是确保不会panic
	assert.NotPanics(t, func() {
		server.printServices()
	})
}

func TestServe(t *testing.T) {
	// 创建服务器
	server := NewGrpcServer()

	// 创建监听器
	ln, err := net.Listen("tcp", ":0") // 随机端口
	assert.NoError(t, err)
	defer ln.Close()

	// 在goroutine中运行服务，因为Serve是阻塞的
	go func() {
		server.Serve(ln)
	}()

	// 给一点时间让服务器启动
	time.Sleep(100 * time.Millisecond)

	// 测试空服务器
	emptyServer := NewEmptyGrpcServer()
	assert.NotPanics(t, func() {
		emptyServer.Serve(ln)
	})
}

func TestServeWithPipe(t *testing.T) {
	// 创建服务器
	server := NewGrpcServer()

	// 创建主监听器
	ln, err := net.Listen("tcp", ":0") // 随机端口
	assert.NoError(t, err)
	defer ln.Close()

	// 创建管道监听器
	pipe, err := net.Listen("tcp", ":0") // 随机端口
	assert.NoError(t, err)
	defer pipe.Close()

	// 在goroutine中运行服务，因为ServeWithPipe是阻塞的
	go func() {
		server.ServeWithPipe(ln, pipe)
	}()

	// 给一点时间让服务器启动
	time.Sleep(100 * time.Millisecond)

	// 测试空服务器
	emptyServer := NewEmptyGrpcServer()
	assert.NotPanics(t, func() {
		emptyServer.ServeWithPipe(ln, pipe)
	})
}
