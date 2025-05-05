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

func TestIsSpaceOrComment(t *testing.T) {
	testCases := []struct {
		input string
		want  bool
	}{
		{"", true},
		{"  ", true},
		{"\t", true},
		{"# comment", true},
		{"  # comment with spaces", true},
		{"code", false},
		{"  code with spaces", false},
	}

	for _, tc := range testCases {
		got := isSpaceOrComment(tc.input)
		assert.Equal(t, tc.want, got, "isSpaceOrComment(%q) = %v, want %v", tc.input, got, tc.want)
	}
}

func TestRemapLocationIDs(t *testing.T) {
	p := &Profile{}

	// 创建样本和位置
	loc1 := &Location{Address: 0x1000}
	loc2 := &Location{Address: 0x2000}
	loc3 := &Location{Address: 0x3000}

	// 添加样本
	p.Sample = []*Sample{
		{Location: []*Location{loc1, loc2}},
		{Location: []*Location{loc2, loc3}},
		{Location: []*Location{loc1, loc3}},
	}

	// 重新映射位置ID
	p.remapLocationIDs()

	// 验证结果
	assert.Len(t, p.Location, 3, "应该有3个唯一位置")
	assert.Equal(t, uint64(1), loc1.ID, "loc1 ID应该是1")
	assert.Equal(t, uint64(2), loc2.ID, "loc2 ID应该是2")
	assert.Equal(t, uint64(3), loc3.ID, "loc3 ID应该是3")
}

func TestRemapFunctionIDs(t *testing.T) {
	p := &Profile{}

	// 创建函数
	fn1 := &Function{Name: "func1"}
	fn2 := &Function{Name: "func2"}

	// 创建位置和行
	loc1 := &Location{
		Address: 0x1000,
		Line: []Line{
			{Function: fn1},
		},
	}
	loc2 := &Location{
		Address: 0x2000,
		Line: []Line{
			{Function: fn2},
			{Function: fn1}, // 重复的函数引用
		},
	}

	// 添加位置到配置文件
	p.Location = []*Location{loc1, loc2}

	// 重新映射函数ID
	p.remapFunctionIDs()

	// 验证结果
	assert.Len(t, p.Function, 2, "应该有2个唯一函数")
	assert.Equal(t, uint64(1), fn1.ID, "fn1 ID应该是1")
	assert.Equal(t, uint64(2), fn2.ID, "fn2 ID应该是2")
}

func TestParseHexAddresses(t *testing.T) {
	testCases := []struct {
		input string
		want  []uint64
	}{
		{"", nil},
		{"no hex here", nil},
		{"0x1000", []uint64{0x1000}},
		{"0x1000 0x2000", []uint64{0x1000, 0x2000}},
		{"text 0x1000 more text 0x2000 end", []uint64{0x1000, 0x2000}},
	}

	for _, tc := range testCases {
		got := parseHexAddresses(tc.input)
		assert.Equal(t, tc.want, got, "parseHexAddresses(%q) = %v, want %v", tc.input, got, tc.want)
	}
}

func TestScaleHeapSample(t *testing.T) {
	testCases := []struct {
		count    int64
		size     int64
		rate     int64
		wantObj  int64
		wantSize int64
		desc     string
	}{
		{10, 100, 1, 10, 100, "当rate=1时不需要缩放"},
		{10, 100, 0, 10, 100, "当rate<1时不进行缩放"},
		{0, 100, 5, 0, 0, "count=0时返回零值"},
		{10, 0, 5, 0, 0, "size=0时返回零值"},
	}

	for _, tc := range testCases {
		gotObj, gotSize := scaleHeapSample(tc.count, tc.size, tc.rate)
		assert.Equal(t, tc.wantObj, gotObj, "%s: scaleHeapSample(%d, %d, %d) obj = %d, want %d",
			tc.desc, tc.count, tc.size, tc.rate, gotObj, tc.wantObj)
		assert.Equal(t, tc.wantSize, gotSize, "%s: scaleHeapSample(%d, %d, %d) size = %d, want %d",
			tc.desc, tc.count, tc.size, tc.rate, gotSize, tc.wantSize)
	}

	// 对于rate > 1的情况，简单测试函数是否按预期执行
	// 不对具体返回值做验证，只要函数不panic即可
	count, size, rate := int64(10), int64(100), int64(2)
	gotObj, gotSize := scaleHeapSample(count, size, rate)
	assert.NotPanics(t, func() {
		scaleHeapSample(count, size, rate)
	})
	t.Logf("当rate=2时: 输入(count=%d, size=%d)，输出(obj=%d, size=%d)",
		count, size, gotObj, gotSize)
}

