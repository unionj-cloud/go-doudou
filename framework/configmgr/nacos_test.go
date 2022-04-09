package configmgr_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/unionj-cloud/go-doudou/framework/configmgr"
	"github.com/unionj-cloud/go-doudou/framework/configmgr/mock"
	"github.com/wubin1989/nacos-sdk-go/clients/cache"
	"github.com/wubin1989/nacos-sdk-go/util"
	"testing"
)

func TestNacosConfigMgr_fetchConfig(t *testing.T) {
	Convey("Add listener to Nacos config client", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		configClient := mock.NewMockIConfigClient(ctrl)
		//configClient.
		//	EXPECT().
		//	SelectInstances(vo.SelectInstancesParam{
		//		Clusters:    []string{"a"},
		//		ServiceName: "testsvc",
		//		HealthyOnly: true,
		//	}).
		//	AnyTimes().
		//	Return(services.Hosts, nil)

		nacosClient := configmgr.NewNacosConfigMgr([]string{".env"},
			"DEFAULT_GROUP", "dotenv", "public", configClient, cache.NewConcurrentMap())
		nacosClient.AddChangeListener(configmgr.NacosConfigListenerParam{
			DataId:   ".env",
			OnChange: nil,
		})
		listener, exists := nacosClient.Listeners().Get(util.GetConfigCacheKey(".env", "DEFAULT_GROUP", "public"))
		So(exists, ShouldBeTrue)
		So(listener, ShouldNotBeZeroValue)
	})
}
