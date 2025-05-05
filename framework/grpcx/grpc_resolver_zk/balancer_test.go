package grpc_resolver_zk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

func TestWeightPickerPick(t *testing.T) {
	// 创建带有权重信息的地址
	addrs := []balancer.SubConn{
		&mockSubConn{id: "conn1"}, // 模拟 SubConn 1
		&mockSubConn{id: "conn2"}, // 模拟 SubConn 2
		&mockSubConn{id: "conn3"}, // 模拟 SubConn 3
	}

	// 创建地址信息，设置不同的权重
	addrInfo := []base.SubConnInfo{
		{
			Address: resolver.Address{
				Addr:               "host1:8080",
				BalancerAttributes: attributes.New(WeightAttributeKey{}, WeightAddrInfo{Weight: 10}),
			},
		},
		{
			Address: resolver.Address{
				Addr:               "host2:8080",
				BalancerAttributes: attributes.New(WeightAttributeKey{}, WeightAddrInfo{Weight: 5}),
			},
		},
		{
			Address: resolver.Address{
				Addr:               "host3:8080",
				BalancerAttributes: attributes.New(WeightAttributeKey{}, WeightAddrInfo{Weight: 1}),
			},
		},
	}

	// 创建选择器
	picker := newWeightPicker(addrs, addrInfo)

	// 测试多次选择，确保每个连接都能被选中
	// 注意：由于权重随机性，这个测试可能偶尔失败，但概率很低
	connCounts := make(map[string]int)
	for i := 0; i < 1000; i++ {
		pickResult, err := picker.Pick(balancer.PickInfo{})
		assert.NoError(t, err)

		// 执行 pickResult.Done 以模拟完成调用
		if pickResult.Done != nil {
			pickResult.Done(balancer.DoneInfo{})
		}

		conn := pickResult.SubConn.(*mockSubConn)
		connCounts[conn.id]++
	}

	// 验证所有连接都被使用
	assert.Contains(t, connCounts, "conn1")
	assert.Contains(t, connCounts, "conn2")
	assert.Contains(t, connCounts, "conn3")

	// 验证权重影响 - conn1(权重10)应该比conn3(权重1)被选择更多次
	assert.Greater(t, connCounts["conn1"], connCounts["conn3"])

	t.Logf("选择次数: %v", connCounts)
}

func TestWeightAddrList(t *testing.T) {
	// 创建测试数据
	list := weightAddrList{
		{addr: "host1:8080", weight: 10, curWeight: 0},
		{addr: "host2:8080", weight: 5, curWeight: 0},
		{addr: "host3:8080", weight: 1, curWeight: 0},
	}

	// 测试 Len 方法
	assert.Equal(t, 3, list.Len())

	// 测试 Swap 方法
	list.Swap(0, 2)
	assert.Equal(t, "host3:8080", list[0].addr)
	assert.Equal(t, "host1:8080", list[2].addr)

	// 测试 Less 方法 - 按权重排序
	list = weightAddrList{
		{addr: "host1:8080", weight: 10, curWeight: 0},
		{addr: "host2:8080", weight: 5, curWeight: 0},
		{addr: "host3:8080", weight: 1, curWeight: 0},
	}
	assert.False(t, list.Less(0, 1)) // host1不小于host2
	assert.False(t, list.Less(1, 2)) // host2不小于host3
	assert.False(t, list.Less(0, 2)) // host1不小于host3

	// 测试按当前权重排序
	list = weightAddrList{
		{addr: "host1:8080", weight: 10, curWeight: 5},
		{addr: "host2:8080", weight: 5, curWeight: 10},
		{addr: "host3:8080", weight: 1, curWeight: 1},
	}
	assert.True(t, list.Less(0, 1))  // host1当前权重小于host2
	assert.False(t, list.Less(1, 2)) // host2当前权重不小于host3
	assert.False(t, list.Less(0, 2)) // host1当前权重不小于host3
}

func TestWeightChooserPick(t *testing.T) {
	// 创建测试数据
	addrs := weightAddrList{
		{addr: "host1:8080", weight: 10, curWeight: 0},
		{addr: "host2:8080", weight: 5, curWeight: 0},
		{addr: "host3:8080", weight: 1, curWeight: 0},
	}

	// 创建选择器
	chooser := newWeightChooser(addrs)

	// 测试多次选择，验证选择的分布
	connCounts := make(map[string]int)
	for i := 0; i < 1000; i++ {
		addr := chooser.pick()
		connCounts[addr]++
	}

	// 验证所有地址都被选择
	assert.Contains(t, connCounts, "host1:8080")
	assert.Contains(t, connCounts, "host2:8080")
	assert.Contains(t, connCounts, "host3:8080")

	// 验证权重影响 - host1(权重10)应该比host3(权重1)被选择更多次
	assert.Greater(t, connCounts["host1:8080"], connCounts["host3:8080"])

	t.Logf("选择次数: %v", connCounts)
}

func TestNewBuilder(t *testing.T) {
	// 测试创建构建器
	builder := newBuilder()
	assert.NotNil(t, builder)
	assert.IsType(t, &weightPickerBuilder{}, builder)
}

// mockSubConn 是用于测试的模拟 SubConn
type mockSubConn struct {
	id string
}

func (m *mockSubConn) UpdateAddresses([]resolver.Address) {
	// 不需要实现
}

func (m *mockSubConn) Connect() {
	// 不需要实现
}

func (m *mockSubConn) GetOrBuildProducer(balancer.ProducerBuilder) (balancer.Producer, func()) {
	return nil, func() {}
}
