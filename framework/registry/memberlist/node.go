package memberlist

import (
	"bytes"
	"fmt"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/hashicorp/go-msgpack/codec"
	"github.com/hashicorp/logutils"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/v2/framework/buildinfo"
	"github.com/unionj-cloud/go-doudou/v2/framework/configmgr"
	"github.com/unionj-cloud/go-doudou/v2/framework/internal/config"
	cons "github.com/unionj-cloud/go-doudou/v2/framework/registry/constants"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/constants"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/memberlist"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	logger "github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mlist memberlist.IMemberlist
var mconf *memberlist.Config
var BroadcastQueue *memberlist.TransmitLimitedQueue
var events = &eventDelegate{}
var shutdownOnce sync.Once
var delegator *delegate

func assertMlistNotNil() {
	if mlist == nil {
		panic("create memberlist first")
	}
}

func init() {
	if _, ok := config.ServiceDiscoveryMap()[cons.SD_MEMBERLIST]; !ok {
		return
	}
	mconf = newConf()
	queue := &memberlist.TransmitLimitedQueue{
		NumNodes:             numNodes,
		RetransmitMultGetter: retransmitMultGetter,
	}
	now := time.Now()
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
	} else if stringutils.IsNotEmpty(config.GddMemWeight.Load()) {
		if w, err := cast.ToIntE(config.GddMemWeight.Load()); err == nil {
			weight = w
		}
	}
	BroadcastQueue = queue
	delegator = &delegate{
		meta: NodeMeta{
			RegisterAt: &now,
			GoVer:      runtime.Version(),
			GddVer:     buildinfo.GddVer,
			BuildUser:  buildinfo.BuildUser,
			BuildTime:  buildTime,
			Weight:     weight,
		},
		queue: queue,
	}
	mconf.Delegate = delegator
	mconf.Events = events
	var err error
	if mlist, err = createMemberlist(mconf); err != nil {
		panic(errors.Wrap(err, "[go-doudou] Failed to create memberlist"))
	}
	if err = join(); err != nil {
		mlist.Shutdown()
		panic(errors.Wrap(err, "[go-doudou] Node register failed"))
	}
	local := mlist.LocalNode()
	logger.Info().Msgf("memberlist created. local node is Node %s, memberlist port %s", local.Name, fmt.Sprint(local.Port))
	registerConfigListener(mconf)
}

func seeds(seedstr string) []string {
	if stringutils.IsEmpty(seedstr) {
		return nil
	}
	s := strings.Split(seedstr, ",")
	for i, seed := range s {
		li := strings.LastIndex(seed, ":")
		if li < 0 {
			s[i] = fmt.Sprintf("%s:%d", seed, config.DefaultGddMemPort)
			continue
		}
		if len(seed) > li+1 {
			if port, err := cast.ToIntE(seed[li+1:]); err != nil {
				s[i] = fmt.Sprintf("%s:%d", seed[:li], config.DefaultGddMemPort)
			} else {
				s[i] = fmt.Sprintf("%s:%d", seed[:li], port)
			}
		}
	}
	return s
}

func join() error {
	assertMlistNotNil()
	seed := config.DefaultGddMemSeed
	if stringutils.IsNotEmpty(config.GddMemSeed.Load()) {
		seed = config.GddMemSeed.Load()
	}
	s := seeds(seed)
	if len(s) == 0 {
		logger.Warn().Msg("No seed found")
		return nil
	}
	_, err := mlist.Join(s)
	if err != nil {
		return errors.Wrap(err, "[go-doudou] Failed to join cluster")
	}
	logger.Info().Msgf("Node %s joined cluster successfully", mlist.LocalNode().FullAddress())
	return nil
}

