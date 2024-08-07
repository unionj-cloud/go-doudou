package restclient

import (
	"net"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/klauspost/compress/gzhttp"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/unionj-cloud/go-doudou/v2/framework/config"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/cast"
)

// RestClient defines service client interface
type RestClient interface {
	SetProvider(provider registry.IServiceProvider)
	SetClient(client *resty.Client)
	SetRootPath(rootPath string)
}

// RestClientOption defines configure function type
type RestClientOption func(RestClient)

// WithProvider sets service provider
func WithProvider(provider registry.IServiceProvider) RestClientOption {
	return func(c RestClient) {
		c.SetProvider(provider)
	}
}

// WithClient sets http client
func WithClient(client *resty.Client) RestClientOption {
	return func(c RestClient) {
		c.SetClient(client)
	}
}

// WithRootPath sets root path for sending http requests
func WithRootPath(rootPath string) RestClientOption {
	return func(c RestClient) {
		c.SetRootPath(rootPath)
	}
}

// ServiceProvider defines an implementation for IServiceProvider
type ServiceProvider struct {
	server string
}

func (s *ServiceProvider) Close() {
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
