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
	"github.com/unionj-cloud/go-doudou/sliceutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
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
	PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, msg error)
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
    github.com/gorilla/mux v1.8.0
	github.com/gorilla/handlers v1.5.1
	github.com/sirupsen/logrus v1.8.1
	github.com/go-resty/resty/v2 v2.6.0
	github.com/unionj-cloud/go-doudou v0.2.7
	github.com/olekukonko/tablewriter v0.0.5
	github.com/ascarter/requestid v0.0.0-20170313220838-5b76ab3d4aee
	github.com/common-nighthawk/go-figure v0.0.0-20200609044655-c4b36f998cf2
)`

const gitignoreTmpl = "# Binaries for programs and plugins\n*.exe\n*.exe~\n*.dll\n*.so\n*.dylib\n\n# Test binary, built with `go test -c`\n*.test\n\n# Output of the go coverage tool, specifically when used with LiteIDE\n*.out\n\n# Dependency directories (remove the comment below to include it)\n# vendor/"

const envTmpl = `GDD_BANNER=on
GDD_BANNERTEXT=Go-doudou
GDD_LOGLEVEL=
GDD_GRACETIMEOUT=15s

DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWD=1234
DB_SCHEMA=test
DB_CHARSET=utf8mb4
DB_DRIVER=mysql

GDD_WRITETIMEOUT=15s
GDD_READTIMEOUT=15s
GDD_IDLETIMEOUT=60s

GDD_NAME={{.SvcName}}
GDD_PORT=6060
GDD_MEM_PORT=
GDD_BASE_URL=
GDD_SEED=192.168.101.6:52634
# Accept 'mono' for monolith mode or 'micro' for microservice mode
GDD_MODE=micro`

const dockerfileTmpl = `FROM golang:1.13.4-alpine AS builder

ENV GO111MODULE=on

WORKDIR /repo

ADD go.mod .
ADD go.sum .

ADD . ./

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache bash tzdata

ENV TZ="Asia/Shanghai"

EXPOSE 6060

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -o api cmd/main.go

ENTRYPOINT ["/repo/api"]
`

func InitSvc(dir string) {
	var (
		err           error
		modName       string
		svcName       string
		gitignorefile string
		svcfile       string
		modfile       string
		vodir         string
		vofile        string
		goVersion     string
		firstLine     string
		f             *os.File
		tpl           *template.Template
		envfile       string
	)
	if stringutils.IsEmpty(dir) {
		if dir, err = os.Getwd(); err != nil {
			panic(err)
		}
	}

	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}

	// git init
	fs := osfs.New(dir)
	dot, _ := fs.Chroot(".git")
	storage := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())

	if _, err = git.Init(storage, fs); err != nil {
		panic("git init error")
	}

	// add .gitignore file
	gitignorefile = filepath.Join(dir, ".gitignore")
	if _, err = os.Stat(gitignorefile); os.IsNotExist(err) {
		if f, err = os.Create(gitignorefile); err != nil {
			panic(err)
		}
		defer f.Close()

		if tpl, err = template.New(".gitignore.tmpl").Parse(gitignoreTmpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(f, nil); err != nil {
			panic(err)
		}
	} else {
		logrus.Warnf("file %s already exists", ".gitignore")
	}

	vnums := sliceutils.StringSlice2InterfaceSlice(strings.Split(strings.TrimPrefix(runtime.Version(), "go"), "."))
	goVersion = fmt.Sprintf("%s.%s%.s", vnums...)
	fmt.Println(goVersion)
	modName = filepath.Base(dir)
	modfile = filepath.Join(dir, "go.mod")
	if _, err = os.Stat(modfile); os.IsNotExist(err) {
		if f, err = os.Create(modfile); err != nil {
			panic(err)
		}
		defer f.Close()

		if tpl, err = template.New("go.mod.tmpl").Parse(modTmpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(f, struct {
			ModName   string
			GoVersion string
		}{
			ModName:   modName,
			GoVersion: goVersion,
		}); err != nil {
			panic(err)
		}
	} else {
		logrus.Warnf("file %s already exists", "go.mod")
	}

	envfile = filepath.Join(dir, ".env")
	if _, err = os.Stat(envfile); os.IsNotExist(err) {
		if f, err = os.Create(envfile); err != nil {
			panic(err)
		}
		defer f.Close()

		if tpl, err = template.New(".env.tmpl").Parse(envTmpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(f, struct {
			SvcName string
		}{
			SvcName: modName,
		}); err != nil {
			panic(err)
		}
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

		if tpl, err = template.New("vo.go.tmpl").Parse(voTmpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(f, nil); err != nil {
			panic(err)
		}
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

		if tpl, err = template.New("svc.go.tmpl").Parse(svcTmpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(f, struct {
			VoPackage string
			SvcName   string
		}{
			VoPackage: modName + "/vo",
			SvcName:   svcName,
		}); err != nil {
			panic(err)
		}
	} else {
		logrus.Warnf("file %s already exists", svcfile)
	}

	dockerfile := filepath.Join(dir, "Dockerfile")
	if _, err = os.Stat(dockerfile); os.IsNotExist(err) {
		if f, err = os.Create(dockerfile); err != nil {
			panic(err)
		}
		defer f.Close()

		if tpl, err = template.New("dockerfile.tmpl").Parse(dockerfileTmpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(f, nil); err != nil {
			panic(err)
		}
	} else {
		logrus.Warnf("file %s already exists", dockerfile)
	}
}
