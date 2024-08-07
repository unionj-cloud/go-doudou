package config_test

import (
	"testing"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/agcache/memory"
	apolloConfig "github.com/apolloconfig/agollo/v4/env/config"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/unionj-cloud/go-doudou/v2/framework/config"
	"github.com/unionj-cloud/go-doudou/v2/framework/configmgr"
	"github.com/unionj-cloud/go-doudou/v2/framework/configmgr/mock"
	"github.com/wubin1989/nacos-sdk-go/v2/clients/cache"
	"github.com/wubin1989/nacos-sdk-go/v2/clients/config_client"
	"github.com/wubin1989/nacos-sdk-go/v2/vo"
)

func Test_envVariable_String(t *testing.T) {
	Convey("Should be on", t, func() {
		config.GddBanner.Write("on")
		So(config.GddBanner.String(), ShouldEqual, "on")
		So(config.GddBanner.Load(), ShouldEqual, "on")
	})
}

func TestMain(m *testing.M) {
	config.GddNacosServerAddr.Write("http://localhost:8848")
	m.Run()
}

func TestLoadConfigFromRemote_Apollo(t *testing.T) {
	Convey("Should not panic to load config from apollo", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		config.GddApolloAddr.Write("http://apollo-config-dev-svc:8080")
		config.GddConfigRemoteType.Write(config.ApolloConfigType)
		config.GddServiceName.Write("configmgr")
		configClient := mock.NewMockClient(ctrl)
		factory := &memory.DefaultCacheFactory{}
		cache := factory.Create()
		cache.Set("gdd.retry.count", "3", 0)
		cache.Set("gdd.weight", "5", 0)
		configClient.
			EXPECT().
			GetConfigCache(config.DefaultGddApolloNamespace).
			AnyTimes().
			Return(cache)

		configmgr.StartWithConfig = func(loadAppConfig func() (*apolloConfig.AppConfig, error)) (agollo.Client, error) {
			_, _ = loadAppConfig()
			return configClient, nil
		}

		if configmgr.ApolloClient != nil {
			configmgr.ApolloClient = configClient
		}

		So(func() {
			config.LoadConfigFromRemote()
		}, ShouldNotPanic)
	})
}

func TestLoadConfigFromRemote_Apollo_Panic(t *testing.T) {
	Convey("Should panic to load config from apollo", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		config.GddApolloAddr.Write("http://apollo-config-dev-svc:8080")
		config.GddConfigRemoteType.Write(config.ApolloConfigType)
		config.GddServiceName.Write("")
		configClient := mock.NewMockClient(ctrl)
		factory := &memory.DefaultCacheFactory{}
		cache := factory.Create()
		cache.Set("gdd.retry.count", "3", 0)
		cache.Set("gdd.weight", "5", 0)
		configClient.
			EXPECT().
			GetConfigCache(config.DefaultGddApolloNamespace).
			AnyTimes().
			Return(cache)

		configmgr.StartWithConfig = func(loadAppConfig func() (*apolloConfig.AppConfig, error)) (agollo.Client, error) {
			_, _ = loadAppConfig()
			return configClient, nil
		}

		if configmgr.ApolloClient != nil {
			configmgr.ApolloClient = configClient
		}

		So(func() {
			config.LoadConfigFromRemote()
		}, ShouldPanic)
	})
}

func TestLoadConfigFromRemote_Apollo_Panic2(t *testing.T) {
	Convey("Should panic to load config from apollo", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		config.GddConfigRemoteType.Write(config.ApolloConfigType)
		config.GddServiceName.Write("configmgr")
		config.GddApolloAddr.Write("")
		configClient := mock.NewMockClient(ctrl)
		factory := &memory.DefaultCacheFactory{}
		cache := factory.Create()
		cache.Set("gdd.retry.count", "3", 0)
		cache.Set("gdd.weight", "5", 0)
		configClient.
			EXPECT().
			GetConfigCache(config.DefaultGddApolloNamespace).
			AnyTimes().
			Return(cache)

		configmgr.StartWithConfig = func(loadAppConfig func() (*apolloConfig.AppConfig, error)) (agollo.Client, error) {
			_, _ = loadAppConfig()
			return configClient, nil
		}

		if configmgr.ApolloClient != nil {
			configmgr.ApolloClient = configClient
		}

		So(func() {
			config.LoadConfigFromRemote()
		}, ShouldPanic)
	})
}

