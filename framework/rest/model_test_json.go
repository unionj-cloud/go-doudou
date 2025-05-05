//go:build !jsontest
// +build !jsontest

package rest

import (
	"testing"
)

// TestJsonMarshalIndent 测试JsonMarshalIndent函数
// 由于包中json变量与包冲突问题，我们将这个测试放在单独的文件中
// 使用build tag隔离，需要单独测试时使用：go test -tags=jsontest
func TestJsonMarshalIndent_Skip(t *testing.T) {
	t.Skip("使用独立的构建标签进行测试，避免json包冲突")
}
