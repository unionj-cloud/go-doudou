package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/v2/framework/config"
)

func TestGetRegisterHost(t *testing.T) {
	// 保存原来的函数以便测试后恢复
	origGetPrivateIP := GetPrivateIP
	defer func() { GetPrivateIP = origGetPrivateIP }()

	// 保存原始环境变量
	origRegisterHost := os.Getenv(string(config.GddRegisterHost))
	defer os.Setenv(string(config.GddRegisterHost), origRegisterHost)

	// 测试1：已设置环境变量的情况
	os.Setenv(string(config.GddRegisterHost), "192.168.1.100")
	result := GetRegisterHost()
	assert.Equal(t, "192.168.1.100", result)

	// 测试2：未设置环境变量，但有默认值的情况
	os.Setenv(string(config.GddRegisterHost), "")

	// 由于DefaultGddRegisterHost是包级变量，无法直接修改，
	// 所以此处我们模拟GetPrivateIP返回的IP地址
	GetPrivateIP = func() (string, error) {
		return "10.0.0.1", nil
	}

	result = GetRegisterHost()
	assert.Equal(t, "10.0.0.1", result)
}
