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
			Name:        "GetDoc",
			Method:      "GET",
			Pattern:     "/go-doudou/doc",
			HandlerFunc: handler.GetDoc,
		},
		{
			Name:        "GetOpenAPI",
			Method:      "GET",
			Pattern:     "/go-doudou/openapi.json",
			HandlerFunc: handler.GetOpenAPI,
		},
	}
}
