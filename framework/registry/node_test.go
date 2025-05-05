package registry

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/unionj-cloud/go-doudou/v2/framework/config"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry/constants"
)

// 由于内部实现和外部系统有交互，我们将在测试中模拟环境变量来控制不同的服务发现模式
func setupTest(t *testing.T) func() {
	// 保存原始环境变量
	origServiceDiscoveryMode := os.Getenv(string(config.GddServiceDiscoveryMode))

	// 清理函数
	return func() {
		// 恢复环境变量
		os.Setenv(string(config.GddServiceDiscoveryMode), origServiceDiscoveryMode)
	}
}

// 在不实际连接到外部系统的情况下测试注册功能
func TestNewRestWithMockEnv(t *testing.T) {
	// 设置环境
	cleanup := setupTest(t)
	defer cleanup()

	// 测试传递空服务发现模式
	os.Setenv(string(config.GddServiceDiscoveryMode), "")
	// 预期不会有错误发生
	NewRest()

	// 测试传递无效的服务发现模式
	os.Setenv(string(config.GddServiceDiscoveryMode), "invalid-mode")
	// 预期会输出警告日志但不会崩溃
	NewRest()

	// 测试传递数据参数
	testData := map[string]interface{}{
		"name": "test-service",
		"port": 8080,
	}
	NewRest(testData)
}

// 在不实际连接到外部系统的情况下测试gRPC注册功能
func TestNewGrpcWithMockEnv(t *testing.T) {
	// 设置环境
	cleanup := setupTest(t)
	defer cleanup()

	// 测试传递空服务发现模式
	os.Setenv(string(config.GddServiceDiscoveryMode), "")
	// 预期不会有错误发生
	NewGrpc()

	// 测试传递无效的服务发现模式
	os.Setenv(string(config.GddServiceDiscoveryMode), "invalid-mode")
	// 预期会输出警告日志但不会崩溃
	NewGrpc()

	// 测试传递数据参数
	testData := map[string]interface{}{
		"name": "test-service",
		"port": 9090,
	}
	NewGrpc(testData)
}

func TestShutdownRest(t *testing.T) {
	// 设置环境
	cleanup := setupTest(t)
	defer cleanup()

	// 测试空服务发现模式
	os.Setenv(string(config.GddServiceDiscoveryMode), "")
	ShutdownRest()

	// 测试无效服务发现模式
	os.Setenv(string(config.GddServiceDiscoveryMode), "invalid-mode")
	ShutdownRest()

	// 测试多个服务发现模式
	os.Setenv(string(config.GddServiceDiscoveryMode), "invalid-mode,another-invalid")
	ShutdownRest()
}

func TestShutdownGrpc(t *testing.T) {
	// 设置环境
	cleanup := setupTest(t)
	defer cleanup()

	// 测试空服务发现模式
	os.Setenv(string(config.GddServiceDiscoveryMode), "")
	ShutdownGrpc()

	// 测试无效服务发现模式
	os.Setenv(string(config.GddServiceDiscoveryMode), "invalid-mode")
	ShutdownGrpc()

	// 测试多个服务发现模式
	os.Setenv(string(config.GddServiceDiscoveryMode), "invalid-mode,another-invalid")
	ShutdownGrpc()
}

func TestServiceDiscoveryMap(t *testing.T) {
	// 保存原始环境变量
	origServiceDiscoveryMode := os.Getenv(string(config.GddServiceDiscoveryMode))
	defer os.Setenv(string(config.GddServiceDiscoveryMode), origServiceDiscoveryMode)

	// 测试空服务发现模式
	os.Setenv(string(config.GddServiceDiscoveryMode), "")
	sdMap := config.ServiceDiscoveryMap()
	assert.Empty(t, sdMap)

	// 测试单一模式
	os.Setenv(string(config.GddServiceDiscoveryMode), constants.SD_NACOS)
	sdMap = config.ServiceDiscoveryMap()
	assert.Contains(t, sdMap, constants.SD_NACOS)
	assert.Len(t, sdMap, 1)

	// 测试多个模式
	os.Setenv(string(config.GddServiceDiscoveryMode), constants.SD_NACOS+","+constants.SD_ETCD)
	sdMap = config.ServiceDiscoveryMap()
	assert.Contains(t, sdMap, constants.SD_NACOS)
	assert.Contains(t, sdMap, constants.SD_ETCD)
	assert.Len(t, sdMap, 2)
}

// 测试IServiceProvider接口
type mockServiceProvider struct{}

func (m *mockServiceProvider) SelectServer() string {
	return "http://localhost:8080"
}

func (m *mockServiceProvider) Close() {
	// 空实现
}

func TestIServiceProvider(t *testing.T) {
	// 测试接口实现
	var provider IServiceProvider = &mockServiceProvider{}

	server := provider.SelectServer()
	assert.Equal(t, "http://localhost:8080", server)

	// 测试关闭方法，不应该有错误
	provider.Close()
}
