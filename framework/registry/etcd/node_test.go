package etcd

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/unionj-cloud/go-doudou/v2/framework/config"
	cons "github.com/unionj-cloud/go-doudou/v2/framework/registry/constants"
)

// 测试前设置模拟环境变量，测试后恢复它们
func setupTestEnv() func() {
	// 保存原始环境变量
	origServiceName := config.GddServiceName.Load()
	origPort := config.GddPort.Load()
	origGrpcPort := config.GddGrpcPort.Load()
	origEtcdEndpoints := config.GddEtcdEndpoints.Load()
	origEtcdLease := config.GddEtcdLease.Load()
	origWeight := config.GddWeight.Load()
	origRouteRootPath := config.GddRouteRootPath.Load()

	// 设置测试环境变量
	config.GddServiceName.Write("etcd-test-service")
	config.GddPort.Write("8080")
	config.GddGrpcPort.Write("50051")
	config.GddEtcdEndpoints.Write("localhost:2379") // 请确保你有一个运行的etcd实例
	config.GddEtcdLease.Write("10")
	config.GddWeight.Write("10")
	config.GddRouteRootPath.Write("/api")

	// 返回清理函数
	return func() {
		config.GddServiceName.Write(origServiceName)
		config.GddPort.Write(origPort)
		config.GddGrpcPort.Write(origGrpcPort)
		config.GddEtcdEndpoints.Write(origEtcdEndpoints)
		config.GddEtcdLease.Write(origEtcdLease)
		config.GddWeight.Write(origWeight)
		config.GddRouteRootPath.Write(origRouteRootPath)

		// 确保关闭etcd客户端
		CloseEtcdClient()
	}
}

// 跳过测试的辅助函数
func skipIfNoEtcd(t *testing.T) {
	if os.Getenv("ETCD_TEST") != "true" {
		t.Skip("跳过需要etcd的测试。设置ETCD_TEST=true环境变量以启用这些测试")
	}
}

func TestPopulateMeta(t *testing.T) {
	// 创建元数据映射
	meta := make(map[string]interface{})

	// 测试填充元数据（REST服务）
	populateMeta(meta, false)

	// 验证元数据字段
	assert.Contains(t, meta, "registerAt")
	assert.Contains(t, meta, "goVer")
	assert.Contains(t, meta, "weight")

	// 测试带用户数据的填充
	meta = make(map[string]interface{})
	userData := map[string]interface{}{
		"customKey": "customValue",
	}
	populateMeta(meta, false, userData)

	// 验证自定义字段
	assert.Equal(t, "customValue", meta["customKey"])

	// 测试gRPC服务标志
	meta = make(map[string]interface{})
	populateMeta(meta, true)

	// 验证gRPC服务没有rootPath
	_, hasRootPath := meta["rootPath"]
	assert.False(t, hasRootPath)
}

func TestInitEtcdCli(t *testing.T) {
	skipIfNoEtcd(t)
	cleanup := setupTestEnv()
	defer cleanup()

	// 测试初始化etcd客户端
	assert.NotPanics(t, func() {
		InitEtcdCli()
	})

	// 验证客户端已初始化
	assert.NotNil(t, EtcdCli)
}

func TestGetLeaseID(t *testing.T) {
	skipIfNoEtcd(t)
	cleanup := setupTestEnv()
	defer cleanup()

	// 确保客户端已初始化
	InitEtcdCli()

	// 测试获取租约ID
	leaseID := getLeaseID()
	assert.NotZero(t, leaseID)
}

func TestRegisterService(t *testing.T) {
	skipIfNoEtcd(t)
	cleanup := setupTestEnv()
	defer cleanup()

	// 确保客户端已初始化
	InitEtcdCli()

	// 获取租约ID
	leaseID := getLeaseID()

	// 测试注册服务
	assert.NotPanics(t, func() {
		registerService("test-service", 8080, leaseID)
	})
}

