package zk

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/unionj-cloud/toolkit/cast"
	"github.com/unionj-cloud/toolkit/constants"
	"github.com/unionj-cloud/toolkit/errorx"
	"github.com/unionj-cloud/toolkit/stringutils"
	"github.com/unionj-cloud/toolkit/zlogger"
	"google.golang.org/grpc"

	"github.com/unionj-cloud/go-doudou/v2/framework/buildinfo"
	"github.com/unionj-cloud/go-doudou/v2/framework/config"
	"github.com/unionj-cloud/go-doudou/v2/framework/grpcx/grpc_resolver_zk"
	cons "github.com/unionj-cloud/go-doudou/v2/framework/registry/constants"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry/interfaces"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry/serversets"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry/utils"
)

var restEndpoint *serversets.Endpoint
var grpcEndpoint *serversets.Endpoint
var providers = map[string]interfaces.IServiceProvider{}

func newServerSet(service string) *serversets.ServerSet {
	zkServers := config.GddZkServers.LoadOrDefault(config.DefaultGddZkServers)
	if stringutils.IsEmpty(zkServers) {
		zlogger.Panic().Msg("[go-doudou] env GDD_ZK_SERVERS is not set")
	}
	zookeepers := strings.Split(zkServers, ",")
	environment := os.Getenv("GDD_ENV")
	if stringutils.IsEmpty(environment) {
		environment = "dev"
	}
	return serversets.New(serversets.Environment(environment), service, zookeepers)
}

func registerService(service string, port uint64, scheme string, userData ...map[string]interface{}) *serversets.Endpoint {
	host := utils.GetRegisterHost()
	metadata := make(map[string]interface{})
	metadata["scheme"] = scheme
	metadata["host"] = host
	metadata["port"] = strconv.Itoa(int(port))
	metadata["service"] = service
	populateMeta(metadata, userData...)
	serverSet := newServerSet(service)
	endpoint, err := serverSet.RegisterEndpointWithMeta(
		host,
		int(port),
		nil,
		metadata)
	if err != nil {
		zlogger.Panic().Err(err).Msgf("[go-doudou] register %s to zookeeper failed", service)
	}
	return endpoint
}

func populateMeta(meta map[string]interface{}, userData ...map[string]interface{}) {
	buildTime := buildinfo.BuildTime
	if stringutils.IsNotEmpty(buildinfo.BuildTime) {
		if t, err := time.Parse(constants.FORMAT15, buildinfo.BuildTime); err == nil {
			buildTime = t.Local().Format(constants.FORMAT8)
		}
	}
	weight := config.DefaultGddWeight
	if stringutils.IsNotEmpty(config.GddWeight.Load()) {
		if w, err := cast.ToIntE(config.GddWeight.Load()); err == nil {
			weight = w
		}
	}
	group := config.GddServiceGroup.LoadOrDefault(config.DefaultGddServiceGroup)
	version := config.GddServiceVersion.LoadOrDefault(config.DefaultGddServiceVersion)
	meta["group"] = group
	meta["version"] = version
	meta["registerAt"] = time.Now().Local().Format(constants.FORMAT8)
	meta["goVer"] = runtime.Version()
	meta["weight"] = weight
	if stringutils.IsNotEmpty(buildinfo.GddVer) {
		meta["gddVer"] = buildinfo.GddVer
	}
	if stringutils.IsNotEmpty(buildinfo.BuildUser) {
		meta["buildUser"] = buildinfo.BuildUser
	}
	if stringutils.IsNotEmpty(buildTime) {
		meta["buildTime"] = buildTime
	}
	meta["rootPath"] = config.GddRouteRootPath.LoadOrDefault(config.DefaultGddRouteRootPath)
	for _, item := range userData {
		for k, v := range item {
			meta[k] = fmt.Sprint(v)
		}
	}
}

func NewRest(data ...map[string]interface{}) {
	service := config.GetServiceName() + "_" + string(cons.REST_TYPE)
	httpPort := config.GetPort()
	restEndpoint = registerService(service, httpPort, "http", data...)
	zlogger.Info().Msgf("[go-doudou] %s registered to zookeeper successfully", service)
}

func NewGrpc(data ...map[string]interface{}) {
	service := config.GetServiceName() + "_" + string(cons.GRPC_TYPE)
	grpcPort := config.GetGrpcPort()
	grpcEndpoint = registerService(service, grpcPort, "grpc", data...)
	zlogger.Info().Msgf("[go-doudou] %s registered to zookeeper successfully", service)
}

func ShutdownRest() {
	if restEndpoint != nil {
		service := config.GetServiceName() + "_" + string(cons.REST_TYPE)
		restEndpoint.Close()
		zlogger.Info().Msgf("[go-doudou] deregistered %s from zookeeper successfully", service)
	}
}

func ShutdownGrpc() {
	if grpcEndpoint != nil {
		service := config.GetServiceName() + "_" + string(cons.GRPC_TYPE)
		grpcEndpoint.Close()
		zlogger.Info().Msgf("[go-doudou] deregistered %s from zookeeper successfully", service)
	}
}

// A Watcher represents how a serverset.Watch is used so it can be stubbed out for tests.
type Watcher interface {
	Endpoints() []string
	Event() <-chan struct{}
	IsClosed() bool
	Close()
}

// RRServiceProvider is a simple round-robin load balance implementation for IServiceProvider
type RRServiceProvider struct {
	current  uint64
	lock     sync.Mutex
	watcher  Watcher
	target   ServiceConfig
	curState atomic.Value
}

type address struct {
	addr          string
	rootPath      string
	weight        int
	currentWeight int
}

type state struct {
	addresses []*address
}

