package prometheus

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/unionj-cloud/go-doudou/svc/http/model"
)

// Routes return route slice for gorilla mux
func Routes() []model.Route {
	return []model.Route{
		{
			"Prometheus",
			"GET",
			"/go-doudou/prometheus",
			promhttp.Handler().ServeHTTP,
		},
	}
}
