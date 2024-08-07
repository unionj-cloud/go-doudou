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
	// protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative transport/grpc/helloworld.proto
	if err := runner.Run("protoc", "--proto_path=.",
		"--go_out=.",
		"--go_opt=paths=source_relative",
		"--go-grpc_out=.",
		"--go-grpc_opt=paths=source_relative",
		protoFile); err != nil {
		panic(err)
	}
	os.Chdir(oldWd)
	GenSvcImplGrpc(dir, ic, grpcSvc)
	FixModGrpc(dir)
	GenMethodAnnotationStore(dir, ic)
}
