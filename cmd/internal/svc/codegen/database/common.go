package database

import (
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/codegen"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/parser"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/validate"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/executils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/gormgen"
	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/protobuf/v3"
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
	svcImplGo()
	svcImplGrpc(v3.Service)
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
	ProtoGenerator      v3.ProtoGenerator
	Omitempty           bool
	AllowGetWithReqBody bool
	Client              bool
	Env                 string
	impl                IOrmGenerator
	runner              executils.Runner
	ConfigPackage       string
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

func (b *AbstractBaseGenerator) svcImplGrpc(grpcService v3.Service) {
	b.impl.svcImplGrpc(grpcService)
}

func (b *AbstractBaseGenerator) svcGo() {
	b.impl.svcGo()
}

func (b *AbstractBaseGenerator) svcImplGo() {
	b.impl.svcImplGo()
}

func (b *AbstractBaseGenerator) Initialize(conf OrmGeneratorConfig) {
	//TODO implement me
	panic("implement me")
}

func (b *AbstractBaseGenerator) GenService() {
	b.orm()
	b.svcGo()
	b.svcImplGo()

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
	wd, _ := os.Getwd()
	os.Chdir(filepath.Join(b.Dir))
	err := b.runner.Run("go", "mod", "tidy")
	if err != nil {
		panic(err)
	}
	os.Chdir(wd)
}
