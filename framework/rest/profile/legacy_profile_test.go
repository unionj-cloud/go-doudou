package profile

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTracebacks(t *testing.T) {
	// 测试有效的Tracebacks数据
	validData := []byte(`TestFunction1 0xabcdef 0x123456
TestFunction2 0x654321 0x789abc
another line with text

@0x123456
TestFunction3 0xdef012

memory map
`)

	p, err := ParseTracebacks(validData)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, "trace", p.PeriodType.Type)
	assert.Equal(t, "count", p.PeriodType.Unit)
	assert.Equal(t, int64(1), p.Period)
	assert.Equal(t, 1, len(p.SampleType))
	assert.Equal(t, "trace", p.SampleType[0].Type)
	assert.Equal(t, "count", p.SampleType[0].Unit)
	assert.NotEmpty(t, p.Location)
	assert.NotEmpty(t, p.Sample)

	// 测试空数据
	emptyData := []byte(``)
	p, err = ParseTracebacks(emptyData)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Empty(t, p.Sample)

	// 测试无有效地址的数据
	noAddressData := []byte(`TestFunction1 
TestFunction2
another line with text
`)
	p, err = ParseTracebacks(noAddressData)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	// 没有有效的十六进制地址，应该没有Location
	assert.Empty(t, p.Location)
}

func TestExtractHexAddresses(t *testing.T) {
	// 测试有地址的字符串
	line := "TestFunction 0xabcdef 0x123456"
	s, addrs := extractHexAddresses(line)
	// 函数会返回提取的十六进制字符串和转换后的uint64值
	assert.Equal(t, []string{"0xabcdef", "0x123456"}, s, "s应该包含两个十六进制字符串")
	assert.Equal(t, []uint64{0xabcdef, 0x123456}, addrs, "addrs应该包含两个地址")

	// 测试只有文本的字符串
	line = "TestFunction no hex addresses here"
	s, addrs = extractHexAddresses(line)
	assert.Empty(t, s, "s应该为空")
	assert.Empty(t, addrs, "addrs应该为空")

	// 测试空字符串
	line = ""
	s, addrs = extractHexAddresses(line)
	assert.Empty(t, s)
	assert.Empty(t, addrs)

	// 测试只有地址的字符串
	line = "0xabcdef 0x123456"
	s, addrs = extractHexAddresses(line)
	assert.Equal(t, []string{"0xabcdef", "0x123456"}, s, "s应该包含两个十六进制字符串")
	assert.Equal(t, []uint64{0xabcdef, 0x123456}, addrs, "addrs应该包含两个地址")

	// 测试混合了不同格式的地址和普通文本
	line = "TestFunction 0xabcdef abcdef 0x123"
	s, addrs = extractHexAddresses(line)
	assert.Equal(t, []string{"0xabcdef", "0x123"}, s, "s应该包含两个十六进制字符串")
	assert.Equal(t, []uint64{0xabcdef, 0x123}, addrs, "addrs应该包含两个地址")
}

func TestAddTracebackSample(t *testing.T) {
	// 创建一个新的Profile
	p := &Profile{
		SampleType: []*ValueType{
			{Type: "trace", Unit: "count"},
		},
	}

	// 创建Location数组
	locs := []*Location{
		{Address: 0x123456},
		{Address: 0xabcdef},
	}

	// 创建sources数组
	sources := []string{"Function1", "Function2"}

	// 添加样本
	addTracebackSample(locs, sources, p)

	// 验证样本是否被正确添加
	assert.Equal(t, 1, len(p.Sample))
	assert.Equal(t, []int64{1}, p.Sample[0].Value)
	assert.Equal(t, locs, p.Sample[0].Location)
	assert.Equal(t, map[string][]string{"source": sources}, p.Sample[0].Label)

	// 再次添加另一个样本
	addTracebackSample(locs, sources, p)

	// 验证是否有两个样本了
	assert.Equal(t, 2, len(p.Sample))
}

func TestParseMemoryMap(t *testing.T) {
	// 创建一个Profile
	p := &Profile{}

	// 创建读入器，包含内存映射数据
	data := bytes.NewBufferString(`memory map:
00400000-00452000: /bin/ls
0054d000-00574000: /lib/x86_64/libc-2.22.so`)

	// 解析内存映射
	err := p.ParseMemoryMap(data)
	assert.NoError(t, err)

	// 验证是否正确解析
	assert.Equal(t, 2, len(p.Mapping))

	assert.Equal(t, uint64(0x400000), p.Mapping[0].Start)
	assert.Equal(t, uint64(0x452000), p.Mapping[0].Limit)
	assert.Equal(t, "/bin/ls", p.Mapping[0].File)

	assert.Equal(t, uint64(0x54d000), p.Mapping[1].Start)
	assert.Equal(t, uint64(0x574000), p.Mapping[1].Limit)
	assert.Equal(t, "/lib/x86_64/libc-2.22.so", p.Mapping[1].File)
}
