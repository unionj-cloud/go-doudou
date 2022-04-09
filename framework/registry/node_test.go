package registry

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/require"
	"github.com/unionj-cloud/go-doudou/framework/internal/config"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	err := NewNode()
	if err != nil {
		panic(err)
	}
	defer Shutdown()
	m.Run()
}

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

func Test_join(t *testing.T) {
	_ = config.GddMemSeed.Write("not exist seed")
	_ = config.GddServiceName.Write("testsvc")
	require.Error(t, join())
}

func TestAllNodes(t *testing.T) {
	Convey("There should be only one node", t, func() {
		nodes, _ := AllNodes()
		So(len(nodes), ShouldEqual, 1)
	})
}

func TestInfo(t *testing.T) {
	Convey("Should not zero value", t, func() {
		info := Info(LocalNode())
		So(info, ShouldNotBeZeroValue)
	})
}

func TestMetaWeight(t *testing.T) {
	Convey("Should not zero value", t, func() {
		weight, _ := MetaWeight(LocalNode())
		So(weight, ShouldNotBeZeroValue)
	})
}

func TestSvcName(t *testing.T) {
	Convey("Should be equal to seed", t, func() {
		So(SvcName(LocalNode()), ShouldEqual, "seed")
	})
}

func TestRegisterServiceProvider(t *testing.T) {
	Convey("", t, func() {
		provider := newMockServiceProvider("TEST")
		RegisterServiceProvider(provider)
		So(len(events.ServiceProviders), ShouldEqual, 1)
		So(len(provider.servers), ShouldEqual, 1)
	})
}
