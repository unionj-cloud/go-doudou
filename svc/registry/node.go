package registry

import (
	"github.com/hashicorp/memberlist"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"os"
	"sync"
)

type NodeMeta struct {
	Service string `json:"service"`
}

type NodeConfig struct {
	memberConf *memberlist.Config
}

type Node struct {
	conf        *NodeConfig
	meta        NodeMeta
	broadcasts  *memberlist.TransmitLimitedQueue
	memberlist  *memberlist.Memberlist
	lock        sync.Mutex
	memberLock  sync.RWMutex
	memberNode  *memberlist.Node
}

func NewNode(conf *NodeConfig) (*Node, error) {
	service := os.Getenv("NODE_SERVICE")
	if stringutils.IsEmpty(service) {
		return nil, errors.New("No env variable NODE_SERVICE found")
	}
	node := &Node{
		conf: conf,
		meta: NodeMeta{
			service,
		},
	}
	conf.memberConf.Delegate = &registry{node}
	list, err := memberlist.Create(conf.memberConf)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create memberlist")
	}
	node.memberlist = list
	node.broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes:       node.NumNodes,
		RetransmitMult: conf.memberConf.RetransmitMult,
	}
	return node, nil
}

func (s *Node) NumNodes() (numNodes int) {
	s.memberLock.RLock()
	numNodes = len(s.memberlist.Members())
	s.memberLock.RUnlock()

	return numNodes
}