func TestNewRest(t *testing.T) {
	skipIfNoEtcd(t)
	cleanup := setupTestEnv()
	defer cleanup()

	// 测试注册REST服务
	assert.NotPanics(t, func() {
		NewRest()
	})

	// 测试带用户数据的注册
	userData := map[string]interface{}{
		"customKey": "customValue",
	}
	assert.NotPanics(t, func() {
		NewRest(userData)
	})
}

func TestNewGrpc(t *testing.T) {
	skipIfNoEtcd(t)
	cleanup := setupTestEnv()
	defer cleanup()

	// 测试注册gRPC服务
	assert.NotPanics(t, func() {
		NewGrpc()
	})

	// 测试带用户数据的注册
	userData := map[string]interface{}{
		"customKey": "customValue",
	}
	assert.NotPanics(t, func() {
		NewGrpc(userData)
	})
}

func TestShutdownRest(t *testing.T) {
	skipIfNoEtcd(t)
	cleanup := setupTestEnv()
	defer cleanup()

	// 先注册服务
	NewRest()

	// 测试关闭REST服务
	assert.NotPanics(t, func() {
		ShutdownRest()
	})
}

func TestShutdownGrpc(t *testing.T) {
	skipIfNoEtcd(t)
	cleanup := setupTestEnv()
	defer cleanup()

	// 先注册服务
	NewGrpc()

	// 测试关闭gRPC服务
	assert.NotPanics(t, func() {
		ShutdownGrpc()
	})
}

func TestCloseEtcdClient(t *testing.T) {
	skipIfNoEtcd(t)
	cleanup := setupTestEnv()
	defer cleanup()

	// 初始化客户端
	InitEtcdCli()

	// 测试关闭客户端
	assert.NotPanics(t, func() {
		CloseEtcdClient()
	})

	// 验证客户端已关闭
	assert.Nil(t, EtcdCli)
}

func TestNewRRServiceProvider(t *testing.T) {
	skipIfNoEtcd(t)
	cleanup := setupTestEnv()
	defer cleanup()

	// 确保有一个REST服务注册
	NewRest()

	// 创建服务提供者
	serviceName := config.GetServiceName() + "_" + string(cons.REST_TYPE)
	provider := NewRRServiceProvider(serviceName)

	// 验证提供者已创建
	assert.NotNil(t, provider)

	// 等待服务发现
	time.Sleep(100 * time.Millisecond)

	// 测试选择服务器
	server := provider.SelectServer()
	assert.NotEmpty(t, server)

	// 关闭提供者
	provider.Close()
}

func TestNewSWRRServiceProvider(t *testing.T) {
	skipIfNoEtcd(t)
	cleanup := setupTestEnv()
	defer cleanup()

	// 确保有一个REST服务注册
	NewRest()

	// 创建服务提供者
	serviceName := config.GetServiceName() + "_" + string(cons.REST_TYPE)
	provider := NewSWRRServiceProvider(serviceName)

	// 验证提供者已创建
	assert.NotNil(t, provider)

	// 等待服务发现
	time.Sleep(100 * time.Millisecond)

	// 测试选择服务器
	server := provider.SelectServer()
	assert.NotEmpty(t, server)

	// 关闭提供者
	provider.Close()
}

func TestConvertToAddress(t *testing.T) {
	// 这是一个内部函数，我们将跳过深入测试
	// 但可以做一个基本的测试确保它不会panic
	assert.NotPanics(t, func() {
		addrs := convertToAddress(nil)
		assert.Empty(t, addrs)
	})
}

func TestNewGrpcClientConn(t *testing.T) {
	skipIfNoEtcd(t)
	cleanup := setupTestEnv()
	defer cleanup()

	// 确保有一个gRPC服务注册
	NewGrpc()

	// 创建gRPC连接选项
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
	}

	// 测试创建gRPC客户端连接 (可能会超时，这是预期的行为)
	serviceName := config.GetServiceName() + "_" + string(cons.GRPC_TYPE)
	conn := NewGrpcClientConn(serviceName, "rr", append(opts, grpc.WithTimeout(100*time.Millisecond))...)
	assert.NotNil(t, conn)
}
