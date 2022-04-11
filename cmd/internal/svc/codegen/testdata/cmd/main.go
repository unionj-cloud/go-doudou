package main

import (
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	service "testdata"
    "testdata/config"
	"testdata/transport/httpsrv"
)

func main() {
	conf := config.LoadFromEnv()
    svc := service.NewUsersvc(conf)
	handler := httpsrv.NewUsersvcHandler(svc)
	srv := ddhttp.NewDefaultHttpSrv()
	srv.AddRoute(httpsrv.Routes(handler)...)
	srv.Run()
}
