package onlinedoc

import (
	"github.com/unionj-cloud/go-doudou/svc/http/model"
	"net/http"
)

type OnlineDocHandler interface {
	GetDoc(w http.ResponseWriter, r *http.Request)
	GetOpenAPI(w http.ResponseWriter, r *http.Request)
}

func Routes() []model.Route {
	handler := NewOnlineDocHandler()
	return []model.Route{
		{
			"GetDoc",
			"GET",
			"/go-doudou/doc",
			handler.GetDoc,
		},
		{
			"GetOpenAPI",
			"GET",
			"/go-doudou/openapi.json",
			handler.GetOpenAPI,
		},
	}
}
