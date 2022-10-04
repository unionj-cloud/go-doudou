package registry

import (
	"bytes"
	"fmt"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/hako/durafmt"
	"github.com/hashicorp/go-msgpack/codec"
	"github.com/hashicorp/logutils"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/framework/buildinfo"
	"github.com/unionj-cloud/go-doudou/framework/configmgr"
	"github.com/unionj-cloud/go-doudou/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/framework/memberlist"
	"github.com/unionj-cloud/go-doudou/framework/registry/nacos"
	"github.com/unionj-cloud/go-doudou/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/toolkit/constants"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	logger "github.com/unionj-cloud/go-doudou/toolkit/zlogger"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var mlist memberlist.IMemberlist
var mconf *memberlist.Config
var BroadcastQueue *memberlist.TransmitLimitedQueue
var events = &eventDelegate{}

type mergedMeta struct {
	Meta nodeMeta               `json:"_meta,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
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
	if mlist == nil {
		return errors.New("mlist is nil")
	}
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
	if mlist == nil {
		return nil, errors.New("mlist is nil")
	}
	var nodes []*memberlist.Node
	for _, node := range mlist.Members() {
		nodes = append(nodes, node)
	}
	return nodes, nil
}

type nodeMeta struct {
	Service       string     `json:"service"`
	RouteRootPath string     `json:"routeRootPath"`
	Port          int        `json:"port"`
	RegisterAt    *time.Time `json:"registerAt"`
	GoVer         string     `json:"goVer"`
	GddVer        string     `json:"gddVer"`
	BuildUser     string     `json:"buildUser"`
	BuildTime     string     `json:"buildTime"`
	Weight        int        `json:"weight"`
}

func newMeta(node *memberlist.Node) (mergedMeta, error) {
	var mm mergedMeta
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
	if mlist == nil {
		return 0
	}
	return mlist.NumMembers()
}

func retransmitMultGetter() int {
	return mconf.RetransmitMult
}

func newNode(data ...map[string]interface{}) error {
	mconf = newConf()
	service := config.DefaultGddServiceName
	if stringutils.IsNotEmpty(config.GddServiceName.Load()) {
		service = config.GddServiceName.Load()
	}
	if stringutils.IsEmpty(service) {
		return errors.New(fmt.Sprintf("NewNode() error: No env variable %s found", string(config.GddServiceName)))
	}
	httpPort := config.DefaultGddPort
	if stringutils.IsNotEmpty(config.GddPort.Load()) {
		if port, err := cast.ToIntE(config.GddPort.Load()); err == nil {
			httpPort = port
		}
	}
	now := time.Now()
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
	mmeta := mergedMeta{
		Meta: nodeMeta{
			Service:       service,
			RouteRootPath: rr,
			Port:          httpPort,
			RegisterAt:    &now,
			GoVer:         runtime.Version(),
			GddVer:        buildinfo.GddVer,
			BuildUser:     buildinfo.BuildUser,
			BuildTime:     buildTime,
			Weight:        weight,
		},
		Data: make(map[string]interface{}),
	}
	if len(data) > 0 {
		mmeta.Data = data[0]
	}
	queue := &memberlist.TransmitLimitedQueue{
		NumNodes:             numNodes,
		RetransmitMultGetter: retransmitMultGetter,
	}
	BroadcastQueue = queue
	mconf.Delegate = &delegate{
		mmeta: mmeta,
		queue: queue,
	}
	mconf.Events = events
	var err error
	if mlist, err = createMemberlist(mconf); err != nil {
		return errors.Wrap(err, "[go-doudou] Failed to create memberlist")
	}
	if err = join(); err != nil {
		mlist.Shutdown()
		return errors.Wrap(err, "[go-doudou] Node register failed")
	}
	local := mlist.LocalNode()
	baseUrl, _ := BaseUrl(local)
	logger.Info().Msgf("memberlist created. local node is Node %s, providing %s service at %s, memberlist port %s",
		local.Name, mmeta.Meta.Service, baseUrl, fmt.Sprint(local.Port))
	registerConfigListener(mconf)
	return nil
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

func getModemap() map[string]struct{} {
	modeStr := config.DefaultGddServiceDiscoveryMode
	if stringutils.IsNotEmpty(config.GddServiceDiscoveryMode.Load()) {
		modeStr = config.GddServiceDiscoveryMode.Load()
	}
	modes := strings.Split(modeStr, ",")
	modemap := make(map[string]struct{})
	for _, mode := range modes {
		modemap[mode] = struct{}{}
	}
	return modemap
}

// NewNode creates a new go-doudou node.
// service related custom data (<= 512 bytes after being marshalled as json format) can be passed into it by data parameter.
// it is made as a variadic function only for backward compatibility purposes,
// only first parameter will be used.
func NewNode(data ...map[string]interface{}) error {
	for mode, _ := range getModemap() {
		switch mode {
		case "nacos":
			if err := nacos.NewNode(data...); err != nil {
				return err
			}
		case "memberlist":
			if err := newNode(data...); err != nil {
				return err
			}
		default:
			logger.Warn().Msgf("[go-doudou] unknown service discovery mode: %s", mode)
		}
	}
	return nil
}

func shutdown() {
	if mlist != nil {
		_ = mlist.Shutdown()
		mlist = nil
		logger.Info().Msg("memberlist shutdown")
	}
}

// Shutdown stops all connections and communications with other nodes in the cluster
func Shutdown() {
	for mode, _ := range getModemap() {
		switch mode {
		case "nacos":
			nacos.Shutdown()
		case "memberlist":
			shutdown()
		default:
			logger.Warn().Msgf("[go-doudou] unknown service discovery mode: %s", mode)
		}
	}
}

// Leave leaves the cluster on purpose
func Leave(timeout time.Duration) {
	if mlist != nil {
		_ = mlist.Leave(timeout)
		logger.Info().Msg("local node left the cluster")
	}
}

// NodeInfo wraps node information
type NodeInfo struct {
	SvcName   string                 `json:"svcName"`
	Hostname  string                 `json:"hostname"`
	BaseUrl   string                 `json:"baseUrl"`
	Status    string                 `json:"status"`
	Uptime    string                 `json:"uptime"`
	GoVer     string                 `json:"goVer"`
	GddVer    string                 `json:"gddVer"`
	BuildUser string                 `json:"buildUser"`
	BuildTime string                 `json:"buildTime"`
	Data      map[string]interface{} `json:"data"`
	Host      string                 `json:"host"`
	SvcPort   int                    `json:"svcPort"`
	MemPort   int                    `json:"memPort"`
}

// Info return node info
func Info(node *memberlist.Node) NodeInfo {
	status := "up"
	if node.State == memberlist.StateSuspect {
		status = "suspect"
	}
	meta, _ := newMeta(node)
	var uptime string
	if meta.Meta.RegisterAt != nil {
		uptime = time.Since(*meta.Meta.RegisterAt).String()
		if duration, err := durafmt.ParseString(uptime); err == nil {
			uptime = duration.LimitFirstN(2).String()
		}
	}
	baseUrl, _ := BaseUrl(node)
	return NodeInfo{
		SvcName:   meta.Meta.Service,
		Hostname:  node.Name,
		BaseUrl:   baseUrl,
		Status:    status,
		Uptime:    uptime,
		GoVer:     meta.Meta.GoVer,
		GddVer:    meta.Meta.GddVer,
		BuildUser: meta.Meta.BuildUser,
		BuildTime: meta.Meta.BuildTime,
		Data:      meta.Data,
		Host:      node.Addr,
		SvcPort:   meta.Meta.Port,
		MemPort:   int(node.Port),
	}
}

func BaseUrl(node *memberlist.Node) (string, error) {
	var (
		mm  mergedMeta
		err error
	)
	if mm, err = newMeta(node); err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s:%d%s", node.Addr, mm.Meta.Port, mm.Meta.RouteRootPath), nil
}

func MetaWeight(node *memberlist.Node) (int, error) {
	var (
		mm  mergedMeta
		err error
	)
	if mm, err = newMeta(node); err != nil {
		return 0, err
	}
	return mm.Meta.Weight, nil
}

func SvcName(node *memberlist.Node) string {
	var (
		mm  mergedMeta
		err error
	)
	if mm, err = newMeta(node); err != nil {
		logger.Error().Err(err).Msg("")
		return ""
	}
	return mm.Meta.Service
}

func RegisterServiceProvider(sp IMemberlistServiceProvider) {
	if mlist == nil {
		return
	}
	for _, node := range mlist.Members() {
		sp.AddNode(node)
	}
	events.ServiceProviders = append(events.ServiceProviders, sp)
}

func LocalNode() *memberlist.Node {
	if mlist == nil {
		return nil
	}
	return mlist.LocalNode()
}
