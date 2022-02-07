package registry

import (
	"encoding/json"
	"fmt"
	"github.com/hako/durafmt"
	"github.com/hashicorp/logutils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/cast"
	"github.com/unionj-cloud/go-doudou/constants"
	"github.com/unionj-cloud/go-doudou/memberlist"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/svc/config"
	"github.com/unionj-cloud/go-doudou/svc/logger"
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var mlist *memberlist.Memberlist
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
		if !strings.Contains(seed, ":") {
			s[i] = seed + ":56199"
		}
	}
	return s
}

func join() error {
	if mlist == nil {
		return errors.New("mlist is nil")
	}
	s := seeds(config.GddMemSeed.Load())
	if len(s) == 0 {
		logger.Warnln("No seed found")
		return nil
	}
	_, err := mlist.Join(s)
	if err != nil {
		return errors.Wrap(err, "Failed to join cluster")
	}
	logger.Infof("Node %s joined cluster successfully", mlist.LocalNode().FullAddress())
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
		if err := json.Unmarshal(node.Meta, &mm); err != nil {
			return mm, errors.Wrap(err, "Unmarshal node meta failed, not a valid json")
		}
	}
	return mm, nil
}

// getFreePort Borrow source code from https://github.com/phayes/freeport/blob/master/freeport.go
// GetFreePort asks the kernel for a free open port that is ready to use.
func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func newConf() *memberlist.Config {
	cfg := memberlist.DefaultWANConfig()
	if indirectChecks, err := cast.ToIntE(config.GddMemIndirectChecks.Load()); err == nil {
		cfg.IndirectChecks = indirectChecks
	}
	minLevel := strings.ToUpper(config.GddLogLevel.Load())
	if minLevel == "ERROR" {
		minLevel = "ERR"
	} else if minLevel == "WARNING" {
		minLevel = "WARN"
	}
	lf := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERR", "INFO"},
		MinLevel: logutils.LogLevel(minLevel),
	}
	if config.GddMemLogDisable.Load() == "true" {
		lf.Writer = ioutil.Discard
	} else {
		lf.Writer = logrus.StandardLogger().Writer()
	}
	cfg.LogOutput = lf
	deadTimeoutStr := config.GddMemDeadTimeout.Load()
	if stringutils.IsNotEmpty(deadTimeoutStr) {
		if deadTimeout, err := strconv.Atoi(deadTimeoutStr); err == nil {
			cfg.GossipToTheDeadTime = time.Duration(deadTimeout) * time.Second
		} else {
			var duration time.Duration
			if duration, err = time.ParseDuration(deadTimeoutStr); err == nil {
				cfg.GossipToTheDeadTime = duration
			}
		}
	}
	cfg.PushPullInterval = 10 * time.Second
	syncIntervalStr := config.GddMemSyncInterval.Load()
	if stringutils.IsNotEmpty(syncIntervalStr) {
		if syncInterval, err := strconv.Atoi(syncIntervalStr); err == nil {
			cfg.PushPullInterval = time.Duration(syncInterval) * time.Second
		} else {
			var duration time.Duration
			if duration, err = time.ParseDuration(syncIntervalStr); err == nil {
				cfg.PushPullInterval = duration
			}
		}
	}
	cfg.DeadNodeReclaimTime = 3 * time.Second
	reclaimTimeoutStr := config.GddMemReclaimTimeout.Load()
	if stringutils.IsNotEmpty(reclaimTimeoutStr) {
		if reclaimTimeout, err := strconv.Atoi(reclaimTimeoutStr); err == nil {
			cfg.DeadNodeReclaimTime = time.Duration(reclaimTimeout) * time.Second
		} else {
			var duration time.Duration
			if duration, err = time.ParseDuration(reclaimTimeoutStr); err == nil {
				cfg.DeadNodeReclaimTime = duration
			}
		}
	}
	probeIntervalStr := config.GddMemProbeInterval.Load()
	if stringutils.IsNotEmpty(probeIntervalStr) {
		if probeInterval, err := strconv.Atoi(probeIntervalStr); err == nil {
			cfg.ProbeInterval = time.Duration(probeInterval) * time.Second
		} else {
			var duration time.Duration
			if duration, err = time.ParseDuration(probeIntervalStr); err == nil {
				cfg.ProbeInterval = duration
			}
		}
	}
	probeTimeoutStr := config.GddMemProbeTimeout.Load()
	if stringutils.IsNotEmpty(probeTimeoutStr) {
		if probeTimeout, err := strconv.Atoi(probeTimeoutStr); err == nil {
			cfg.ProbeTimeout = time.Duration(probeTimeout) * time.Second
		} else {
			var duration time.Duration
			if duration, err = time.ParseDuration(probeTimeoutStr); err == nil {
				cfg.ProbeTimeout = duration
			}
		}
	}
	suspicionMultStr := config.GddMemSuspicionMult.Load()
	if stringutils.IsNotEmpty(suspicionMultStr) {
		if suspicionMult, err := strconv.Atoi(suspicionMultStr); err == nil {
			cfg.SuspicionMult = suspicionMult
		}
	}
	gossipNodesStr := config.GddMemGossipNodes.Load()
	if stringutils.IsNotEmpty(gossipNodesStr) {
		if gossipNodes, err := strconv.Atoi(gossipNodesStr); err == nil {
			cfg.GossipNodes = gossipNodes
		}
	}
	gossipIntervalStr := config.GddMemGossipInterval.Load()
	if stringutils.IsNotEmpty(gossipIntervalStr) {
		if gossipInterval, err := strconv.Atoi(gossipIntervalStr); err == nil {
			cfg.GossipInterval = time.Duration(gossipInterval) * time.Millisecond
		} else {
			var duration time.Duration
			if duration, err = time.ParseDuration(gossipIntervalStr); err == nil {
				cfg.GossipInterval = duration
			}
		}
	}
	// if env GDD_MEM_WEIGHT is set to > 0, then disable weight calculation, client will always use the same weight
	if cast.ToInt(config.GddMemWeight.Load()) > 0 {
		cfg.WeightInterval = 0
	} else {
		weightIntervalStr := config.GddMemWeightInterval.Load()
		if stringutils.IsNotEmpty(weightIntervalStr) {
			if weightInterval, err := strconv.Atoi(weightIntervalStr); err == nil {
				cfg.WeightInterval = time.Duration(weightInterval) * time.Millisecond
			} else {
				var duration time.Duration
				if duration, err = time.ParseDuration(weightIntervalStr); err == nil {
					cfg.WeightInterval = duration
				}
			}
		}
	}
	tcpTimeoutStr := config.GddMemTCPTimeout.Load()
	if stringutils.IsNotEmpty(tcpTimeoutStr) {
		if tcpTimeout, err := strconv.Atoi(tcpTimeoutStr); err == nil {
			cfg.TCPTimeout = time.Duration(tcpTimeout) * time.Second
		} else {
			var duration time.Duration
			if duration, err = time.ParseDuration(tcpTimeoutStr); err == nil {
				cfg.TCPTimeout = duration
			}
		}
	}
	nodename := config.GddMemName.Load()
	if stringutils.IsNotEmpty(nodename) {
		cfg.Name = nodename
	}
	memport := cast.ToInt(config.GddMemPort.Load())
	if memport == 0 {
		memport, _ = getFreePort()
	}
	if memport > 0 {
		cfg.BindPort = memport
		cfg.AdvertisePort = memport
	}
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

