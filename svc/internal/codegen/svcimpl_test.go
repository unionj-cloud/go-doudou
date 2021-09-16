package codegen

import (
	"github.com/unionj-cloud/go-doudou/astutils"
	"io/ioutil"
	"os"
	"testing"
)

func TestGenSvcImpl(t *testing.T) {
	dir := testDir + "svcimpl"
	InitSvc(dir)
	defer os.RemoveAll(dir)
	ic := astutils.BuildInterfaceCollector(dir+"/svc.go", astutils.ExprString)
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
	file := dir + "/svcimpl.go"
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
