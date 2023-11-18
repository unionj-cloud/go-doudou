package codegen

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/copier"
	"path/filepath"
	"testing"
)

func Test_unimplementedMethods(t *testing.T) {
	ic := astutils.BuildInterfaceCollector(filepath.Join(testDir, "svc.go"), astutils.ExprString)
	var meta astutils.InterfaceMeta
	_ = copier.DeepCopy(ic.Interfaces[0], &meta)
	unimplementedMethods(&meta, filepath.Join(testDir, "transport/httpsrv"), meta.Name+"HandlerImpl")
	fmt.Println(len(meta.Methods))
}
