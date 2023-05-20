package database

import (
	"github.com/iancoleman/strcase"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/codegen"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/validate"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/executils"
	"gorm.io/gen"
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
	Driver string
	Dsn    string
	Dir    string
	Soft   string
}

type IOrmGenerator interface {
	svcGo()
	svcImplGo()
	dto()
	Initialize(conf OrmGeneratorConfig)
	GenService()
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
	Env                 string
	impl                IOrmGenerator
	runner              executils.Runner
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

	wd, _ := os.Getwd()
	os.Chdir(filepath.Join(b.Dir, "dto"))
	err := b.runner.Run("go", "generate", "./...")
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
}
