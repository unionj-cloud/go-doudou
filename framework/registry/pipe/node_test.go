package pipe

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/toolkit/pipeconn"
)

// 模拟管道连接函数
func mockDialContextFunc(ctx context.Context) (net.Conn, error) {
	// 创建一对连接
	clientConn, _ := net.Pipe()

	// 返回客户端连接
	return clientConn, nil
}

func TestNewGrpcClientConn(t *testing.T) {
	// 创建模拟管道连接函数
	dialCtx := pipeconn.DialContextFunc(mockDialContextFunc)

	// 测试创建gRPC客户端连接
	assert.NotPanics(t, func() {
		conn := NewGrpcClientConn(dialCtx)
		assert.NotNil(t, conn)

		// 关闭连接
		defer conn.Close()

		// 验证连接状态
		state := conn.GetState()
		t.Logf("连接状态: %v", state)

		// 确保连接可以使用
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		// 等待连接就绪
		conn.WaitForStateChange(ctx, state)
	})
}
