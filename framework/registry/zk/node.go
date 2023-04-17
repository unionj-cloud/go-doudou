package zk

import (
	"context"
	"fmt"
	"github.com/unionj-cloud/go-doudou/v2/framework/buildinfo"
	"github.com/unionj-cloud/go-doudou/v2/framework/internal/config"
	cons "github.com/unionj-cloud/go-doudou/v2/framework/registry/constants"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry/serversets"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry/utils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/constants"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/errorx"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var restEndpoint *serversets.Endpoint
var grpcEndpoint *serversets.Endpoint

func newServerSet() *serversets.ServerSet {
	zkServers := config.GddZkServers.LoadOrDefault(config.DefaultGddZkServers)
	if stringutils.IsEmpty(zkServers) {
		zlogger.Panic().Msg("[go-doudou] env GDD_ZK_SERVERS is not set")
	}
	zookeepers := strings.Split(zkServers, ",")
	environment := os.Getenv("GDD_ENV")
	if stringutils.IsEmpty(environment) {
		environment = "dev"
	}
	return serversets.New(serversets.Environment(environment), config.GetServiceName(), zookeepers)
}

func registerService(service string, port uint64, scheme string, userData ...map[string]interface{}) *serversets.Endpoint {
	host := utils.GetRegisterHost()
	metadata := make(map[string]interface{})
	metadata["scheme"] = scheme
	metadata["host"] = host
	metadata["port"] = strconv.Itoa(int(port))
	metadata["service"] = service
	populateMeta(metadata, strings.HasPrefix(scheme, "http"), userData...)
	pingFunction := func() error {
		c := http.Client{Timeout: time.Duration(3) * time.Second}
		_, err := c.Get(fmt.Sprintf("%s://%s:%s%s", metadata["scheme"], metadata["host"], metadata["port"], metadata["rootPath"]))
		return errorx.Handle(err)
	}
	serverSet := newServerSet()
	endpoint, err := serverSet.RegisterEndpointWithMeta(
		host,
		int(port),
		pingFunction,
		metadata)
	if err != nil {
		zlogger.Panic().Err(err).Msgf("[go-doudou] register %s to zookeeper failed", service)
	}
	return endpoint
}

func populateMeta(meta map[string]interface{}, isRest bool, userData ...map[string]interface{}) {
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
	rr := config.DefaultGddRouteRootPath
	if stringutils.IsNotEmpty(config.GddRouteRootPath.Load()) {
		rr = config.GddRouteRootPath.Load()
	}
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
	if stringutils.IsNotEmpty(rr) && isRest {
		meta["rootPath"] = rr
	}
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

// RRServiceProvider is a simple round-robin load balance implementation for IServiceProvider
type RRServiceProvider struct {
	current  uint64
	lock     sync.Mutex
	c        *clientv3.Client
	target   string
	wch      endpoints.WatchChannel
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
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

func (r *RRServiceProvider) watch() {
	defer r.wg.Done()

	allUps := make(map[string]*endpoints.Update)
	for {
		select {
		case <-r.ctx.Done():
			return
		case ups, ok := <-r.wch:
			if !ok {
				return
			}

			for _, up := range ups {
				switch up.Op {
				case endpoints.Add:
					allUps[up.Key] = up
				case endpoints.Delete:
					delete(allUps, up.Key)
				}
			}

			addrs := convertToAddress(allUps)
			r.curState.Store(state{addresses: addrs})
		}
	}
}

func convertToAddress(ups map[string]*endpoints.Update) (addrs []*address) {
	for _, up := range ups {
		weight := 1
		var rootPath string
		if metadata, ok := up.Endpoint.Metadata.(map[string]interface{}); !ok {
			zlogger.Error().Msg("[go-doudou] zookeeper endpoint metadata is not map[string]string type")
		} else {
			weight = int(metadata["weight"].(float64))
			rootPath = metadata["rootPath"].(string)
		}
		addr := &address{
			addr:     up.Endpoint.Addr,
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

// NewRRServiceProvider creates new RRServiceProvider instance
func NewRRServiceProvider(serviceName string) *RRServiceProvider {
	onceZk.Do(func() {
		InitZkCli()
	})
	r := &RRServiceProvider{
		c:      ZkCli,
		target: serviceName,
	}
	r.ctx, r.cancel = context.WithCancel(context.Background())
	em, err := endpoints.NewManager(r.c, r.target)
	if err != nil {
		zlogger.Panic().Err(err).Msg("[go-doudou] failed to create endpoint manager")
	}
	r.wch, err = em.NewWatchChannel(r.ctx)
	if err != nil {
		zlogger.Panic().Err(err).Msg("[go-doudou] failed to create watch channel")
	}
	r.wg.Add(1)
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
func NewSWRRServiceProvider(serviceName string) *SWRRServiceProvider {
	return &SWRRServiceProvider{
		RRServiceProvider: NewRRServiceProvider(serviceName),
	}
}

func NewSWRRGrpcClientConn(service string, dialOptions ...grpc.DialOption) *grpc.ClientConn {
	return NewGrpcClientConn(service, "zookeeper_weight_balancer", dialOptions...)
}

func NewRRGrpcClientConn(service string, dialOptions ...grpc.DialOption) *grpc.ClientConn {
	return NewGrpcClientConn(service, "round_robin", dialOptions...)
}

func NewGrpcClientConn(service string, lb string, dialOptions ...grpc.DialOption) *grpc.ClientConn {
	onceZk.Do(func() {
		InitZkCli()
	})
	zookeeperResolver, err := resolver.NewBuilder(ZkCli)
	if err != nil {
		zlogger.Panic().Err(err).Msg("[go-doudou] failed to create zookeeper resolver")
	}
	dialOptions = append(dialOptions,
		grpc.WithBlock(),
		grpc.WithResolvers(zookeeperResolver),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "`+lb+`"}`),
	)
	serverAddr := fmt.Sprintf("zookeeper:///%s", service)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	grpcConn, err := grpc.DialContext(ctx, serverAddr, dialOptions...)
	if err != nil {
		zlogger.Panic().Err(err).Msgf("[go-doudou] failed to connect to server %s", serverAddr)
	}
	return grpcConn
}
