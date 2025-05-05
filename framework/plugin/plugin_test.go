package plugin

import (
	"testing"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/v2/framework/grpcx"
	"github.com/unionj-cloud/go-doudou/v2/framework/rest"
	"github.com/unionj-cloud/toolkit/pipeconn"
)

// 创建一个用于测试的mock ServicePlugin实现
type mockServicePlugin struct {
	name string
}

func (m *mockServicePlugin) Initialize(restServer *rest.RestServer, grpcServer *grpcx.GrpcServer, dialCtx pipeconn.DialContextFunc) {
	// 空实现
}

func (m *mockServicePlugin) GetName() string {
	return m.name
}

func (m *mockServicePlugin) Close() {
	// 空实现
}

func (m *mockServicePlugin) GoDoudouServicePlugin() {
	// 空实现
}

func TestRegisterServicePlugin(t *testing.T) {
	// 清空已有插件，以便测试
	servicePlugins = orderedmap.NewOrderedMap[string, ServicePlugin]()

	// 创建测试插件
	plugin1 := &mockServicePlugin{name: "plugin1"}
	plugin2 := &mockServicePlugin{name: "plugin2"}

	// 注册插件
	RegisterServicePlugin(plugin1)
	RegisterServicePlugin(plugin2)

	// 验证插件是否正确注册
	plugins := GetServicePlugins()
	assert.Equal(t, 2, plugins.Len())

	// 验证插件名称
	keys := plugins.Keys()
	assert.Contains(t, keys, "plugin1")
	assert.Contains(t, keys, "plugin2")

	// 验证可以通过名称获取插件
	p1, ok := plugins.Get("plugin1")
	assert.True(t, ok)
	assert.Equal(t, "plugin1", p1.GetName())

	p2, ok := plugins.Get("plugin2")
	assert.True(t, ok)
	assert.Equal(t, "plugin2", p2.GetName())
}

func TestGetServicePlugins(t *testing.T) {
	// 清空已有插件，以便测试
	servicePlugins = orderedmap.NewOrderedMap[string, ServicePlugin]()

	// 创建测试插件
	plugin1 := &mockServicePlugin{name: "plugin1"}
	plugin2 := &mockServicePlugin{name: "plugin2"}

	// 注册插件
	servicePlugins.Set(plugin1.GetName(), plugin1)
	servicePlugins.Set(plugin2.GetName(), plugin2)

	// 获取插件
	plugins := GetServicePlugins()

	// 验证返回的插件集合
	assert.Equal(t, 2, plugins.Len())
	assert.Same(t, servicePlugins, plugins, "GetServicePlugins应返回servicePlugins变量")

	// 验证插件内容
	p1, ok := plugins.Get("plugin1")
	assert.True(t, ok)
	assert.Equal(t, plugin1, p1)

	p2, ok := plugins.Get("plugin2")
	assert.True(t, ok)
	assert.Equal(t, plugin2, p2)
}

// 测试覆盖插件注册
func TestRegisterServicePlugin_Override(t *testing.T) {
	// 清空已有插件，以便测试
	servicePlugins = orderedmap.NewOrderedMap[string, ServicePlugin]()

	// 创建测试插件
	plugin1 := &mockServicePlugin{name: "plugin"}
	plugin2 := &mockServicePlugin{name: "plugin"} // 相同名称的插件

	// 注册插件
	RegisterServicePlugin(plugin1)
	RegisterServicePlugin(plugin2)

	// 验证插件是否被覆盖
	plugins := GetServicePlugins()
	assert.Equal(t, 1, plugins.Len())

	// 验证获取的是新插件
	p, ok := plugins.Get("plugin")
	assert.True(t, ok)
	assert.Equal(t, plugin2, p)
}
