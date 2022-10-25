package registry

import (
	"github.com/unionj-cloud/go-doudou/v2/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry/etcd"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry/nacos"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	logger "github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	"strings"
)

type IServiceProvider interface {
	SelectServer() string
}

func getModemap() map[string]struct{} {
	modeStr := config.GddServiceDiscoveryMode.LoadOrDefault(config.DefaultGddServiceDiscoveryMode)
	if stringutils.IsEmpty(modeStr) {
		return nil
	}
	modes := strings.Split(modeStr, ",")
	modemap := make(map[string]struct{})
	for _, mode := range modes {
		modemap[mode] = struct{}{}
	}
	return modemap
}

func NewRest(data ...map[string]interface{}) {
	for mode, _ := range getModemap() {
		switch mode {
		case "nacos":
			nacos.NewRest(data...)
		case "etcd":
			etcd.NewRest(data...)
		default:
			logger.Warn().Msgf("[go-doudou] unknown service discovery mode: %s", mode)
		}
	}
}

func NewGrpc(data ...map[string]interface{}) {
	for mode, _ := range getModemap() {
		switch mode {
		case "nacos":
			nacos.NewGrpc(data...)
		case "etcd":
			etcd.NewGrpc(data...)
		default:
			logger.Warn().Msgf("[go-doudou] unknown service discovery mode: %s", mode)
		}
	}
}

func ShutdownRest() {
	for mode, _ := range getModemap() {
		switch mode {
		case "nacos":
			nacos.ShutdownRest()
		case "etcd":
			etcd.ShutdownRest()
		default:
			logger.Warn().Msgf("[go-doudou] unknown service discovery mode: %s", mode)
		}
	}
}

func ShutdownGrpc() {
	for mode, _ := range getModemap() {
		switch mode {
		case "nacos":
			nacos.ShutdownGrpc()
		case "etcd":
			etcd.ShutdownGrpc()
		default:
			logger.Warn().Msgf("[go-doudou] unknown service discovery mode: %s", mode)
		}
	}
}
