package codegen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	v3 "github.com/unionj-cloud/go-doudou/v2/cmd/internal/protobuf/v3"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	"path/filepath"
	"testing"
)

func TestGenGrpcProto(t *testing.T) {
	svcfile := filepath.Join(testDir, "svc.go")
	ic := astutils.BuildInterfaceCollector(svcfile, astutils.ExprString)
	p := v3.NewProtoGenerator(v3.WithFieldNamingFunc(strcase.ToLowerCamel))
	ParseDtoGrpc(testDir, p, "dto")
	_, gotProtoFile := GenGrpcProto(testDir, ic, p)
	fmt.Println(gotProtoFile)
}
