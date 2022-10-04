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

func getPort() uint64 {
	httpPort := config.DefaultGddPort
	if stringutils.IsNotEmpty(config.GddPort.Load()) {
		if port, err := cast.ToIntE(config.GddPort.Load()); err == nil {
			httpPort = port
		}
	}
	return uint64(httpPort)
}

func getServiceName() string {
	service := config.DefaultGddServiceName
	if stringutils.IsNotEmpty(config.GddServiceName.Load()) {
		service = config.GddServiceName.Load()
	}
	if stringutils.IsEmpty(service) {
		logger.Panic().Msgf("[go-doudou] no value for environment variable %s found", config.GddServiceName)
	}
	return service
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

func NewNode(data ...map[string]interface{}) error {
	onceNacos.Do(func() {
		InitialiseNacosNamingClient()
	})
	registerHost := getRegisterHost()
	httpPort := getPort()
	service := getServiceName()
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
		return errors.Errorf("[go-doudou] failed to register to nacos server: %s", err)
	}
	if success {
		logger.Info().Msg("[go-doudou] registered to nacos server successfully")
	}
	return nil
}

func Shutdown() {
	if NamingClient != nil {
		registerHost := getRegisterHost()
		httpPort := getPort()
		service := getServiceName()
		success, err := NamingClient.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          registerHost,
			Port:        httpPort,
			ServiceName: service,
			Ephemeral:   true,
		})
		NamingClient = nil
		if err != nil {
			logger.Error().Err(err).Msg("[go-doudou] failed to deregister from nacos server")
			return
		}
		if !success {
			logger.Error().Msg("[go-doudou] failed to deregister from nacos server")
			return
		}
		logger.Info().Msg("[go-doudou] deregistered from nacos server successfully")
	}
}
