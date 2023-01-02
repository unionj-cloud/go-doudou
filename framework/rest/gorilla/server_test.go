package gorilla

import (
	"github.com/unionj-cloud/go-doudou/v2/framework/rest"
	"testing"
)

func TestDefaultHttpSrv_printRoutes(t *testing.T) {
	srv := NewRestServer()
	srv.gddRoutes = append(srv.gddRoutes, []rest.Route{
		{
			Name:    "GetStatsvizWs",
			Method:  "GET",
			Pattern: gddPathPrefix + "statsviz/ws",
		},
		{
			Name:    "GetStatsviz",
			Method:  "GET",
			Pattern: gddPathPrefix + "statsviz/",
		},
	}...)
	srv.printRoutes()
}
