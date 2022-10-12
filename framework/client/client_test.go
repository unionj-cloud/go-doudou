package client_test

import (
	"bytes"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/go-msgpack/codec"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/require"
	"github.com/unionj-cloud/go-doudou/framework/client"
	"github.com/unionj-cloud/go-doudou/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/framework/memberlist"
	"github.com/unionj-cloud/go-doudou/framework/registry"
	nmock "github.com/unionj-cloud/go-doudou/framework/registry/nacos/mock"
	"github.com/wubin1989/nacos-sdk-go/v2/common/constant"
	"github.com/wubin1989/nacos-sdk-go/v2/model"
	"github.com/wubin1989/nacos-sdk-go/v2/vo"
	"os"
	"testing"
	"time"
)

var clientConfigTest = *constant.NewClientConfig(
	constant.WithTimeoutMs(10*1000),
	constant.WithBeatInterval(5*1000),
	constant.WithNotLoadCacheAtStart(true),
)

var serverConfigTest = *constant.NewServerConfig("console.nacos.io", 80, constant.WithContextPath("/nacos"))

var services = model.Service{
	Name:        "DEFAULT_GROUP@@DEMO",
	CacheMillis: 1000,
	Hosts: []model.Instance{
		{
			InstanceId: "10.10.10.10-80-a-DEMO",
			Port:       80,
			Ip:         "10.10.10.10",
			Weight:     10,
			Metadata: map[string]string{
				"rootPath": "/api",
			},
			ClusterName: "a",
			ServiceName: "DEMO",
			Enable:      true,
			Healthy:     true,
		},
		{
			InstanceId:  "10.10.10.11-80-a-DEMO",
			Port:        80,
			Ip:          "10.10.10.11",
			Weight:      9,
			Metadata:    map[string]string{},
			ClusterName: "a",
			ServiceName: "DEMO",
			Enable:      true,
			Healthy:     true,
		},
		{
			InstanceId:  "10.10.10.12-80-a-DEMO",
			Port:        80,
			Ip:          "10.10.10.12",
			Weight:      8,
			Metadata:    map[string]string{},
			ClusterName: "a",
			ServiceName: "DEMO",
			Enable:      true,
			Healthy:     false,
		},
		{
			InstanceId:  "10.10.10.13-80-a-DEMO",
			Port:        80,
			Ip:          "10.10.10.13",
			Weight:      7,
			Metadata:    map[string]string{},
			ClusterName: "a",
			ServiceName: "DEMO",
			Enable:      false,
			Healthy:     true,
		},
		{
			InstanceId:  "10.10.10.14-80-a-DEMO",
			Port:        80,
			Ip:          "10.10.10.14",
			Weight:      6,
			Metadata:    map[string]string{},
			ClusterName: "a",
			ServiceName: "DEMO",
			Enable:      true,
			Healthy:     true,
		},
	},
	Checksum:    "3bbcf6dd1175203a8afdade0e77a27cd1528787794594",
	LastRefTime: 1528787794594,
	Clusters:    "a",
}

func TestNacosRRServiceProvider_SelectServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	namingClient := nmock.NewMockINamingClient(ctrl)
	namingClient.
		EXPECT().
		SelectInstances(vo.SelectInstancesParam{
			Clusters:    []string{"a"},
			ServiceName: "testsvc",
			HealthyOnly: true,
		}).
		AnyTimes().
		Return(services.Hosts, nil)

	n := client.NewNacosRRServiceProvider("testsvc",
		client.WithNacosNamingClient(namingClient),
		client.WithNacosClusters([]string{"a"}))
	for i := 0; i < len(services.Hosts)*2; i++ {
		got := n.SelectServer()
		fmt.Println(got)
	}
}

func TestNacosWRRServiceProvider_SelectServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	namingClient := nmock.NewMockINamingClient(ctrl)
	namingClient.
		EXPECT().
		SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
			Clusters:    []string{"a"},
			ServiceName: "testsvc",
		}).
		AnyTimes().
		Return(&services.Hosts[0], nil)

	n := client.NewNacosWRRServiceProvider("testsvc",
		client.WithNacosNamingClient(namingClient),
		client.WithNacosClusters([]string{"a"}))
	got := n.SelectServer()
	require.Equal(t, got, "http://10.10.10.10:80/api")
}

type MockDdClient struct {
	provider registry.IServiceProvider
	client   *resty.Client
	rootPath string
}

func (receiver *MockDdClient) SetRootPath(rootPath string) {
	receiver.rootPath = rootPath
}

func (receiver *MockDdClient) SetProvider(provider registry.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *MockDdClient) SetClient(client *resty.Client) {
	receiver.client = client
}

func NewMockDdClient(opts ...client.DdClientOption) *MockDdClient {
	defaultProvider := client.NewServiceProvider("MOCKDDCLIENT")
	defaultClient := client.NewClient()

	svcClient := &MockDdClient{
		provider: defaultProvider,
		client:   defaultClient,
	}

	for _, opt := range opts {
		opt(svcClient)
	}

	return svcClient
}

func TestWithProvider(t *testing.T) {
	Convey("Create a DdClient instance with custom provider", t, func() {
		m := NewMockDdClient(client.WithProvider(client.NewMemberlistServiceProvider("mock-svc")))
		So(m.provider, ShouldNotBeZeroValue)
	})
}

func TestWithClient(t *testing.T) {
	Convey("Create a DdClient instance with custom client", t, func() {
		m := NewMockDdClient(client.WithClient(resty.New()))
		So(m.client, ShouldNotBeZeroValue)
	})
}

func TestWithRootPath(t *testing.T) {
	Convey("Create a DdClient instance with custom rootPath", t, func() {
		m := NewMockDdClient(client.WithRootPath("/v1"))
		So(m.rootPath, ShouldEqual, "/v1")
	})
}

