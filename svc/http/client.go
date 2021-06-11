package ddhttp

import (
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"net"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
)

type IServiceProvider interface {
	SelectServer() (string, error)
}

type ServiceProvider struct {
	Name string
}

func (s *ServiceProvider) SelectServer() (string, error) {
	re := regexp.MustCompile(`\s+`)
	address := os.Getenv(strings.ToUpper(re.ReplaceAllString(s.Name, "_")))
	if stringutils.IsEmpty(address) {
		return "", errors.New("No service address for Usersvc found!")
	}
	return address, nil
}

type ServiceProviderOption func(IServiceProvider)

func NewServiceProvider(name string, opts ...ServiceProviderOption) IServiceProvider {
	provider := &ServiceProvider{
		Name: name,
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
