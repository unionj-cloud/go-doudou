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

type DdClient interface {
	SetProvider(provider IServiceProvider)
	SetClient(client *resty.Client)
}

type DdClientOption func(DdClient)

func WithProvider(provider IServiceProvider) DdClientOption {
	return func(c DdClient) {
		c.SetProvider(provider)
	}
}

func WithClient(client *resty.Client) DdClientOption {
	return func(c DdClient) {
		c.SetClient(client)
	}
}

type IServiceProvider interface {
	SelectServer() (string, error)
}

type ServiceProvider struct {
	Env string
}

func (s *ServiceProvider) SelectServer() (string, error) {
	address := os.Getenv(s.Env)
	if stringutils.IsEmpty(address) {
		return "", errors.Errorf("no service address found from environment variable %s", s.Env)
	}
	return address, nil
}

type ServiceProviderOption func(IServiceProvider)

func NewServiceProvider(env string, opts ...ServiceProviderOption) IServiceProvider {
	provider := &ServiceProvider{
		Env: env,
	}

	for _, opt := range opts {
		opt(provider)
	}

	return provider
}

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

type MemberlistServiceProvider struct {
	// Name of the service that dependent on
	name     string
	registry registry.IRegistry
	current  uint64
}

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

type MemberlistProviderOption func(IServiceProvider)

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