func TestLoadConfigFromRemote_Apollo_Log(t *testing.T) {
	Convey("Should panic to load config from apollo", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		config.GddConfigRemoteType.Write(config.ApolloConfigType)
		config.GddServiceName.Write("configmgr")
		config.GddApolloAddr.Write("http://apollo-config-dev-svc:8080")
		config.GddApolloLogEnable.Write("true")
		configClient := mock.NewMockClient(ctrl)
		factory := &memory.DefaultCacheFactory{}
		cache := factory.Create()
		cache.Set("gdd.retry.count", "3", 0)
		cache.Set("gdd.weight", "5", 0)
		configClient.
			EXPECT().
			GetConfigCache(config.DefaultGddApolloNamespace).
			AnyTimes().
			Return(cache)

		configmgr.StartWithConfig = func(loadAppConfig func() (*apolloConfig.AppConfig, error)) (agollo.Client, error) {
			_, _ = loadAppConfig()
			return configClient, nil
		}

		if configmgr.ApolloClient != nil {
			configmgr.ApolloClient = configClient
		}

		So(func() {
			config.LoadConfigFromRemote()
		}, ShouldNotPanic)
	})
}

func TestLoadConfigFromRemote_Nacos(t *testing.T) {
	Convey("Should not panic to load config from nacos", t, func() {
		config.GddConfigRemoteType.Write(config.NacosConfigType)
		config.GddNacosConfigDataid.Write(".env")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		dataId := ".env"
		configClient := mock.NewMockIConfigClient(ctrl)
		configClient.
			EXPECT().
			GetConfig(vo.ConfigParam{
				DataId: dataId,
				Group:  config.DefaultGddNacosConfigGroup,
			}).
			AnyTimes().
			Return("GDD_SERVICE_NAME=configmgr\n\nGDD_READ_TIMEOUT=60s\nGDD_WRITE_TIMEOUT=60s\nGDD_IDLE_TIMEOUT=120s", nil)

		configClient.
			EXPECT().
			ListenConfig(gomock.Any()).
			AnyTimes().
			Return(nil)

		configmgr.NewConfigClient = func(param vo.NacosClientParam) (iClient config_client.IConfigClient, err error) {
			return configClient, nil
		}

		if configmgr.NacosClient != nil {
			configmgr.NacosClient = configmgr.NewNacosConfigMgr([]string{dataId},
				config.DefaultGddNacosConfigGroup, configmgr.DotenvConfigFormat, config.DefaultGddNacosNamespaceId, configClient, cache.NewConcurrentMap())
		}

		So(func() {
			config.LoadConfigFromRemote()
		}, ShouldNotPanic)
	})
}

func TestLoadConfigFromRemote_Nacos_Panic(t *testing.T) {
	Convey("Should panic to load config from nacos", t, func() {
		config.GddConfigRemoteType.Write(config.NacosConfigType)
		config.GddNacosConfigDataid.Write("")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		dataId := ".env"
		configClient := mock.NewMockIConfigClient(ctrl)
		configClient.
			EXPECT().
			GetConfig(vo.ConfigParam{
				DataId: dataId,
				Group:  config.DefaultGddNacosConfigGroup,
			}).
			AnyTimes().
			Return("GDD_SERVICE_NAME=configmgr\n\nGDD_READ_TIMEOUT=60s\nGDD_WRITE_TIMEOUT=60s\nGDD_IDLE_TIMEOUT=120s", nil)

		configClient.
			EXPECT().
			ListenConfig(gomock.Any()).
			AnyTimes().
			Return(nil)

		configmgr.NewConfigClient = func(param vo.NacosClientParam) (iClient config_client.IConfigClient, err error) {
			return configClient, nil
		}

		if configmgr.NacosClient != nil {
			configmgr.NacosClient = configmgr.NewNacosConfigMgr([]string{dataId},
				config.DefaultGddNacosConfigGroup, configmgr.DotenvConfigFormat, config.DefaultGddNacosNamespaceId, configClient, cache.NewConcurrentMap())
		}

		So(func() {
			config.LoadConfigFromRemote()
		}, ShouldPanic)
	})
}

