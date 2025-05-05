package grpc_resolver_zk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseURL(t *testing.T) {
	// 保存原始配置，以便在测试后恢复
	originalConfigs := ZkConfigs
	defer func() {
		ZkConfigs = originalConfigs
	}()

	// 设置测试配置
	ZkConfigs = make(map[string]*ZkConfig)
	mockWatcher := &mockWatcher{}
	ZkConfigs["service-name"] = &ZkConfig{
		Label:       "service-name",
		ServiceName: "service-name",
		Watcher:     mockWatcher,
		Group:       "test-group",
		Version:     "test-version",
	}

	tests := []struct {
		name    string
		target  string
		wantOK  bool
		wantSvc string
	}{
		{
			name:    "有效的 URL",
			target:  "zk://service-name/",
			wantOK:  true,
			wantSvc: "service-name",
		},
		{
			name:    "无效的 scheme",
			target:  "invalid://service-name/",
			wantOK:  false,
			wantSvc: "",
		},
		{
			name:    "缺少服务名称",
			target:  "zk:///",
			wantOK:  false,
			wantSvc: "",
		},
		{
			name:    "格式错误的 URL",
			target:  "zk:/service-name",
			wantOK:  false,
			wantSvc: "",
		},
		{
			name:    "未配置的服务",
			target:  "zk://unknown-service/",
			wantOK:  false,
			wantSvc: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := parseURL(tt.target)
			if tt.wantOK {
				assert.NoError(t, err, "parseURL() 不应该返回错误")
				assert.Equal(t, tt.wantSvc, config.ServiceName, "解析的服务名称不正确")
			} else {
				assert.Error(t, err, "parseURL() 应该返回错误")
			}
		})
	}
}

func TestAddZkConfig(t *testing.T) {
	// 保存原始配置，以便在测试后恢复
	originalConfigs := ZkConfigs
	defer func() {
		ZkConfigs = originalConfigs
	}()

	// 重置配置映射
	ZkConfigs = make(map[string]*ZkConfig)

	// 添加配置
	mockWatcher := &mockWatcher{}
	config := ZkConfig{
		Label:       "test-label",
		ServiceName: "test-service",
		Watcher:     mockWatcher,
		Group:       "test-group",
		Version:     "test-version",
	}

	AddZkConfig(config)

	// 验证配置已添加
	assert.Len(t, ZkConfigs, 1, "应该有一个配置")
	assert.Contains(t, ZkConfigs, "test-label", "配置应该以标签为键")

	// 验证配置内容
	addedConfig := ZkConfigs["test-label"]
	assert.Equal(t, "test-service", addedConfig.ServiceName)
	assert.Equal(t, mockWatcher, addedConfig.Watcher)
	assert.Equal(t, "test-group", addedConfig.Group)
	assert.Equal(t, "test-version", addedConfig.Version)
}

func TestDelZkConfig(t *testing.T) {
	// 保存原始配置，以便在测试后恢复
	originalConfigs := ZkConfigs
	defer func() {
		ZkConfigs = originalConfigs
	}()

	// 重置配置映射并添加测试数据
	ZkConfigs = make(map[string]*ZkConfig)
	mockWatcher := &mockWatcher{}
	ZkConfigs["test-label"] = &ZkConfig{
		Label:       "test-label",
		ServiceName: "test-service",
		Watcher:     mockWatcher,
		Group:       "test-group",
		Version:     "test-version",
	}

	// 删除配置
	DelZkConfig("test-label")

	// 验证配置已删除
	assert.Len(t, ZkConfigs, 0, "配置应该已被删除")
	assert.NotContains(t, ZkConfigs, "test-label", "配置键不应该存在")

	// 测试删除不存在的键
	DelZkConfig("non-existent-label")
	// 应该不会发生错误
}

// 实现模拟的 Watcher 接口，用于测试
type mockWatcher struct {
	endpoints []string
	eventCh   chan struct{}
	closed    bool
}

func (m *mockWatcher) Endpoints() []string {
	return m.endpoints
}

func (m *mockWatcher) Event() <-chan struct{} {
	if m.eventCh == nil {
		m.eventCh = make(chan struct{})
	}
	return m.eventCh
}

func (m *mockWatcher) IsClosed() bool {
	return m.closed
}

func (m *mockWatcher) Close() {
	m.closed = true
	if m.eventCh != nil {
		close(m.eventCh)
	}
}
