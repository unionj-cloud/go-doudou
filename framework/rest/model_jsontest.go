//go:build jsontest
// +build jsontest

// 这个文件仅在使用jsontest构建标签时编译
// 运行测试时使用: go test -tags=jsontest

package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestJsonMarshalIndent 测试JsonMarshalIndent函数
// 由于包中json变量与编码器冲突，所以单独隔离测试
func TestJsonMarshalIndent(t *testing.T) {
	// 测试数据
	data := map[string]interface{}{
		"name": "John",
		"age":  30,
		"html": "<script>alert('XSS')</script>",
	}

	// 不转义HTML
	result, err := JsonMarshalIndent(data, "", "  ", true)
	assert.NoError(t, err)
	assert.Contains(t, result, `<script>alert('XSS')</script>`)

	// 转义HTML
	result, err = JsonMarshalIndent(data, "", "  ", false)
	assert.NoError(t, err)
	assert.Contains(t, result, `\u003cscript\u003ealert('XSS')\u003c/script\u003e`)
}
