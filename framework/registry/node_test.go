package registry

import (
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/agcache/memory"
	apolloConfig "github.com/apolloconfig/agollo/v4/env/config"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/unionj-cloud/go-doudou/v2/framework/configmgr"
	"github.com/unionj-cloud/go-doudou/v2/framework/configmgr/mock"
	"github.com/unionj-cloud/go-doudou/v2/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry/nacos"
	nmock "github.com/unionj-cloud/go-doudou/v2/framework/registry/nacos/mock"
	"github.com/wubin1989/nacos-sdk-go/v2/clients/cache"
	"github.com/wubin1989/nacos-sdk-go/v2/clients/config_client"
	"github.com/wubin1989/nacos-sdk-go/v2/clients/naming_client"
	"github.com/wubin1989/nacos-sdk-go/v2/vo"
	"testing"
)

func setup() {
	_ = config.GddServiceName.Write("seed")
	_ = config.GddLogLevel.Write("debug")
	_ = config.GddPort.Write("8088")
	_ = config.GddRouteRootPath.Write("/v1")
	_ = config.GddApolloAddr.Write("http://apollo-config-dev-svc:8080")
	_ = config.GddNacosServerAddr.Write("http://localhost:8848")
}

func TestNewNode_NacosConfigType(t *testing.T) {
	Convey("Should not have error", t, func() {
		setup()
		_ = config.GddConfigRemoteType.Write("nacos")
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
			Return("GDD_READ_TIMEOUT=60s\nGDD_WRITE_TIMEOUT=60s\nGDD_IDLE_TIMEOUT=120s", nil)

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

		err := configmgr.LoadFromNacos(config.GetNacosClientParam(), dataId, string(config.DefaultGddNacosConfigFormat), config.DefaultGddNacosConfigGroup)
		So(err, ShouldBeNil)
		So(func() {
			NewRest()
		}, ShouldNotPanic)
		defer ShutdownRest()
	})
}

func TestNewNode_ApolloConfigType(t *testing.T) {
	Convey("Should not have error", t, func() {
		setup()
		_ = config.GddConfigRemoteType.Write("apollo")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
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

		configClient.
			EXPECT().
			AddChangeListener(gomock.Any()).
			AnyTimes().
			Return()

		configmgr.StartWithConfig = func(loadAppConfig func() (*apolloConfig.AppConfig, error)) (agollo.Client, error) {
			_, _ = loadAppConfig()
			return configClient, nil
		}

		if configmgr.ApolloClient != nil {
			configmgr.ApolloClient = configClient
		}

		apolloCluster := config.DefaultGddApolloCluster
		apolloAddr := config.GddApolloAddr.Load()
		apolloNamespace := config.DefaultGddApolloNamespace
		apolloBackupPath := config.DefaultGddApolloBackupPath
		c := &apolloConfig.AppConfig{
			AppID:            config.GddServiceName.Load(),
			Cluster:          apolloCluster,
			IP:               apolloAddr,
			NamespaceName:    apolloNamespace,
			IsBackupConfig:   false,
			BackupConfigPath: apolloBackupPath,
			MustStart:        false,
		}
		So(func() {
			configmgr.LoadFromApollo(c)
		}, ShouldNotPanic)

		So(func() {
			NewRest()
		}, ShouldNotPanic)
		defer ShutdownRest()
	})
}

func TestNewNode_Nacos(t *testing.T) {
	Convey("Should return nil", t, func() {
		setup()
		_ = config.GddServiceDiscoveryMode.Write("nacos")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		defer ShutdownRest()

		namingClient := nmock.NewMockINamingClient(ctrl)
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
			NewRest()
		}, ShouldNotPanic)
	})
}

func TestNewNode_InvalidServiceDiscoveryMode(t *testing.T) {
	Convey("Should return nil", t, func() {
		setup()
		_ = config.GddServiceDiscoveryMode.Write("invalid")
		So(func() {
			NewRest()
		}, ShouldNotPanic)
	})
}
