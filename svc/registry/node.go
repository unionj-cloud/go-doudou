package registry

import (
	"encoding/json"
	"fmt"
	"github.com/hako/durafmt"
	"github.com/hashicorp/logutils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/cast"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/svc/config"
	"github.com/unionj-cloud/memberlist"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// IRegistry wraps service registry behaviors
type IRegistry interface {
	Register() error
	Discover(svc string) ([]*Node, error)
}

type registry struct {
	memberConf *memberlist.Config
	broadcasts *memberlist.TransmitLimitedQueue
	memberlist *memberlist.Memberlist
	lock       sync.Mutex
	memberLock sync.RWMutex
	members    []*memberlist.Node
}

func seeds(seedstr string) []string {
	if stringutils.IsEmpty(seedstr) {
		return nil
	}
	seeds := strings.Split(seedstr, ",")
	for i, seed := range seeds {
		if !strings.Contains(seed, ":") {
			seeds[i] = seed + ":56199"
		}
	}
	return seeds
}

// Register registers local node to cluster
func (r *registry) Register() error {
	if r.memberlist == nil {
		return errors.New("Memberlist is nil")
	}
	seeds := seeds(config.GddMemSeed.Load())
	if len(seeds) == 0 {
		logrus.Warnln("No seed found")
		return nil
	}
	_, err := r.memberlist.Join(seeds)
	if err != nil {
		return errors.Wrap(err, "Failed to join cluster")
	}
	logrus.Infof("Node %s joined cluster successfully", r.memberlist.LocalNode().FullAddress())
	return nil
}

// Discover finds nodes which supplying specified service
func (r *registry) Discover(svc string) ([]*Node, error) {
	if r.memberlist == nil {
		return nil, errors.New("Memberlist is nil")
	}
	var nodes []*Node
	for _, member := range r.memberlist.Members() {
		if member.State != memberlist.StateAlive {
			continue
		}
		logrus.Debugf("Member: %s %s\n", member.Name, member.Addr)
		var mmeta mergedMeta
		if err := json.Unmarshal(member.Meta, &mmeta); err != nil {
			return nil, errors.Wrap(err, "")
		}
		if stringutils.IsEmpty(svc) || mmeta.Meta.Service == svc {
			nodes = append(nodes, &Node{
				mmeta:      mmeta,
				memberNode: member,
				remote:     true,
			})
		}
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
}

func newMeta(mnode *memberlist.Node) (mergedMeta, error) {
	var mm mergedMeta
	if len(mnode.Meta) > 0 {
		if err := json.Unmarshal(mnode.Meta, &mm); err != nil {
			return mm, errors.Wrap(err, "Unmarshal node meta failed, not a valid json")
		}
	}
	return mm, nil
}

type mergedMeta struct {
	Meta nodeMeta    `json:"_meta,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

// Node represents go-doudou node
type Node struct {
	mmeta      mergedMeta
	memberNode *memberlist.Node
	*registry
	// check the node is a local node or remote node
	remote bool
}

// LocalNode store local node globally
var LocalNode *Node

// NodeOption sets node properties
type NodeOption func(*Node)

// WithData sets data that local node carrying
func WithData(data interface{}) NodeOption {
	return func(node *Node) {
		node.mmeta = mergedMeta{
			Data: data,
		}
	}
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
	mconf := memberlist.DefaultWANConfig()
	// if both udp and tcp ping failed, the node should be suspected,
	// no need to send indirect ping message for RESTFul microservice use case
	mconf.IndirectChecks = 0
	minLevel := strings.ToUpper(config.GddLogLevel.Load())
	if minLevel == "ERROR" {
		minLevel = "ERR"
	} else if minLevel == "WARNING" {
		minLevel = "WARN"
	}
	mconf.LogOutput = &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERR", "INFO"},
		MinLevel: logutils.LogLevel(minLevel),
		Writer:   logrus.StandardLogger().Writer(),
	}
	mconf.GossipToTheDeadTime = 30 * time.Second
	deadTimeoutStr := config.GddMemDeadTimeout.Load()
	if stringutils.IsNotEmpty(deadTimeoutStr) {
		if deadTimeout, err := strconv.Atoi(deadTimeoutStr); err == nil {
			mconf.GossipToTheDeadTime = time.Duration(deadTimeout) * time.Second
		}
	}
	mconf.PushPullInterval = 5 * time.Second
	syncIntervalStr := config.GddMemSyncInterval.Load()
	if stringutils.IsNotEmpty(syncIntervalStr) {
		if syncInterval, err := strconv.Atoi(syncIntervalStr); err == nil {
			mconf.PushPullInterval = time.Duration(syncInterval) * time.Second
		}
	}
	mconf.DeadNodeReclaimTime = 3 * time.Second
	reclaimTimeoutStr := config.GddMemReclaimTimeout.Load()
	if stringutils.IsNotEmpty(reclaimTimeoutStr) {
		if reclaimTimeout, err := strconv.Atoi(reclaimTimeoutStr); err == nil {
			mconf.DeadNodeReclaimTime = time.Duration(reclaimTimeout) * time.Second
		}
	}
	mconf.ProbeInterval = 1 * time.Second
	probeIntervalStr := config.GddMemProbeInterval.Load()
	if stringutils.IsNotEmpty(probeIntervalStr) {
		if probeInterval, err := strconv.Atoi(probeIntervalStr); err == nil {
			mconf.ProbeInterval = time.Duration(probeInterval) * time.Second
		}
	}
	nodename := config.GddMemName.Load()
	if stringutils.IsNotEmpty(nodename) {
		mconf.Name = nodename
	}
	memport := cast.ToInt(config.GddMemPort.Load())
	if memport == 0 {
		memport, _ = getFreePort()
	}
	if memport > 0 {
		mconf.BindPort = memport
		mconf.AdvertisePort = memport
	}
	memhost := config.GddMemHost.Load()
	if stringutils.IsNotEmpty(memhost) {
		if strings.HasPrefix(memhost, ".") {
			hostname, _ := os.Hostname()
			mconf.AdvertiseAddr = hostname + memhost
		} else {
			mconf.AdvertiseAddr = memhost
		}
	}
	return mconf
}

// NewNode creates new go-doudou node
func NewNode(opts ...NodeOption) (*Node, error) {
	mconf := newConf()
	service := config.GddServiceName.Load()
	if stringutils.IsEmpty(service) {
		return nil, errors.New(fmt.Sprintf("NewNode() error: No env variable %s found", config.GddServiceName))
	}
	node := &Node{
		registry: &registry{
			memberConf: mconf,
		},
	}
	for _, opt := range opts {
		opt(node)
	}
	port := cast.ToInt(config.GddPort.Load())
	if port == 0 {
		port, _ = getFreePort()
	}
	config.GddPort.Write(fmt.Sprint(port))
	now := time.Now()
	node.mmeta.Meta = nodeMeta{
		Service:       service,
		RouteRootPath: config.GddRouteRootPath.Load(),
		Port:          port,
		RegisterAt:    &now,
		GoVer:         runtime.Version(),
		GddVer:        config.GddVer,
		BuildUser:     config.BuildUser,
		BuildTime:     config.BuildTime,
	}
	mconf.Delegate = &delegate{node}
	mconf.Events = &eventDelegate{node}
	list, err := memberlist.Create(mconf)
	if err != nil {
		return nil, errors.Wrap(err, "NewNode() error: Failed to create memberlist")
	}
	node.registry.memberlist = list
	node.registry.broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes:       node.NumNodes,
		RetransmitMult: mconf.RetransmitMult,
	}
	if err = node.Register(); err != nil {
		node.registry.memberlist.Shutdown()
		return nil, errors.Wrap(err, "NewNode() error: Node register failed")
	}
	node.memberNode = list.LocalNode()
	LocalNode = node
	return node, nil
}

// NumNodes return node number in the cluster
func (n *Node) NumNodes() (numNodes int) {
	n.memberLock.RLock()
	numNodes = len(n.memberlist.Members())
	n.memberLock.RUnlock()

	return numNodes
}

// Shutdown stops all connections and communications with other nodes in the cluster
func (n *Node) Shutdown() {
	if err := n.memberlist.Shutdown(); err != nil {
		logrus.Errorf("memberlist shutdown fail: %+v\n", err)
	}
	return
}

// NodeInfo wraps node information
type NodeInfo struct {
	SvcName   string `json:"svcName"`
	Hostname  string `json:"hostname"`
	BaseUrl   string `json:"baseUrl"`
	Status    string `json:"status"`
	Uptime    string `json:"uptime"`
	GoVer     string `json:"goVer"`
	GddVer    string `json:"gddVer"`
	BuildUser string `json:"buildUser"`
	BuildTime string `json:"buildTime"`
	Data      string `json:"data"`
}

// Info return node info
func (n *Node) Info() NodeInfo {
	status := "up"
	if n.memberNode.State == memberlist.StateSuspect {
		status = "suspect"
	}
	var data string
	if n.mmeta.Data != nil {
		if b, err := json.Marshal(n.mmeta.Data); err == nil {
			data = string(b)
		}
	}
	var uptime string
	if n.mmeta.Meta.RegisterAt != nil {
		uptime = time.Since(*n.mmeta.Meta.RegisterAt).String()
		if duration, err := durafmt.ParseString(uptime); err == nil {
			uptime = duration.LimitFirstN(2).String()
		}
	}
	return NodeInfo{
		SvcName:   n.mmeta.Meta.Service,
		Hostname:  n.memberNode.Name,
		BaseUrl:   n.BaseUrl(),
		Status:    status,
		Uptime:    uptime,
		GoVer:     n.mmeta.Meta.GoVer,
		GddVer:    n.mmeta.Meta.GddVer,
		BuildUser: n.mmeta.Meta.BuildUser,
		BuildTime: n.mmeta.Meta.BuildTime,
		Data:      data,
	}
}

// BaseUrl return base url for restful service
func (n *Node) BaseUrl() string {
	return fmt.Sprintf("http://%s:%d%s", n.memberNode.Addr, n.mmeta.Meta.Port, n.mmeta.Meta.RouteRootPath)
}

// String return string representation
func (n *Node) String() string {
	return fmt.Sprintf("Node %s, providing %s service at %s, memberlist port %s",
		n.memberNode.Name, n.mmeta.Meta.Service, n.BaseUrl(), fmt.Sprint(n.memberNode.Port))
}
