package main

import (
	"github.com/ascarter/requestid"
	"github.com/gorilla/handlers"
	ddhttp "github.com/unionj-cloud/go-doudou/svc/http"
	service "testsvc"
	"testsvc/config"
	"testsvc/transport/httpsrv"
)

func main() {
	conf := config.LoadFromEnv()
	svc := service.NewTestsvc(conf)

	handler := httpsrv.NewTestsvcHandler(svc)
	srv := ddhttp.NewDefaultHttpSrv()
	srv.AddMiddleware(ddhttp.Metrics, requestid.RequestIDHandler, handlers.CompressHandler, handlers.ProxyHeaders, ddhttp.Rest)
	srv.AddRoute(httpsrv.Routes(handler)...)
	srv.Run()
}
