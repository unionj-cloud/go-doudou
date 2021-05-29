package codegen

import (
	"github.com/unionj-cloud/go-doudou/astutils"
	"path/filepath"
	"testing"
)

func TestGenHttpHandlerImplWithImpl(t *testing.T) {
	dir := "/Users/wubin1989/workspace/cloud/usersvc"
	svcfile := filepath.Join(dir, "svc.go")
	ic := buildIc(svcfile)

	type args struct {
		dir string
		ic  astutils.InterfaceCollector
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				dir: dir,
				ic:  ic,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenHttpHandlerImplWithImpl(tt.args.dir, tt.args.ic)
		})
	}
}
