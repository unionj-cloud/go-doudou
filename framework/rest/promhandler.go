package rest

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func promRoutes() []Route {
	return []Route{
		{
			Name:        "Prometheus",
			Method:      "GET",
			Pattern:     "/go-doudou/prometheus",
			HandlerFunc: promhttp.Handler().ServeHTTP,
		},
	}
}
