package codegen

import (
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/validate"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/templates"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/common"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/executils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"github.com/unionj-cloud/go-doudou/v2/version"
	"os"
	"path/filepath"
	"text/template"
)

const svcTmpl = templates.EditableHeaderTmpl + `package service

import (
	"context"
	"{{.DtoPackage}}"
)

//go:generate go-doudou svc http -c
//go:generate go-doudou svc grpc

type {{.SvcName}} interface {
	// You can define your service methods as your need. Below is an example.
	// You can also add annotations here like @role(admin) to add meta data to routes for 
	// implementing your own middlewares
	PostUser(ctx context.Context, user dto.GddUser) (data int32, err error)
	GetUser_Id(ctx context.Context, id int32) (data dto.GddUser, err error)
	PutUser(ctx context.Context, user dto.GddUser) error
	DeleteUser_Id(ctx context.Context, id int32) error
	GetUsers(ctx context.Context, parameter dto.Parameter) (data dto.Page, err error)
}
`

const dtoTmpl = templates.EditableHeaderTmpl + `package dto

//go:generate go-doudou name --file $GOFILE --form

type GddUser struct {
	Id    int32
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
	Page    string
	Size    string
	Sort    string
	Order   string
	Fields  string
	Filters []interface{}
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
	github.com/jmoiron/sqlx v1.3.1
	github.com/kelseyhightower/envconfig v1.4.0
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
	Dir      string
	ModName  string
	Runner   executils.Runner
	GenSvcGo bool
	Module   bool
}

// InitProj inits a service project
// dir is root path
// modName is module name
func InitProj(conf InitProjConfig) {
	var (
		err       error
		svcName   string
		svcfile   string
		dtodir    string
		dtofile   string
		goVersion string
		f         *os.File
		tpl       *template.Template
		envfile   string
	)
	dir, modName, runner, genSvcGo, module := conf.Dir, conf.ModName, conf.Runner, conf.GenSvcGo, conf.Module
	if stringutils.IsEmpty(dir) {
		dir, _ = os.Getwd()
	}
	_ = os.MkdirAll(dir, os.ModePerm)

	common.InitGitRepo(dir)
	common.GitIgnore(dir)

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
		logrus.Warnf("file %s already exists", modfile)
	}

	envfile = filepath.Join(dir, ".env")
	if _, err = os.Stat(envfile); os.IsNotExist(err) {
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
	} else {
		logrus.Warnf("file %s already exists", envfile)
	}

	dtodir = filepath.Join(dir, "dto")
	if err = os.MkdirAll(dtodir, os.ModePerm); err != nil {
		panic(err)
	}
	dtofile = filepath.Join(dtodir, "dto.go")
	if _, err = os.Stat(dtofile); os.IsNotExist(err) {
		if f, err = os.Create(dtofile); err != nil {
			panic(err)
		}
		defer f.Close()

		tpl, _ = template.New("dto.go.tmpl").Parse(dtoTmpl)
		_ = tpl.Execute(f, struct {
			Version string
		}{
			Version: version.Release,
		})
	} else {
		logrus.Warnf("file %s already exists", dtofile)
	}

	if genSvcGo {
		svcName = strcase.ToCamel(filepath.Base(dir))
		svcfile = filepath.Join(dir, "svc.go")
		if _, err = os.Stat(svcfile); os.IsNotExist(err) {
			if f, err = os.Create(svcfile); err != nil {
				panic(err)
			}
			defer f.Close()

			tpl, _ = template.New("svc.go.tmpl").Parse(svcTmpl)
			_ = tpl.Execute(f, struct {
				DtoPackage string
				SvcName    string
				Version    string
			}{
				DtoPackage: filepath.Join(modName, "dto"),
				SvcName:    svcName,
				Version:    version.Release,
			})
		} else {
			logrus.Warnf("file %s already exists", svcfile)
		}
	}

	dockerfile := filepath.Join(dir, "Dockerfile")
	if _, err = os.Stat(dockerfile); os.IsNotExist(err) {
		if f, err = os.Create(dockerfile); err != nil {
			panic(err)
		}
		defer f.Close()

		tpl, _ = template.New("dockerfile.tmpl").Parse(dockerfileTmpl)
		_ = tpl.Execute(f, nil)
	} else {
		logrus.Warnf("file %s already exists", dockerfile)
	}

	dockerignorefile := filepath.Join(dir, ".dockerignore")
	if _, err = os.Stat(dockerignorefile); os.IsNotExist(err) {
		if f, err = os.Create(dockerignorefile); err != nil {
			panic(err)
		}
		defer f.Close()

		tpl, _ = template.New("dockerignorefile.tmpl").Parse(dockerignorefileTmpl)
		_ = tpl.Execute(f, nil)
	} else {
		logrus.Warnf("file %s already exists", dockerignorefile)
	}

	if module {
		ParseDto(dir, "vo")
		ParseDto(dir, "dto")
		validate.ValidateDataType(dir)
		ic := astutils.BuildInterfaceCollector(filepath.Join(dir, "svc.go"), astutils.ExprString)
		validate.ValidateRestApi(dir, ic)
		genHttp(dir, ic)

		// TODO
		// genGrpc(dir, ic)
		genPlugin(dir, ic)
		genMainModule(dir)
	}
}

// InitSvc inits a service project, test purpose only
func InitSvc(dir string) {
	var (
		err       error
		svcName   string
		svcfile   string
		dtodir    string
		dtofile   string
		goVersion string
		f         *os.File
		tpl       *template.Template
		envfile   string
	)
	if stringutils.IsEmpty(dir) {
		dir, _ = os.Getwd()
	}
	_ = os.MkdirAll(dir, os.ModePerm)

	common.InitGitRepo(dir)
	common.GitIgnore(dir)

	goVersion, err = common.GetGoVersionNum(executils.CmdRunner{})
	if err != nil {
		panic(err)
	}
	modName := filepath.Base(dir)
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
		logrus.Warnf("file %s already exists", "go.mod")
	}

	envfile = filepath.Join(dir, ".env")
	if _, err = os.Stat(envfile); os.IsNotExist(err) {
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
	} else {
		logrus.Warnf("file %s already exists", dtofile)
	}

	dtodir = filepath.Join(dir, "dto")
	if err = os.MkdirAll(dtodir, os.ModePerm); err != nil {
		panic(err)
	}
	dtofile = filepath.Join(dtodir, "dto.go")
	if _, err = os.Stat(dtofile); os.IsNotExist(err) {
		if f, err = os.Create(dtofile); err != nil {
			panic(err)
		}
		defer f.Close()

		tpl, _ = template.New("dto.go.tmpl").Parse(dtoTmpl)
		_ = tpl.Execute(f, struct {
			Version string
		}{
			Version: version.Release,
		})
	} else {
		logrus.Warnf("file %s already exists", dtofile)
	}

	svcName = strcase.ToCamel(filepath.Base(dir))
	svcfile = filepath.Join(dir, "svc.go")
	if _, err = os.Stat(svcfile); os.IsNotExist(err) {
		if f, err = os.Create(svcfile); err != nil {
			panic(err)
		}
		defer f.Close()

		tpl, _ = template.New("svc.go.tmpl").Parse(svcTmpl)
		_ = tpl.Execute(f, struct {
			DtoPackage string
			SvcName    string
			Version    string
		}{
			DtoPackage: filepath.Join(modName, "dto"),
			SvcName:    svcName,
			Version:    version.Release,
		})
	} else {
		logrus.Warnf("file %s already exists", svcfile)
	}

	dockerfile := filepath.Join(dir, "Dockerfile")
	if _, err = os.Stat(dockerfile); os.IsNotExist(err) {
		if f, err = os.Create(dockerfile); err != nil {
			panic(err)
		}
		defer f.Close()

		tpl, _ = template.New("dockerfile.tmpl").Parse(dockerfileTmpl)
		_ = tpl.Execute(f, nil)
	} else {
		logrus.Warnf("file %s already exists", dockerfile)
	}
}
