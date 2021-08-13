package registry

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/logutils"
	"github.com/hashicorp/memberlist"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/cast"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/svc/config"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

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

func (r *registry) Register() error {
	if r.memberlist == nil {
		return errors.New("Memberlist is nil")
	}
	seed := config.GddMemSeed.Load()
	if stringutils.IsEmpty(seed) {
		logrus.Warnln("No seed found")
		return nil
	}
	_, err := r.memberlist.Join([]string{seed})
	if err != nil {
		return errors.Wrap(err, "Failed to join cluster")
	}
	logrus.Infof("Node %s joined cluster successfully", r.memberlist.LocalNode().FullAddress())
	return nil
}

func (r *registry) Discover(svc string) ([]*Node, error) {
	if r.memberlist == nil {
		return nil, errors.New("Memberlist is nil")
	}
	var nodes []*Node
	for _, member := range r.memberlist.Members() {
		logrus.Debugf("Member: %s %s\n", member.Name, member.Addr)
		var mmeta mergedMeta
		if err := json.Unmarshal(member.Meta, &mmeta); err != nil {
			return nil, errors.Wrap(err, "")
		}
		if mmeta.Meta.Service == svc {
			nodes = append(nodes, &Node{
				mmeta:      mmeta,
				state:      Alive,
				memberNode: member,
				remote:     true,
			})
		}
	}
	return nodes, nil
}

type nodeMeta struct {
	Service string `json:"service"`
	BaseUrl string `json:"baseUrl"`
	Port    int    `json:"port"`
	Host    string `json:"host"`
}

func newMeta(mnode *memberlist.Node) (mergedMeta, error) {
	var mm mergedMeta
	if err := json.Unmarshal(mnode.Meta, &mm); err != nil {
		return mm, errors.Wrap(err, "Unmarshal node meta failed, not a valid json")
	}
	return mm, nil
}

type NodeState int

const (
	Alive NodeState = iota
	Leaving
	Left
	Shutdown
)

func (s NodeState) String() string {
	switch s {
	case Alive:
		return "alive"
	case Leaving:
		return "leaving"
	case Left:
		return "left"
	case Shutdown:
		return "shutdown"
	default:
		return "unknown"
	}
}

type mergedMeta struct {
	Meta nodeMeta    `json:"_meta,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type Node struct {
	mmeta      mergedMeta
	state      NodeState
	memberNode *memberlist.Node
	*registry
	// check the node is a local node or remote node
	remote bool
}

type NodeOption func(*Node)

func WithData(data interface{}) NodeOption {
	return func(node *Node) {
		node.mmeta = mergedMeta{
			Data: data,
		}
	}
}

// Borrow source code from https://github.com/phayes/freeport/blob/master/freeport.go
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

func NewNode(opts ...NodeOption) (*Node, error) {
	mconf := memberlist.DefaultWANConfig()
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
	memport := cast.ToInt(config.GddMemPort.Load())
	if memport == 0 {
		memport, _ = getFreePort()
	}
	if memport > 0 {
		mconf.BindPort = memport
		mconf.AdvertisePort = memport
	}
	nodeName := config.GddMemNodeName.Load()
	if stringutils.IsNotEmpty(nodeName) {
		mconf.Name = nodeName
	}
	service := config.GddServiceName.Load()
	if stringutils.IsEmpty(service) {
		return nil, errors.New(fmt.Sprintf("NewNode() error: No env variable %s found", config.GddServiceName))
	}
	node := &Node{
		state: -1,
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
	baseUrl := config.GddBaseUrl.Load()
	node.mmeta.Meta = nodeMeta{
		Service: service,
		Port:    port,
		BaseUrl: baseUrl,
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
	node.state = Alive
	node.memberNode = list.LocalNode()
	return node, nil
}

func (n *Node) NumNodes() (numNodes int) {
	n.memberLock.RLock()
	numNodes = len(n.memberlist.Members())
	n.memberLock.RUnlock()

	return numNodes
}

func (n *Node) BaseUrl() string {
	if stringutils.IsNotEmpty(n.mmeta.Meta.BaseUrl) {
		return n.mmeta.Meta.BaseUrl
	}
	return fmt.Sprintf("http://%s:%d", n.memberNode.Addr.String(), n.mmeta.Meta.Port)
}

func (n *Node) String() string {
	if stringutils.IsNotEmpty(n.mmeta.Meta.Service) {
		return fmt.Sprintf("Node %s, providing %s service at %s, memberlist port %s, service port %d",
			n.memberNode.Name, n.mmeta.Meta.Service, n.memberNode.Addr, fmt.Sprint(n.memberNode.Port), n.mmeta.Meta.Port)
	}
	return fmt.Sprintf("Node %s", n.memberNode.Name)
}
