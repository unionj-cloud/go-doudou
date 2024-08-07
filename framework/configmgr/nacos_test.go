package configmgr_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/unionj-cloud/go-doudou/v2/framework/config"
	"github.com/unionj-cloud/go-doudou/v2/framework/configmgr"
	"github.com/unionj-cloud/go-doudou/v2/framework/configmgr/mock"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"github.com/wubin1989/nacos-sdk-go/v2/clients/cache"
	"github.com/wubin1989/nacos-sdk-go/v2/clients/config_client"
	"github.com/wubin1989/nacos-sdk-go/v2/util"
	"github.com/wubin1989/nacos-sdk-go/v2/vo"
)

func TestMain(m *testing.M) {
	config.GddServiceName.Write("configmgr")
	config.GddApolloAddr.Write("http://apollo-config-dev-svc:8080")
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

func TestNacosConfigMgr_CallbackOnChange_Dotenv(t *testing.T) {
	Convey("Should react to dotenv config change", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		dataId := ".env"
		configClient := mock.NewMockIConfigClient(ctrl)
		nacosClient := configmgr.NewNacosConfigMgr([]string{dataId},
			config.DefaultGddNacosConfigGroup, configmgr.DotenvConfigFormat, config.DefaultGddNacosNamespaceId, configClient, cache.NewConcurrentMap())
		nacosClient.AddChangeListener(configmgr.NacosConfigListenerParam{
			DataId: dataId,
			OnChange: func(event *configmgr.NacosChangeEvent) {
				So(len(event.Changes), ShouldEqual, 4)
			},
		})
		nacosClient.CallbackOnChange(config.DefaultGddNacosNamespaceId, config.DefaultGddNacosConfigGroup, dataId, "GDD_SERVICE_NAME=configmgr\n\nGDD_READ_TIMEOUT=60s\nGDD_WRITE_TIMEOUT=60s\nGDD_IDLE_TIMEOUT=120s", "")
	})
}

func TestNacosConfigMgr_CallbackOnChange_Yaml(t *testing.T) {
	Convey("Should react to yaml config change", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		dataId := "app.yml"
		configClient := mock.NewMockIConfigClient(ctrl)
		nacosClient := configmgr.NewNacosConfigMgr([]string{dataId},
			config.DefaultGddNacosConfigGroup, configmgr.YamlConfigFormat, config.DefaultGddNacosNamespaceId, configClient, cache.NewConcurrentMap())
		nacosClient.AddChangeListener(configmgr.NacosConfigListenerParam{
			DataId: dataId,
			OnChange: func(event *configmgr.NacosChangeEvent) {
				So(len(event.Changes), ShouldEqual, 2)
			},
		})
		nacosClient.CallbackOnChange(config.DefaultGddNacosNamespaceId, config.DefaultGddNacosConfigGroup, dataId, "gdd:\n  port: 6060\n  tracing:\n    metrics:\n      root: \"go-doudou\"", "")
	})
}

func ErrReader(err error) io.Reader {
	return &errReader{err: err}
}

type errReader struct {
	err error
}

func (r *errReader) Read(p []byte) (int, error) {
	return 0, r.err
}

func TestNacosConfigMgr_CallbackOnChange_Yaml_Error(t *testing.T) {
	Convey("Should fail to react to yaml config change", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		dataId := "app.yml"
		configClient := mock.NewMockIConfigClient(ctrl)
		nacosClient := configmgr.NewNacosConfigMgr([]string{dataId},
			config.DefaultGddNacosConfigGroup, configmgr.YamlConfigFormat, config.DefaultGddNacosNamespaceId, configClient, cache.NewConcurrentMap())
		configmgr.StringReader = func(s string) io.Reader {
			return ErrReader(errors.New("mock read error"))
		}
		nacosClient.CallbackOnChange(config.DefaultGddNacosNamespaceId, config.DefaultGddNacosConfigGroup, dataId, "gdd:\n  port: 6060\n  tracing:\n    metrics:\n      root: \"go-doudou\"", "")
	})
}

func TestNacosConfigMgr_CallbackOnChange_Yaml_Error2(t *testing.T) {
	Convey("Should fail to react to yaml config change", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		dataId := "app.yml"
		configClient := mock.NewMockIConfigClient(ctrl)
		nacosClient := configmgr.NewNacosConfigMgr([]string{dataId},
			config.DefaultGddNacosConfigGroup, configmgr.YamlConfigFormat, config.DefaultGddNacosNamespaceId, configClient, cache.NewConcurrentMap())
		configmgr.StringReader = func(s string) io.Reader {
			if stringutils.IsEmpty(s) {
				return ErrReader(errors.New("mock read error"))
			} else {
				return strings.NewReader(s)
			}
		}
		nacosClient.CallbackOnChange(config.DefaultGddNacosNamespaceId, config.DefaultGddNacosConfigGroup, dataId, "gdd:\n  port: 6060\n  tracing:\n    metrics:\n      root: \"go-doudou\"", "")
	})
}

