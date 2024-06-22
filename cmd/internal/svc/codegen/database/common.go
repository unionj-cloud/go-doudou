package database

import (
	"fmt"
	"github.com/docker/docker/pkg/fileutils"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/codegen"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/parser"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/validate"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/constants"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/executils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/gormgen"
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
	Driver           string
	Dsn              string
	TablePrefix      string
	TableGlob        string
	TableExcludeGlob string
	GenGenGo         bool
	CaseConverter    func(string) string
	Dir              string
	Soft             string
	Grpc             bool
	Omitempty        bool
}

type IOrmGenerator interface {
	svcGo()
	svcImplGo()
	svcImplGrpc(v3.Service)
	dto()
	orm()
	fix()
	Initialize(conf OrmGeneratorConfig)
	GenService()
	GenDao()
}

var _ IOrmGenerator = (*AbstractBaseGenerator)(nil)

type AbstractBaseGenerator struct {
	Driver              string
	Dsn                 string
	TablePrefix         string
	TableGlob           string
	TableExcludeGlob    string
	GenGenGo            bool
	Dir                 string
	g                   *gormgen.Generator
	CaseConverter       func(string) string
	Omitempty           bool
	AllowGetWithReqBody bool
	Client              bool
	Grpc                bool
	Env                 string
	impl                IOrmGenerator
	runner              executils.Runner
	ConfigPackage       string
}

func (b *AbstractBaseGenerator) GenDao() {
	b.dto()
	b.orm()
	b.goModTidy()
}

func (b *AbstractBaseGenerator) fix() {
	b.impl.fix()
}

func (b *AbstractBaseGenerator) orm() {
	b.impl.orm()
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
	b.dto()
	b.svcGo()

	validate.DataType(b.Dir)
	ic := astutils.BuildInterfaceCollector(filepath.Join(b.Dir, "svc.go"), astutils.ExprString)
	validate.RestApi(b.Dir, ic)

	serviceName := ic.Interfaces[0].Name
	envfile := filepath.Join(b.Dir, ".env")
	if _, err := os.Stat(envfile); err != nil {
		fileutils.CreateIfNotExists(envfile, false)
	}
	envSource, err := ioutil.ReadFile(envfile)
	if err != nil {
		panic(err)
	}
	envContent := string(envSource)
	if !strings.Contains(envContent, strings.ToUpper(serviceName)+"_DB_DRIVER") {
		envContent += fmt.Sprintf(`%s_DB_DRIVER=%s`, strings.ToUpper(serviceName), b.Driver)
		envContent += constants.LineBreak
	}
	if !strings.Contains(envContent, strings.ToUpper(serviceName)+"_DB_DSN") {
		envContent += fmt.Sprintf(`%s_DB_DSN=%s`, strings.ToUpper(serviceName), b.Dsn)
		envContent += constants.LineBreak
	}
	ioutil.WriteFile(envfile, []byte(envContent), os.ModePerm)

	codegen.GenConfig(b.Dir, ic)

	cfgPkg := astutils.GetPkgPath(filepath.Join(b.Dir, "config"))
	b.g.ConfigPackage = cfgPkg

	b.orm()
	b.svcImplGo()

	codegen.GenHttpMiddleware(b.Dir)
	codegen.GenMain(b.Dir, ic)
	codegen.GenHttpHandler(b.Dir, ic, 0)
	codegen.GenHttpHandlerImpl(b.Dir, ic, codegen.GenHttpHandlerImplConfig{
		Omitempty:           b.Omitempty,
		AllowGetWithReqBody: b.AllowGetWithReqBody,
		CaseConvertor:       b.CaseConverter,
	})
	if b.Client {
		codegen.GenGoIClient(b.Dir, ic)
		codegen.GenGoClient(b.Dir, ic, codegen.GenGoClientConfig{
			Env:                 b.Env,
			AllowGetWithReqBody: b.AllowGetWithReqBody,
			CaseConvertor:       b.CaseConverter,
		})
		codegen.GenGoClientProxy(b.Dir, ic)
	}
	parser.GenDoc(b.Dir, ic, parser.GenDocConfig{
		AllowGetWithReqBody: b.AllowGetWithReqBody,
	})

	if b.Grpc {
		p := v3.NewProtoGenerator(v3.WithFieldNamingFunc(b.CaseConverter))
		parser.ParseDtoGrpc(b.Dir, p, "dto")
		grpcSvc, protoFile := codegen.GenGrpcProto(b.Dir, ic, p)
		protoFile, _ = filepath.Rel(b.Dir, protoFile)
		wd, _ := os.Getwd()
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

	b.fix()
	b.goModTidy()
}

func (b *AbstractBaseGenerator) goModTidy() {
	wd, _ := os.Getwd()
	os.Chdir(filepath.Join(b.Dir))
	err := b.runner.Run("go", "mod", "tidy")
	if err != nil {
		panic(err)
	}
	os.Chdir(wd)
}
