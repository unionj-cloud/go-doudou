package nacos_test

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/require"
	"github.com/unionj-cloud/go-doudou/v2/framework/buildinfo"
	"github.com/unionj-cloud/go-doudou/v2/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry/nacos"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry/nacos/mock"
	"github.com/wubin1989/nacos-sdk-go/v2/clients/naming_client"
	"github.com/wubin1989/nacos-sdk-go/v2/model"
	"github.com/wubin1989/nacos-sdk-go/v2/vo"
	"testing"
)

func setup() {
	_ = config.GddServiceName.Write("seed")
	_ = config.GddLogLevel.Write("debug")
	_ = config.GddPort.Write("8088")
	_ = config.GddRouteRootPath.Write("/v1")
	_ = config.GddNacosServerAddr.Write("http://localhost:8848")
	_ = config.GddWeight.Write("5")
}

func TestNewRest(t *testing.T) {
	Convey("Should not have error", t, func() {
		setup()
		_ = config.GddNacosRegisterHost.Write("seed")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buildinfo.BuildTime = "Mon Jan 2 15:04:05 MST 2006"
		namingClient := mock.NewMockINamingClient(ctrl)
		namingClient.
			EXPECT().
			RegisterInstance(gomock.Any()).
			AnyTimes().
			Return(true, nil)

		namingClient.
			EXPECT().
			DeregisterInstance(gomock.Any()).
			AnyTimes().
			Return(true, nil)

		nacos.NewNamingClient = func(param vo.NacosClientParam) (iClient naming_client.INamingClient, err error) {
			return namingClient, nil
		}

		if nacos.NamingClient == nil {
			nacos.NamingClient = namingClient
		}

		So(func() {
			nacos.NewRest(map[string]interface{}{
				"foo": "bar",
			})
		}, ShouldNotPanic)

		nacos.ShutdownRest()
	})
}

func TestNewRest2(t *testing.T) {
	Convey("Should not have error", t, func() {
		setup()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buildinfo.BuildTime = "Mon Jan 2 15:04:05 MST 2006"
		namingClient := mock.NewMockINamingClient(ctrl)
		namingClient.
			EXPECT().
			RegisterInstance(gomock.Any()).
			AnyTimes().
			Return(true, nil)

		namingClient.
			EXPECT().
			DeregisterInstance(gomock.Any()).
			AnyTimes().
			Return(true, nil)

		nacos.NewNamingClient = func(param vo.NacosClientParam) (iClient naming_client.INamingClient, err error) {
			return namingClient, nil
		}

		if nacos.NamingClient == nil {
			nacos.NamingClient = namingClient
		}

		_ = config.GddNacosRegisterHost.Write("")

		So(func() {
			nacos.NewRest(map[string]interface{}{
				"foo": "bar",
			})
		}, ShouldNotPanic)

		nacos.ShutdownRest()

	})
}

func TestShutdownRestFail(t *testing.T) {
	Convey("Should fail", t, func() {
		setup()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buildinfo.BuildTime = "Mon Jan 2 15:04:05 MST 2006"
		namingClient := mock.NewMockINamingClient(ctrl)
		namingClient.
			EXPECT().
			RegisterInstance(gomock.Any()).
			AnyTimes().
			Return(true, nil)

		namingClient.
			EXPECT().
			DeregisterInstance(gomock.Any()).
			AnyTimes().
			Return(false, errors.New("mock test error"))

		nacos.NewNamingClient = func(param vo.NacosClientParam) (iClient naming_client.INamingClient, err error) {
			return namingClient, nil
		}

		if nacos.NamingClient == nil {
			nacos.NamingClient = namingClient
		}

		_ = config.GddNacosRegisterHost.Write("")

		So(func() {
			nacos.NewRest(map[string]interface{}{
				"foo": "bar",
			})
		}, ShouldNotPanic)

		nacos.ShutdownRest()

	})
}

func TestShutdownRestFail2(t *testing.T) {
	Convey("Should fail", t, func() {
		setup()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buildinfo.BuildTime = "Mon Jan 2 15:04:05 MST 2006"
		namingClient := mock.NewMockINamingClient(ctrl)
		namingClient.
			EXPECT().
			RegisterInstance(gomock.Any()).
			AnyTimes().
			Return(true, nil)

		namingClient.
			EXPECT().
			DeregisterInstance(gomock.Any()).
			AnyTimes().
			Return(false, nil)

		nacos.NewNamingClient = func(param vo.NacosClientParam) (iClient naming_client.INamingClient, err error) {
			return namingClient, nil
		}

		if nacos.NamingClient == nil {
			nacos.NamingClient = namingClient
		}

		_ = config.GddNacosRegisterHost.Write("")

		So(func() {
			nacos.NewRest(map[string]interface{}{
				"foo": "bar",
			})
		}, ShouldNotPanic)

		nacos.ShutdownRest()
	})
}

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

	n := nacos.NewRRServiceProvider("testsvc",
		nacos.WithNacosNamingClient(namingClient),
		nacos.WithNacosClusters([]string{"a"}))
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

	n := nacos.NewWRRServiceProvider("testsvc",
		nacos.WithNacosNamingClient(namingClient),
		nacos.WithNacosClusters([]string{"a"}))
	got := n.SelectServer()
	require.Equal(t, got, "http://10.10.10.10:80/api")
}
