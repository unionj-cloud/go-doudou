package etcd

import (
	"context"
	"fmt"
	"github.com/unionj-cloud/go-doudou/v2/framework/buildinfo"
	"github.com/unionj-cloud/go-doudou/v2/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry/utils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/constants"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var onceEtcd sync.Once
var EtcdCli *clientv3.Client
var restLease clientv3.LeaseID
var grpcLease clientv3.LeaseID

func InitEtcdCli() {
	etcdEndpoints := config.GddEtcdEndpoints.LoadOrDefault(config.DefaultGddEtcdEndpoints)
	if stringutils.IsEmpty(etcdEndpoints) {
		zlogger.Panic().Msg("[go-doudou] env GDD_ETCD_ENDPOINTS is not set")
	}
	endpoints := strings.Split(etcdEndpoints, ",")
	var err error
	if EtcdCli, err = clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	}); err != nil {
		zlogger.Panic().Err(err).Msg("[go-doudou] register to etcd failed")
	}
}

func getLeaseID() clientv3.LeaseID {
	// grant lease time
	tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	lease := config.DefaultGddEtcdLease
	leaseStr := config.GddEtcdLease.Load()
	if stringutils.IsNotEmpty(leaseStr) {
		if value, err := cast.ToInt64E(leaseStr); err != nil {
			zlogger.Error().Err(err).Msgf("[go-doudou] cast %s to int failed", leaseStr)
		} else {
			lease = value
		}
	}
	leaseResp, err := EtcdCli.Grant(tctx, lease)
	if err != nil {
		zlogger.Panic().Err(err).Msgf("[go-doudou] get etcd lease ID failed")
	}
	return leaseResp.ID
}

func registerService(service string, port uint64, lease clientv3.LeaseID, userData ...map[string]interface{}) {
	em, err := endpoints.NewManager(EtcdCli, service)
	if err != nil {
		zlogger.Panic().Err(err).Msgf("[go-doudou] register %s to etcd failed", service)
	}
	host := utils.GetRegisterHost()
	addr := host + ":" + strconv.Itoa(int(port))
	metadata := make(map[string]interface{})
	populateMeta(metadata, strings.HasSuffix(service, "_grpc"), userData...)
	tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = em.AddEndpoint(tctx, service+"/"+addr, endpoints.Endpoint{Addr: addr, Metadata: metadata}, clientv3.WithLease(lease)); err != nil {
		zlogger.Panic().Err(err).Msgf("[go-doudou] register %s to etcd failed", service)
	}
	// set keep-alive logic
	leaseRespChan, err := EtcdCli.KeepAlive(context.Background(), lease)
	if err != nil {
		zlogger.Panic().Err(err).Msgf("[go-doudou] register %s to etcd failed", service)
	}
	go func() {
		for leaseKeepResp := range leaseRespChan {
			zlogger.Debug().Msgf("[go-doudou] %#v", *leaseKeepResp)
		}
	}()
}

func populateMeta(meta map[string]interface{}, isGrpc bool, userData ...map[string]interface{}) {
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
	if stringutils.IsNotEmpty(rr) && !isGrpc {
		meta["rootPath"] = rr
	}
	for _, item := range userData {
		for k, v := range item {
			meta[k] = fmt.Sprint(v)
		}
	}
}

func NewRest(data ...map[string]interface{}) {
	onceEtcd.Do(func() {
		InitEtcdCli()
	})
	service := config.GetServiceName() + "_rest"
	port := config.GetPort()
	restLease = getLeaseID()
	registerService(service, port, restLease, data...)
	zlogger.Info().Msgf("[go-doudou] %s registered to etcd successfully", service)
}

func NewGrpc(data ...map[string]interface{}) {
	onceEtcd.Do(func() {
		InitEtcdCli()
	})
	service := config.GetServiceName() + "_grpc"
	port := config.GetGrpcPort()
	grpcLease = getLeaseID()
	registerService(service, port, grpcLease, data...)
	zlogger.Info().Msgf("[go-doudou] %s registered to etcd successfully", service)
}

