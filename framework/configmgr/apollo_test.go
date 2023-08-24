package configmgr_test

import (
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/agcache/memory"
	apolloConfig "github.com/apolloconfig/agollo/v4/env/config"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/unionj-cloud/go-doudou/v2/framework/configmgr"
	"github.com/unionj-cloud/go-doudou/v2/framework/configmgr/mock"
	"github.com/unionj-cloud/go-doudou/v2/framework/config"
	"testing"
)

func TestLoadFromApollo(t *testing.T) {
	Convey("Should not have error", t, func() {
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
	})
}

func TestInitialiseApolloConfig(t *testing.T) {
	Convey("Should have error", t, func() {
		configmgr.StartWithConfig = func(loadAppConfig func() (*apolloConfig.AppConfig, error)) (agollo.Client, error) {
			return nil, errors.New("mock test error")
		}
		So(func() {
			configmgr.InitialiseApolloConfig(nil)
		}, ShouldPanic)
	})
}
