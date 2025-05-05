package lazyregexp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// 保存原始标志位，以便测试后恢复
	origInTest := inTest
	defer func() { inTest = origInTest }()

	// 测试非测试环境
	inTest = false
	re := New("hello")
	assert.NotNil(t, re)
	assert.Equal(t, "hello", re.str)
	assert.Nil(t, re.rx)

	// 测试测试环境
	inTest = true
	re = New("world")
	assert.NotNil(t, re)
	// 在测试环境中，会立即编译，所以str会被清空
	assert.Empty(t, re.str)
	assert.NotNil(t, re.rx)
}

func TestRegexp_FindStringSubmatch(t *testing.T) {
	re := New("h(e)llo")
	matches := re.FindStringSubmatch("hello world")
	assert.Equal(t, 2, len(matches))
	assert.Equal(t, "hello", matches[0])
	assert.Equal(t, "e", matches[1])
}

func TestRegexp_FindAllString(t *testing.T) {
	re := New("\\w+")
	matches := re.FindAllString("hello world", -1)
	assert.Equal(t, 2, len(matches))
	assert.Equal(t, "hello", matches[0])
	assert.Equal(t, "world", matches[1])
}

func TestRegexp_ReplaceAllString(t *testing.T) {
	re := New("(hello)")
	result := re.ReplaceAllString("hello world", "hi")
	assert.Equal(t, "hi world", result)
}

func TestRegexp_MatchString(t *testing.T) {
	re := New("hello")
	assert.True(t, re.MatchString("hello world"))
	assert.False(t, re.MatchString("hi world"))
}

func TestRegexp_FindString(t *testing.T) {
	re := New("h(e)llo")
	result := re.FindString("hello world")
	assert.Equal(t, "hello", result)
}

func TestRegexp_FindSubmatch(t *testing.T) {
	re := New("h(e)llo")
	matches := re.FindSubmatch([]byte("hello world"))
	assert.Equal(t, 2, len(matches))
	assert.Equal(t, []byte("hello"), matches[0])
	assert.Equal(t, []byte("e"), matches[1])
}

func TestRegexp_FindStringSubmatchIndex(t *testing.T) {
	re := New("h(e)llo")
	index := re.FindStringSubmatchIndex("hello world")
	assert.Equal(t, 4, len(index))
	assert.Equal(t, 0, index[0]) // 匹配开始
	assert.Equal(t, 5, index[1]) // 匹配结束
	assert.Equal(t, 1, index[2]) // 捕获组1开始
	assert.Equal(t, 2, index[3]) // 捕获组1结束
}

func TestRegexp_SubexpNames(t *testing.T) {
	re := New("h(?P<letter>e)llo")
	names := re.SubexpNames()
	assert.Equal(t, 2, len(names))
	assert.Equal(t, "", names[0])       // 整个表达式没有名称
	assert.Equal(t, "letter", names[1]) // 第一个捕获组的名称
}

func TestLazyBuild(t *testing.T) {
	// 创建非测试环境下的延迟编译正则表达式
	saved := inTest
	inTest = false
	defer func() { inTest = saved }()

	re := New("hello")
	assert.NotNil(t, re)
	// 在非测试环境中，不会立即编译
	assert.Equal(t, "hello", re.str)
	assert.Nil(t, re.rx)

	// 执行操作会触发编译
	re.MatchString("hello")
	assert.NotNil(t, re.rx)
	assert.Empty(t, re.str) // 编译后会清空原始字符串
}

func TestBuild(t *testing.T) {
	re := &Regexp{str: "hello"}
	re.build()
	assert.NotNil(t, re.rx)
	assert.Empty(t, re.str)
}
