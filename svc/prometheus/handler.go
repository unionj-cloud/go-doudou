package prometheus

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/unionj-cloud/go-doudou/svc/config"
	ddhttp "github.com/unionj-cloud/go-doudou/svc/http"
)

func Routes() []ddhttp.Route {
	rootPath := config.GddRouteRootPath.Load()
	return []ddhttp.Route{
		{
			"Prometheus",
			"GET",
			rootPath + "/go-doudou/prometheus",
			promhttp.Handler().ServeHTTP,
		},
	}
}
