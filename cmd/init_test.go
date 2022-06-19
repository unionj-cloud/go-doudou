package cmd_test

import (
	"github.com/unionj-cloud/go-doudou/cmd"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestInitCmd(t *testing.T) {
	dir := testDir + "/initcmd"
	// go-doudou svc init ordersvc
	_, _, err := ExecuteCommandC(cmd.GetRootCmd(), []string{"svc", "init", dir}...)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	expect := `package service

import (
	"context"
	"initcmd/vo"
)

//go:generate go-doudou svc http --handler -c --doc

type Initcmd interface {
	// You can define your service methods as your need. Below is an example.
	// You can also add annotations here like @role(admin) to add meta data to routes for 
	// implementing your own middlewares
	PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, err error)
}
`
	svcfile := filepath.Join(dir, "svc.go")
	f, err := os.Open(svcfile)
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
