package nacos

import (
	"fmt"
	"github.com/hashicorp/go-sockaddr"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/framework/buildinfo"
	"github.com/unionj-cloud/go-doudou/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/toolkit/constants"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	logger "github.com/unionj-cloud/go-doudou/toolkit/zlogger"
	"github.com/wubin1989/nacos-sdk-go/clients"
	"github.com/wubin1989/nacos-sdk-go/clients/naming_client"
	"github.com/wubin1989/nacos-sdk-go/vo"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var NamingClient naming_client.INamingClient

var GetPrivateIP = sockaddr.GetPrivateIP

func getRegisterHost() string {
	registerHost := config.DefaultGddNacosRegisterHost
	if stringutils.IsNotEmpty(config.GddNacosRegisterHost.Load()) {
		registerHost = config.GddNacosRegisterHost.Load()
	}
	if stringutils.IsEmpty(registerHost) {
		var err error
		registerHost, err = GetPrivateIP()
		if err != nil {
			logger.Panic().Err(err).Msg("[go-doudou] failed to get interface addresses")
		}
		if stringutils.IsEmpty(registerHost) {
			logger.Panic().Msg("[go-doudou] no private IP address found, and explicit IP not provided")
		}
	}
	return registerHost
}

var onceNacos sync.Once
var NewNamingClient = clients.NewNamingClient

func InitialiseNacosNamingClient() {
	var err error
	NamingClient, err = NewNamingClient(config.GetNacosClientParam())
	if err != nil {
		logger.Panic().Err(err).Msg("[go-doudou] failed to create nacos discovery client")
	}
}

func NewNode(data ...map[string]interface{}) {
	onceNacos.Do(func() {
		InitialiseNacosNamingClient()
	})
	registerHost := getRegisterHost()
	httpPort := config.GetPort()
	service := config.GetServiceName()
	weight := config.DefaultGddWeight
	if stringutils.IsNotEmpty(config.GddWeight.Load()) {
		if w, err := cast.ToIntE(config.GddWeight.Load()); err == nil {
			weight = w
		}
	}
	buildTime := buildinfo.BuildTime
	if stringutils.IsNotEmpty(buildinfo.BuildTime) {
		if t, err := time.Parse(constants.FORMAT15, buildinfo.BuildTime); err == nil {
			buildTime = t.Local().Format(constants.FORMAT8)
		}
	}
	rr := config.DefaultGddRouteRootPath
	if stringutils.IsNotEmpty(config.GddRouteRootPath.Load()) {
		rr = config.GddRouteRootPath.Load()
	}
	metadata := make(map[string]string)
	metadata["registerAt"] = time.Now().Local().Format(constants.FORMAT8)
	metadata["goVer"] = runtime.Version()
	metadata["gddVer"] = buildinfo.GddVer
	metadata["buildUser"] = buildinfo.BuildUser
	metadata["buildTime"] = buildTime
	metadata["weight"] = strconv.Itoa(weight)
	metadata["rootPath"] = rr
	for _, item := range data {
		for k, v := range item {
			metadata[k] = fmt.Sprint(v)
		}
	}
	success, err := NamingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          registerHost,
		Port:        httpPort,
		Weight:      float64(weight),
		Enable:      true,
		Healthy:     true,
		Metadata:    metadata,
		ClusterName: config.GddNacosClusterName.LoadOrDefault(config.DefaultGddNacosClusterName),
		ServiceName: service,
		GroupName:   config.GddNacosGroupName.LoadOrDefault(config.DefaultGddNacosGroupName),
		Ephemeral:   true,
	})
	if err != nil {
		panic(errors.Errorf("[go-doudou] %s failed to register to nacos server: %s", service, err))
	}
	if success {
		logger.Info().Msgf("[go-doudou] %s registered to nacos server successfully", service)
	}
}

func NewGrpc(data ...map[string]interface{}) {
	onceNacos.Do(func() {
		InitialiseNacosNamingClient()
	})
	registerHost := getRegisterHost()
	grpcPort := config.GetGrpcPort()
	service := config.GetServiceName() + "_grpc"
	weight := config.DefaultGddWeight
	if stringutils.IsNotEmpty(config.GddWeight.Load()) {
		if w, err := cast.ToIntE(config.GddWeight.Load()); err == nil {
			weight = w
		}
	}
	buildTime := buildinfo.BuildTime
	if stringutils.IsNotEmpty(buildinfo.BuildTime) {
		if t, err := time.Parse(constants.FORMAT15, buildinfo.BuildTime); err == nil {
			buildTime = t.Local().Format(constants.FORMAT8)
		}
	}
	metadata := make(map[string]string)
	metadata["registerAt"] = time.Now().Local().Format(constants.FORMAT8)
	metadata["goVer"] = runtime.Version()
	metadata["gddVer"] = buildinfo.GddVer
	metadata["buildUser"] = buildinfo.BuildUser
	metadata["buildTime"] = buildTime
	metadata["weight"] = strconv.Itoa(weight)
	for _, item := range data {
		for k, v := range item {
			metadata[k] = fmt.Sprint(v)
		}
	}
	success, err := NamingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          registerHost,
		Port:        grpcPort,
		Weight:      float64(weight),
		Enable:      true,
		Healthy:     true,
		Metadata:    metadata,
		ClusterName: config.GddNacosClusterName.LoadOrDefault(config.DefaultGddNacosClusterName),
		ServiceName: service,
		GroupName:   config.GddNacosGroupName.LoadOrDefault(config.DefaultGddNacosGroupName),
		Ephemeral:   true,
	})
	if err != nil {
		panic(errors.Errorf("[go-doudou] %s failed to register to nacos server: %s", service, err))
	}
	if success {
		logger.Info().Msgf("[go-doudou] %s registered to nacos server successfully", service)
	}
}

func Shutdown() {
	if NamingClient != nil {
		registerHost := getRegisterHost()
		httpPort := config.GetPort()
		service := config.GetServiceName()
		success, err := NamingClient.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          registerHost,
			Port:        httpPort,
			ServiceName: service,
			Ephemeral:   true,
		})
		NamingClient = nil
		if err != nil {
			logger.Error().Err(err).Msgf("[go-doudou] failed to deregister %s from nacos server", service)
			return
		}
		if !success {
			logger.Error().Msgf("[go-doudou] failed to deregister %s from nacos server", service)
			return
		}
		logger.Info().Msgf("[go-doudou] deregistered %s from nacos server successfully", service)
	}
}

func ShutdownGrpc() {
	if NamingClient != nil {
		registerHost := getRegisterHost()
		grpcPort := config.GetGrpcPort()
		service := config.GetServiceName() + "_grpc"
		success, err := NamingClient.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          registerHost,
			Port:        grpcPort,
			ServiceName: service,
			Ephemeral:   true,
		})
		NamingClient = nil
		if err != nil {
			logger.Error().Err(err).Msgf("[go-doudou] failed to deregister %s from nacos server", service)
			return
		}
		if !success {
			logger.Error().Msgf("[go-doudou] failed to deregister %s from nacos server", service)
			return
		}
		logger.Info().Msgf("[go-doudou] deregistered %s from nacos server successfully", service)
	}
}
