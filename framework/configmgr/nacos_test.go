package configmgr_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/unionj-cloud/go-doudou/framework/configmgr"
	"github.com/unionj-cloud/go-doudou/framework/configmgr/mock"
	"github.com/unionj-cloud/go-doudou/framework/internal/config"
	"github.com/wubin1989/nacos-sdk-go/clients/cache"
	"github.com/wubin1989/nacos-sdk-go/clients/config_client"
	"github.com/wubin1989/nacos-sdk-go/util"
	"github.com/wubin1989/nacos-sdk-go/vo"
	"testing"
)

func TestMain(m *testing.M) {
	config.GddNacosServerAddr.Write("http://localhost:8848")
	m.Run()
}

func TestNacosConfigMgr_AddChangeListener(t *testing.T) {
	Convey("Add listener to Nacos config client", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		configClient := mock.NewMockIConfigClient(ctrl)
		nacosClient := configmgr.NewNacosConfigMgr([]string{".env"},
			"DEFAULT_GROUP", "dotenv", "public", configClient, cache.NewConcurrentMap())
		nacosClient.AddChangeListener(configmgr.NacosConfigListenerParam{
			DataId:   ".env",
			OnChange: nil,
		})
		nacosClient.AddChangeListener(configmgr.NacosConfigListenerParam{
			DataId:   ".env",
			OnChange: nil,
		})
		listener, exists := nacosClient.Listeners().Get(util.GetConfigCacheKey(".env", "DEFAULT_GROUP", "public"))
		So(exists, ShouldBeTrue)
		So(listener, ShouldNotBeZeroValue)
	})
}

func TestNacosConfigMgr_LoadFromNacos_Dotenv(t *testing.T) {
	Convey("Should not have error", t, func() {
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

		err := configmgr.LoadFromNacos(config.GetNacosClientParam(), dataId, string(config.DefaultGddNacosConfigFormat), config.DefaultGddNacosConfigGroup)
		So(err, ShouldBeNil)
	})
}

func TestNacosConfigMgr_LoadFromNacos_Yaml(t *testing.T) {
	Convey("Should not have error", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		dataId := "app.yml"
		configClient := mock.NewMockIConfigClient(ctrl)
		configClient.
			EXPECT().
			GetConfig(vo.ConfigParam{
				DataId: dataId,
				Group:  config.DefaultGddNacosConfigGroup,
			}).
			AnyTimes().
			Return("gdd:\n  port: 8088", nil)

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
				config.DefaultGddNacosConfigGroup, configmgr.YamlConfigFormat, config.DefaultGddNacosNamespaceId, configClient, cache.NewConcurrentMap())
		}

		err := configmgr.LoadFromNacos(config.GetNacosClientParam(), dataId, string(configmgr.YamlConfigFormat), config.DefaultGddNacosConfigGroup)
		So(err, ShouldBeNil)
	})
}