// AllNodes return all memberlist nodes except dead and left nodes
func AllNodes() ([]*memberlist.Node, error) {
	assertMlistNotNil()
	var nodes []*memberlist.Node
	for _, node := range mlist.Members() {
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func ParseMeta(node *memberlist.Node) (NodeMeta, error) {
	var mm NodeMeta
	if len(node.Meta) > 0 {
		r := bytes.NewReader(node.Meta)
		dec := codec.NewDecoder(r, &codec.MsgpackHandle{})
		if err := dec.Decode(&mm); err != nil {
			logger.Panic().Err(errors.Wrap(err, "[go-doudou] parse node meta data error")).Msg("")
		}
	}
	return mm, nil
}

func newConf() *memberlist.Config {
	cfg := memberlist.DefaultWANConfig()
	cidrs := config.GddMemCIDRsAllowed.LoadOrDefault(config.DefaultGddMemCIDRsAllowed)
	if stringutils.IsNotEmpty(cidrs) {
		var err error
		if cfg.CIDRsAllowed, err = memberlist.ParseCIDRs(strings.Split(cidrs, ",")); err != nil {
			logger.Error().Msgf("call ParseCIDRs error: %s\n", err.Error())
		}
	}
	setGddMemIndirectChecks(cfg)
	minLevel := config.DefaultGddLogLevel
	if stringutils.IsNotEmpty(config.GddLogLevel.Load()) {
		minLevel = strings.ToUpper(config.GddLogLevel.Load())
		if minLevel == "ERROR" {
			minLevel = "ERR"
		} else if minLevel == "WARNING" {
			minLevel = "WARN"
		}
	}
	lf := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERR", "INFO"},
		MinLevel: logutils.LogLevel(minLevel),
	}
	disable := cast.ToBoolOrDefault(config.GddMemLogDisable.Load(), config.DefaultGddMemLogDisable)
	if disable {
		lf.Writer = ioutil.Discard
	} else {
		lf.Writer = logger.Logger
	}
	cfg.LogOutput = lf
	setGddMemDeadTimeout(cfg)
	setGddMemSyncInterval(cfg)
	setGddMemReclaimTimeout(cfg)
	setGddMemProbeInterval(cfg)
	setGddMemProbeTimeout(cfg)
	setGddMemSuspicionMult(cfg)
	setGddMemRetransmitMult(cfg)
	setGddMemGossipNodes(cfg)
	setGddMemGossipInterval(cfg)
	// if env GDD_MEM_WEIGHT is set to > 0, then disable weight calculation, client will always use the same weight
	weight := config.DefaultGddWeight
	if stringutils.IsNotEmpty(config.GddWeight.Load()) {
		if w, err := cast.ToIntE(config.GddWeight.Load()); err == nil {
			weight = w
		}
	} else if stringutils.IsNotEmpty(config.GddMemWeight.Load()) {
		if w, err := cast.ToIntE(config.GddMemWeight.Load()); err == nil {
			weight = w
		}
	}
	if weight > 0 {
		cfg.WeightInterval = 0
	} else {
		cfg.WeightInterval = config.DefaultGddMemWeightInterval
		weightIntervalStr := config.GddMemWeightInterval.Load()
		if stringutils.IsNotEmpty(weightIntervalStr) {
			if weightInterval, err := strconv.Atoi(weightIntervalStr); err == nil {
				cfg.WeightInterval = time.Duration(weightInterval) * time.Millisecond
			} else {
				if duration, err := time.ParseDuration(weightIntervalStr); err == nil {
					cfg.WeightInterval = duration
				}
			}
		}
	}
	cfg.TCPTimeout, _ = time.ParseDuration(config.DefaultGddMemTCPTimeout)
	tcpTimeoutStr := config.GddMemTCPTimeout.Load()
	if stringutils.IsNotEmpty(tcpTimeoutStr) {
		if tcpTimeout, err := strconv.Atoi(tcpTimeoutStr); err == nil {
			cfg.TCPTimeout = time.Duration(tcpTimeout) * time.Second
		} else {
			if duration, err := time.ParseDuration(tcpTimeoutStr); err == nil {
				cfg.TCPTimeout = duration
			}
		}
	}
	if stringutils.IsNotEmpty(config.GddMemName.Load()) {
		cfg.Name = config.GddMemName.Load()
	}
	memport := config.DefaultGddMemPort
	if m, err := cast.ToIntE(config.GddMemPort.Load()); err == nil {
		memport = m
	}
	cfg.BindPort = memport
	cfg.AdvertisePort = memport
	memhost := config.GddMemHost.Load()
	if stringutils.IsNotEmpty(memhost) {
		if strings.HasPrefix(memhost, ".") {
			hostname, _ := os.Hostname()
			cfg.AdvertiseAddr = hostname + memhost
		} else {
			cfg.AdvertiseAddr = memhost
		}
	}
	return cfg
}

var createMemberlist = memberlist.Create

func numNodes() int {
	assertMlistNotNil()
	return mlist.NumMembers()
}

func retransmitMultGetter() int {
	return mconf.RetransmitMult
}

func NewRest(data ...map[string]interface{}) {
	assertMlistNotNil()
	service := config.GetServiceName() + "_" + string(cons.REST_TYPE)
	httpPort := config.GetPort()
	rr := config.DefaultGddRouteRootPath
	if stringutils.IsNotEmpty(config.GddRouteRootPath.Load()) {
		rr = config.GddRouteRootPath.Load()
	}
	si := Service{
		Name:          service,
		Host:          mlist.AdvertiseAddr(),
		Port:          int(httpPort),
		RouteRootPath: rr,
		Type:          cons.REST_TYPE,
	}
	if len(data) > 0 {
		si.Data = data[0]
	}
	delegator.AddService(si)
	if err := mlist.UpdateNode(mlist.Config().TCPTimeout); err != nil {
		panic(errors.Wrapf(err, "[go-doudou] failed to register %s service to memberlist", service))
	}
	logger.Info().Msgf("[go-doudou] registered %s service to memberlist successfully", service)
}

func NewGrpc(data ...map[string]interface{}) {
	assertMlistNotNil()
	service := config.GetServiceName() + "_" + string(cons.GRPC_TYPE)
	grpcPort := config.GetGrpcPort()
	si := Service{
		Name: service,
		Host: mlist.AdvertiseAddr(),
		Port: int(grpcPort),
		Type: cons.GRPC_TYPE,
	}
	if len(data) > 0 {
		si.Data = data[0]
	}
	delegator.AddService(si)
	if err := mlist.UpdateNode(mlist.Config().TCPTimeout); err != nil {
		panic(errors.Wrapf(err, "[go-doudou] failed to register %s service to memberlist", service))
	}
	logger.Info().Msgf("[go-doudou] registered %s service to memberlist successfully", service)
}

type memConfigListener struct {
	configmgr.BaseApolloListener
	memConf *memberlist.Config
}

func (c *memConfigListener) OnChange(event *storage.ChangeEvent) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	if !c.SkippedFirstEvent {
		c.SkippedFirstEvent = true
		return
	}
	for key, value := range event.Changes {
		upperKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
		if strings.HasPrefix(upperKey, "GDD_MEM_") {
			if value.ChangeType == storage.DELETED {
				_ = os.Unsetenv(upperKey)
				continue
			}
			_ = os.Setenv(upperKey, fmt.Sprint(value.NewValue))
		}
	}
	setGddMemDeadTimeout(c.memConf)
	setGddMemSyncInterval(c.memConf)
	setGddMemReclaimTimeout(c.memConf)
	setGddMemProbeInterval(c.memConf)
	setGddMemGossipInterval(c.memConf)
	setGddMemProbeTimeout(c.memConf)
	setGddMemSuspicionMult(c.memConf)
	setGddMemRetransmitMult(c.memConf)
	setGddMemGossipNodes(c.memConf)
	setGddMemIndirectChecks(c.memConf)
}

