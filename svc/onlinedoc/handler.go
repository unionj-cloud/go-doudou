package onlinedoc

import (
	"github.com/unionj-cloud/go-doudou/svc/config"
	ddhttp "github.com/unionj-cloud/go-doudou/svc/http"
	"net/http"
)

type OnlineDocHandler interface {
	GetDoc(w http.ResponseWriter, r *http.Request)
	GetOpenAPI(w http.ResponseWriter, r *http.Request)
}

func Routes() []ddhttp.Route {
	rootPath := config.GddRouteRootPath.Load()
	handler := NewOnlineDocHandler()
	return []ddhttp.Route{
		{
			"GetDoc",
			"GET",
			rootPath + "/go-doudou/doc",
			handler.GetDoc,
		},
		{
			"GetOpenAPI",
			"GET",
			rootPath + "/go-doudou/openapi.json",
			handler.GetOpenAPI,
		},
	}
}
