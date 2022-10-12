package client

import (
	"github.com/go-resty/resty/v2"
	"github.com/klauspost/compress/gzhttp"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/unionj-cloud/go-doudou/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/framework/registry"
	"github.com/unionj-cloud/go-doudou/toolkit/cast"
	"net"
	"net/http"
	"os"
	"runtime"
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
	client.SetTransport(gzhttp.Transport(&nethttp.Transport{
		RoundTripper: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
			MaxConnsPerHost:       10000,
		},
	}))
	retryCnt := config.DefaultGddRetryCount
	if cnt, err := cast.ToIntE(config.GddRetryCount.Load()); err == nil {
		retryCnt = cnt
	}
	client.SetRetryCount(retryCnt)
	return client
}
