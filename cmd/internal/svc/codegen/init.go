package codegen

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/openapi/v3/codegen/server"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/parser"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/templates"
	"github.com/unionj-cloud/toolkit/astutils"
	"github.com/unionj-cloud/toolkit/common"
	"github.com/unionj-cloud/toolkit/executils"
	v3 "github.com/unionj-cloud/toolkit/protobuf/v3"
	"github.com/unionj-cloud/toolkit/stringutils"
	"github.com/unionj-cloud/go-doudou/v2/version"
)

const svcTmpl = templates.EditableHeaderTmpl + `package service

import (
	"context"
	"{{.DtoPackage}}"
)

{{- if eq .ProjectType "rest" }}
//go:generate go-doudou svc http --case {{ .JsonCase }}
{{- else }}
//go:generate go-doudou svc grpc --http2grpc --case {{ .JsonCase }}
{{- end }}

type {{.SvcName}} interface {
	// You can define your service methods as your need. Below is an example.
	// You can also add annotations here like @role(admin) to add meta data to routes for 
	// implementing your own middlewares
	PostUser(ctx context.Context, user dto.GddUser) (data dto.GddUser, err error)
	PutUser(ctx context.Context, user dto.GddUser) error
	DeleteUser(ctx context.Context, user dto.GddUser) error
	GetUsers(ctx context.Context, parameter dto.Parameter) (data dto.Page, err error)
}
`

const dtoTmpl = templates.EditableHeaderTmpl + `package dto

//go:generate go-doudou name --file $GOFILE --form --case {{ .JsonCase }}

type GddUser struct {
	Id    int64
	Name  string
	Phone string
	Dept  string
}

// Page result wrapper
type Page struct {
	Items      []interface{}
	Page       int64
	Size       int64
	MaxPage    int64
	TotalPages int64
	Total      int64
	Last       bool
	First      bool
	Visible    int64
}

// Parameter struct
type Parameter struct {
	Page    int64
	Size    int64
	Sort    string
	Order   string
	Fields  string
	Filters string
}

func (receiver Parameter) GetPage() int64 {
	return receiver.Page
}

func (receiver Parameter) GetSize() int64 {
	return receiver.Size
}

func (receiver Parameter) GetSort() string {
	return receiver.Sort
}

func (receiver Parameter) GetOrder() string {
	return receiver.Order
}

func (receiver Parameter) GetFields() string {
	return receiver.Fields
}

func (receiver Parameter) GetFilters() interface{} {
	return receiver.Filters
}

func (receiver Parameter) IParameterInstance() {

}
`

const modTmpl = `module {{.ModName}}

go {{.GoVersion}}

require (
	github.com/ascarter/requestid v0.0.0-20170313220838-5b76ab3d4aee
	github.com/brianvoe/gofakeit/v6 v6.10.0
	github.com/go-resty/resty/v2 v2.6.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gorilla/handlers v1.5.1
	github.com/iancoleman/strcase v0.1.3
	github.com/opentracing-contrib/go-stdlib v1.0.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.28.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2 v2.0.0-rc.2
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.0.0-rc.2
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/prometheus/client_golang v1.11.0
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
	github.com/bytedance/sonic v1.12.3
	github.com/unionj-cloud/toolkit v0.0.1
	github.com/wubin1989/gen v0.0.2
	github.com/unionj-cloud/go-doudou/v2 ` + version.Release + `
)`

const envTmpl = ``

const dockerignorefileTmpl = `**/*.local
`

const dockerfileTmpl = `FROM devopsworks/golang-upx:1.18 AS builder

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
ARG user
ENV HOST_USER=$user

WORKDIR /repo

# all the steps are cached
ADD go.mod .
ADD go.sum .
# if go.mod/go.sum not changed, this step is also cached
RUN go mod download

ADD . ./
RUN go mod vendor

RUN export GDD_VER=$(go list -mod=vendor -m -f '{{` + "`" + `{{` + "`" + `}} .Version {{` + "`" + `}}` + "`" + `}}' github.com/unionj-cloud/go-doudou/v2) && \
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags="-s -w -X 'github.com/unionj-cloud/go-doudou/v2/framework/buildinfo.BuildUser=$HOST_USER' -X 'github.com/unionj-cloud/go-doudou/v2/framework/buildinfo.BuildTime=$(date)' -X 'github.com/unionj-cloud/go-doudou/v2/framework/buildinfo.GddVer=$GDD_VER'" -mod vendor -o api cmd/main.go && \
strip api && /usr/local/bin/upx api

FROM alpine:3.14

COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ="Asia/Shanghai"

WORKDIR /repo

COPY --from=builder /repo/api ./

COPY .env* ./

ENTRYPOINT ["/repo/api"]
`

type InitProjConfig struct {
	Dir            string
	ModName        string
	Runner         executils.Runner
	Module         bool
	ProtoGenerator v3.ProtoGenerator
	JsonCase       string
	DocPath        string
	ProjectType    string
}

