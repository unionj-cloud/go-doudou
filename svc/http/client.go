package ddhttp

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/svc/registry"
	"net"
	"net/http"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

// DdClient defines service client interface
type DdClient interface {
	SetProvider(provider IServiceProvider)
	SetClient(client *resty.Client)
}

// DdClientOption defines configure function type
type DdClientOption func(DdClient)

// WithProvider sets service provider
func WithProvider(provider IServiceProvider) DdClientOption {
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

// IServiceProvider defines service provider interface for server discovery
type IServiceProvider interface {
	SelectServer() (string, error)
}

// ServiceProvider defines an implementation for IServiceProvider
type ServiceProvider struct {
	Env string
}

// SelectServer return service address from environment variable
func (s *ServiceProvider) SelectServer() (string, error) {
	address := os.Getenv(s.Env)
	if stringutils.IsEmpty(address) {
		return "", errors.Errorf("no service address found from environment variable %s", s.Env)
	}
	return address, nil
}

// ServiceProviderOption sets properties of ServiceProvider
type ServiceProviderOption func(IServiceProvider)

// NewServiceProvider creates new ServiceProvider instance
func NewServiceProvider(env string, opts ...ServiceProviderOption) IServiceProvider {
	provider := &ServiceProvider{
		Env: env,
	}

	for _, opt := range opts {
		opt(provider)
	}

	return provider
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

// MemberlistServiceProvider defines an implementation for IServiceProvider. Recommend to use.
type MemberlistServiceProvider struct {
	// Name of the service that dependent on
	name     string
	registry registry.IRegistry
	current  uint64
}

// SelectServer selects a node which is supplying service specified by name property from cluster
func (m *MemberlistServiceProvider) SelectServer() (string, error) {
	nodes, err := m.registry.Discover(m.name)
	if err != nil {
		return "", errors.Wrap(err, "SelectServer() fail")
	}
	if len(nodes) == 0 {
		return "", errors.Wrap(errors.New(fmt.Sprintf("no service %s supplier found", m.name)), "SelectServer() fail")
	}
	next := int(atomic.AddUint64(&m.current, uint64(1)) % uint64(len(nodes)))
	m.current = uint64(next)
	selected := nodes[next]
	return selected.BaseUrl(), nil
}

// MemberlistProviderOption defines a function for setting properties of MemberlistServiceProvider
type MemberlistProviderOption func(IServiceProvider)

// NewMemberlistServiceProvider create an NewMemberlistServiceProvider instance
func NewMemberlistServiceProvider(name string, registry registry.IRegistry, opts ...MemberlistProviderOption) IServiceProvider {
	provider := &MemberlistServiceProvider{
		name:     name,
		registry: registry,
	}

	for _, opt := range opts {
		opt(provider)
	}

	return provider
}
