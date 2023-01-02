package rest

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var PromRoutes = promRoutes

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
