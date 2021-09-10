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
    github.com/gorilla/mux v1.8.0
	github.com/gorilla/handlers v1.5.1
	github.com/sirupsen/logrus v1.8.1
	github.com/go-resty/resty/v2 v2.6.0
	github.com/unionj-cloud/go-doudou v0.6.7
	github.com/olekukonko/tablewriter v0.0.5
	github.com/ascarter/requestid v0.0.0-20170313220838-5b76ab3d4aee
	github.com/common-nighthawk/go-figure v0.0.0-20200609044655-c4b36f998cf2
	github.com/unionj-cloud/cast v1.3.2
)`

const gitignoreTmpl = "# Binaries for programs and plugins\n*.exe\n*.exe~\n*.dll\n*.so\n*.dylib\n\n# Test binary, built with `go test -c`\n*.test\n\n# Output of the go coverage tool, specifically when used with LiteIDE\n*.out\n\n# Dependency directories (remove the comment below to include it)\n# vendor/"

const envTmpl = `GDD_BANNER=on
GDD_BANNER_TEXT=Go-doudou
# GddLogLevel accept values are panic, fatal, error, warn, warning, info, debug, trace
GDD_LOG_LEVEL=info

DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWD=1234
DB_SCHEMA=test
DB_CHARSET=utf8mb4
DB_DRIVER=mysql

GDD_GRACE_TIMEOUT=15s
GDD_WRITE_TIMEOUT=15s
GDD_READ_TIMEOUT=15s
GDD_IDLE_TIMEOUT=60s

# GDD_ROUTE_ROOT_PATH add prefix path to all routes
GDD_ROUTE_ROOT_PATH=

# GDD_MANAGE_ENABLE if true, it will add built-in apis with /go-doudou path prefix for online api document and service status monitor etc.
# if you don't' need the feature, just set it false or remove it
GDD_MANAGE_ENABLE=true
# GDD_MANAGE_USER if you want to disable http basic auth for management api endpoints, just set GDD_MANAGE_USER and GDD_MANAGE_PASS empty 
# or remove them
GDD_MANAGE_USER=admin
GDD_MANAGE_PASS=admin

GDD_SERVICE_NAME={{.SvcName}}
GDD_PORT=6060
# GDD_MODE accept 'mono' for monolith mode or 'micro' for microservice mode
GDD_MODE=micro

# GDD_MEM_PORT if empty or not set, an available port will be chosen randomly. recommend specifying a port
GDD_MEM_PORT=
GDD_MEM_SEED=localhost:56199
# GDD_MEM_DEAD_TIMEOUT dead node will be removed from node map if not received refute messages from it in GDD_MEM_DEAD_TIMEOUT second
GDD_MEM_DEAD_TIMEOUT=30
# GDD_MEM_SYNC_INTERVAL local node will synchronize states from other random node every GDD_MEM_SYNC_INTERVAL second
GDD_MEM_SYNC_INTERVAL=5
# GDD_MEM_RECLAIM_TIMEOUT dead node will be replaced with new node with the same name but different full address in GDD_MEM_RECLAIM_TIMEOUT second
GDD_MEM_RECLAIM_TIMEOUT=3
# GDD_MEM_NAME unique name of this node in cluster. if not provided, hostname will be used instead
GDD_MEM_NAME=
# GDD_MEM_HOST specify AdvertiseAddr attribute of memberlist config struct.
# if GDD_MEM_HOST starts with dot such as .seed-svc-headless.default.svc.cluster.local,
# it will be prefixed by hostname such as seed-2.seed-svc-headless.default.svc.cluster.local
# for supporting k8s stateful service
GDD_MEM_HOST=`

const dockerfileTmpl = `FROM golang:1.13.4-alpine AS builder

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
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags="-X 'github.com/unionj-cloud/go-doudou/svc/config.BuildUser=$HOST_USER' -X 'github.com/unionj-cloud/go-doudou/svc/config.BuildTime=$(date)' -X 'github.com/unionj-cloud/go-doudou/svc/config.GddVer=$GDD_VER'" -mod vendor -o api cmd/main.go

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
		dir, _ = os.Getwd()
	}
	_ = os.MkdirAll(dir, os.ModePerm)

	// git init
	fs := osfs.New(dir)
	dot, _ := fs.Chroot(".git")
	storage := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())

	_, _ = git.Init(storage, fs)

	// add .gitignore file
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
