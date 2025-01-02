package plugin

import (
	"github.com/elliotchance/orderedmap/v2"
	"github.com/unionj-cloud/go-doudou/v2/framework/grpcx"
	"github.com/unionj-cloud/go-doudou/v2/framework/rest"
	"github.com/unionj-cloud/toolkit/pipeconn"
)

var servicePlugins = orderedmap.NewOrderedMap[string, ServicePlugin]()

type ServicePlugin interface {
	Initialize(restServer *rest.RestServer, grpcServer *grpcx.GrpcServer, dialCtx pipeconn.DialContextFunc)
	GetName() string
	Close()
	GoDoudouServicePlugin()
}

func RegisterServicePlugin(plugin ServicePlugin) {
	servicePlugins.Set(plugin.GetName(), plugin)
}

func GetServicePlugins() *orderedmap.OrderedMap[string, ServicePlugin] {
	return servicePlugins
}