func TestLoadConfigFromRemote_Nacos_Panic2(t *testing.T) {
	Convey("Should panic to load config from nacos", t, func() {
		config.GddConfigRemoteType.Write(config.NacosConfigType)
		config.GddNacosConfigDataid.Write(".env")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		dataId := ".env"
		configClient := mock.NewMockIConfigClient(ctrl)
		configClient.
			EXPECT().
			GetConfig(vo.ConfigParam{
				DataId: dataId,
				Group:  config.DefaultGddNacosConfigGroup,
			}).
			AnyTimes().
			Return("", errors.New("mock error from GetConfig"))

		configClient.
			EXPECT().
			ListenConfig(gomock.Any()).
			AnyTimes().
			Return(nil)

		configmgr.NewConfigClient = func(param vo.NacosClientParam) (iClient config_client.IConfigClient, err error) {
			return configClient, nil
		}

		if configmgr.NacosClient != nil {
			configmgr.NacosClient = configmgr.NewNacosConfigMgr([]string{dataId},
				config.DefaultGddNacosConfigGroup, configmgr.DotenvConfigFormat, config.DefaultGddNacosNamespaceId, configClient, cache.NewConcurrentMap())
		}

		So(func() {
			config.LoadConfigFromRemote()
		}, ShouldPanic)
	})
}

func TestLoadConfigFromRemote_Panic(t *testing.T) {
	Convey("Should panic to load config from remote", t, func() {
		config.GddConfigRemoteType.Write("Unknown remote config type")
		So(func() {
			config.LoadConfigFromRemote()
		}, ShouldPanic)
	})
}

func TestGetNacosClientParam(t *testing.T) {
	Convey("Should panic because of invalid url", t, func() {
		config.GddNacosNamespaceId.Write("test namespace")
		config.GddNacosTimeoutMs.Write("300")
		config.GddNacosNotLoadCacheAtStart.Write("true")
		config.GddNacosLogDir.Write("/tmp")
		config.GddNacosCacheDir.Write("/tmp")
		config.GddNacosLogLevel.Write("debug")
		So(func() {
			config.GetNacosClientParam()
		}, ShouldNotPanic)
	})
}

func TestGetNacosClientParam_Panic(t *testing.T) {
	Convey("Should panic because of invalid url", t, func() {
		config.GddNacosNamespaceId.Write("test namespace")
		config.GddNacosTimeoutMs.Write("300")
		config.GddNacosNotLoadCacheAtStart.Write("true")
		config.GddNacosLogDir.Write("/tmp")
		config.GddNacosCacheDir.Write("/tmp")
		config.GddNacosLogLevel.Write("debug")
		config.GddNacosServerAddr.Write("invalid url")
		So(func() {
			config.GetNacosClientParam()
		}, ShouldPanic)
	})
}

func TestGetNacosClientParam_Panic1(t *testing.T) {
	Convey("Should panic because of invalid url", t, func() {
		config.GddNacosNamespaceId.Write("test namespace")
		config.GddNacosTimeoutMs.Write("300")
		config.GddNacosNotLoadCacheAtStart.Write("true")
		config.GddNacosLogDir.Write("/tmp")
		config.GddNacosCacheDir.Write("/tmp")
		config.GddNacosLogLevel.Write("debug")
		config.GddNacosServerAddr.Write("#$@$%^&$@@")
		So(func() {
			config.GetNacosClientParam()
		}, ShouldPanic)
	})
}

func Test_envVariable_MarshalJSON(t *testing.T) {
	Convey("Should be equal to ", t, func() {
		config.GddPort.Write("8080")
		data, err := config.GddPort.MarshalJSON()
		So(err, ShouldBeNil)
		So(string(data), ShouldEqual, `"8080"`)
	})
}
