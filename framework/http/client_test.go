package ddhttp_test

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/stretchr/testify/require"
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	"github.com/unionj-cloud/go-doudou/framework/http/mock"
	"testing"
)

var clientConfigTest = *constant.NewClientConfig(
	constant.WithTimeoutMs(10*1000),
	constant.WithBeatInterval(5*1000),
	constant.WithNotLoadCacheAtStart(true),
)

var serverConfigTest = *constant.NewServerConfig("console.nacos.io", 80, constant.WithContextPath("/nacos"))

var services = model.Service{
	Name:            "DEFAULT_GROUP@@DEMO",
	CacheMillis:     1000,
	UseSpecifiedURL: false,
	Hosts: []model.Instance{
		{
			Valid:      true,
			Marked:     false,
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
			Valid:       true,
			Marked:      false,
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
			Valid:       true,
			Marked:      false,
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
			Valid:       true,
			Marked:      false,
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
			Valid:       true,
			Marked:      false,
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
	LastRefTime: 1528787794594, Env: "", Clusters: "a",
	Metadata: map[string]string(nil),
}

func TestNacosRRServiceProvider_SelectServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	namingClient := mock.NewMockINamingClient(ctrl)
	namingClient.
		EXPECT().
		SelectInstances(vo.SelectInstancesParam{
			Clusters:    []string{"a"},
			ServiceName: "testsvc",
			HealthyOnly: true,
		}).
		AnyTimes().
		Return(services.Hosts, nil)

	n := ddhttp.NewNacosRRServiceProvider(namingClient, "testsvc", ddhttp.WithNacosClusters([]string{"a"}))
	for i := 0; i < len(services.Hosts)*2; i++ {
		got := n.SelectServer()
		fmt.Println(got)
	}
}

func TestNacosWRRServiceProvider_SelectServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	namingClient := mock.NewMockINamingClient(ctrl)
	namingClient.
		EXPECT().
		SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
			Clusters:    []string{"a"},
			ServiceName: "testsvc",
		}).
		AnyTimes().
		Return(&services.Hosts[0], nil)

	n := ddhttp.NewNacosWRRServiceProvider(namingClient, "testsvc", ddhttp.WithNacosClusters([]string{"a"}))
	got := n.SelectServer()
	require.Equal(t, got, "http://10.10.10.10:80/api")
}