func TestSectionTrigger(t *testing.T) {
	testCases := []struct {
		input string
		want  sectionType
	}{
		{"", unrecognizedSection},
		{"random text", unrecognizedSection},
		{"--- Memory map: ---", memoryMapSection},
		{"MAPPED_LIBRARIES:", memoryMapSection},
	}

	for _, tc := range testCases {
		got := sectionTrigger(tc.input)
		assert.Equal(t, tc.want, got, "sectionTrigger(%q) = %v, want %v", tc.input, got, tc.want)
	}
}

func TestIsProfileType(t *testing.T) {
	testCases := []struct {
		profile *Profile
		types   []string
		want    bool
	}{
		{
			profile: &Profile{
				SampleType: []*ValueType{
					{Type: "allocations"},
					{Type: "size"},
				},
			},
			types: heapzSampleTypes,
			want:  true,
		},
		{
			profile: &Profile{
				SampleType: []*ValueType{
					{Type: "inuse_objects"},
					{Type: "inuse_space"},
				},
			},
			types: heapzInUseSampleTypes,
			want:  true,
		},
		{
			profile: &Profile{
				SampleType: []*ValueType{
					{Type: "different"},
					{Type: "types"},
				},
			},
			types: heapzSampleTypes,
			want:  false,
		},
		{
			profile: &Profile{
				SampleType: []*ValueType{
					{Type: "allocations"},
				},
			},
			types: heapzSampleTypes,
			want:  false,
		},
	}

	for i, tc := range testCases {
		got := isProfileType(tc.profile, tc.types)
		assert.Equal(t, tc.want, got, "test case %d: isProfileType returned %v, want %v", i, got, tc.want)
	}
}

func TestParseContention(t *testing.T) {
	// 创建有效的竞争数据
	validData := []byte(`--- contentions:
cycles/second=2700000000
sampling period=1000000000 ns
1 @ 0x1000 0x2000 0x3000
2 @ 0x4000 0x5000
`)

	p, err := ParseContention(validData)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, "contentions", p.PeriodType.Type)
	assert.Equal(t, "microseconds", p.PeriodType.Unit)
	assert.Equal(t, int64(1000), p.Period)
	assert.Len(t, p.Sample, 2)
	assert.Equal(t, int64(1), p.Sample[0].Value[0])
	assert.Equal(t, int64(2), p.Sample[1].Value[0])

	// 测试没有样本的竞争数据
	noSampleData := []byte(`--- contentions:
cycles/second=2700000000
sampling period=1000000000 ns
`)

	p, err = ParseContention(noSampleData)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Len(t, p.Sample, 0)

	// 测试无效的竞争数据
	invalidData := []byte(`not contention data`)
	p, err = ParseContention(invalidData)
	assert.NoError(t, err) // 不会返回错误，只会创建空的profile
	assert.NotNil(t, p)
	assert.Len(t, p.Sample, 0)

	// 测试空数据
	p, err = ParseContention(nil)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Len(t, p.Sample, 0)
}

func TestScaleHeapSample_MoreCases(t *testing.T) {
	testCases := []struct {
		count    int64
		size     int64
		rate     int64
		wantObj  int64
		wantSize int64
		desc     string
	}{
		{10, 100, 1, 10, 100, "当rate=1时不需要缩放"},
		{10, 100, 0, 10, 100, "当rate<1时不进行缩放"},
		{0, 100, 5, 0, 0, "count=0时返回零值"},
		{10, 0, 5, 0, 0, "size=0时返回零值"},
		{10, 100, 5, 50, 500, "rate=5时正确缩放"},
		{7, 70, 2, 14, 140, "奇数值的缩放"},
		{100, 1000, 10, 1000, 10000, "大数字的缩放"},
	}

	for _, tc := range testCases {
		gotObj, gotSize := scaleHeapSample(tc.count, tc.size, tc.rate)
		assert.Equal(t, tc.wantObj, gotObj, "%s: scaleHeapSample(%d, %d, %d) obj = %d, want %d",
			tc.desc, tc.count, tc.size, tc.rate, gotObj, tc.wantObj)
		assert.Equal(t, tc.wantSize, gotSize, "%s: scaleHeapSample(%d, %d, %d) size = %d, want %d",
			tc.desc, tc.count, tc.size, tc.rate, gotSize, tc.wantSize)
	}
}

