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

// MemberlistServiceProvider defines an implementation for IServiceProvider
type MemberlistServiceProvider struct {
	// name is the service name
	name    string
	current uint64
	nodes   []*memberlist.Node
	nodeMap map[string]struct{}
	lock    sync.RWMutex
}

// AddNode add or update node providing the service
func (m *MemberlistServiceProvider) AddNode(node *memberlist.Node) {
	m.lock.Lock()
	defer m.lock.Unlock()
	svcName := registry.SvcName(node)
	if svcName != m.name {
		return
	}
	if _, exists := m.nodeMap[node.Name]; !exists {
		m.nodes = append(m.nodes, node)
		m.nodeMap[node.Name] = struct{}{}
		logrus.Infof("Node %s joined, supplying %s service", node.Name, svcName)
	}
}

func (m *MemberlistServiceProvider) RemoveNode(node *memberlist.Node) {
	m.lock.Lock()
	defer m.lock.Unlock()
	svcName := registry.SvcName(node)
	if svcName != m.name {
		return
	}
	if _, exists := m.nodeMap[node.Name]; exists {
		var idx int
		for i, n := range m.nodes {
			if n.Name == node.Name {
				idx = i
			}
		}
		m.nodes = append(m.nodes[:idx], m.nodes[idx+1:]...)
		delete(m.nodeMap, node.Name)
		logrus.Infof("Node %s left, supplying %s service", node.Name, svcName)
	}
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
	return registry.BaseUrl(selected)
}

// NewMemberlistServiceProvider create an NewMemberlistServiceProvider instance
func NewMemberlistServiceProvider(name string) *MemberlistServiceProvider {
	sp := &MemberlistServiceProvider{
		name:    name,
		nodeMap: make(map[string]struct{}),
	}
	registry.RegisterServiceProvider(sp)
	return sp
}

//// SmoothWeightedRoundRobinProvider is a smooth weighted round-robin algo implementation for IServiceProvider
//type SmoothWeightedRoundRobinProvider struct {
//	name string
//}
//
//// SelectServer selects a node which is supplying service specified by name property from cluster
//func (m *SmoothWeightedRoundRobinProvider) SelectServer() (string, error) {
//	nodes, err := m.registry.Discover(m.name)
//	if err != nil {
//		return "", errors.Wrap(err, "SelectServer() fail")
//	}
//	if len(nodes) == 0 {
//		return "", errors.Wrap(errors.New(fmt.Sprintf("no service %s supplier found", m.name)), "SelectServer() fail")
//	}
//	next := int(atomic.AddUint64(&m.current, uint64(1)) % uint64(len(nodes)))
//	m.current = uint64(next)
//	selected := nodes[next]
//	return selected.BaseUrl(), nil
//}
//
//// NewSmoothWeightedRoundRobinProvider create an SmoothWeightedRoundRobinProvider instance
//func NewSmoothWeightedRoundRobinProvider(name string, registry registry.IRegistry) *SmoothWeightedRoundRobinProvider {
//	provider := &SmoothWeightedRoundRobinProvider{
//		name: name,
//	}
//	return provider
//}
