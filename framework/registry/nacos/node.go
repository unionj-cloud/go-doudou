package nacos

import (
	"context"
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/v2/framework/buildinfo"
	"github.com/unionj-cloud/go-doudou/v2/framework/config"
	"github.com/unionj-cloud/go-doudou/v2/framework/grpcx/grpc_resolver_nacos"
	cons "github.com/unionj-cloud/go-doudou/v2/framework/registry/constants"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry/utils"
	"github.com/unionj-cloud/toolkit/cast"
	"github.com/unionj-cloud/toolkit/constants"
	"github.com/unionj-cloud/toolkit/stringutils"
	logger "github.com/unionj-cloud/toolkit/zlogger"
	"github.com/wubin1989/nacos-sdk-go/v2/clients"
	"github.com/wubin1989/nacos-sdk-go/v2/clients/naming_client"
	"github.com/wubin1989/nacos-sdk-go/v2/model"
	"github.com/wubin1989/nacos-sdk-go/v2/vo"
	"google.golang.org/grpc"
)

var NamingClient naming_client.INamingClient
var onceNacos sync.Once
var NewNamingClient = clients.NewNamingClient

func InitialiseNacosNamingClient() {
	var err error
	NamingClient, err = NewNamingClient(config.GetNacosClientParam())
	if err != nil {
		logger.Panic().Err(err).Msg("[go-doudou] failed to create nacos discovery client")
	}
}

func NewRest(data ...map[string]interface{}) {
	onceNacos.Do(func() {
		InitialiseNacosNamingClient()
	})
	registerHost := utils.GetRegisterHost()
	httpPort := config.GetPort()
	service := config.GetServiceName() + "_" + string(cons.REST_TYPE)
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
	registerHost := utils.GetRegisterHost()
	grpcPort := config.GetGrpcPort()
	service := config.GetServiceName() + "_" + string(cons.GRPC_TYPE)
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

func ShutdownRest() {
	if NamingClient != nil {
		registerHost := utils.GetRegisterHost()
		httpPort := config.GetPort()
		service := config.GetServiceName() + "_" + string(cons.REST_TYPE)
		success, err := NamingClient.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          registerHost,
			Port:        httpPort,
			ServiceName: service,
			Ephemeral:   true,
		})
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
		registerHost := utils.GetRegisterHost()
		grpcPort := config.GetGrpcPort()
		service := config.GetServiceName() + "_" + string(cons.GRPC_TYPE)
		success, err := NamingClient.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          registerHost,
			Port:        grpcPort,
			ServiceName: service,
			Ephemeral:   true,
		})
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

var shutdownOnce sync.Once

func CloseNamingClient() {
	shutdownOnce.Do(func() {
		if NamingClient != nil {
			NamingClient.CloseClient()
			NamingClient = nil
			logger.Info().Msg("[go-doudou] nacos naming client closed")
		}
	})
}

type nacosBase struct {
	clusters     []string //optional,default:DEFAULT
	serviceName  string   //required
	groupName    string   //optional,default:DEFAULT_GROUP
	lock         sync.Mutex
	namingClient naming_client.INamingClient
}

func (b *nacosBase) SetClusters(clusters []string) {
	b.clusters = clusters
}

func (b *nacosBase) SetGroupName(groupName string) {
	b.groupName = groupName
}

func (b *nacosBase) SetNamingClient(namingClient naming_client.INamingClient) {
	b.namingClient = namingClient
}

type INacosServiceProvider interface {
	SetClusters(clusters []string)
	SetGroupName(groupName string)
	SetNamingClient(namingClient naming_client.INamingClient)
}

type NacosProviderOption func(INacosServiceProvider)

func WithNacosClusters(clusters []string) NacosProviderOption {
	return func(provider INacosServiceProvider) {
		provider.SetClusters(clusters)
	}
}

func WithNacosGroupName(groupName string) NacosProviderOption {
	return func(provider INacosServiceProvider) {
		provider.SetGroupName(groupName)
	}
}

func WithNacosNamingClient(namingClient naming_client.INamingClient) NacosProviderOption {
	return func(provider INacosServiceProvider) {
		provider.SetNamingClient(namingClient)
	}
}

type instance []model.Instance

func (a instance) Len() int {
	return len(a)
}

func (a instance) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a instance) Less(i, j int) bool {
	return a[i].InstanceId < a[j].InstanceId
}

// RRServiceProvider is a simple round-robin load balance implementation for IServiceProvider
type RRServiceProvider struct {
	nacosBase
	current uint64
}