func ShutdownRest() {
	if EtcdCli != nil {
		service := config.GetServiceName() + "_rest"
		em, err := endpoints.NewManager(EtcdCli, service)
		if err != nil {
			zlogger.Error().Err(err).Msgf("[go-doudou] failed to deregister %s from etcd", service)
			return
		}
		addr := utils.GetRegisterHost() + ":" + strconv.Itoa(int(config.GetPort()))
		tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err = em.DeleteEndpoint(tctx, service+"/"+addr); err != nil {
			zlogger.Error().Err(err).Msgf("[go-doudou] failed to deregister %s from etcd", service)
			return
		}
		zlogger.Info().Msgf("[go-doudou] deregistered %s from etcd successfully", service)
	}
}

func ShutdownGrpc() {
	if EtcdCli != nil {
		service := config.GetServiceName() + "_grpc"
		em, err := endpoints.NewManager(EtcdCli, service)
		if err != nil {
			zlogger.Error().Err(err).Msgf("[go-doudou] failed to deregister %s from etcd", service)
			return
		}
		addr := utils.GetRegisterHost() + ":" + strconv.Itoa(int(config.GetGrpcPort()))
		tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err = em.DeleteEndpoint(tctx, service+"/"+addr); err != nil {
			zlogger.Error().Err(err).Msgf("[go-doudou] failed to deregister %s from etcd", service)
			return
		}
		zlogger.Info().Msgf("[go-doudou] deregistered %s from etcd successfully", service)
	}
}

var shutdownOnce sync.Once

func CloseEtcdClient() {
	shutdownOnce.Do(func() {
		if EtcdCli != nil {
			EtcdCli.Close()
			EtcdCli = nil
			zlogger.Info().Msg("[go-doudou] etcd client closed")
		}
	})
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
			zlogger.Error().Msg("[go-doudou] etcd endpoint metadata is not map[string]string type")
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
	onceEtcd.Do(func() {
		InitEtcdCli()
	})
	r := &RRServiceProvider{
		c:      EtcdCli,
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

type SmoothWeightedRoundRobinProvider struct {
	*RRServiceProvider
}

// SelectServer selects a node which is supplying service specified by name property from cluster
func (n *SmoothWeightedRoundRobinProvider) SelectServer() string {
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

// NewSWRRServiceProvider creates new SmoothWeightedRoundRobinProvider instance
func NewSWRRServiceProvider(serviceName string) *SmoothWeightedRoundRobinProvider {
	return &SmoothWeightedRoundRobinProvider{
		RRServiceProvider: NewRRServiceProvider(serviceName),
	}
}

func NewWRRGrpcClientConn(service string, dialOptions ...grpc.DialOption) *grpc.ClientConn {
	return NewGrpcClientConn(service, "etcd_weight_balancer", dialOptions...)
}

func NewRRGrpcClientConn(service string, dialOptions ...grpc.DialOption) *grpc.ClientConn {
	return NewGrpcClientConn(service, "round_robin", dialOptions...)
}

func NewGrpcClientConn(service string, lb string, dialOptions ...grpc.DialOption) *grpc.ClientConn {
	onceEtcd.Do(func() {
		InitEtcdCli()
	})
	etcdResolver, err := resolver.NewBuilder(EtcdCli)
	if err != nil {
		zlogger.Panic().Err(err).Msg("[go-doudou] failed to create etcd resolver")
	}
	dialOptions = append(dialOptions,
		grpc.WithBlock(),
		grpc.WithResolvers(etcdResolver),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "`+lb+`"}`),
	)
	serverAddr := fmt.Sprintf("etcd:///%s", service)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	grpcConn, err := grpc.DialContext(ctx, serverAddr, dialOptions...)
	if err != nil {
		zlogger.Panic().Err(err).Msgf("[go-doudou] failed to connect to server %s", serverAddr)
	}
	return grpcConn
}