func TestRetryCount(t *testing.T) {
	Convey("Create a DdClient instance with 10 retry count", t, func() {
		os.Setenv("GDD_RETRY_COUNT", "10")
		m := NewMockDdClient()
		So(m.client.RetryCount, ShouldEqual, 10)
	})
}

func TestMain(m *testing.M) {
	setup()
	registry.NewNode()
	defer registry.Shutdown()
	m.Run()
}

func setup() {
	_ = config.GddMemSeed.Write("")
	_ = config.GddServiceName.Write("ddhttp")
	_ = config.GddMemName.Write("ddhttp")
	_ = config.GddMemWeight.Write("8")
	_ = config.GddMemDeadTimeout.Write("8s")
	_ = config.GddMemSyncInterval.Write("8s")
	_ = config.GddMemReclaimTimeout.Write("8s")
	_ = config.GddMemProbeInterval.Write("8s")
	_ = config.GddMemProbeTimeout.Write("8s")
	_ = config.GddMemSuspicionMult.Write("8")
	_ = config.GddMemGossipNodes.Write("8")
	_ = config.GddMemGossipInterval.Write("8s")
	_ = config.GddMemWeightInterval.Write("8s")
	_ = config.GddMemTCPTimeout.Write("8s")
	_ = config.GddMemHost.Write("localhost")
	_ = config.GddMemIndirectChecks.Write("8")
	_ = config.GddLogLevel.Write("debug")
	_ = config.GddPort.Write("8088")
	_ = config.GddRouteRootPath.Write("/v1")
}

func Test_base_AddNode(t *testing.T) {
	Convey("Should select one node", t, func() {
		provider := client.NewMemberlistServiceProvider("ddhttp")
		provider.AddNode(registry.LocalNode())
		So(provider.SelectServer(), ShouldEqual, fmt.Sprintf("http://%s:%d%s", "localhost", 8088, "/v1"))

		Convey("Node should be removed", func() {
			provider.RemoveNode(registry.LocalNode())
			So(provider.GetServer("ddhttp"), ShouldBeNil)
		})
	})
}

func Test_base_AddNode_Fail(t *testing.T) {
	Convey("Should select one node", t, func() {
		provider := client.NewMemberlistServiceProvider("test")
		provider.AddNode(registry.LocalNode())
		So(provider.SelectServer(), ShouldBeZeroValue)

		Convey("Node should not be removed", func() {
			provider.RemoveNode(registry.LocalNode())
			So(provider.GetServer("ddhttp"), ShouldBeNil)
		})
	})
}

type nodeMeta struct {
	Service       string     `json:"service"`
	RouteRootPath string     `json:"routeRootPath"`
	Port          int        `json:"port"`
	RegisterAt    *time.Time `json:"registerAt"`
	GoVer         string     `json:"goVer"`
	GddVer        string     `json:"gddVer"`
	BuildUser     string     `json:"buildUser"`
	BuildTime     string     `json:"buildTime"`
	Weight        int        `json:"weight"`
}

type mergedMeta struct {
	Meta nodeMeta               `json:"_meta,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
}

func setMetaWeight(node *memberlist.Node, weight int) error {
	var mm mergedMeta
	if len(node.Meta) > 0 {
		r := bytes.NewReader(node.Meta)
		dec := codec.NewDecoder(r, &codec.MsgpackHandle{})
		if err := dec.Decode(&mm); err != nil {
			return err
		}
	}
	mm.Meta.Weight = weight
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.MsgpackHandle{})
	if err := enc.Encode(mm); err != nil {
		return err
	}
	node.Meta = buf.Bytes()
	return nil
}

func Test_base_UpdateWeight(t *testing.T) {
	_ = setMetaWeight(registry.LocalNode(), 0)
	Convey("Weight should change", t, func() {
		provider := client.NewMemberlistServiceProvider("ddhttp")
		provider.AddNode(registry.LocalNode())
		registry.LocalNode().Weight = 4
		provider.UpdateWeight(registry.LocalNode())
		So(provider.GetServer("ddhttp").Weight(), ShouldEqual, 4)
	})
}

func Test_base_UpdateWeight_Fail(t *testing.T) {
	_ = setMetaWeight(registry.LocalNode(), 8)
	Convey("Weight should not change as meta weight greater than 0", t, func() {
		provider := client.NewMemberlistServiceProvider("ddhttp")
		provider.AddNode(registry.LocalNode())
		registry.LocalNode().Weight = 4
		provider.UpdateWeight(registry.LocalNode())
		So(provider.GetServer("ddhttp").Weight(), ShouldEqual, 8)
	})
}

func Test_base_UpdateWeight_Fail_SvcName_Mismatch(t *testing.T) {
	_ = setMetaWeight(registry.LocalNode(), 0)
	Convey("Weight should not change as service name mismatch", t, func() {
		provider := client.NewMemberlistServiceProvider("test")
		provider.UpdateWeight(registry.LocalNode())
		So(provider.GetServer("ddhttp"), ShouldBeNil)
	})
}

func Test_SMRR_AddNode(t *testing.T) {
	Convey("Should select one node", t, func() {
		provider := client.NewSmoothWeightedRoundRobinProvider("ddhttp")
		provider.AddNode(registry.LocalNode())
		So(provider.SelectServer(), ShouldEqual, fmt.Sprintf("http://%s:%d%s", "localhost", 8088, "/v1"))
	})
}

func Test_SMRR_SelectServer(t *testing.T) {
	Convey("Should select none node", t, func() {
		provider := client.NewSmoothWeightedRoundRobinProvider("ddhttp")
		provider.RemoveNode(registry.LocalNode())
		So(provider.SelectServer(), ShouldBeZeroValue)
	})
}
