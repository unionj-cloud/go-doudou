package prometheus

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/unionj-cloud/go-doudou/framework/http/model"
)

// Routes return route slice for gorilla mux
func Routes() []model.Route {
	return []model.Route{
		{
			Name:        "Prometheus",
			Method:      "GET",
			Pattern:     "/go-doudou/prometheus",
			HandlerFunc: promhttp.Handler().ServeHTTP,
		},
	}
}
