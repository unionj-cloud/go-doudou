package database

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/codegen"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/validate"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/constants"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/executils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/gorm_gen"
	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/protobuf/v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var OrmGeneratorRegistry = make(map[OrmKind]IOrmGenerator)

type OrmKind string

func RegisterOrmGenerator(kind OrmKind, instance IOrmGenerator) {
	OrmGeneratorRegistry[kind] = instance
}

func GetOrmGenerator(kind OrmKind) IOrmGenerator {
	if gen, ok := OrmGeneratorRegistry[kind]; ok {
		return gen
	}
	return nil
}

type OrmGeneratorConfig struct {
	Driver string
	Dsn    string
	Dir    string
	Soft   string
	Grpc   bool
}

type IOrmGenerator interface {
	svcGo()
	svcImplGo()
	svcImplGrpc(v3.Service)
	dto()
	Initialize(conf OrmGeneratorConfig)
	GenService()
	ProtoFieldNamingFn() func(string) string
}

var _ IOrmGenerator = (*AbstractBaseGenerator)(nil)

type AbstractBaseGenerator struct {
	Driver              string
	Dsn                 string
	Dir                 string
	g                   *gen.Generator
	Jsonattrcase        string
	Omitempty           bool
	AllowGetWithReqBody bool
	Client              bool
	Grpc                bool
	Env                 string
	impl                IOrmGenerator
	runner              executils.Runner
}

func (b *AbstractBaseGenerator) ProtoFieldNamingFn() func(string) string {
	return b.impl.ProtoFieldNamingFn()
}

func (b *AbstractBaseGenerator) svcImplGrpc(grpcService v3.Service) {
	b.impl.svcImplGrpc(grpcService)
}

func (b *AbstractBaseGenerator) svcGo() {
	b.impl.svcGo()
}

func (b *AbstractBaseGenerator) svcImplGo() {
	b.impl.svcImplGo()
}

func (b *AbstractBaseGenerator) dto() {
	b.impl.dto()
}

func (b *AbstractBaseGenerator) Initialize(conf OrmGeneratorConfig) {
	//TODO implement me
	panic("implement me")
}

func (b *AbstractBaseGenerator) GenService() {
	envfile := filepath.Join(b.Dir, ".env")
	envSource, err := ioutil.ReadFile(envfile)
	if err != nil {
		panic(err)
	}
	envContent := string(envSource)
	if !strings.Contains(envContent, "GDD_DB_DRIVER") {
		envContent += fmt.Sprintf(`GDD_DB_DRIVER=%s`, b.Driver)
		envContent += constants.LineBreak
	}
	if !strings.Contains(envContent, "GDD_DB_DSN") {
		envContent += fmt.Sprintf(`GDD_DB_DSN=%s`, b.Dsn)
		envContent += constants.LineBreak
	}
	ioutil.WriteFile(envfile, []byte(envContent), os.ModePerm)

	b.dto()

	wd, _ := os.Getwd()
	os.Chdir(filepath.Join(b.Dir, "dto"))
	err = b.runner.Run("go", "generate", "./...")
	if err != nil {
		panic(err)
	}
	os.Chdir(wd)

	b.svcGo()

	validate.ValidateDataType(b.Dir)
	ic := astutils.BuildInterfaceCollector(filepath.Join(b.Dir, "svc.go"), astutils.ExprString)
	validate.ValidateRestApi(b.Dir, ic)

	codegen.GenConfig(b.Dir)

	b.svcImplGo()

	codegen.GenHttpMiddleware(b.Dir)
	codegen.GenMain(b.Dir, ic)
	codegen.GenHttpHandler(b.Dir, ic, 0)
	var caseConvertor func(string) string
	switch b.Jsonattrcase {
	case "snake":
		caseConvertor = strcase.ToSnake
	default:
		caseConvertor = strcase.ToLowerCamel
	}
	codegen.GenHttpHandlerImpl(b.Dir, ic, codegen.GenHttpHandlerImplConfig{
		Omitempty:           b.Omitempty,
		AllowGetWithReqBody: b.AllowGetWithReqBody,
		CaseConvertor:       caseConvertor,
	})
	if b.Client {
		codegen.GenGoIClient(b.Dir, ic)
		codegen.GenGoClient(b.Dir, ic, codegen.GenGoClientConfig{
			Env:                 b.Env,
			AllowGetWithReqBody: b.AllowGetWithReqBody,
			CaseConvertor:       caseConvertor,
		})
		codegen.GenGoClientProxy(b.Dir, ic)
	}
	codegen.GenDoc(b.Dir, ic, codegen.GenDocConfig{
		AllowGetWithReqBody: b.AllowGetWithReqBody,
	})

	if b.Grpc {
		p := v3.NewProtoGenerator(v3.WithFieldNamingFunc(b.ProtoFieldNamingFn()))
		codegen.ParseDtoGrpc(b.Dir, p, "dto")
		grpcSvc, protoFile := codegen.GenGrpcProto(b.Dir, ic, p)
		protoFile, _ = filepath.Rel(b.Dir, protoFile)
		fmt.Println(protoFile)
		wd, _ = os.Getwd()
		os.Chdir(filepath.Join(b.Dir))
		// protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative transport/grpc/helloworld.proto
		if err = b.runner.Run("protoc", "--proto_path=.",
			"--go_out=.",
			"--go_opt=paths=source_relative",
			"--go-grpc_out=.",
			"--go-grpc_opt=paths=source_relative",
			protoFile); err != nil {
			panic(err)
		}
		os.Chdir(wd)
		b.svcImplGrpc(grpcSvc)
		codegen.GenMainGrpc(b.Dir, ic, grpcSvc)
		codegen.FixModGrpc(b.Dir)
		codegen.GenMethodAnnotationStore(b.Dir, ic)
	}

	wd, _ = os.Getwd()
	os.Chdir(filepath.Join(b.Dir))
	err = b.runner.Run("go", "mod", "tidy")
	if err != nil {
		panic(err)
	}
	os.Chdir(wd)
}
