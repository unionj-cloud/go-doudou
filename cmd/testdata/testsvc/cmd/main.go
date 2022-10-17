package main

import (
	"github.com/ascarter/requestid"
	"github.com/gorilla/handlers"
	"github.com/unionj-cloud/go-doudou/v2/framework/rest"
	service "testsvc"
	"testsvc/config"
	"testsvc/transport/httpsrv"
)

func main() {
	conf := config.LoadFromEnv()
	svc := service.NewTestsvc(conf)

	handler := httpsrv.NewTestsvcHandler(svc)
	srv := rest.NewRestServer()
	srv.AddMiddleware(rest.Tracing, rest.Metrics, requestid.RequestIDHandler, handlers.CompressHandler, handlers.ProxyHeaders, rest.Logger, rest.Rest, rest.Recover)
	srv.AddRoute(httpsrv.Routes(handler)...)
	srv.Run()
}
