package registry

import (
	"github.com/unionj-cloud/go-doudou/v2/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry/constants"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry/etcd"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry/memberlist"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry/nacos"
	logger "github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
)

type IServiceProvider interface {
	SelectServer() string
}

func NewRest(data ...map[string]interface{}) {
	for mode, _ := range config.ServiceDiscoveryMap() {
		switch mode {
		case constants.SD_NACOS:
			nacos.NewRest(data...)
		case constants.SD_ETCD:
			etcd.NewRest(data...)
		case constants.SD_MEMBERLIST:
			memberlist.NewRest(data...)
		default:
			logger.Warn().Msgf("[go-doudou] unknown service discovery mode: %s", mode)
		}
	}
}

func NewGrpc(data ...map[string]interface{}) {
	for mode, _ := range config.ServiceDiscoveryMap() {
		switch mode {
		case constants.SD_NACOS:
			nacos.NewGrpc(data...)
		case constants.SD_ETCD:
			etcd.NewGrpc(data...)
		case constants.SD_MEMBERLIST:
			memberlist.NewGrpc(data...)
		default:
			logger.Warn().Msgf("[go-doudou] unknown service discovery mode: %s", mode)
		}
	}
}

func ShutdownRest() {
	for mode, _ := range config.ServiceDiscoveryMap() {
		switch mode {
		case constants.SD_NACOS:
			nacos.ShutdownRest()
		case constants.SD_ETCD:
			etcd.ShutdownRest()
		case constants.SD_MEMBERLIST:
			memberlist.Shutdown()
		default:
			logger.Warn().Msgf("[go-doudou] unknown service discovery mode: %s", mode)
		}
	}
}

func ShutdownGrpc() {
	for mode, _ := range config.ServiceDiscoveryMap() {
		switch mode {
		case constants.SD_NACOS:
			nacos.ShutdownGrpc()
		case constants.SD_ETCD:
			etcd.ShutdownGrpc()
		case constants.SD_MEMBERLIST:
			memberlist.Shutdown()
		default:
			logger.Warn().Msgf("[go-doudou] unknown service discovery mode: %s", mode)
		}
	}
}
