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

func TestNewRest(t *testing.T) {
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

	// 注意：我们无法直接测试真实的服务注册情况，因为那需要外部系统如nacos/etcd等
	// 实际项目中应该考虑使用mock库来模拟这些依赖
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

// 后续NewGrpc和ShutdownGrpc的测试类似，省略
// 需要注意的是完整测试应该考虑使用mocking库来模拟外部依赖
