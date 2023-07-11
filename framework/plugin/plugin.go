package plugin

import (
	"github.com/unionj-cloud/go-doudou/v2/framework/grpcx"
	"github.com/unionj-cloud/go-doudou/v2/framework/rest"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/pipeconn"
)

var servicePlugins = map[string]ServicePlugin{}

type ServicePlugin interface {
	Initialize(restServer *rest.RestServer, grpcServer *grpcx.GrpcServer, dialCtx pipeconn.DialContextFunc)
	GetName() string
	Close()
	GoDoudouServicePlugin()
}

func RegisterServicePlugin(plugin ServicePlugin) {
	servicePlugins[plugin.GetName()] = plugin
}

func GetServicePlugins() map[string]ServicePlugin {
	return servicePlugins
}
