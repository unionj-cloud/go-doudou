package codegen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/astutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/copier"
	"os"
	"path/filepath"
	"testing"
)

func TestGenHttpHandlerImplWithImpl(t *testing.T) {
	dir := "testdata"
	defer os.RemoveAll(filepath.Join(dir, "transport"))
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
			GenHttpHandlerImpl(tt.args.dir, tt.args.ic, true, strcase.ToLowerCamel)
		})
	}
}

func TestGenHttpHandlerImplWithImpl2(t *testing.T) {
	svcfile := testDir + "/svc.go"
	ic := astutils.BuildInterfaceCollector(svcfile, astutils.ExprString)
	GenHttpHandlerImpl(testDir, ic, true, strcase.ToLowerCamel)
}

func Test_unimplementedMethods(t *testing.T) {
	ic := astutils.BuildInterfaceCollector(filepath.Join(testDir, "svc.go"), astutils.ExprString)
	var meta astutils.InterfaceMeta
	_ = copier.DeepCopy(ic.Interfaces[0], &meta)
	unimplementedMethods(&meta, filepath.Join(testDir, "transport/httpsrv"))
	fmt.Println(len(meta.Methods))
}
