package onlinedoc

import (
	"github.com/unionj-cloud/go-doudou/framework/http/model"
	"net/http"
)

// OnlineDocHandler define http handler interface
type OnlineDocHandler interface {
	GetDoc(w http.ResponseWriter, r *http.Request)
	GetOpenAPI(w http.ResponseWriter, r *http.Request)
}

// Routes return route slice for gorilla mux
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
