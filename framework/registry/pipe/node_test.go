package pipe

import (
	"context"
	"errors"
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

// 模拟失败的管道连接函数
func mockFailDialContextFunc(ctx context.Context) (net.Conn, error) {
	return nil, errors.New("mock dial error")
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

func TestNewGrpcClientConnFail(t *testing.T) {
	// 创建失败的模拟管道连接函数
	dialCtx := pipeconn.DialContextFunc(mockFailDialContextFunc)

	// 测试创建失败的gRPC客户端连接
	assert.NotPanics(t, func() {
		conn := NewGrpcClientConn(dialCtx)
		assert.NotNil(t, conn)

		// 关闭连接
		defer conn.Close()

		// 验证连接初始状态可能是CONNECTING或TRANSIENT_FAILURE
		state := conn.GetState()
		assert.Equal(t, "CONNECTING", state.String())
	})
}

func TestServerClientConnection(t *testing.T) {
	// 创建一对管道连接
	clientConn, serverConn := net.Pipe()

	dialCtx := func(ctx context.Context) (net.Conn, error) {
		return clientConn, nil
	}

	// 测试创建gRPC客户端连接
	assert.NotPanics(t, func() {
		conn := NewGrpcClientConn(pipeconn.DialContextFunc(dialCtx))
		assert.NotNil(t, conn)
		defer conn.Close()

		// 模拟服务器端写入数据
		go func() {
			_, err := serverConn.Write([]byte("hello from server"))
			assert.NoError(t, err)
		}()

		// 等待一段时间，确保数据写入
		time.Sleep(100 * time.Millisecond)

		// 关闭服务器连接
		serverConn.Close()
	})
}
