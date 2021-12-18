package ddhttp

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/svc/registry"
	"github.com/unionj-cloud/memberlist"
	"net"
	"net/http"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// DdClient defines service client interface
type DdClient interface {
	SetProvider(provider registry.IServiceProvider)
	SetClient(client *resty.Client)
}

// DdClientOption defines configure function type
type DdClientOption func(DdClient)

// WithProvider sets service provider
func WithProvider(provider registry.IServiceProvider) DdClientOption {
	return func(c DdClient) {
		c.SetProvider(provider)
	}
}

// WithClient sets http client
func WithClient(client *resty.Client) DdClientOption {
	return func(c DdClient) {
		c.SetClient(client)
	}
}

// ServiceProvider defines an implementation for IServiceProvider
type ServiceProvider struct {
	Env string
}

func (s *ServiceProvider) AddNode(node *memberlist.Node) {
}

func (s *ServiceProvider) RemoveNode(node *memberlist.Node) {
}

func (s *ServiceProvider) UpdateWeight(node *memberlist.Node) {
}

// SelectServer return service address from environment variable
func (s *ServiceProvider) SelectServer() (string, error) {
	address := os.Getenv(s.Env)
	if stringutils.IsEmpty(address) {
		return "", errors.Errorf("no service address found from environment variable %s", s.Env)
	}
	return address, nil
}

// NewServiceProvider creates new ServiceProvider instance
func NewServiceProvider(env string) *ServiceProvider {
	return &ServiceProvider{
		Env: env,
	}
}

// NewClient creates new resty Client instance
func NewClient() *resty.Client {
	client := resty.New()
	client.SetTimeout(1 * time.Minute)

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	client.SetTransport(&http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		MaxConnsPerHost:       100,
	})
	return client
}

type server struct {
	service       string
	node          string
	baseUrl       string
	weight        int
	currentWeight int
}

type base struct {
	name    string
	nodes   []*server
	nodeMap map[string]*server
	lock    sync.RWMutex
}

// AddNode add or update node providing the service
func (m *base) AddNode(node *memberlist.Node) {
	m.lock.Lock()
	defer m.lock.Unlock()
	svcName := registry.SvcName(node)
	if svcName != m.name {
		return
	}
	baseUrl, _ := registry.BaseUrl(node)
	weight, _ := registry.MetaWeight(node)
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
		logrus.Infof("node %s joined, supplying %s service", node.Name, svcName)
	} else {
		old := *s
		s.baseUrl = baseUrl
		s.weight = weight
		logrus.Infof("node %s update, supplying %s service, old: %+v, new: %+v", node.Name, svcName, old, *s)
	}
}

func (m *base) UpdateWeight(node *memberlist.Node) {
	weight, _ := registry.MetaWeight(node)
	if weight > 0 {
		return
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	svcName := registry.SvcName(node)
	if svcName != m.name {
		return
	}
	if s, exists := m.nodeMap[node.Name]; exists {
		old := *s
		s.weight = node.Weight
		logrus.Infof("weight of node %s update, old: %d, new: %d", node.Name, old.weight, s.weight)
	}
}

func (m *base) RemoveNode(node *memberlist.Node) {
	m.lock.Lock()
	defer m.lock.Unlock()
	svcName := registry.SvcName(node)
	if svcName != m.name {
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
		logrus.Infof("node %s left, supplying %s service", node.Name, svcName)
	}
}

// MemberlistServiceProvider defines an implementation for IServiceProvider
type MemberlistServiceProvider struct {
	base
	current uint64
}

// SelectServer selects a node which is supplying service specified by name property from cluster
func (m *MemberlistServiceProvider) SelectServer() (string, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if len(m.nodes) == 0 {
		return "", errors.Wrap(errors.New(fmt.Sprintf("no service %s supplier found", m.name)), "SelectServer() fail")
	}
	next := int(atomic.AddUint64(&m.current, uint64(1)) % uint64(len(m.nodes)))
	m.current = uint64(next)
	selected := m.nodes[next]
	return selected.baseUrl, nil
}

// NewMemberlistServiceProvider create an NewMemberlistServiceProvider instance
func NewMemberlistServiceProvider(name string) *MemberlistServiceProvider {
	sp := &MemberlistServiceProvider{
		base: base{
			name:    name,
			nodeMap: make(map[string]*server),
		},
	}
	registry.RegisterServiceProvider(sp)
	return sp
}

// SmoothWeightedRoundRobinProvider is a smooth weighted round-robin algo implementation for IServiceProvider
// https://github.com/nginx/nginx/commit/52327e0627f49dbda1e8db695e63a4b0af4448b1
type SmoothWeightedRoundRobinProvider struct {
	base
}

// SelectServer selects a node which is supplying service specified by name property from cluster
func (m *SmoothWeightedRoundRobinProvider) SelectServer() (string, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if len(m.nodes) == 0 {
		return "", errors.Wrap(errors.New(fmt.Sprintf("no service %s supplier found", m.name)), "SelectServer() fail")
	}
	var selected *server
	total := 0
	for i := 0; i < len(m.nodes); i++ {
		s := m.nodes[i]
		if s == nil {
			continue
		}
		s.currentWeight += s.weight
		total += s.weight
		if selected == nil || s.currentWeight > selected.currentWeight {
			selected = s
		}
	}
	if selected == nil {
		return "", errors.Wrap(errors.New(fmt.Sprintf("no service %s supplier found", m.name)), "SelectServer() fail")
	}
	selected.currentWeight -= total
	return selected.baseUrl, nil
}

// NewSmoothWeightedRoundRobinProvider create an SmoothWeightedRoundRobinProvider instance
func NewSmoothWeightedRoundRobinProvider(name string) *SmoothWeightedRoundRobinProvider {
	sp := &SmoothWeightedRoundRobinProvider{
		base: base{
			name:    name,
			nodeMap: make(map[string]*server),
		},
	}
	registry.RegisterServiceProvider(sp)
	return sp
}
