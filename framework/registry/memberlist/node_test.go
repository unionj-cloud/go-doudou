package memberlist

import (
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/require"
	"github.com/unionj-cloud/go-doudou/v2/framework/buildinfo"
	"github.com/unionj-cloud/go-doudou/v2/framework/configmgr"
	"github.com/unionj-cloud/go-doudou/v2/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/maputils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/memberlist"
	memmock "github.com/unionj-cloud/go-doudou/v2/toolkit/memberlist/mock"
	"os"
	"reflect"
	"testing"
	"time"
)

func setup() {
	_ = config.GddMemSeed.Write("")
	_ = config.GddServiceName.Write("seed")
	_ = config.GddMemName.Write("seed")
	_ = config.GddMemPort.Write("56199")
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
	_ = config.GddMemHost.Write("seed.seed-svc-headless.default.svc.cluster.local")
	_ = config.GddMemIndirectChecks.Write("8")
	_ = config.GddLogLevel.Write("debug")
	_ = config.GddPort.Write("8088")
	_ = config.GddRouteRootPath.Write("/v1")
	_ = config.GddApolloAddr.Write("http://apollo-config-dev-svc:8080")
	_ = config.GddNacosServerAddr.Write("http://localhost:8848")
}

func Test_seeds(t *testing.T) {
	type args struct {
		seedstr string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "",
			args: args{
				seedstr: "seed-01,seed-02,seed-03",
			},
			want: []string{"seed-01:7946", "seed-02:7946", "seed-03:7946"},
		},
		{
			name: "",
			args: args{
				seedstr: "",
			},
			want: nil,
		},
		{
			name: "",
			args: args{
				seedstr: "seed-01:56199,seed-02,seed-03:03,seed-04:abc",
			},
			want: []string{"seed-01:56199", "seed-02:7946", "seed-03:3", "seed-04:7946"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := seeds(tt.args.seedstr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("seeds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllNodes(t *testing.T) {
	setup()
	NewRest()
	defer Shutdown()
	Convey("There should be only one node", t, func() {
		nodes, _ := AllNodes()
		So(len(nodes), ShouldEqual, 1)
	})
}

func TestAllNodesError(t *testing.T) {
	Convey("There should be only one node", t, func() {
		_, err := AllNodes()
		So(err.Error(), ShouldEqual, "mlist is nil")
	})
}

func TestRegisterServiceProvider(t *testing.T) {
	setup()
	NewRest()
	defer Shutdown()
	Convey("", t, func() {
		provider := newMockServiceProvider("TEST")
		RegisterServiceProvider(provider)
		So(len(events.ServiceProviders), ShouldEqual, 1)
		So(len(provider.servers), ShouldEqual, 1)
	})
}

func Test_memConfigListener_OnChange(t *testing.T) {
	Convey("Test OnChange callback", t, func() {
		c := &memConfigListener{
			memConf: &memberlist.Config{},
		}
		c.OnChange(&storage.ChangeEvent{
			Changes: map[string]*storage.ConfigChange{
				"gdd.mem.dead.timeout": {
					OldValue:   "60s",
					NewValue:   "30s",
					ChangeType: storage.MODIFIED,
				},
			},
		})
		Convey("Should equal to 8s", func() {
			So(os.Getenv("GDD_MEM_DEAD_TIMEOUT"), ShouldEqual, "8s")
		})

		c.OnChange(&storage.ChangeEvent{
			Changes: map[string]*storage.ConfigChange{
				"gdd.mem.dead.timeout": {
					OldValue:   "8s",
					NewValue:   "30s",
					ChangeType: storage.MODIFIED,
				},
			},
		})
		Convey("Should equal to 30s", func() {
			So(os.Getenv("GDD_MEM_DEAD_TIMEOUT"), ShouldEqual, "30s")
			So(c.memConf.GossipToTheDeadTime, ShouldEqual, 30*time.Second)
		})

		c.OnChange(&storage.ChangeEvent{
			Changes: map[string]*storage.ConfigChange{
				"gdd.mem.dead.timeout": {
					OldValue:   "30s",
					NewValue:   "",
					ChangeType: storage.DELETED,
				},
			},
		})
		Convey("Should equal to 60s", func() {
			So(os.Getenv("GDD_MEM_DEAD_TIMEOUT"), ShouldEqual, "")
			So(c.memConf.GossipToTheDeadTime, ShouldEqual, 60*time.Second)
		})

	})
}

func TestNewRest_InvalidConfigType(t *testing.T) {
	Convey("Should panic", t, func() {
		setup()
		_ = config.GddConfigRemoteType.Write("invalid")
		defer Shutdown()
		So(func() {
			NewRest()
		}, ShouldPanic)
	})
}

func TestCallbackOnChange(t *testing.T) {
	Convey("Should equal to 30s", t, func() {
		listener := &memConfigListener{
			memConf: &memberlist.Config{},
		}
		CallbackOnChange(listener)(&configmgr.NacosChangeEvent{
			Namespace: config.DefaultGddNacosNamespaceId,
			Group:     config.DefaultGddNacosConfigGroup,
			DataId:    ".env",
			Changes: map[string]maputils.Change{
				"gdd.mem.dead.timeout": {
					OldValue:   "8s",
					NewValue:   "30s",
					ChangeType: maputils.MODIFIED,
				},
			},
		})
		CallbackOnChange(listener)(&configmgr.NacosChangeEvent{
			Namespace: config.DefaultGddNacosNamespaceId,
			Group:     config.DefaultGddNacosConfigGroup,
			DataId:    ".env",
			Changes: map[string]maputils.Change{
				"gdd.mem.dead.timeout": {
					OldValue:   "8s",
					NewValue:   "30s",
					ChangeType: maputils.MODIFIED,
				},
			},
		})
		So(os.Getenv("GDD_MEM_DEAD_TIMEOUT"), ShouldEqual, "30s")
		So(listener.memConf.GossipToTheDeadTime, ShouldEqual, 30*time.Second)
	})
}

func Test_join_fail1(t *testing.T) {
	mlist = nil
	require.Error(t, join())
}

func Test_join(t *testing.T) {
	Convey("Join should be successful", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mem := memmock.NewMockIMemberlist(ctrl)
		mlist = mem
		seedAddr := "seed:7946"
		config.GddMemSeed.Write(seedAddr)
		mem.
			EXPECT().
			Join([]string{seedAddr}).
			AnyTimes().
			Return(0, nil)

		mem.
			EXPECT().
			LocalNode().
			AnyTimes().
			Return(&memberlist.Node{
				Name: "",
				Addr: "test",
				Port: 7946,
			})

		So(join(), ShouldBeNil)
	})
}

func Test_newConf_GddMemCIDRsAllowed(t *testing.T) {
	config.GddMemCIDRsAllowed.Write("172.28.0.0/16")
	newConf()
}

func Test_newConf_GddMemCIDRsAllowed_error(t *testing.T) {
	config.GddMemCIDRsAllowed.Write("invalid")
	newConf()
}

func Test_newConf(t *testing.T) {
	defer os.Clearenv()
	config.GddLogLevel.Write("ERROR")
	newConf()

	config.GddLogLevel.Write("WARNING")
	newConf()

	config.GddMemLogDisable.Write("true")
	newConf()

	config.GddWeight.Write("5")
	newConf()

	config.GddWeight.Write("0")
	newConf()

	config.GddWeight.Write("0")
	config.GddMemWeightInterval.Write("200")
	newConf()

	config.GddMemTCPTimeout.Write("10")
	newConf()

	config.GddMemHost.Write(".seed-headless")
	newConf()
}

func Test_newNode(t *testing.T) {
	Convey("Test newNode", t, func() {
		Convey("Should return error as join failed", func() {
			setup()
			buildinfo.BuildTime = "Mon Jan 2 15:04:05 MST 2006"
			config.GddWeight.Write("8")
			config.GddMemSeed.Write("invalid seed")

			Convey("Should return error as memberlist create failed", func() {
				defer func() {
					createMemberlist = memberlist.Create
				}()
				createMemberlist = func(conf *memberlist.Config) (*memberlist.Memberlist, error) {
					return nil, errors.New("mock test error")
				}
			})
		})

		Convey("Should return nil", func() {
			setup()
			So(numNodes(), ShouldEqual, 1)
			So(retransmitMultGetter(), ShouldEqual, 4)
		})
	})
}

func Test_numNodes(t *testing.T) {
	Convey("Should return 0", t, func() {
		mlist = nil
		So(numNodes(), ShouldEqual, 0)
	})
}

func TestShutdownRest(t *testing.T) {
	config.GddServiceDiscoveryMode.Write("invalid")
	Shutdown()
}

func TestLeave(t *testing.T) {
	Convey("Should leave", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mem := memmock.NewMockIMemberlist(ctrl)
		mlist = mem
		mem.
			EXPECT().
			Leave(10 * time.Second).
			AnyTimes().
			Return(nil)

		Leave(10 * time.Second)
	})
}
