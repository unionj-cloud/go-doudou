package codegen

import (
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/cmd/internal/executils"
	"os"
	"path/filepath"
	"testing"
)

func TestInitProj(t *testing.T) {
	dir := filepath.Join("testdata", "init")
	os.MkdirAll(dir, os.ModePerm)
	defer os.RemoveAll(dir)
	assert.NotPanics(t, func() {
		InitProj(dir, "testinit", executils.CmdRunner{})
	})
}

func Test_getGoVersion(t *testing.T) {
	type args struct {
		runtimeVersion string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				runtimeVersion: "go1.17",
			},
			want: "1.17",
		},
		{
			name: "",
			args: args{
				runtimeVersion: "go1.17.8",
			},
			want: "1.17",
		},
		{
			name: "",
			args: args{
				runtimeVersion: "go1.13",
			},
			want: "1.13",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getGoVersionNum(tt.args.runtimeVersion), "getGoVersionNum(%v)", tt.args.runtimeVersion)
		})
	}
}
