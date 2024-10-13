package codegen

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/unionj-cloud/toolkit/astutils"
	"github.com/unionj-cloud/toolkit/copier"
)

func Test_unimplementedMethods(t *testing.T) {
	ic := astutils.BuildInterfaceCollector(filepath.Join(testDir, "svc.go"), astutils.ExprString)
	var meta astutils.InterfaceMeta
	_ = copier.DeepCopy(ic.Interfaces[0], &meta)
	unimplementedMethods(&meta, filepath.Join(testDir, "transport/httpsrv"), meta.Name+"HandlerImpl")
	fmt.Println(len(meta.Methods))
}
