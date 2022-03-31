package codegen

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/cmd/internal/executils"
	"github.com/unionj-cloud/go-doudou/toolkit/sliceutils"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

const svcTmpl = `package service

import (
	"context"
	"{{.VoPackage}}"
)

type {{.SvcName}} interface {
	// You can define your service methods as your need. Below is an example.
	PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, err error)
}
`

const voTmpl = `package vo

//go:generate go-doudou name --file $GOFILE

type PageFilter struct {
	// 真实姓名，前缀匹配
	Name string
	// 所属部门ID
	Dept int
}

type Order struct {
	Col  string
	Sort string
}

type Page struct {
	// 排序规则
	Orders []Order
	// 页码
	PageNo int
	// 每页行数
	Size int
}

// 分页筛选条件
type PageQuery struct {
	Filter PageFilter
	Page   Page
}

type PageRet struct {
	Items    interface{}
	PageNo   int
	PageSize int
	Total    int
	HasNext  bool
}

type UserVo struct {
	Id    int
	Name  string
	Phone string
	Dept  string
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
	github.com/sirupsen/logrus v1.8.1
	github.com/unionj-cloud/go-doudou v1.0.5
)`

const gitignoreTmpl = `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/
**/*.local
.DS_Store
.idea`

const envTmpl = ``

const dockerfileTmpl = `FROM golang:1.16.6-alpine AS builder

ENV GO111MODULE=on
ARG user
ENV HOST_USER=$user

WORKDIR /repo

ADD go.mod .
ADD go.sum .

ADD . ./

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache bash tzdata

ENV TZ="Asia/Shanghai"

RUN export GDD_VER=$(go list -mod=vendor -m -f '{{` + "`" + `{{ .Version }}` + "`" + `}}' github.com/unionj-cloud/go-doudou) && \
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags="-X 'github.com/unionj-cloud/go-doudou/framework/buildinfo.BuildUser=$HOST_USER' -X 'github.com/unionj-cloud/go-doudou/framework/buildinfo.BuildTime=$(date)' -X 'github.com/unionj-cloud/go-doudou/framework/buildinfo.GddVer=$GDD_VER'" -mod vendor -o api cmd/main.go

ENTRYPOINT ["/repo/api"]
`

func getGoVersionNum(goVersion string) string {
	vnums := sliceutils.StringSlice2InterfaceSlice(strings.Split(strings.TrimPrefix(strings.TrimSpace(goVersion), "go"), "."))
	nums := make([]interface{}, 2)
	copy(nums, vnums)
	return fmt.Sprintf("%s.%s", nums...)
}