func (r *RRServiceProvider) updateState() {
	addrs := r.convertToAddress(r.watcher.Endpoints())
	r.curState.Store(state{addresses: addrs})
}

func (r *RRServiceProvider) watch() {
	r.updateState()
	for {
		select {
		case _, ok := <-r.watcher.Event():
			if !ok {
				return
			}
			r.updateState()
		}

		if r.watcher.IsClosed() {
			return
		}
	}
}

func (r *RRServiceProvider) convertToAddress(ups []string) (addrs []*address) {
	for _, up := range ups {
		unescaped, _ := url.QueryUnescape(up)
		u, _ := url.Parse(unescaped)
		weight := cast.ToIntOrDefault(u.Query().Get("weight"), 1)
		rootPath := u.Query().Get("rootPath")
		group := u.Query().Get("group")
		version := u.Query().Get("version")
		if group != r.target.Group || version != r.target.Version {
			continue
		}
		addr := &address{
			addr:     u.Host,
			rootPath: rootPath,
			weight:   weight,
		}
		addrs = append(addrs, addr)
	}
	return
}

// SelectServer return service address from environment variable
func (n *RRServiceProvider) SelectServer() string {
	n.lock.Lock()
	defer n.lock.Unlock()
	if n.curState.Load() == nil {
		return ""
	}
	instances := n.curState.Load().(state).addresses
	if len(instances) == 0 {
		zlogger.Error().Msgf("[go-doudou] %s server not found", n.target)
		return ""
	}
	sort.SliceStable(instances, func(i, j int) bool {
		return instances[i].addr < instances[j].addr
	})
	next := int(atomic.AddUint64(&n.current, uint64(1)) % uint64(len(instances)))
	n.current = uint64(next)
	selected := instances[next]
	return fmt.Sprintf("http://%s%s", selected.addr, selected.rootPath)
}

func (r *RRServiceProvider) Close() {
	r.watcher.Close()
}

type ServiceConfig struct {
	Name    string
	Group   string
	Version string
}

// NewRRServiceProvider creates new RRServiceProvider instance.
// If you don't need it, You should call Close to release resource.
// You can also call CloseProviders to close all at one shot
// NewRRServiceProvider is not thread-safe
func NewRRServiceProvider(conf ServiceConfig) *RRServiceProvider {
	serverSet := newServerSet(conf.Name)
	watcher, err := serverSet.Watch()
	if err != nil {
		errorx.Panic(err.Error())
	}
	r := &RRServiceProvider{
		watcher: watcher,
		target:  conf,
	}
	defer func() {
		providers[conf.Name] = r
	}()
	go r.watch()
	return r
}

// SWRRServiceProvider is a smooth weighted round-robin service provider
type SWRRServiceProvider struct {
	*RRServiceProvider
}

// SelectServer selects a node which is supplying service specified by name property from cluster
func (n *SWRRServiceProvider) SelectServer() string {
	n.lock.Lock()
	defer n.lock.Unlock()
	if n.curState.Load() == nil {
		return ""
	}
	instances := n.curState.Load().(state).addresses
	if len(instances) == 0 {
		zlogger.Error().Msgf("[go-doudou] %s server not found", n.target)
		return ""
	}
	var selected *address
	total := 0
	for i := 0; i < len(instances); i++ {
		s := instances[i]
		s.currentWeight += s.weight
		total += s.weight
		if selected == nil || s.currentWeight > selected.currentWeight {
			selected = s
		}
	}
	selected.currentWeight -= total
	return fmt.Sprintf("http://%s%s", selected.addr, selected.rootPath)
}

// NewSWRRServiceProvider creates new SWRRServiceProvider instance
func NewSWRRServiceProvider(conf ServiceConfig) *SWRRServiceProvider {
	return &SWRRServiceProvider{
		RRServiceProvider: NewRRServiceProvider(conf),
	}
}

// NewSWRRGrpcClientConn is not thread-safe
func NewSWRRGrpcClientConn(conf ServiceConfig, dialOptions ...grpc.DialOption) *grpc.ClientConn {
	return NewGrpcClientConn(conf, "zk_weight_balancer", dialOptions...)
}

// NewRRGrpcClientConn is not thread-safe
func NewRRGrpcClientConn(conf ServiceConfig, dialOptions ...grpc.DialOption) *grpc.ClientConn {
	return NewGrpcClientConn(conf, "round_robin", dialOptions...)
}

func NewGrpcClientConn(conf ServiceConfig, lb string, dialOptions ...grpc.DialOption) *grpc.ClientConn {
	serverSet := newServerSet(conf.Name)
	watcher, err := serverSet.Watch()
	if err != nil {
		errorx.Panic(err.Error())
	}
	grpc_resolver_zk.AddZkConfig(grpc_resolver_zk.ZkConfig{
		Label:       conf.Name,
		ServiceName: conf.Name,
		Watcher:     watcher,
		Group:       conf.Group,
		Version:     conf.Version,
	})
	serverAddr := fmt.Sprintf("zk://%s/", conf.Name)
	dialOptions = append(dialOptions, grpc.WithBlock(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "`+lb+`"}`))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	grpcConn, err := grpc.DialContext(ctx, serverAddr, dialOptions...)
	if err != nil {
		zlogger.Panic().Err(err).Msgf("[go-doudou] failed to connect to server %s", serverAddr)
	}
	return grpcConn
}

var shutdownOnce sync.Once

// CloseProviders you must call CloseProviders when program is shutting down, otherwise goroutine will leak
func CloseProviders() {
	shutdownOnce.Do(func() {
		for _, p := range providers {
			p.Close()
		}
	})
}
