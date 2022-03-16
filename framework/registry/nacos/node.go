package nacos

import (
	"fmt"
	"github.com/hashicorp/go-sockaddr"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/unionj-cloud/go-doudou/framework/buildinfo"
	"github.com/unionj-cloud/go-doudou/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/framework/logger"
	"github.com/unionj-cloud/go-doudou/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/toolkit/constants"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"runtime"
	"strconv"
	"time"
)

var NamingClient naming_client.INamingClient

func getRegisterHost() string {
	registerHost := config.DefaultGddNacosRegisterHost
	if stringutils.IsNotEmpty(config.GddNacosRegisterHost.Load()) {
		registerHost = config.GddNacosRegisterHost.Load()
	}
	if stringutils.IsEmpty(registerHost) {
		var err error
		registerHost, err = sockaddr.GetPrivateIP()
		if err != nil {
			logger.Panic(fmt.Errorf("[go-doudou] failed to get interface addresses: %v", err))
		}
		if stringutils.IsEmpty(registerHost) {
			logger.Panic(fmt.Errorf("[go-doudou] no private IP address found, and explicit IP not provided"))
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
		logger.Panic(fmt.Sprintf("[go-doudou] no value for environment variable %s found", config.GddServiceName))
	}
	return service
}

func NewNode(data ...map[string]interface{}) {
	var err error
	NamingClient, err = clients.NewNamingClient(config.GetNacosClientParam())
	if err != nil {
		logger.Panic(fmt.Errorf("[go-doudou] failed to create nacos discovery client: %v", err))
	}
	registerHost := getRegisterHost()
	httpPort := getPort()
	service := getServiceName()
	weight := config.DefaultGddWeight
	if stringutils.IsNotEmpty(config.GddWeight.Load()) {
		if w, err := cast.ToIntE(config.GddWeight.Load()); err == nil {
			weight = w
		}
	}
	var buildTime string
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
		Port:        httpPort,
		ServiceName: service,
		Weight:      float64(weight),
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    metadata,
	})
	if err != nil {
		logger.Panic(fmt.Sprintf("[go-doudou] failed to register to nacos server: %s", err))
	}
	if success {
		logger.Info("[go-doudou] registered to nacos server successfully")
	}
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
			logger.Error(fmt.Sprintf("[go-doudou] failed to deregister from nacos server: %s", err))
			return
		}
		if !success {
			logger.Error("[go-doudou] failed to deregister from nacos server")
			return
		}
		logger.Info("[go-doudou] deregistered from nacos server successfully")
	}
}
