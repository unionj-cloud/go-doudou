package zk

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/unionj-cloud/go-doudou/v2/framework/config"
)

// 模拟 Watcher 实现
type mockWatcher struct {
	endpoints []string
	eventCh   chan struct{}
	closed    bool
	mu        sync.Mutex
}

func newMockWatcher(endpoints []string) *mockWatcher {
	return &mockWatcher{
		endpoints: endpoints,
		eventCh:   make(chan struct{}, 1),
		closed:    false,
	}
}

func (m *mockWatcher) Endpoints() []string {
	return m.endpoints
}

func (m *mockWatcher) Event() <-chan struct{} {
	return m.eventCh
}

func (m *mockWatcher) IsClosed() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.closed
}

func (m *mockWatcher) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.closed {
		m.closed = true
		close(m.eventCh)
	}
}

func (m *mockWatcher) UpdateEndpoints(endpoints []string) {
	m.endpoints = endpoints
	m.eventCh <- struct{}{}
}

// 测试设置
func setupTest() func() {
	// 保存原始环境变量
	origZkServers := config.GddZkServers.Load()
	origServiceName := config.GddServiceName.Load()
	origServiceGroup := config.GddServiceGroup.Load()
	origServiceVersion := config.GddServiceVersion.Load()
	origPort := config.GddPort.Load()
	origGrpcPort := config.GddGrpcPort.Load()
	origRouteRootPath := config.GddRouteRootPath.Load()
	origWeight := config.GddWeight.Load()

	// 设置测试环境变量
	config.GddZkServers.Write("localhost:2181")
	config.GddServiceName.Write("test-service")
	config.GddServiceGroup.Write("test-group")
	config.GddServiceVersion.Write("v1.0.0")
	config.GddPort.Write("8080")
	config.GddGrpcPort.Write("50051")
	config.GddRouteRootPath.Write("/api")
	config.GddWeight.Write("10")

	// 返回清理函数
	return func() {
		config.GddZkServers.Write(origZkServers)
		config.GddServiceName.Write(origServiceName)
		config.GddServiceGroup.Write(origServiceGroup)
		config.GddServiceVersion.Write(origServiceVersion)
		config.GddPort.Write(origPort)
		config.GddGrpcPort.Write(origGrpcPort)
		config.GddRouteRootPath.Write(origRouteRootPath)
		config.GddWeight.Write(origWeight)
	}
}

// 跳过需要外部zk服务的测试
func skipIfNoZk(t *testing.T) {
	if os.Getenv("ZK_TEST") != "true" {
		t.Skip("跳过需要Zookeeper的测试。设置ZK_TEST=true环境变量以启用这些测试")
	}
}

func TestPopulateMeta(t *testing.T) {
	cleanup := setupTest()
	defer cleanup()

	// 创建元数据映射
	meta := make(map[string]interface{})

	// 测试填充元数据
	populateMeta(meta)

	// 验证元数据字段
	assert.Contains(t, meta, "registerAt")
	assert.Contains(t, meta, "goVer")
	assert.Contains(t, meta, "weight")
	assert.Contains(t, meta, "group")
	assert.Contains(t, meta, "version")
	assert.Contains(t, meta, "rootPath")
	assert.Equal(t, "test-group", meta["group"])
	assert.Equal(t, "v1.0.0", meta["version"])
	assert.Equal(t, "/api", meta["rootPath"])
	assert.Equal(t, 10, meta["weight"])

	// 测试带用户数据的填充
	meta = make(map[string]interface{})
	userData := map[string]interface{}{
		"customKey": "customValue",
	}
	populateMeta(meta, userData)

	// 验证自定义字段
	assert.Equal(t, "customValue", meta["customKey"])
}

func TestRRServiceProviderConvertToAddress(t *testing.T) {
	// 创建一个ServiceProvider
	target := ServiceConfig{
		Name:    "test-service",
		Group:   "test-group",
		Version: "v1.0.0",
	}
	provider := &RRServiceProvider{
		target: target,
	}

	// 测试地址转换
	endpoints := []string{
		"http://host1:8080?weight=10&rootPath=/api&group=test-group&version=v1.0.0",
		"http://host2:8080?weight=20&rootPath=/api&group=test-group&version=v1.0.0",
		"http://host3:8080?weight=5&rootPath=/api&group=other-group&version=v1.0.0", // 不同的组
		"http://host4:8080?weight=5&rootPath=/api&group=test-group&version=v2.0.0",  // 不同的版本
	}

	addrs := provider.convertToAddress(endpoints)

	// 应该只有匹配组和版本的地址被包含
	assert.Equal(t, 2, len(addrs))
	assert.Equal(t, "host1:8080", addrs[0].addr)
	assert.Equal(t, 10, addrs[0].weight)
	assert.Equal(t, "/api", addrs[0].rootPath)
	assert.Equal(t, "host2:8080", addrs[1].addr)
	assert.Equal(t, 20, addrs[1].weight)
}

func TestRRServiceProviderSelectServer(t *testing.T) {
	// 创建服务配置
	target := ServiceConfig{
		Name:    "test-service",
		Group:   "test-group",
		Version: "v1.0.0",
	}

	// 创建模拟观察者
	mockWatch := newMockWatcher([]string{
		"http://host1:8080?weight=1&rootPath=/api&group=test-group&version=v1.0.0",
		"http://host2:8080?weight=1&rootPath=/api&group=test-group&version=v1.0.0",
	})

	// 创建服务提供者
	provider := &RRServiceProvider{
		target:  target,
		watcher: mockWatch,
	}

	// 初始化状态
	provider.updateState()

	// 测试轮询选择逻辑 - 确保不会panic
	assert.NotPanics(t, func() {
		server1 := provider.SelectServer()
		assert.NotEmpty(t, server1, "第一次选择的服务器不应为空")

		server2 := provider.SelectServer()
		assert.NotEmpty(t, server2, "第二次选择的服务器不应为空")

		t.Logf("选择的服务器: %s, %s", server1, server2)
	})

	// 关闭提供者
	provider.Close()
	assert.True(t, mockWatch.IsClosed())
}

func TestSWRRServiceProviderSelectServer(t *testing.T) {
	// 注意：由于加权轮询算法依赖于内部状态，如果没有足够的调用，可能不总是表现出权重效果
	// 我们在这里只测试函数是否正常工作，而不测试具体的加权效果

	// 创建服务配置
	target := ServiceConfig{
		Name:    "test-service",
		Group:   "test-group",
		Version: "v1.0.0",
	}

	// 创建模拟观察者
	mockWatch := newMockWatcher([]string{
		"http://host1:8080?weight=10&rootPath=/api&group=test-group&version=v1.0.0",
		"http://host2:8080?weight=5&rootPath=/api&group=test-group&version=v1.0.0",
	})

	// 创建RR服务提供者
	rrProvider := &RRServiceProvider{
		target:  target,
		watcher: mockWatch,
	}

	// 创建SWRR服务提供者
	provider := &SWRRServiceProvider{
		RRServiceProvider: rrProvider,
	}

	// 初始化状态
	rrProvider.updateState()

	// 测试函数不应该panic
	assert.NotPanics(t, func() {
		for i := 0; i < 5; i++ {
			server := provider.SelectServer()
			assert.NotEmpty(t, server)
			t.Logf("选择的服务器: %s", server)
		}
	})

	// 关闭提供者
	provider.Close()
}
