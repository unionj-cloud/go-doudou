package svc

import (
	"bufio"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/sliceutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
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
	Create()
	Update()
}

type Svc struct {
	Dir string
}

func (receiver Svc) Update() {
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

}

func (receiver Svc) Create() {
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
		dir       string
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

	modName = filepath.Base(dir)
	vnums := sliceutils.StringSlice2InterfaceSlice(strings.Split(strings.TrimPrefix(runtime.Version(), "go"), "."))
	goVersion = fmt.Sprintf("%s.%s%.s", vnums...)
	fmt.Println(goVersion)
	modfile = filepath.Join(dir, "go.mod")
	if _, err = os.Stat(modfile); os.IsNotExist(err) {
		if f, err = os.Create(modfile); err != nil {
			panic(err)
		}
		defer f.Close()

		if tpl, err = template.New("go.mod.tmpl").Parse(modTempl); err != nil {
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

var modTempl = `module {{.ModName}}

go {{.GoVersion}}

require (
	github.com/unionj-cloud/go-doudou v0.1.3
)`
