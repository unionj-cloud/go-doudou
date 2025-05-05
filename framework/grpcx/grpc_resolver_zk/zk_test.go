package grpc_resolver_zk

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
)

// 模拟 ClientConn 接口，用于测试
type mockClientConn struct {
	state      resolver.State
	updateErr  error
	reportErr  error
	parsedAddr []resolver.Address
}

func (m *mockClientConn) UpdateState(state resolver.State) error {
	m.state = state
	m.parsedAddr = state.Addresses
	return m.updateErr
}

func (m *mockClientConn) ReportError(err error) {
	m.reportErr = err
}

func (m *mockClientConn) NewAddress(addresses []resolver.Address) {
	m.parsedAddr = addresses
}

func (m *mockClientConn) NewServiceConfig(serviceConfig string) {
	// 不需要实现
}

func (m *mockClientConn) ParseServiceConfig(serviceConfigJSON string) *resolver.ServiceConfig {
	return nil
}

func TestScheme(t *testing.T) {
	resolver := &zkResolver{}
	assert.Equal(t, "zk", resolver.Scheme())
}

func TestResolveNow(t *testing.T) {
	mockWatcher := &mockWatcher{
		endpoints: []string{
			"http://host1:8080?weight=10&group=test-group&version=v1.0.0",
			"http://host2:8080?weight=20&group=test-group&version=v1.0.0",
		},
	}

	// 保存原始配置，以便在测试后恢复
	originalConfigs := ZkConfigs
	defer func() {
		ZkConfigs = originalConfigs
	}()

	// 创建测试配置
	ZkConfigs = make(map[string]*ZkConfig)
	ZkConfigs["test-service"] = &ZkConfig{
		Label:       "test-service",
		ServiceName: "test-service",
		Watcher:     mockWatcher,
		Group:       "test-group",
		Version:     "v1.0.0",
	}

	// 创建模拟连接
	cc := &mockClientConn{}

	// 创建解析器
	zkr := &zkResolver{}
	r, err := zkr.Build(resolver.Target{URL: *mustParseURL("zk://test-service/")}, cc, resolver.BuildOptions{})
	assert.NoError(t, err)
	assert.NotNil(t, r)

	// 延迟检查，等待 updateState 完成
	time.Sleep(100 * time.Millisecond)

	// 测试 ResolveNow
	r.ResolveNow(resolver.ResolveNowOptions{})

	// 验证地址已正确更新
	assert.NotEmpty(t, cc.parsedAddr)

	// 清理
	r.Close()
}

func TestUpdateState(t *testing.T) {
	// 创建测试数据
	mockWatcher := &mockWatcher{
		endpoints: []string{
			"http://host1:8080?weight=10&group=test-group&version=v1.0.0",
			"http://host2:8080?weight=20&group=test-group&version=v1.0.0",
			"http://host3:8080?weight=5&group=other-group&version=v1.0.0", // 不同的组
			"http://host4:8080?weight=5&group=test-group&version=v2.0.0",  // 不同的版本
		},
	}

	// 创建 zkBuilder 和 zkResolver
	cc := &mockClientConn{}
	builder := &zkBuilder{}
	config := &ZkConfig{
		ServiceName: "test-service",
		Watcher:     mockWatcher,
		Group:       "test-group",
		Version:     "v1.0.0",
	}

	res := &zkResolver{
		target:  "test-service",
		cc:      cc,
		watcher: mockWatcher,
		config:  config,
	}

	// 调用 updateState 方法
	res.updateState()

	// 验证只有符合条件的地址被添加
	assert.Equal(t, 2, len(cc.parsedAddr))

	// 验证地址格式正确
	addrMap := make(map[string]bool)
	for _, addr := range cc.parsedAddr {
		addrMap[addr.Addr] = true
	}
	assert.Contains(t, addrMap, "host1:8080")
	assert.Contains(t, addrMap, "host2:8080")
	assert.NotContains(t, addrMap, "host3:8080") // 不同组
	assert.NotContains(t, addrMap, "host4:8080") // 不同版本
}

func TestConvertToAddress(t *testing.T) {
	// 创建测试数据
	endpoints := []string{
		"http://host1:8080?weight=10&group=test-group&version=v1.0.0",
		"http://host2:8080?weight=20&group=test-group&version=v1.0.0",
		"http://host3:8080?weight=5&group=other-group&version=v1.0.0", // 不同的组
		"http://host4:8080?weight=5&group=test-group&version=v2.0.0",  // 不同的版本
	}

	config := &ZkConfig{
		ServiceName: "test-service",
		Group:       "test-group",
		Version:     "v1.0.0",
	}

	// 调用 convertToAddress 方法
	addrs := convertToAddress(endpoints, config)

	// 验证只有符合条件的地址被转换
	assert.Equal(t, 2, len(addrs))

	// 验证地址信息正确
	addrMap := make(map[string]resolver.Address)
	for _, addr := range addrs {
		addrMap[addr.Addr] = addr
	}

	assert.Contains(t, addrMap, "host1:8080")
	assert.Contains(t, addrMap, "host2:8080")
	assert.NotContains(t, addrMap, "host3:8080") // 不同组
	assert.NotContains(t, addrMap, "host4:8080") // 不同版本
}

func TestWatchZkService(t *testing.T) {
	// 创建测试数据
	mockWatcher := &mockWatcher{
		endpoints: []string{
			"http://host1:8080?weight=10&group=test-group&version=v1.0.0",
		},
		eventCh: make(chan struct{}, 1),
	}

	// 创建 zkResolver
	cc := &mockClientConn{}
	config := &ZkConfig{
		ServiceName: "test-service",
		Watcher:     mockWatcher,
		Group:       "test-group",
		Version:     "v1.0.0",
	}

	res := &zkResolver{
		target:  "test-service",
		cc:      cc,
		watcher: mockWatcher,
		config:  config,
	}

	// 启动监视协程
	ctx, cancel := context.WithCancel(context.Background())
	go res.watchZkService(ctx)

	// 初始化后应该有地址
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 1, len(cc.parsedAddr))

	// 模拟端点更新
	mockWatcher.endpoints = []string{
		"http://host1:8080?weight=10&group=test-group&version=v1.0.0",
		"http://host2:8080?weight=20&group=test-group&version=v1.0.0",
	}
	mockWatcher.eventCh <- struct{}{}

	// 等待更新
	time.Sleep(100 * time.Millisecond)

	// 验证地址已更新
	assert.Equal(t, 2, len(cc.parsedAddr))

	// 取消上下文，停止监视
	cancel()

	// 模拟再次更新，应该不会再处理
	mockWatcher.endpoints = []string{
		"http://host3:8080?weight=30&group=test-group&version=v1.0.0",
	}

	// 由于上下文已取消，不应该有新的更新
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 2, len(cc.parsedAddr))
}

// 辅助函数
func mustParseURL(u string) *resolver.URL {
	url, err := resolver.Parse(u)
	if err != nil {
		panic(err)
	}
	return url
}
