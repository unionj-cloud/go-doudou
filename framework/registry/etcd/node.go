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
	logger "github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var onceEtcd sync.Once
var EtcdCli *clientv3.Client
var restLease clientv3.LeaseID
var grpcLease clientv3.LeaseID

func InitEtcdCli() {
	etcdEndpoints := config.GddEtcdEndpoints.LoadOrDefault(config.DefaultGddEtcdEndpoints)
	if stringutils.IsEmpty(etcdEndpoints) {
		logger.Panic().Msg("env GDD_ETCD_ENDPOINTS is not set")
	}
	endpoints := strings.Split(etcdEndpoints, ",")
	var err error
	if EtcdCli, err = clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	}); err != nil {
		logger.Panic().Err(err).Msg("register to etcd failed")
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
			logger.Error().Err(err).Msgf("cast %s to int failed", leaseStr)
		} else {
			lease = value
		}
	}
	leaseResp, err := EtcdCli.Grant(tctx, lease)
	if err != nil {
		logger.Panic().Err(err).Msgf("get etcd lease ID failed")
	}
	return leaseResp.ID
}

func registerService(service string, port uint64, lease clientv3.LeaseID, userData ...map[string]interface{}) {
	em, err := endpoints.NewManager(EtcdCli, service)
	if err != nil {
		logger.Panic().Err(err).Msgf("register %s to etcd failed", service)
	}
	host := utils.GetRegisterHost()
	addr := host + ":" + strconv.Itoa(int(port))
	metadata := make(map[string]string)
	populateMeta(metadata, userData...)
	tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = em.AddEndpoint(tctx, service+"/"+addr, endpoints.Endpoint{Addr: addr, Metadata: metadata}, clientv3.WithLease(lease)); err != nil {
		logger.Panic().Err(err).Msgf("register %s to etcd failed", service)
	}
	// set keep-alive logic
	leaseRespChan, err := EtcdCli.KeepAlive(context.Background(), lease)
	if err != nil {
		logger.Panic().Err(err).Msgf("register %s to etcd failed", service)
	}
	go func() {
		for leaseKeepResp := range leaseRespChan {
			logger.Debug().Msgf("%#v", *leaseKeepResp)
		}
	}()
}

func populateMeta(meta map[string]string, userData ...map[string]interface{}) {
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
	meta["weight"] = strconv.Itoa(weight)
	if stringutils.IsNotEmpty(buildinfo.GddVer) {
		meta["gddVer"] = buildinfo.GddVer
	}
	if stringutils.IsNotEmpty(buildinfo.BuildUser) {
		meta["buildUser"] = buildinfo.BuildUser
	}
	if stringutils.IsNotEmpty(buildTime) {
		meta["buildTime"] = buildTime
	}
	if stringutils.IsNotEmpty(rr) {
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
	service := config.GetServiceName()
	port := config.GetPort()
	restLease = getLeaseID()
	registerService(service, port, restLease, data...)
	logger.Info().Msgf("[go-doudou] %s registered to etcd successfully", service)
}

func NewGrpc(data ...map[string]interface{}) {
	onceEtcd.Do(func() {
		InitEtcdCli()
	})
	service := config.GetServiceName() + "_grpc"
	port := config.GetGrpcPort()
	grpcLease = getLeaseID()
	registerService(service, port, grpcLease, data...)
	logger.Info().Msgf("[go-doudou] %s registered to etcd successfully", service)
}

func ShutdownRest() {
	if EtcdCli != nil {
		service := config.GetServiceName()
		em, err := endpoints.NewManager(EtcdCli, service)
		if err != nil {
			logger.Error().Err(err).Msgf("[go-doudou] failed to deregister %s from etcd", service)
			return
		}
		addr := utils.GetRegisterHost() + ":" + strconv.Itoa(int(config.GetPort()))
		tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err = em.DeleteEndpoint(tctx, service+"/"+addr, clientv3.WithLease(restLease)); err != nil {
			logger.Error().Err(err).Msgf("[go-doudou] failed to deregister %s from etcd", service)
			return
		}
		logger.Info().Msgf("[go-doudou] deregistered %s from etcd successfully", service)
	}
}

func ShutdownGrpc() {
	if EtcdCli != nil {
		service := config.GetServiceName() + "_grpc"
		em, err := endpoints.NewManager(EtcdCli, service)
		if err != nil {
			logger.Error().Err(err).Msgf("[go-doudou] failed to deregister %s from etcd", service)
			return
		}
		addr := utils.GetRegisterHost() + ":" + strconv.Itoa(int(config.GetGrpcPort()))
		tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err = em.DeleteEndpoint(tctx, service+"/"+addr, clientv3.WithLease(grpcLease)); err != nil {
			logger.Error().Err(err).Msgf("[go-doudou] failed to deregister %s from etcd", service)
			return
		}
		logger.Info().Msgf("[go-doudou] deregistered %s from etcd successfully", service)
	}
}

var shutdownOnce sync.Once

func CloseEtcdClient() {
	shutdownOnce.Do(func() {
		if EtcdCli != nil {
			EtcdCli.Close()
			EtcdCli = nil
			logger.Info().Msg("[go-doudou] etcd client closed")
		}
	})
}
