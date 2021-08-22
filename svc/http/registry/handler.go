package registry

import (
	"github.com/unionj-cloud/go-doudou/svc/http/model"
	"net/http"
)

type RegistryHandler interface {
	GetRegistry(w http.ResponseWriter, r *http.Request)
}

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
