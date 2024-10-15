package database

import (
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/codegen"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/parser"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/validate"
	"github.com/unionj-cloud/toolkit/astutils"
	"github.com/unionj-cloud/toolkit/executils"
	v3 "github.com/unionj-cloud/toolkit/protobuf/v3"
	"github.com/wubin1989/gen"
	"os"
	"path/filepath"
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
	ProtoGenerator   v3.ProtoGenerator
	Omitempty        bool
}

type IOrmGenerator interface {
	svcGo()
	svcImplGrpc()
	svcImplRest()
	orm()
	fix()
	Initialize(conf OrmGeneratorConfig)
	GenGrpc()
	GenRest()
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
	ProtoGenerator      v3.ProtoGenerator
	Omitempty           bool
	AllowGetWithReqBody bool
	Client              bool
	Env                 string
	impl                IOrmGenerator
	runner              executils.Runner
	ConfigPackage       string
}

func (b *AbstractBaseGenerator) GenRest() {
	b.orm()
	b.svcGo()
	b.svcImplRest()

	validate.DataType(b.Dir, parser.DEFAULT_DTO_PKGS...)
	ic := astutils.BuildInterfaceCollector(filepath.Join(b.Dir, "svc.go"), astutils.ExprString)
	validate.RestApi(b.Dir, ic)

	codegen.GenHttpMiddleware(b.Dir)
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

	b.fix()
	b.goModTidy()
}

func (b *AbstractBaseGenerator) GenDao() {
	b.orm()
	b.goModTidy()
}

func (b *AbstractBaseGenerator) fix() {
	b.impl.fix()
}

func (b *AbstractBaseGenerator) orm() {
	b.impl.orm()
}

func (b *AbstractBaseGenerator) svcImplRest() {
	b.impl.svcImplRest()
}

func (b *AbstractBaseGenerator) svcGo() {
	b.impl.svcGo()
}

func (b *AbstractBaseGenerator) svcImplGrpc() {
	b.impl.svcImplGrpc()
}

func (b *AbstractBaseGenerator) Initialize(conf OrmGeneratorConfig) {
	//TODO implement me
	panic("implement me")
}

func (b *AbstractBaseGenerator) GenGrpc() {
	b.orm()
	b.svcGo()
	b.svcImplGrpc()

	validate.DataType(b.Dir, parser.DEFAULT_DTO_PKGS...)
	ic := astutils.BuildInterfaceCollector(filepath.Join(b.Dir, "svc.go"), astutils.ExprString)
	validate.GrpcApi(b.Dir, ic, true)
	parser.ParseDtoGrpc(b.Dir, b.ProtoGenerator, parser.DEFAULT_DTO_PKGS...)
	_, protoFile := codegen.GenGrpcProto(b.Dir, ic, b.ProtoGenerator)
	protoFile, _ = filepath.Rel(b.Dir, protoFile)
	wd, _ := os.Getwd()
	os.Chdir(filepath.Join(b.Dir))
	if err := b.ProtoGenerator.Generate(protoFile, b.runner); err != nil {
		panic(err)
	}
	os.Chdir(wd)
	codegen.GenHttpHandler(b.Dir, ic, 0)
	codegen.GenHttp2Grpc(b.Dir, ic, codegen.GenHttp2GrpcConfig{
		AllowGetWithReqBody: b.AllowGetWithReqBody,
		CaseConvertor:       b.CaseConverter,
		Omitempty:           b.Omitempty,
	})
	parser.GenDoc(b.Dir, ic, parser.GenDocConfig{
		RoutePatternStrategy: 0,
		AllowGetWithReqBody:  b.AllowGetWithReqBody,
	})
	codegen.GenMethodAnnotationStore(b.Dir, ic)

	b.fix()
	b.goModTidy()
}

func (b *AbstractBaseGenerator) goModTidy() {
	// here go mod tidy cause performance issue on some computer
	//wd, _ := os.Getwd()
	//os.Chdir(filepath.Join(b.Dir))
	//err := b.runner.Run("go", "mod", "tidy")
	//if err != nil {
	//	panic(err)
	//}
	//os.Chdir(wd)
}
