package registry

import (
	"github.com/unionj-cloud/go-doudou/svc/http/model"
	"net/http"
)

// ConfigHandler define http handler interface
type ConfigHandler interface {
	GetConfig(w http.ResponseWriter, r *http.Request)
}

// Routes return route slice for gorilla mux
func Routes() []model.Route {
	handler := NewConfigHandler()
	return []model.Route{
		{
			"GetConfig",
			"GET",
			"/go-doudou/config",
			handler.GetConfig,
		},
	}
}