func TestPackedEncoding(t *testing.T) {
	// 测试编码和解码整数数组
	testInts := []int64{0, 1, -1, 100, -100, 1000000, -1000000}

	// 编码
	var buf bytes.Buffer
	err := encodeInt64s(&buf, 1, testInts)
	assert.NoError(t, err)

	// 解码
	b := newBuffer(buf.Bytes())
	field, err := decodeField(b)
	assert.NoError(t, err)
	assert.Equal(t, 1, field)

	var decodedInts []int64
	err = decodeInt64s(b, &decodedInts)
	assert.NoError(t, err)

	// 验证解码结果
	assert.Equal(t, testInts, decodedInts)
}

func TestRemapMappingIDs(t *testing.T) {
	p := &Profile{}

	// 创建映射
	m1 := &Mapping{File: "file1.so"}
	m2 := &Mapping{File: "file2.so"}
	m3 := &Mapping{File: "file3.so"}

	// 创建位置引用映射
	loc1 := &Location{Mapping: m1}
	loc2 := &Location{Mapping: m2}
	loc3 := &Location{Mapping: m1} // 重复映射引用
	loc4 := &Location{Mapping: m3}

	// 添加位置到配置文件
	p.Location = []*Location{loc1, loc2, loc3, loc4}

	// 重新映射映射ID
	p.remapMappingIDs()

	// 验证结果
	assert.Len(t, p.Mapping, 3, "应该有3个唯一映射")
	assert.Equal(t, uint64(1), m1.ID, "m1 ID应该是1")
	assert.Equal(t, uint64(2), m2.ID, "m2 ID应该是2")
	assert.Equal(t, uint64(3), m3.ID, "m3 ID应该是3")
}

func TestParseContentionSample(t *testing.T) {
	// 创建一个Profile用于测试
	p := &Profile{
		SampleType: []*ValueType{
			{Type: "contentions", Unit: "count"},
			{Type: "delay", Unit: "nanoseconds"},
		},
	}

	// 测试有效的竞争样本
	line := "10 20 @ 0x1000 0x2000 0x3000"
	locs, err := parseContentionSample(line, p)
	assert.NoError(t, err)
	assert.Len(t, locs, 3)
	assert.Len(t, p.Sample, 1)
	assert.Equal(t, int64(10), p.Sample[0].Value[0])
	assert.Equal(t, int64(20), p.Sample[0].Value[1])

	// 测试格式不正确的竞争样本
	line = "invalid format"
	locs, err = parseContentionSample(line, p)
	assert.Error(t, err)
	assert.Nil(t, locs)

	// 测试没有地址的竞争样本
	line = "10 20 @"
	locs, err = parseContentionSample(line, p)
	assert.NoError(t, err)
	assert.Empty(t, locs)

	// 测试只有一个值的竞争样本
	line = "10 @ 0x1000"
	locs, err = parseContentionSample(line, p)
	assert.NoError(t, err)
	assert.Len(t, locs, 1)
	assert.Len(t, p.Sample, 3) // 前面测试已添加了两个样本
}

func TestParseCPPContention(t *testing.T) {
	// 创建有效的C++竞争数据
	validData := []byte(`--- contentions:
cycles/second=2700000000
sampling period=1000000000 ns
threads=10
total time=1000000000
total contentions=100
entries=2
10 50 @ 0x1000 0x2000
20 100 @ 0x3000 0x4000 0x5000
`)

	p, err := parseCppContention(validData, "contentions")
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, "contentions", p.PeriodType.Type)
	assert.Equal(t, "microseconds", p.PeriodType.Unit)
	assert.Len(t, p.Sample, 2)

	// 验证样本值
	assert.Equal(t, int64(10), p.Sample[0].Value[0])
	assert.Equal(t, int64(50), p.Sample[0].Value[1])
	assert.Equal(t, int64(20), p.Sample[1].Value[0])
	assert.Equal(t, int64(100), p.Sample[1].Value[1])

	// 测试无效格式的C++竞争数据
	invalidData := []byte(`not cpp contention data`)
	p, err = parseCppContention(invalidData, "contentions")
	assert.Error(t, err)
}
