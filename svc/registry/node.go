package registry

import (
	"encoding/json"
	"github.com/hashicorp/memberlist"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"os"
	"sync"
)

type IRegistry interface {
	Register() error
	Discovery(svc string) ([]*Node, error)
}

type NodeMeta struct {
	Service string `json:"service"`
}

type NodeConfig struct {
	memberConf *memberlist.Config
}

type Node struct {
	conf       *NodeConfig
	meta       NodeMeta
	broadcasts *memberlist.TransmitLimitedQueue
	memberlist *memberlist.Memberlist
	lock       sync.Mutex
	memberLock sync.RWMutex
	memberNode *memberlist.Node
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

func (n *Node) NumNodes() (numNodes int) {
	n.memberLock.RLock()
	numNodes = len(n.memberlist.Members())
	n.memberLock.RUnlock()

	return numNodes
}

func (n *Node) Register() error {
	seed := os.Getenv("SEED")
	if stringutils.IsEmpty(seed) {
		return errors.New("No seed found, register failed")
	}
	_, err := n.memberlist.Join([]string{seed})
	if err != nil {
		return errors.Wrap(err, "Failed to join cluster")
	}
	return nil
}

func (n *Node) Discovery(svc string) ([]*Node, error) {
	var nodes []*Node
	for _, member := range n.memberlist.Members() {
		logrus.Infof("Member: %s %s\n", member.Name, member.Addr)
		if member.State == memberlist.StateAlive {
			var nmeta NodeMeta
			if err := json.Unmarshal(member.Meta, &nmeta); err != nil {
				return nil, errors.Wrap(err, "")
			}
			if nmeta.Service == svc {
				nodes = append(nodes, &Node{
					meta:       nmeta,
					memberNode: member,
				})
			}
		}
	}
	return nodes, nil
}
