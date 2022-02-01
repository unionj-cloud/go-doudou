package codegen

import (
	"github.com/iancoleman/strcase"
	"github.com/unionj-cloud/go-doudou/astutils"
	"os"
	"path/filepath"
	"testing"
)

func TestGenGoClient(t *testing.T) {
	dir := testDir + "client1"
	InitSvc(dir)
	defer os.RemoveAll(dir)
	svcfile := filepath.Join(dir, "svc.go")
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
				dir: dir,
				ic:  ic,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenGoClient(tt.args.dir, tt.args.ic, "", 1, strcase.ToLowerCamel)
		})
	}
}

func TestGenGoClient2(t *testing.T) {
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
			GenGoClient(tt.args.dir, tt.args.ic, "", 1, strcase.ToLowerCamel)
		})
	}
}
