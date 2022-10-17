package restclient_test

import (
	"github.com/go-resty/resty/v2"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/unionj-cloud/go-doudou/v2/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry"
	"github.com/unionj-cloud/go-doudou/v2/framework/restclient"
	"github.com/wubin1989/nacos-sdk-go/v2/common/constant"
	"os"
	"testing"
)

var clientConfigTest = *constant.NewClientConfig(
	constant.WithTimeoutMs(10*1000),
	constant.WithBeatInterval(5*1000),
	constant.WithNotLoadCacheAtStart(true),
)

var serverConfigTest = *constant.NewServerConfig("console.nacos.io", 80, constant.WithContextPath("/nacos"))

type MockRestClient struct {
	provider registry.IServiceProvider
	client   *resty.Client
	rootPath string
}

func (receiver *MockRestClient) SetRootPath(rootPath string) {
	receiver.rootPath = rootPath
}

func (receiver *MockRestClient) SetProvider(provider registry.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *MockRestClient) SetClient(client *resty.Client) {
	receiver.client = client
}

func NewMockRestClient(opts ...restclient.RestClientOption) *MockRestClient {
	defaultProvider := restclient.NewServiceProvider("MOCKRESTCLIENT")
	defaultClient := restclient.NewClient()

	svcClient := &MockRestClient{
		provider: defaultProvider,
		client:   defaultClient,
	}

	for _, opt := range opts {
		opt(svcClient)
	}

	return svcClient
}

func TestWithClient(t *testing.T) {
	Convey("Create a RestClient instance with custom client", t, func() {
		m := NewMockRestClient(restclient.WithClient(resty.New()))
		So(m.client, ShouldNotBeZeroValue)
	})
}

func TestWithRootPath(t *testing.T) {
	Convey("Create a RestClient instance with custom rootPath", t, func() {
		m := NewMockRestClient(restclient.WithRootPath("/v1"))
		So(m.rootPath, ShouldEqual, "/v1")
	})
}

func TestRetryCount(t *testing.T) {
	Convey("Create a RestClient instance with 10 retry count", t, func() {
		os.Setenv("GDD_RETRY_COUNT", "10")
		m := NewMockRestClient()
		So(m.client.RetryCount, ShouldEqual, 10)
	})
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
}

func setup() {
	_ = config.GddServiceName.Write("awesome-service")
	_ = config.GddLogLevel.Write("debug")
	_ = config.GddPort.Write("8088")
	_ = config.GddRouteRootPath.Write("/v1")
}
