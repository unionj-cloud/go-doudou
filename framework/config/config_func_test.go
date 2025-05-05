package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/unionj-cloud/go-doudou/v2/framework/config"
)

func TestCheckDev(t *testing.T) {
	// 保存原始环境变量
	origEnv := os.Getenv("GDD_ENV")
	defer os.Setenv("GDD_ENV", origEnv)

	// 测试未设置环境变量的情况
	os.Unsetenv("GDD_ENV")
	assert.True(t, config.CheckDev())

	// 测试设置为dev的情况
	os.Setenv("GDD_ENV", "dev")
	assert.True(t, config.CheckDev())

	// 测试设置为其他值的情况
	os.Setenv("GDD_ENV", "prod")
	assert.False(t, config.CheckDev())
}

func TestServiceDiscoveryMap(t *testing.T) {
	// 保存原始值
	origMode := config.GddServiceDiscoveryMode.Load()
	defer config.GddServiceDiscoveryMode.Write(origMode)

	// 测试空模式
	config.GddServiceDiscoveryMode.Write("")
	sdMap := config.ServiceDiscoveryMap()
	assert.Nil(t, sdMap)

	// 测试单一模式
	config.GddServiceDiscoveryMode.Write("nacos")
	sdMap = config.ServiceDiscoveryMap()
	assert.NotNil(t, sdMap)
	assert.Len(t, sdMap, 1)
	_, exists := sdMap["nacos"]
	assert.True(t, exists)

	// 测试多个模式
	config.GddServiceDiscoveryMode.Write("nacos,etcd,zk")
	sdMap = config.ServiceDiscoveryMap()
	assert.NotNil(t, sdMap)
	assert.Len(t, sdMap, 3)
	_, exists = sdMap["nacos"]
	assert.True(t, exists)
	_, exists = sdMap["etcd"]
	assert.True(t, exists)
	_, exists = sdMap["zk"]
	assert.True(t, exists)
}

func TestShutdown(t *testing.T) {
	// 测试文件为空时的Shutdown
	config.G_LogFile = nil
	config.Shutdown()
	assert.Nil(t, config.G_LogFile)

	// 创建临时日志文件进行测试
	tempFile, err := os.CreateTemp("", "test-log-*.log")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name()) // 清理

	// 设置日志文件并测试Shutdown
	config.G_LogFile = tempFile
	assert.NotNil(t, config.G_LogFile)
	config.Shutdown()
	assert.Nil(t, config.G_LogFile)

	// 检查文件是否已关闭（写入已关闭的文件会出错）
	_, err = tempFile.Write([]byte("test"))
	assert.Error(t, err)
}

func TestLoadConfigFromLocal(t *testing.T) {
	// 保存原始环境变量
	origEnv := os.Getenv("GDD_ENV")
	defer os.Setenv("GDD_ENV", origEnv)

	// 测试未设置环境变量的情况
	os.Unsetenv("GDD_ENV")
	assert.NotPanics(t, func() {
		config.LoadConfigFromLocal()
	})

	// 测试设置环境变量的情况
	os.Setenv("GDD_ENV", "test")
	assert.NotPanics(t, func() {
		config.LoadConfigFromLocal()
	})
}
