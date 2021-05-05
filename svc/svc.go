package svc

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/sliceutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/svc/codegen"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

type SvcCmd interface {
	Init()
	Http()
}

type Svc struct {
	Dir string
}

func (receiver Svc) Http() {
	var (
		err     error
		svcfile string
		dir     string
	)
	dir = receiver.Dir
	if stringutils.IsEmpty(dir) {
		if dir, err = os.Getwd(); err != nil {
			panic(err)
		}
	}
	fmt.Println(dir)

	codegen.GenConfig(dir)
	codegen.GenDotenv(dir)
	codegen.GenDb(dir)
	codegen.GenHttpMiddleware(dir)

	svcfile = filepath.Join(dir, "svc.go")
	if _, err = os.Stat(svcfile); os.IsNotExist(err) {
		panic("Svc.go file cannot be found. Execute command go-doudou svc create first!")
	}

	var ic astutils.InterfaceCollector
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, svcfile, nil, parser.ParseComments)
	if err != nil {
		logrus.Panicln(err)
	}
	ast.Walk(&ic, root)

	fmt.Printf("%+v\n", ic)

	if len(ic.Interfaces) > 0 {
		codegen.GenHttpServer(dir, ic)
		codegen.GenMain(dir, ic)
		codegen.GenHttpHandler(dir, ic)
		codegen.GenHttpHandlerImpl(dir, ic)
		codegen.GenSvcImpl(dir, ic)
	}
}

func (receiver Svc) Init() {
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
		dir           string
		envfile       string
	)
	dir = receiver.Dir
	if stringutils.IsEmpty(dir) {
		if dir, err = os.Getwd(); err != nil {
			panic(err)
		}
	}
	fmt.Println(dir)

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
		if err = tpl.Execute(f, nil); err != nil {
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
}

var voTmpl = `package vo

import "github.com/unionj-cloud/go-doudou/ddl/query"

//go:generate go-doudou name --file $GOFILE

type Ret struct {
	Code int
	Data interface{}
	Msg  string
}

type PageFilter struct {
	// 真实姓名，前缀匹配
	Name string
	// 所属部门ID
	Dept int
}

// 分页筛选条件
type PageQuery struct {
	filter PageFilter
	page   query.Page
}

type UserVo struct {
	Id    int
	Name  string
	Phone string
	Dept  string
}
`

var svcTmpl = `package service

import (
	"context"
	"{{.VoPackage}}"
	"github.com/unionj-cloud/go-doudou/ddl/query"
)

type {{.SvcName}} interface {
	// You can define your service methods as your need. Below is an example.
	PageUsers(ctx context.Context, query vo.PageQuery) (query.PageRet, error)
}
`

var modTmpl = `module {{.ModName}}

go {{.GoVersion}}

require (
    github.com/gorilla/mux v1.8.0
	github.com/sirupsen/logrus v1.8.1
	github.com/unionj-cloud/go-doudou v0.1.8
	github.com/olekukonko/tablewriter v0.0.5
	github.com/common-nighthawk/go-figure v0.0.0-20200609044655-c4b36f998cf2
)`

var gitignoreTmpl = "# Binaries for programs and plugins\n*.exe\n*.exe~\n*.dll\n*.so\n*.dylib\n\n# Test binary, built with `go test -c`\n*.test\n\n# Output of the go coverage tool, specifically when used with LiteIDE\n*.out\n\n# Dependency directories (remove the comment below to include it)\n# vendor/"

var envTmpl = `APP_BANNER=on
APP_BANNERTEXT=Go-doudou
APP_LOGLEVEL=
APP_GRACETIMEOUT=15s

DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWD=1234
DB_SCHEMA=test
DB_CHARSET=utf8mb4
DB_DRIVER=mysql

SRV_HOST=
SRV_PORT=6060
SRV_WRITETIMEOUT=15s
SRV_READTIMEOUT=15s
SRV_IDLETIMEOUT=60s`