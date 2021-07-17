package onlinedoc

import (
	"net/http"

	"github.com/unionj-cloud/go-doudou/svc/config"
	ddhttp "github.com/unionj-cloud/go-doudou/svc/http"
)

type OnlineDocHandler interface {
	GetDoc(w http.ResponseWriter, r *http.Request)
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
	}
}
