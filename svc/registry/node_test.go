package registry

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unionj-cloud/go-doudou/svc/config"
	"reflect"
	"testing"
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
	_ = config.GddMemName.Write("test00")
	_ = config.GddMemHost.Write(".seed-svc-headless.default.svc.cluster.local")
	_ = config.GddMemPort.Write("56199")
	_ = config.GddMemIndirectChecks.Write("8")
	_ = config.GddLogLevel.Write("debug")
}

func setup1() {
	_ = config.GddMemSeed.Write("")
	_ = config.GddServiceName.Write("seed")
	_ = config.GddMemName.Write("seed")
	_ = config.GddMemPort.Write("56199")
	_ = config.GddMemWeight.Write("8")
	_ = config.GddMemDeadTimeout.Write("8")
	_ = config.GddMemSyncInterval.Write("8")
	_ = config.GddMemReclaimTimeout.Write("8")
	_ = config.GddMemProbeInterval.Write("8")
	_ = config.GddMemProbeTimeout.Write("8")
	_ = config.GddMemSuspicionMult.Write("8")
	_ = config.GddMemGossipNodes.Write("8")
	_ = config.GddMemGossipInterval.Write("8")
	_ = config.GddMemWeightInterval.Write("8")
	_ = config.GddMemTCPTimeout.Write("8")
	_ = config.GddMemName.Write("test00")
	_ = config.GddMemHost.Write(".seed-svc-headless.default.svc.cluster.local")
	_ = config.GddMemPort.Write("56199")
	_ = config.GddMemIndirectChecks.Write("8")
	_ = config.GddLogLevel.Write("debug")
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

func Test_registry_Register2(t *testing.T) {
	setup()
	err := NewNode()
	if err != nil {
		panic(err)
	}
	defer mlist.Shutdown()
	_ = config.GddMemSeed.Write("not exist seed")
	_ = config.GddServiceName.Write("testsvc")
	require.Error(t, NewNode())
}

func Test_join(t *testing.T) {
	setup()
	err := NewNode()
	if err != nil {
		panic(err)
	}
	defer mlist.Shutdown()
	_ = config.GddMemSeed.Write("not exist seed")
	_ = config.GddServiceName.Write("testsvc")
	require.Error(t, join())
}

func TestAllNodes(t *testing.T) {
	setup()
	err := NewNode()
	if err != nil {
		panic(err)
	}
	defer mlist.Shutdown()
	nodes, _ := AllNodes()
	if got := len(nodes); got != 1 {
		t.Errorf("got is not equal to 1, got is %d", got)
	}
}

func TestShutdown(t *testing.T) {
	setup()
	err := NewNode()
	if err != nil {
		panic(err)
	}
	defer mlist.Shutdown()
	assert.NotPanics(t, func() {
		Shutdown()
	})
}

func TestInfo(t *testing.T) {
	setup()
	err := NewNode()
	if err != nil {
		panic(err)
	}
	defer mlist.Shutdown()
	info := Info(LocalNode())
	assert.NotZero(t, info)
}

func TestMetaWeight(t *testing.T) {
	setup()
	err := NewNode()
	if err != nil {
		panic(err)
	}
	defer mlist.Shutdown()
	weight, _ := MetaWeight(LocalNode())
	assert.NotZero(t, weight)
}

func TestSvcName(t *testing.T) {
	setup()
	err := NewNode()
	if err != nil {
		panic(err)
	}
	defer mlist.Shutdown()
	assert.Equalf(t, "seed", SvcName(LocalNode()), "SvcName(%v)", LocalNode())
}

func TestRegisterServiceProvider(t *testing.T) {
	setup()
	err := NewNode()
	if err != nil {
		panic(err)
	}
	defer mlist.Shutdown()
	RegisterServiceProvider(newMockServiceProvider("TEST"))
}

func TestNewNode(t *testing.T) {
	setup1()
	err := NewNode()
	if err != nil {
		panic(err)
	}
	defer mlist.Shutdown()
}
