package nacos_test

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/unionj-cloud/go-doudou/framework/buildinfo"
	"github.com/unionj-cloud/go-doudou/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/framework/registry/nacos"
	"github.com/unionj-cloud/go-doudou/framework/registry/nacos/mock"
	"github.com/wubin1989/nacos-sdk-go/clients/naming_client"
	"github.com/wubin1989/nacos-sdk-go/vo"
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

func TestNewNode(t *testing.T) {
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
			nacos.NewNode(map[string]interface{}{
				"foo": "bar",
			})
		}, ShouldNotPanic)

		nacos.Shutdown()
		So(nacos.NamingClient, ShouldBeNil)
	})
}

func TestNewNode2(t *testing.T) {
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
			nacos.NewNode(map[string]interface{}{
				"foo": "bar",
			})
		}, ShouldNotPanic)

		nacos.Shutdown()
		So(nacos.NamingClient, ShouldBeNil)
	})
}

func TestNewNodePanic(t *testing.T) {
	Convey("Should panic as service name not set", t, func() {
		setup()
		_ = config.GddServiceName.Write("")

		So(func() {
			nacos.NewNode()
		}, ShouldPanic)
	})
}

func TestNewNodePanic2(t *testing.T) {
	Convey("Should panic as register failed", t, func() {
		setup()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		namingClient := mock.NewMockINamingClient(ctrl)
		namingClient.
			EXPECT().
			RegisterInstance(gomock.Any()).
			AnyTimes().
			Return(false, errors.New("test error"))

		nacos.NewNamingClient = func(param vo.NacosClientParam) (iClient naming_client.INamingClient, err error) {
			return namingClient, nil
		}

		if nacos.NamingClient == nil {
			nacos.NamingClient = namingClient
		}

		So(func() {
			nacos.NewNode()
		}, ShouldPanic)
	})
}

func TestInitialiseNacosNamingClient(t *testing.T) {
	Convey("Should panic as service name not set", t, func() {
		nacos.NewNamingClient = func(param vo.NacosClientParam) (iClient naming_client.INamingClient, err error) {
			return nil, errors.New("test error")
		}
		So(func() {
			nacos.InitialiseNacosNamingClient()
		}, ShouldPanic)
	})
}

func TestShutdownFail(t *testing.T) {
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
			nacos.NewNode(map[string]interface{}{
				"foo": "bar",
			})
		}, ShouldNotPanic)

		nacos.Shutdown()
		So(nacos.NamingClient, ShouldBeNil)
	})
}

func TestShutdownFail2(t *testing.T) {
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
			nacos.NewNode(map[string]interface{}{
				"foo": "bar",
			})
		}, ShouldNotPanic)

		nacos.Shutdown()
		So(nacos.NamingClient, ShouldBeNil)
	})
}

func TestNewNodePanic3(t *testing.T) {
	Convey("Should panic as fail to get private ip", t, func() {
		setup()
		_ = config.GddNacosRegisterHost.Write("")
		nacos.GetPrivateIP = func() (string, error) {
			return "", errors.New("mock test error")
		}
		So(func() {
			nacos.NewNode(map[string]interface{}{
				"foo": "bar",
			})
		}, ShouldPanic)
	})
}

func TestNewNodePanic4(t *testing.T) {
	Convey("Should panic as register host is empty", t, func() {
		setup()
		_ = config.GddNacosRegisterHost.Write("")
		nacos.GetPrivateIP = func() (string, error) {
			return "", nil
		}
		So(func() {
			nacos.NewNode(map[string]interface{}{
				"foo": "bar",
			})
		}, ShouldPanic)
	})
}
