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
	ic := astutils.BuildInterfaceCollector(dir + "/svc.go", astutils.ExprString)
	GenSvcImpl(dir, ic)
	expect := `package service

import (
	"context"
	"testfilessvcimpl/config"
	"testfilessvcimpl/vo"

	"github.com/jmoiron/sqlx"
)

type TestfilessvcimplImpl struct {
	conf *config.Config
}

func (receiver *TestfilessvcimplImpl) PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, msg error) {
	panic("implement me")
}

func NewTestfilessvcimpl(conf *config.Config, db *sqlx.DB) Testfilessvcimpl {
	return &TestfilessvcimplImpl{
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
