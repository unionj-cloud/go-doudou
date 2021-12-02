package codegen

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestInitProj(t *testing.T) {
	dir := filepath.Join("testdata", "init")
	os.MkdirAll(dir, os.ModePerm)
	defer os.RemoveAll(dir)
	type args struct {
		dir     string
		modName string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				dir:     filepath.Join("testdata", "init"),
				modName: "testinit",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				InitProj(tt.args.dir, tt.args.modName)
			})
		})
	}
}
