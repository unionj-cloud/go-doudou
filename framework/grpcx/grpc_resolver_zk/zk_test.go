package grpc_resolver_zk

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
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

func (m *mockClientConn) ParseServiceConfig(serviceConfigJSON string) *serviceconfig.ParseResult {
	return nil
}

func TestScheme(t *testing.T) {
	resolver := &ZkResolver{}
	assert.Equal(t, "zk", resolver.Scheme())
}

func TestResolveNow(t *testing.T) {
	mockWatcher := &mockWatcher{
		endpoints: []string{
			"http://host1:8080?weight=10&group=test-group&version=v1.0.0",
			"http://host2:8080?weight=20&group=test-group&version=v1.0.0",
		},
		eventCh: make(chan struct{}, 1),
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
	zkr := &ZkResolver{
		ZkConfig: ZkConfigs["test-service"],
	}

	// 创建 Target
	testURL, _ := url.Parse("zk://test-service/")
	target := resolver.Target{URL: *testURL}

	// 测试构建解析器
	r, err := zkr.Build(target, cc, resolver.BuildOptions{})
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
		eventCh: make(chan struct{}, 1),
	}

	// 创建 ZkResolver
	cc := &mockClientConn{}
	config := &ZkConfig{
		ServiceName: "test-service",
		Watcher:     mockWatcher,
		Group:       "test-group",
		Version:     "v1.0.0",
	}

	res := &ZkResolver{
		ZkConfig: config,
	}

	// 调用 updateState 方法
	res.updateState(cc)

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

	// 创建解析器
	resolver := &ZkResolver{
		ZkConfig: config,
	}

	// 调用 convertToAddress 方法
	addrs := resolver.convertToAddress(endpoints)

	// 验证只有符合条件的地址被转换
	assert.Equal(t, 2, len(addrs))

	// 验证地址信息正确
	addrMap := make(map[string]serviceInfo)
	for _, addr := range addrs {
		addrMap[addr.Address] = addr
	}

	assert.Contains(t, addrMap, "host1:8080")
	assert.Equal(t, 10, addrMap["host1:8080"].Weight)
	assert.Contains(t, addrMap, "host2:8080")
	assert.Equal(t, 20, addrMap["host2:8080"].Weight)
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

	// 创建 ZkResolver
	cc := &mockClientConn{}
	config := &ZkConfig{
		ServiceName: "test-service",
		Watcher:     mockWatcher,
		Group:       "test-group",
		Version:     "v1.0.0",
	}

	res := &ZkResolver{
		ZkConfig: config,
	}

	// 启动监视协程
	done := make(chan struct{})
	go func() {
		res.watchZkService(cc)
		close(done)
	}()

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

	// 关闭watcher，确保协程退出
	mockWatcher.Close()

	// 等待协程退出
	select {
	case <-done:
		// 协程已退出
	case <-time.After(500 * time.Millisecond):
		t.Fatal("watchZkService goroutine did not exit in time")
	}
}
