package ddhttp

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/unionj-cloud/go-doudou/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/framework/logger"
	"github.com/unionj-cloud/go-doudou/framework/memberlist"
	"github.com/unionj-cloud/go-doudou/framework/registry"
	"github.com/unionj-cloud/go-doudou/framework/registry/nacos"
	"github.com/unionj-cloud/go-doudou/toolkit/cast"
	"github.com/wubin1989/nacos-sdk-go/clients/naming_client"
	"github.com/wubin1989/nacos-sdk-go/model"
	"github.com/wubin1989/nacos-sdk-go/vo"
	"net"
	"net/http"
	"os"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// DdClient defines service client interface
type DdClient interface {
	SetProvider(provider registry.IServiceProvider)
	SetClient(client *resty.Client)
	SetRootPath(rootPath string)
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

// WithRootPath sets root path for sending http requests
func WithRootPath(rootPath string) DdClientOption {
	return func(c DdClient) {
		c.SetRootPath(rootPath)
	}
}

// ServiceProvider defines an implementation for IServiceProvider
type ServiceProvider struct {
	server string
}

func (s *ServiceProvider) AddNode(node *memberlist.Node) {
}

func (s *ServiceProvider) RemoveNode(node *memberlist.Node) {
}

func (s *ServiceProvider) UpdateWeight(node *memberlist.Node) {
}

// SelectServer return service address from environment variable
func (s *ServiceProvider) SelectServer() string {
	return s.server
}

// NewServiceProvider creates new ServiceProvider instance
func NewServiceProvider(env string) *ServiceProvider {
	return &ServiceProvider{
		server: os.Getenv(env),
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
	client.SetTransport(&nethttp.Transport{
		RoundTripper: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
			MaxConnsPerHost:       100,
		},
	})
	retryCnt := config.DefaultGddRetryCount
	if cnt, err := cast.ToIntE(config.GddRetryCount.Load()); err == nil {
		retryCnt = cnt
	}
	client.SetRetryCount(retryCnt)
	return client
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
		logger.Infof("[go-doudou] add node %s to load balancer, supplying %s service", node.Name, svcName)
	} else {
		old := *s
		s.baseUrl = baseUrl
		s.weight = weight
		logger.Infof("[go-doudou] node %s update, supplying %s service, old: %+v, new: %+v", node.Name, svcName, old, *s)
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
		logger.Infof("[go-doudou] weight of node %s update, old: %d, new: %d", node.Name, old.weight, s.weight)
	}
}

func (m *base) GetServer(nodeName string) *server {
	return m.nodeMap[nodeName]
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
		logger.Infof("[go-doudou] remove node %s from load balancer, supplying %s service", node.Name, svcName)
	}
}

// MemberlistServiceProvider defines an implementation for IServiceProvider
type MemberlistServiceProvider struct {
	base
	current uint64
}

