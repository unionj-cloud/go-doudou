package grpc_resolver_zk

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

// 基础的测试
func TestNewBuilder(t *testing.T) {
	builder := newBuilder()
	assert.NotNil(t, builder)
	// 验证构建器名称
	assert.Equal(t, Name, builder.Name())
}

// 测试 WeightAttributeKey 和 WeightAddrInfo
func TestWeightAttributes(t *testing.T) {
	info := WeightAddrInfo{Weight: 10}
	assert.Equal(t, 10, info.Weight)

	// 验证 WeightAttributeKey 是一个唯一的类型
	key1 := WeightAttributeKey{}
	key2 := WeightAttributeKey{}
	assert.Equal(t, key1, key2)

	// 检查它们的类型
	keyType := reflect.TypeOf(key1)
	assert.Equal(t, "grpc_resolver_zk.WeightAttributeKey", keyType.String())
}

// 测试 conns 的排序功能
func TestConnsSort(t *testing.T) {
	// 创建测试数据
	list := conns{
		{Weight: 10},
		{Weight: 5},
		{Weight: 1},
	}

	// 测试 Len 方法
	assert.Equal(t, 3, list.Len())

	// 测试 Swap 方法
	list.Swap(0, 2)
	assert.Equal(t, 1, list[0].Weight)
	assert.Equal(t, 10, list[2].Weight)

	// 测试 Less 方法 - 按权重排序
	list = conns{
		{Weight: 10},
		{Weight: 5},
		{Weight: 1},
	}
	assert.False(t, list.Less(0, 1)) // 10不小于5
	assert.False(t, list.Less(1, 2)) // 5不小于1
	assert.False(t, list.Less(0, 2)) // 10不小于1

	// 测试反向排序
	list = conns{
		{Weight: 1},
		{Weight: 5},
		{Weight: 10},
	}
	assert.True(t, list.Less(0, 1)) // 1小于5
	assert.True(t, list.Less(1, 2)) // 5小于10
	assert.True(t, list.Less(0, 2)) // 1小于10
}

// 测试 wPickerBuilder
func TestWPickerBuilder(t *testing.T) {
	builder := &wPickerBuilder{}

	// 测试没有准备好的连接时返回错误选择器
	emptyInfo := base.PickerBuildInfo{
		ReadySCs: map[balancer.SubConn]base.SubConnInfo{},
	}
	picker := builder.Build(emptyInfo)
	_, err := picker.Pick(balancer.PickInfo{})
	assert.Error(t, err)
	assert.Equal(t, balancer.ErrNoSubConnAvailable, err)
}

// 没有必要测试实际选择逻辑，因为无法轻易模拟 SubConn
// 但是可以测试 Chooser 的初始化
func TestChooserInit(t *testing.T) {
	// 创建测试数据
	testConns := conns{
		{Weight: 10},
		{Weight: 5},
		{Weight: 1},
	}

	// 创建选择器
	chooser := newChooser(testConns)

	// 验证总权重计算正确
	assert.Equal(t, 16, chooser.max) // 10 + 5 + 1 = 16

	// 验证累计权重数组计算正确
	assert.Equal(t, []int{10, 15, 16}, chooser.totals)
}