// SelectServer return service address from environment variable
func (n *RRServiceProvider) SelectServer() string {
	n.lock.Lock()
	defer n.lock.Unlock()
	if n.namingClient == nil {
		logger.Error().Msg("[go-doudou] nacos discovery client has not been initialized")
		return ""
	}
	instances, err := n.namingClient.SelectInstances(vo.SelectInstancesParam{
		Clusters:    n.clusters,
		ServiceName: n.serviceName,
		GroupName:   n.groupName,
		HealthyOnly: true,
	})
	if err != nil {
		logger.Error().Err(err).Msgf("[go-doudou] %s server not found", n.serviceName)
		return ""
	}
	if len(instances) == 0 {
		logger.Error().Msgf("[go-doudou] %s server not found", n.serviceName)
		return ""
	}
	sort.Sort(instance(instances))
	next := int(atomic.AddUint64(&n.current, uint64(1)) % uint64(len(instances)))
	n.current = uint64(next)
	selected := instances[next]
	return fmt.Sprintf("http://%s:%d%s", selected.Ip, selected.Port, selected.Metadata["rootPath"])
}

func (n *RRServiceProvider) Close() {
}

// NewRRServiceProvider creates new ServiceProvider instance
func NewRRServiceProvider(serviceName string, opts ...NacosProviderOption) *RRServiceProvider {
	onceNacos.Do(func() {
		InitialiseNacosNamingClient()
	})
	provider := &RRServiceProvider{
		nacosBase: nacosBase{
			serviceName:  serviceName,
			namingClient: NamingClient,
		},
	}
	for _, opt := range opts {
		opt(provider)
	}
	return provider
}

// WRRServiceProvider is a WRR load balance implementation for IServiceProvider
type WRRServiceProvider struct {
	nacosBase
}

// SelectServer return service address from environment variable
func (n *WRRServiceProvider) SelectServer() string {
	n.lock.Lock()
	defer n.lock.Unlock()
	if n.namingClient == nil {
		logger.Error().Msg("[go-doudou] nacos discovery client has not been initialized")
		return ""
	}
	instance, err := n.namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		Clusters:    n.clusters,
		ServiceName: n.serviceName,
		GroupName:   n.groupName,
	})
	if err != nil {
		logger.Error().Err(err).Msgf("[go-doudou] %s server not found", n.serviceName)
		return ""
	}
	return fmt.Sprintf("http://%s:%d%s", instance.Ip, instance.Port, instance.Metadata["rootPath"])
}

func (n *WRRServiceProvider) Close() {
}

// NewWRRServiceProvider creates new ServiceProvider instance
func NewWRRServiceProvider(serviceName string, opts ...NacosProviderOption) *WRRServiceProvider {
	onceNacos.Do(func() {
		InitialiseNacosNamingClient()
	})
	provider := &WRRServiceProvider{
		nacosBase{
			serviceName:  serviceName,
			namingClient: NamingClient,
		},
	}
	for _, opt := range opts {
		opt(provider)
	}
	return provider
}

type NacosConfig struct {
	ServiceName string
	Clusters    []string
	GroupName   string
}

func NewWRRGrpcClientConn(config NacosConfig, dialOptions ...grpc.DialOption) *grpc.ClientConn {
	return NewGrpcClientConn(config, "nacos_weight_balancer", dialOptions...)
}

func NewRRGrpcClientConn(config NacosConfig, dialOptions ...grpc.DialOption) *grpc.ClientConn {
	return NewGrpcClientConn(config, "round_robin", dialOptions...)
}

func NewGrpcClientConn(config NacosConfig, lb string, dialOptions ...grpc.DialOption) *grpc.ClientConn {
	onceNacos.Do(func() {
		InitialiseNacosNamingClient()
	})
	grpc_resolver_nacos.AddNacosConfig(grpc_resolver_nacos.NacosConfig{
		Label:       config.ServiceName,
		ServiceName: config.ServiceName,
		Clusters:    config.Clusters,
		GroupName:   config.GroupName,
		NacosClient: NamingClient,
	})
	serverAddr := fmt.Sprintf("nacos://%s/", config.ServiceName)
	dialOptions = append(dialOptions, grpc.WithBlock(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "`+lb+`"}`))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	grpcConn, err := grpc.DialContext(ctx, serverAddr, dialOptions...)
	if err != nil {
		logger.Panic().Err(err).Msgf("[go-doudou] failed to connect to server %s", serverAddr)
	}
	return grpcConn
}