// SelectServer selects a node which is supplying service specified by name property from cluster
func (m *MemberlistServiceProvider) SelectServer() string {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if len(m.nodes) == 0 {
		return ""
	}
	next := int(atomic.AddUint64(&m.current, uint64(1)) % uint64(len(m.nodes)))
	m.current = uint64(next)
	selected := m.nodes[next]
	return selected.baseUrl
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
func (m *SmoothWeightedRoundRobinProvider) SelectServer() string {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if len(m.nodes) == 0 {
		return ""
	}
	var selected *server
	total := 0
	for i := 0; i < len(m.nodes); i++ {
		s := m.nodes[i]
		s.currentWeight += s.weight
		total += s.weight
		if selected == nil || s.currentWeight > selected.currentWeight {
			selected = s
		}
	}
	selected.currentWeight -= total
	return selected.baseUrl
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

type iNacosServiceProvider interface {
	SetClusters(clusters []string)
	SetGroupName(groupName string)
	SetNamingClient(namingClient naming_client.INamingClient)
}

type nacosBase struct {
	clusters    []string //optional,default:DEFAULT
	serviceName string   //required
	groupName   string   //optional,default:DEFAULT_GROUP

	lock         sync.RWMutex
	namingClient naming_client.INamingClient
}

func (b *nacosBase) SetClusters(clusters []string) {
	b.clusters = clusters
}

func (b *nacosBase) SetGroupName(groupName string) {
	b.groupName = groupName
}

func (b *nacosBase) SetNamingClient(namingClient naming_client.INamingClient) {
	b.namingClient = namingClient
}

func (b *nacosBase) AddNode(node *memberlist.Node) {
}

func (b *nacosBase) RemoveNode(node *memberlist.Node) {
}

func (b *nacosBase) UpdateWeight(node *memberlist.Node) {
}

type NacosProviderOption func(iNacosServiceProvider)

func WithNacosClusters(clusters []string) NacosProviderOption {
	return func(provider iNacosServiceProvider) {
		provider.SetClusters(clusters)
	}
}

func WithNacosGroupName(groupName string) NacosProviderOption {
	return func(provider iNacosServiceProvider) {
		provider.SetGroupName(groupName)
	}
}

func WithNacosNamingClient(namingClient naming_client.INamingClient) NacosProviderOption {
	return func(provider iNacosServiceProvider) {
		provider.SetNamingClient(namingClient)
	}
}

type instance []model.Instance

func (a instance) Len() int {
	return len(a)
}

func (a instance) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a instance) Less(i, j int) bool {
	return a[i].InstanceId < a[j].InstanceId
}

// NacosRRServiceProvider is a simple round-robin load balance implementation for IServiceProvider
type NacosRRServiceProvider struct {
	nacosBase
	current uint64
}

// SelectServer return service address from environment variable
func (n *NacosRRServiceProvider) SelectServer() string {
	n.lock.RLock()
	defer n.lock.RUnlock()
	if n.namingClient == nil {
		logger.Error("[go-doudou] nacos discovery client has not been initialized")
		return ""
	}
	instances, err := n.namingClient.SelectInstances(vo.SelectInstancesParam{
		Clusters:    n.clusters,
		ServiceName: n.serviceName,
		GroupName:   n.groupName,
		HealthyOnly: true,
	})
	if err != nil {
		logger.Error(fmt.Sprintf("[go-doudou] error:%s", err))
		return ""
	}
	sort.Sort(instance(instances))
	next := int(atomic.AddUint64(&n.current, uint64(1)) % uint64(len(instances)))
	n.current = uint64(next)
	selected := instances[next]
	return fmt.Sprintf("http://%s:%d%s", selected.Ip, selected.Port, selected.Metadata["rootPath"])
}

// NewNacosRRServiceProvider creates new ServiceProvider instance
func NewNacosRRServiceProvider(serviceName string, opts ...NacosProviderOption) *NacosRRServiceProvider {
	provider := &NacosRRServiceProvider{
		nacosBase: nacosBase{
			serviceName:  serviceName,
			namingClient: nacos.NamingClient,
		},
	}
	for _, opt := range opts {
		opt(provider)
	}
	return provider
}

// NacosWRRServiceProvider is a WRR load balance implementation for IServiceProvider
type NacosWRRServiceProvider struct {
	nacosBase
}

// SelectServer return service address from environment variable
func (n *NacosWRRServiceProvider) SelectServer() string {
	n.lock.RLock()
	defer n.lock.RUnlock()
	if n.namingClient == nil {
		logger.Error("[go-doudou] nacos discovery client has not been initialized")
		return ""
	}
	instance, err := n.namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		Clusters:    n.clusters,
		ServiceName: n.serviceName,
		GroupName:   n.groupName,
	})
	if err != nil {
		logger.Error(fmt.Sprintf("[go-doudou] failed to select one healthy instance:%s", err))
		return ""
	}
	return fmt.Sprintf("http://%s:%d%s", instance.Ip, instance.Port, instance.Metadata["rootPath"])
}

// NewNacosWRRServiceProvider creates new ServiceProvider instance
func NewNacosWRRServiceProvider(serviceName string, opts ...NacosProviderOption) *NacosWRRServiceProvider {
	provider := &NacosWRRServiceProvider{
		nacosBase{
			serviceName:  serviceName,
			namingClient: nacos.NamingClient,
		},
	}
	for _, opt := range opts {
		opt(provider)
	}
	return provider
}
