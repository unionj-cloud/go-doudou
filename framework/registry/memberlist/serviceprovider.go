package memberlist

import (
	"context"
	"fmt"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/memberlist"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	logger "github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	"google.golang.org/grpc"
	"sync"
	"sync/atomic"
	"time"
)

// IMemberlistServiceProvider defines service provider interface for server discovery
type IMemberlistServiceProvider interface {
	AddNode(node *memberlist.Node)
	UpdateWeight(node *memberlist.Node)
	RemoveNode(node *memberlist.Node)
}

type server struct {
	service       string
	node          string
	baseUrl       string
	weight        int
	currentWeight int
}

func (s *server) Weight() int {
	return s.weight
}

type base struct {
	name    string
	nodes   []*server
	nodeMap map[string]*server
}

func (m *base) GetService(meta NodeMeta) Service {
	for _, service := range meta.Services {
		if service.Name == m.name {
			return service
		}
	}
	return Service{}
}

// AddNode add or update node providing the service
func (m *base) AddNode(node *memberlist.Node) {
	meta, _ := ParseMeta(node)
	service := m.GetService(meta)
	if stringutils.IsEmpty(service.Name) {
		return
	}
	baseUrl := service.BaseUrl()
	weight := meta.Weight
	if s, exists := m.nodeMap[node.Name]; !exists {
		s = &server{
			service:       m.name,
			node:          node.Name,
			baseUrl:       baseUrl,
			weight:        weight,
			currentWeight: 0,
		}
		m.nodes = append(m.nodes, s)
		m.nodeMap[node.Name] = s
		logger.Info().Msgf("[go-doudou] add node %s to load balancer, supplying %s service", node.Name, service.Name)
	} else {
		old := *s
		s.baseUrl = baseUrl
		s.weight = weight
		logger.Info().Msgf("[go-doudou] node %s update, supplying %s service, old: %+v, new: %+v", node.Name, service.Name, old, *s)
	}
}

func (m *base) UpdateWeight(node *memberlist.Node) {
	meta, _ := ParseMeta(node)
	if meta.Weight > 0 {
		return
	}
	service := m.GetService(meta)
	if stringutils.IsEmpty(service.Name) {
		return
	}
	if s, exists := m.nodeMap[node.Name]; exists {
		old := *s
		s.weight = node.Weight
		logger.Info().Msgf("[go-doudou] weight of node %s update, old: %d, new: %d", node.Name, old.weight, s.weight)
	}
}

func (m *base) GetServer(nodeName string) *server {
	return m.nodeMap[nodeName]
}

func (m *base) RemoveNode(node *memberlist.Node) {
	meta, _ := ParseMeta(node)
	service := m.GetService(meta)
	if stringutils.IsEmpty(service.Name) {
		return
	}
	if _, exists := m.nodeMap[node.Name]; exists {
		var idx int
		for i, n := range m.nodes {
			if n.node == node.Name {
				idx = i
			}
		}
		m.nodes = append(m.nodes[:idx], m.nodes[idx+1:]...)
		delete(m.nodeMap, node.Name)
		logger.Info().Msgf("[go-doudou] remove node %s from load balancer, supplying %s service", node.Name, service.Name)
	}
}

var _ IMemberlistServiceProvider = (*RRServiceProvider)(nil)

// RRServiceProvider defines an implementation for IMemberlistServiceProvider
type RRServiceProvider struct {
	base    base
	current uint64
	lock    sync.RWMutex
}

func (m *RRServiceProvider) AddNode(node *memberlist.Node) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.base.AddNode(node)
}

func (m *RRServiceProvider) UpdateWeight(node *memberlist.Node) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.base.UpdateWeight(node)
}

func (m *RRServiceProvider) RemoveNode(node *memberlist.Node) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.base.RemoveNode(node)
}

// SelectServer selects a node which is supplying service specified by name property from cluster
func (m *RRServiceProvider) SelectServer() string {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if len(m.base.nodes) == 0 {
		return ""
	}
	next := int(atomic.AddUint64(&m.current, uint64(1)) % uint64(len(m.base.nodes)))
	m.current = uint64(next)
	selected := m.base.nodes[next]
	return selected.baseUrl
}

// NewRRServiceProvider create an RRServiceProvider instance
func NewRRServiceProvider(name string) *RRServiceProvider {
	sp := &RRServiceProvider{
		base: base{
			name:    name,
			nodeMap: make(map[string]*server),
		},
	}
	RegisterServiceProvider(sp)
	return sp
}

var _ IMemberlistServiceProvider = (*SWRRServiceProvider)(nil)

// SWRRServiceProvider is a smooth weighted round-robin algo implementation for IMemberlistServiceProvider
// https://github.com/nginx/nginx/commit/52327e0627f49dbda1e8db695e63a4b0af4448b1
type SWRRServiceProvider struct {
	base base
	lock sync.RWMutex
}

func (m *SWRRServiceProvider) AddNode(node *memberlist.Node) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.base.AddNode(node)
}

func (m *SWRRServiceProvider) UpdateWeight(node *memberlist.Node) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.base.UpdateWeight(node)
}

func (m *SWRRServiceProvider) RemoveNode(node *memberlist.Node) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.base.RemoveNode(node)
}

// SelectServer selects a node which is supplying service specified by name property from cluster
func (m *SWRRServiceProvider) SelectServer() string {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if len(m.base.nodes) == 0 {
		return ""
	}
	var selected *server
	total := 0
	for i := 0; i < len(m.base.nodes); i++ {
		s := m.base.nodes[i]
		s.currentWeight += s.weight
		total += s.weight
		if selected == nil || s.currentWeight > selected.currentWeight {
			selected = s
		}
	}
	selected.currentWeight -= total
	return selected.baseUrl
}

// NewSWRRServiceProvider create an SWRRServiceProvider instance
func NewSWRRServiceProvider(name string) *SWRRServiceProvider {
	sp := &SWRRServiceProvider{
		base: base{
			name:    name,
			nodeMap: make(map[string]*server),
		},
	}
	RegisterServiceProvider(sp)
	return sp
}

func NewSWRRGrpcClientConn(service string, dialOptions ...grpc.DialOption) *grpc.ClientConn {
	return NewGrpcClientConn(service, "memberlist_weight_balancer", dialOptions...)
}

func NewRRGrpcClientConn(service string, dialOptions ...grpc.DialOption) *grpc.ClientConn {
	return NewGrpcClientConn(service, "round_robin", dialOptions...)
}

func NewGrpcClientConn(service string, lb string, dialOptions ...grpc.DialOption) *grpc.ClientConn {
	serverAddr := fmt.Sprintf(schemeName+"://%s/", service)
	dialOptions = append(dialOptions, grpc.WithBlock(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "`+lb+`"}`))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	grpcConn, err := grpc.DialContext(ctx, serverAddr, dialOptions...)
	if err != nil {
		logger.Panic().Err(err).Msgf("[go-doudou] failed to connect to server %s", serverAddr)
	}
	return grpcConn
}