// NewNode creates a new go-doudou node.
// service related custom data (<= 512 bytes after being marshalled as json format) can be passed into it by data parameter.
// it is made as a variadic function only for backward compatibility purposes,
// only first parameter will be used.
func NewNode(data ...map[string]interface{}) error {
	mconf := newConf()
	service := config.GddServiceName.Load()
	if stringutils.IsEmpty(service) {
		return errors.New(fmt.Sprintf("NewNode() error: No env variable %s found", config.GddServiceName))
	}
	port := cast.ToInt(config.GddPort.Load())
	if port == 0 {
		port, _ = getFreePort()
	}
	config.GddPort.Write(fmt.Sprint(port))
	now := time.Now()
	var buildTime string
	if stringutils.IsNotEmpty(config.BuildTime) {
		if t, err := time.Parse(constants.FORMAT15, config.BuildTime); err == nil {
			buildTime = t.Local().Format(constants.FORMAT)
		}
	}
	mmeta := mergedMeta{
		Meta: nodeMeta{
			Service:       service,
			RouteRootPath: config.GddRouteRootPath.Load(),
			Port:          port,
			RegisterAt:    &now,
			GoVer:         runtime.Version(),
			GddVer:        config.GddVer,
			BuildUser:     config.BuildUser,
			BuildTime:     buildTime,
			Weight:        cast.ToInt(config.GddMemWeight.Load()),
		},
		Data: make(map[string]interface{}),
	}
	if len(data) > 0 {
		mmeta.Data = data[0]
	}
	queue := &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			if mlist == nil {
				return 0
			}
			return mlist.NumMembers()
		},
		RetransmitMult: mconf.RetransmitMult,
	}
	BroadcastQueue = queue
	mconf.Delegate = &delegate{
		mmeta: mmeta,
		queue: queue,
	}
	mconf.Events = events
	var err error
	if mlist, err = memberlist.Create(mconf); err != nil {
		return errors.Wrap(err, "NewNode() error: Failed to create memberlist")
	}
	if err = join(); err != nil {
		mlist.Shutdown()
		return errors.Wrap(err, "NewNode() error: Node register failed")
	}
	local := mlist.LocalNode()
	baseUrl, _ := BaseUrl(local)
	logger.Infof("memberlist created. local node is Node %s, providing %s service at %s, memberlist port %s",
		local.Name, mmeta.Meta.Service, baseUrl, fmt.Sprint(local.Port))
	return nil
}

// Shutdown stops all connections and communications with other nodes in the cluster
func Shutdown() {
	if mlist != nil {
		_ = mlist.Shutdown()
		mlist = nil
		logger.Info("memberlist shutdown")
	}
}

// Leave leaves the cluster on purpose
func Leave(timeout time.Duration) {
	if mlist != nil {
		_ = mlist.Leave(timeout)
		logger.Info("local node left the cluster")
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
		logger.Errorln(fmt.Sprintf("%+v", err))
		return ""
	}
	return mm.Meta.Service
}

func RegisterServiceProvider(sp IServiceProvider) {
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