// InitProj inits a service project
// dir is root path
// modName is module name
func InitProj(conf InitProjConfig) {
	var (
		err         error
		svcName     string
		svcfile     string
		dtodir      string
		dtofile     string
		goVersion   string
		f           *os.File
		tpl         *template.Template
		envfile     string
		docPath     string
		projectType string
	)
	dir, modName, runner, module, jsonCase, docPath, projectType := conf.Dir, conf.ModName, conf.Runner, conf.Module, conf.JsonCase, conf.DocPath, conf.ProjectType
	if stringutils.IsEmpty(dir) {
		dir, _ = os.Getwd()
	}
	_ = os.MkdirAll(dir, os.ModePerm)

	goVersion, err = common.GetGoVersionNum(runner)
	if err != nil {
		panic(err)
	}
	if stringutils.IsEmpty(modName) {
		modName = filepath.Base(dir)
	}
	modfile := filepath.Join(dir, "go.mod")
	if _, err = os.Stat(modfile); os.IsNotExist(err) {
		if f, err = os.Create(modfile); err != nil {
			panic(err)
		}
		defer f.Close()

		tpl, _ = template.New("go.mod.tmpl").Parse(modTmpl)
		_ = tpl.Execute(f, struct {
			ModName   string
			GoVersion string
		}{
			ModName:   modName,
			GoVersion: goVersion,
		})
	} else {
		logrus.Warn("Project has already been initialized, it is not safe to be reinitialized")
		return
	}
	if module {
		if err = runner.Run("go", "work", "use", filepath.Base(dir)); err != nil {
			panic(err)
		}
	}

	envfile = filepath.Join(dir, ".env")
	if f, err = os.Create(envfile); err != nil {
		panic(err)
	}
	defer f.Close()

	tpl, _ = template.New(".env.tmpl").Parse(envTmpl)
	_ = tpl.Execute(f, struct {
		SvcName string
	}{
		SvcName: modName,
	})

	if stringutils.IsNotEmpty(docPath) {
		server.GenSvcGo(dir, docPath)
	} else {
		dtodir = filepath.Join(dir, "dto")
		if err = os.MkdirAll(dtodir, os.ModePerm); err != nil {
			panic(err)
		}
		dtofile = filepath.Join(dtodir, "dto.go")
		if f, err = os.Create(dtofile); err != nil {
			panic(err)
		}
		defer f.Close()

		tpl, _ = template.New(dtoTmpl).Parse(dtoTmpl)
		_ = tpl.Execute(f, struct {
			Version  string
			JsonCase string
		}{
			Version:  version.Release,
			JsonCase: jsonCase,
		})

		svcName = strcase.ToCamel(filepath.Base(dir))
		svcfile = filepath.Join(dir, "svc.go")
		if f, err = os.Create(svcfile); err != nil {
			panic(err)
		}
		defer f.Close()

		tpl, _ = template.New(svcTmpl).Parse(svcTmpl)
		_ = tpl.Execute(f, struct {
			DtoPackage  string
			SvcName     string
			Version     string
			JsonCase    string
			ProjectType string
		}{
			DtoPackage:  strings.ReplaceAll(filepath.Join(modName, "dto"), string(os.PathSeparator), "/"),
			SvcName:     svcName,
			Version:     version.Release,
			JsonCase:    jsonCase,
			ProjectType: projectType,
		})
	}

	oldWd, _ := os.Getwd()
	os.Chdir(dir)
	if err = runner.Run("go", "generate", "./..."); err != nil {
		panic(err)
	}
	os.Chdir(oldWd)

	dockerfile := filepath.Join(dir, "Dockerfile")
	if f, err = os.Create(dockerfile); err != nil {
		panic(err)
	}
	defer f.Close()

	tpl, _ = template.New("dockerfile.tmpl").Parse(dockerfileTmpl)
	_ = tpl.Execute(f, nil)

	dockerignorefile := filepath.Join(dir, ".dockerignore")
	if f, err = os.Create(dockerignorefile); err != nil {
		panic(err)
	}
	defer f.Close()

	tpl, _ = template.New("dockerignorefile.tmpl").Parse(dockerignorefileTmpl)
	_ = tpl.Execute(f, nil)

	if module {
		parser.ParseDto(dir, parser.DEFAULT_DTO_PKGS...)
		ic := astutils.BuildInterfaceCollector(filepath.Join(dir, "svc.go"), astutils.ExprString)
		genPlugin(dir, ic, CodeGenConfig{
			ProjectType: projectType,
		})
		genMain(dir, CodeGenConfig{
			ProjectType: projectType,
		})
		mainMainFile := filepath.Join(filepath.Dir(dir), "main", "cmd", "main.go")
		fileContent, err := ioutil.ReadFile(mainMainFile)
		if err != nil {
			panic(err)
		}
		pluginPkg := astutils.GetPkgPath(filepath.Join(dir, "plugin"))
		original := astutils.AppendImportStatements(fileContent, []byte(fmt.Sprintf(`_ "%s"`, pluginPkg)))
		astutils.FixImport(original, mainMainFile)
		// Comment below code due to performance issue
		//if err = runner.Run("go", "work", "sync"); err != nil {
		//	panic(err)
		//}
	}
}

// InitSvc inits a service project, test purpose only
func InitSvc(dir string) {
}
