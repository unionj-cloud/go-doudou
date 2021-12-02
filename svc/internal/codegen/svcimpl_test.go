package codegen

import (
	"github.com/unionj-cloud/go-doudou/astutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGenSvcImpl(t *testing.T) {
	dir := testDir + "svcimpl"
	InitSvc(dir)
	defer os.RemoveAll(dir)
	ic := astutils.BuildInterfaceCollector(filepath.Join(dir, "svc.go"), astutils.ExprString)
	GenSvcImpl(dir, ic)
	expect := `package service

import (
	"context"
	"testdatasvcimpl/config"
	"testdatasvcimpl/vo"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jmoiron/sqlx"
)

type TestdatasvcimplImpl struct {
	conf *config.Config
}

func (receiver *TestdatasvcimplImpl) PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, err error) {
	var _result struct {
		Code int
		Data vo.PageRet
	}
	_ = gofakeit.Struct(&_result)
	return _result.Code, _result.Data, nil
}

func NewTestdatasvcimpl(conf *config.Config, db *sqlx.DB) Testdatasvcimpl {
	return &TestdatasvcimplImpl{
		conf,
	}
}
`
	file := filepath.Join(dir, "svcimpl.go")
	f, err := os.Open(file)
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != expect {
		t.Errorf("want %s, got %s\n", expect, string(content))
	}
}

func TestGenSvcImplAppend(t *testing.T) {
	dir := filepath.Join(testDir, "svcimplappend")
	ic := astutils.BuildInterfaceCollector(filepath.Join(dir, "svc.go"), astutils.ExprString)
	GenSvcImpl(dir, ic)
	file := filepath.Join(dir, "svcimpl.go")
	original := `package service

import (
	"context"
	"svcimplappend/config"
	"svcimplappend/vo"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jmoiron/sqlx"
)

type SvcimplappendImpl struct {
	conf *config.Config
}

func (receiver *SvcimplappendImpl) PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, err error) {
	var _result struct {
		Code int
		Data vo.PageRet
	}
	_ = gofakeit.Struct(&_result)
	return _result.Code, _result.Data, nil
}

func NewSvcimplappend(conf *config.Config, db *sqlx.DB) Svcimplappend {
	return &SvcimplappendImpl{
		conf,
	}
}
`
	defer func() {
		os.WriteFile(file, []byte(original), os.ModePerm)
	}()
	expect := `package service

import (
	"context"
	"svcimplappend/config"
	"svcimplappend/vo"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jmoiron/sqlx"
)

type SvcimplappendImpl struct {
	conf *config.Config
}

func (receiver *SvcimplappendImpl) PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, err error) {
	var _result struct {
		Code int
		Data vo.PageRet
	}
	_ = gofakeit.Struct(&_result)
	return _result.Code, _result.Data, nil
}

func NewSvcimplappend(conf *config.Config, db *sqlx.DB) Svcimplappend {
	return &SvcimplappendImpl{
		conf,
	}
}

func (receiver *SvcimplappendImpl) GetUser(ctx context.Context, userId string, photo string) (code int, data string, msg error) {
	var _result struct {
		Code int
		Data string
	}
	_ = gofakeit.Struct(&_result)
	return _result.Code, _result.Data, nil
}
`
	f, err := os.Open(file)
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != expect {
		t.Errorf("want %s, got %s\n", expect, string(content))
	}
}
