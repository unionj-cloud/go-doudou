package ddhttp

import (
	"github.com/unionj-cloud/go-doudou/framework/http/model"
	"testing"
)

func TestDefaultHttpSrv_printRoutes(t *testing.T) {
	srv := NewDefaultHttpSrv()
	srv.gddRoutes = append(srv.gddRoutes, []model.Route{
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