func TestNacosConfigMgr_CallbackOnChange_Dotenv_Error(t *testing.T) {
	Convey("Should fail to react to dotenv config change", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		dataId := ".env"
		configClient := mock.NewMockIConfigClient(ctrl)
		nacosClient := configmgr.NewNacosConfigMgr([]string{dataId},
			config.DefaultGddNacosConfigGroup, configmgr.DotenvConfigFormat, config.DefaultGddNacosNamespaceId, configClient, cache.NewConcurrentMap())
		configmgr.StringReader = func(s string) io.Reader {
			return ErrReader(errors.New("mock read error"))
		}
		nacosClient.CallbackOnChange(config.DefaultGddNacosNamespaceId, config.DefaultGddNacosConfigGroup, dataId, "GDD_SERVICE_NAME=configmgr\n\nGDD_READ_TIMEOUT=60s\nGDD_WRITE_TIMEOUT=60s\nGDD_IDLE_TIMEOUT=120s", "")
	})
}

func TestNacosConfigMgr_CallbackOnChange_Dotenv_Error2(t *testing.T) {
	Convey("Should fail to react to dotenv config change", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		dataId := ".env"
		configClient := mock.NewMockIConfigClient(ctrl)
		nacosClient := configmgr.NewNacosConfigMgr([]string{dataId},
			config.DefaultGddNacosConfigGroup, configmgr.DotenvConfigFormat, config.DefaultGddNacosNamespaceId, configClient, cache.NewConcurrentMap())
		configmgr.StringReader = func(s string) io.Reader {
			if stringutils.IsEmpty(s) {
				return ErrReader(errors.New("mock read error"))
			} else {
				return strings.NewReader(s)
			}
		}
		nacosClient.CallbackOnChange(config.DefaultGddNacosNamespaceId, config.DefaultGddNacosConfigGroup, dataId, "GDD_SERVICE_NAME=configmgr\n\nGDD_READ_TIMEOUT=60s\nGDD_WRITE_TIMEOUT=60s\nGDD_IDLE_TIMEOUT=120s", "")
	})
}

func TestNacosConfigMgr_LoadFromNacos_Panic(t *testing.T) {
	Convey("Should panic from listenConfig", t, func() {
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
			Return(errors.New("mock returned error"))

		configmgr.NewConfigClient = func(param vo.NacosClientParam) (iClient config_client.IConfigClient, err error) {
			return configClient, nil
		}

		if configmgr.NacosClient != nil {
			configmgr.NacosClient = configmgr.NewNacosConfigMgr([]string{dataId},
				config.DefaultGddNacosConfigGroup, configmgr.YamlConfigFormat, config.DefaultGddNacosNamespaceId, configClient, cache.NewConcurrentMap())
		}

		So(func() {
			configmgr.LoadFromNacos(config.GetNacosClientParam(), dataId, string(configmgr.YamlConfigFormat), config.DefaultGddNacosConfigGroup)
		}, ShouldPanic)
	})
}

func TestNacosConfigMgr_LoadFromNacos_UnknownFormat(t *testing.T) {
	Convey("Should return error", t, func() {
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

		unknownFormat := "Unknown format"

		if configmgr.NacosClient != nil {
			configmgr.NacosClient = configmgr.NewNacosConfigMgr([]string{dataId},
				config.DefaultGddNacosConfigGroup, "Unknown format", config.DefaultGddNacosNamespaceId, configClient, cache.NewConcurrentMap())
		}

		err := configmgr.LoadFromNacos(config.GetNacosClientParam(), dataId, unknownFormat, config.DefaultGddNacosConfigGroup)
		So(err, ShouldResemble, fmt.Errorf("[go-doudou] unknown config format: %s\n", unknownFormat))
	})
}

func TestNacosConfigMgr_LoadFromNacos_Yaml_Error(t *testing.T) {
	Convey("Should return error from GetConfig", t, func() {
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
				config.DefaultGddNacosConfigGroup, configmgr.YamlConfigFormat, config.DefaultGddNacosNamespaceId, configClient, cache.NewConcurrentMap())
		}

		err := configmgr.LoadFromNacos(config.GetNacosClientParam(), dataId, string(configmgr.YamlConfigFormat), config.DefaultGddNacosConfigGroup)
		So(err, ShouldNotBeNil)
	})
}

func TestNacosConfigMgr_LoadFromNacos_Dotenv_Error(t *testing.T) {
	Convey("Should return error from GetConfig", t, func() {
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

		err := configmgr.LoadFromNacos(config.GetNacosClientParam(), dataId, string(configmgr.DotenvConfigFormat), config.DefaultGddNacosConfigGroup)
		So(err, ShouldNotBeNil)
	})
}

func TestNacosConfigMgr_LoadFromNacos_Dotenv_Error2(t *testing.T) {
	Convey("Should return error from GetConfig", t, func() {
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

		configmgr.StringReader = func(s string) io.Reader {
			return ErrReader(errors.New("mock read error"))
		}

		err := configmgr.LoadFromNacos(config.GetNacosClientParam(), dataId, string(configmgr.DotenvConfigFormat), config.DefaultGddNacosConfigGroup)
		So(err, ShouldNotBeNil)
	})
}

func TestNacosConfigMgr_InitialiseNacosConfig(t *testing.T) {
	Convey("Should panic", t, func() {
		dataId := ".env"
		configmgr.NewConfigClient = func(param vo.NacosClientParam) (iClient config_client.IConfigClient, err error) {
			return nil, errors.New("mock error from NewConfigClient")
		}
		So(func() {
			configmgr.InitialiseNacosConfig(config.GetNacosClientParam(), dataId, string(config.DefaultGddNacosConfigFormat), config.DefaultGddNacosConfigGroup)
		}, ShouldPanic)
	})
}