// InitProj inits a service project
// dir is root path
// modName is module name
func InitProj(dir string, modName string, runner executils.Runner) {
	var (
		err       error
		svcName   string
		svcfile   string
		modfile   string
		vodir     string
		vofile    string
		goVersion string
		firstLine string
		f         *os.File
		tpl       *template.Template
		envfile   string
		out       []byte
	)
	if stringutils.IsEmpty(dir) {
		dir, _ = os.Getwd()
	}
	_ = os.MkdirAll(dir, os.ModePerm)

	gitInit(dir)
	gitIgnore(dir)

	if out, err = runner.Output("go", "version"); err != nil {
		panic(err)
	}
	// go version go1.13 darwin/amd64
	goVersion = getGoVersionNum(strings.Split(strings.TrimSpace(string(out)), " ")[2])
	if stringutils.IsEmpty(modName) {
		modName = filepath.Base(dir)
	}
	modfile = filepath.Join(dir, "go.mod")
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

	vodir = filepath.Join(dir, "vo")
	if err = os.MkdirAll(vodir, os.ModePerm); err != nil {
		panic(err)
	}
	vofile = filepath.Join(vodir, "vo.go")
	if _, err = os.Stat(vofile); os.IsNotExist(err) {
		if f, err = os.Create(vofile); err != nil {
			panic(err)
		}
		defer f.Close()

		tpl, _ = template.New("vo.go.tmpl").Parse(voTmpl)
		_ = tpl.Execute(f, nil)
	} else {
		logrus.Warnf("file %s already exists", vofile)
	}

	svcName = strcase.ToCamel(filepath.Base(dir))
	svcfile = filepath.Join(dir, "svc.go")
	if _, err = os.Stat(svcfile); os.IsNotExist(err) {
		if f, err = os.Open(modfile); err != nil {
			panic(err)
		}
		reader := bufio.NewReader(f)
		if firstLine, err = reader.ReadString('\n'); err != nil {
			panic(err)
		}
		modName = strings.TrimSpace(strings.TrimPrefix(firstLine, "module"))

		if f, err = os.Create(svcfile); err != nil {
			panic(err)
		}
		defer f.Close()

		tpl, _ = template.New("svc.go.tmpl").Parse(svcTmpl)
		_ = tpl.Execute(f, struct {
			VoPackage string
			SvcName   string
		}{
			VoPackage: modName + "/vo",
			SvcName:   svcName,
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

// InitSvc inits a service project, test purpose only
func InitSvc(dir string) {
	var (
		err       error
		modName   string
		svcName   string
		svcfile   string
		modfile   string
		vodir     string
		vofile    string
		goVersion string
		firstLine string
		f         *os.File
		tpl       *template.Template
		envfile   string
	)
	if stringutils.IsEmpty(dir) {
		dir, _ = os.Getwd()
	}
	_ = os.MkdirAll(dir, os.ModePerm)

	gitInit(dir)
	gitIgnore(dir)

	goVersion = getGoVersionNum(runtime.Version())
	modName = filepath.Base(dir)
	modfile = filepath.Join(dir, "go.mod")
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
		logrus.Warnf("file %s already exists", vofile)
	}

	vodir = filepath.Join(dir, "vo")
	if err = os.MkdirAll(vodir, os.ModePerm); err != nil {
		panic(err)
	}
	vofile = filepath.Join(vodir, "vo.go")
	if _, err = os.Stat(vofile); os.IsNotExist(err) {
		if f, err = os.Create(vofile); err != nil {
			panic(err)
		}
		defer f.Close()

		tpl, _ = template.New("vo.go.tmpl").Parse(voTmpl)
		_ = tpl.Execute(f, nil)
	} else {
		logrus.Warnf("file %s already exists", vofile)
	}

	svcName = strcase.ToCamel(filepath.Base(dir))
	svcfile = filepath.Join(dir, "svc.go")
	if _, err = os.Stat(svcfile); os.IsNotExist(err) {
		if f, err = os.Open(modfile); err != nil {
			panic(err)
		}
		reader := bufio.NewReader(f)
		if firstLine, err = reader.ReadString('\n'); err != nil {
			panic(err)
		}
		modName = strings.TrimSpace(strings.TrimPrefix(firstLine, "module"))
		fmt.Println(modName)

		if f, err = os.Create(svcfile); err != nil {
			panic(err)
		}
		defer f.Close()

		tpl, _ = template.New("svc.go.tmpl").Parse(svcTmpl)
		_ = tpl.Execute(f, struct {
			VoPackage string
			SvcName   string
		}{
			VoPackage: modName + "/vo",
			SvcName:   svcName,
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

// gitIgnore adds .gitignore file
func gitIgnore(dir string) {
	var (
		gitignorefile string
		err           error
		f             *os.File
		tpl           *template.Template
	)
	gitignorefile = filepath.Join(dir, ".gitignore")
	if _, err = os.Stat(gitignorefile); os.IsNotExist(err) {
		if f, err = os.Create(gitignorefile); err != nil {
			panic(err)
		}
		defer f.Close()

		tpl, _ = template.New(".gitignore.tmpl").Parse(gitignoreTmpl)
		_ = tpl.Execute(f, nil)
	} else {
		logrus.Warnf("file %s already exists", ".gitignore")
	}
}

// gitInit inits git repository
func gitInit(dir string) {
	fs := osfs.New(dir)
	dot, _ := fs.Chroot(".git")
	storage := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())

	_, _ = git.Init(storage, fs)
}
