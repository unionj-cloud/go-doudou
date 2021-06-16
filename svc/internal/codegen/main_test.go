package codegen

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGenMain(t *testing.T) {
	dir := testDir + "main"
	InitSvc(dir)
	defer os.RemoveAll(dir)
	ic := BuildIc(dir + "/svc.go")
	GenMain(dir, ic)
	expect := `package main

import (
	"github.com/ascarter/requestid"
	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/pathutils"
	ddhttp "github.com/unionj-cloud/go-doudou/svc/http"
	service "testfilesmain"
    "testfilesmain/config"
	"testfilesmain/db"
	"testfilesmain/transport/httpsrv"
)

func main() {
	env := config.NewDotenv(pathutils.Abs("../.env"))
	conf := env.Get()

	conn, err := db.NewDb(conf.DbConf)
	if err != nil {
		panic(err)
	}
	defer func() {
		if conn == nil {
			return
		}
		if err := conn.Close(); err == nil {
			logrus.Infoln("Database connection is closed")
		} else {
			logrus.Warnln("Failed to close database connection")
		}
	}()

    svc := service.NewTestfilesmain(conf, conn)

	handler := httpsrv.NewTestfilesmainHandler(svc)
	srv := ddhttp.NewDefaultHttpSrv()
	srv.AddMiddleware(ddhttp.Metrics, requestid.RequestIDHandler, handlers.CompressHandler, handlers.ProxyHeaders, ddhttp.Logger, ddhttp.Rest)
	srv.AddRoute(httpsrv.Routes(handler)...)
	srv.Run()
}
`
	file := dir + "/cmd/main.go"
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
