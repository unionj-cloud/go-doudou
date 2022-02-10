package registry

import (
	"github.com/unionj-cloud/go-doudou/framework/http/model"
	"net/http"
)

// RegistryHandler define http handler interface
type RegistryHandler interface {
	GetRegistry(w http.ResponseWriter, r *http.Request)
}

// Routes return route slice for gorilla mux
func Routes() []model.Route {
	handler := NewRegistryHandler()
	return []model.Route{
		{
			"GetRegistry",
			"GET",
			"/go-doudou/registry",
			handler.GetRegistry,
		},
	}
}