func CallbackOnChange(listener *memConfigListener) func(event *configmgr.NacosChangeEvent) {
	return func(event *configmgr.NacosChangeEvent) {
		changes := make(map[string]*storage.ConfigChange)
		for k, v := range event.Changes {
			changes[k] = &storage.ConfigChange{
				OldValue:   v.OldValue,
				NewValue:   v.NewValue,
				ChangeType: storage.ConfigChangeType(v.ChangeType),
			}
		}
		changeEvent := &storage.ChangeEvent{
			Changes: changes,
		}
		listener.OnChange(changeEvent)
	}
}

func registerConfigListener(memConf *memberlist.Config) {
	listener := &memConfigListener{
		memConf: memConf,
	}
	configType := config.GddConfigRemoteType.LoadOrDefault(config.DefaultGddConfigRemoteType)
	switch configType {
	case "":
		return
	case config.NacosConfigType:
		dataIdStr := config.GddNacosConfigDataid.LoadOrDefault(config.DefaultGddNacosConfigDataid)
		dataIds := strings.Split(dataIdStr, ",")
		listener.SkippedFirstEvent = true
		for _, dataId := range dataIds {
			configmgr.NacosClient.AddChangeListener(configmgr.NacosConfigListenerParam{
				DataId:   "__" + dataId + "__" + "registry",
				OnChange: CallbackOnChange(listener),
			})
		}
	case config.ApolloConfigType:
		configmgr.ApolloClient.AddChangeListener(listener)
	default:
		panic(fmt.Errorf("[go-doudou] from registry pkg: unknown config type: %s\n", configType))
	}
}

func Shutdown() {
	shutdownOnce.Do(func() {
		if mlist != nil {
			_ = mlist.Shutdown()
			mlist = nil
			logger.Info().Msg("memberlist shutdown")
		}
	})
}

// Leave leaves the cluster on purpose
func Leave(timeout time.Duration) {
	if mlist != nil {
		_ = mlist.Leave(timeout)
		logger.Info().Msg("local node left the cluster")
	}
}

func RegisterServiceProvider(sp IMemberlistServiceProvider) {
	if mlist != nil {
		for _, node := range mlist.Members() {
			sp.AddNode(node)
		}
	}
	events.ServiceProviders = append(events.ServiceProviders, sp)
}

func LocalNode() *memberlist.Node {
	assertMlistNotNil()
	return mlist.LocalNode()
}
