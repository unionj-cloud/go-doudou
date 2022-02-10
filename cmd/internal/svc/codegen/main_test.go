package codegen

import (
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGenMain(t *testing.T) {
	dir := testDir + "main"
	InitSvc(dir)
	defer os.RemoveAll(dir)
	ic := astutils.BuildInterfaceCollector(filepath.Join(dir, "svc.go"), astutils.ExprString)
	GenMain(dir, ic)
	expect := `package main

import (
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	service "testdatamain"
    "testdatamain/config"
	"testdatamain/transport/httpsrv"
)

func main() {
	conf := config.LoadFromEnv()
    svc := service.NewTestdatamain(conf)
	handler := httpsrv.NewTestdatamainHandler(svc)
	srv := ddhttp.NewDefaultHttpSrv()
	srv.AddRoute(httpsrv.Routes(handler)...)
	srv.Run()
}
`
	file := filepath.Join(dir, "cmd", "main.go")
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
