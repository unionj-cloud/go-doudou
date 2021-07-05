package main

import (
	"fmt"
	"github.com/ascarter/requestid"
	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/pathutils"
	ddconfig "github.com/unionj-cloud/go-doudou/svc/config"
	ddhttp "github.com/unionj-cloud/go-doudou/svc/http"
	"github.com/unionj-cloud/go-doudou/svc/registry"
	service "testfiles3"
    "testfiles3/config"
	"testfiles3/db"
	"testfiles3/transport/httpsrv"
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

	if ddconfig.GddMode.Load() == "micro" {
		node, err := registry.NewNode()
		if err != nil {
			logrus.Panicln(fmt.Sprintf("%+v", err))
		}
		logrus.Infof("Memberlist created. Local node is %s\n", node)
	}

    svc := service.NewTestfiles3(conf, conn)

	handler := httpsrv.NewTestfiles3Handler(svc)
	srv := ddhttp.NewDefaultHttpSrv()
	srv.AddMiddleware(ddhttp.Metrics, requestid.RequestIDHandler, handlers.CompressHandler, handlers.ProxyHeaders, ddhttp.Logger, ddhttp.Rest)
	srv.AddRoute(httpsrv.Routes(handler)...)
	srv.Run()
}
