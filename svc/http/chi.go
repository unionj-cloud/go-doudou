package ddhttp

import (
	"github.com/go-chi/chi/v5"
	"github.com/unionj-cloud/go-doudou/svc/http/model"
)

// TODO https://github.com/go-chi/chi
type ChiHttpSrv struct {
	*chi.Mux
	routes []model.Route
}
