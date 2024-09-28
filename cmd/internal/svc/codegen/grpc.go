package codegen

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/parser"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/executils"
	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/protobuf/v3"
)

func genGrpc(dir string, ic astutils.InterfaceCollector, runner executils.Runner, protoGenerator v3.ProtoGenerator) {
	parser.ParseDtoGrpc(dir, protoGenerator, "dto")
	grpcSvc, protoFile := GenGrpcProto(dir, ic, protoGenerator)
	oldWd, _ := os.Getwd()
	os.Chdir(dir)
	protoFile = strings.TrimPrefix(protoFile, dir+string(filepath.Separator))
	if err := protoGenerator.Generate(protoFile, runner); err != nil {
		panic(err)
	}
	os.Chdir(oldWd)
	GenSvcImplGrpc(dir, ic, grpcSvc)
	FixModGrpc(dir)
	GenMethodAnnotationStore(dir, ic)
}
