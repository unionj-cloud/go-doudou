package httpsrv

import (
	"github.com/go-chi/chi/v5"
)

// TODO https://github.com/go-chi/chi
type ChiHttpSrv struct {
	*chi.Mux
	routes []Route
}
