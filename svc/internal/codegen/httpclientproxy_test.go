package codegen

import (
	"github.com/unionj-cloud/go-doudou/astutils"
	"path/filepath"
	"testing"
)

func TestGenGoClientProxy(t *testing.T) {
	svcfile := filepath.Join(testDir, "svc.go")
	ic := astutils.BuildInterfaceCollector(svcfile, astutils.ExprString)

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
				dir: testDir,
				ic:  ic,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenGoClientProxy(tt.args.dir, tt.args.ic)
		})
	}
}
