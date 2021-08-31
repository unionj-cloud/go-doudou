package cmd

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestInitCmd(t *testing.T) {
	dir := testDir + "initcmd"
	// go-doudou svc init ordersvc
	_, _, err := ExecuteCommandC(rootCmd, []string{"svc", "init", dir}...)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	expect := `package service

import (
	"context"
	"testfilesinitcmd/vo"
)

type Testfilesinitcmd interface {
	// You can define your service methods as your need. Below is an example.
	PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, err error)
}
`
	svcfile := dir + "/svc.go"
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
