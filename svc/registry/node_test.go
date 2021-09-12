package registry

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/unionj-cloud/go-doudou/svc/config"
	"github.com/unionj-cloud/memberlist"
	"reflect"
	"testing"
	"time"
)

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
			want: []string{"seed-01:56199", "seed-02:56199", "seed-03:56199"},
		},
		{
			name: "",
			args: args{
				seedstr: "",
			},
			want: nil,
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

func newSeed(port string) *Node {
	_ = config.GddMemSeed.Write("")
	_ = config.GddServiceName.Write("seed")
	_ = config.GddMemName.Write("seed" + fmt.Sprint(time.Now().UnixNano()))
	_ = config.GddMemPort.Write(port)
	seed, err := NewNode()
	if err != nil {
		panic(err)
	}
	return seed
}

func Test_registry_Register1(t *testing.T) {
	conf := memberlist.DefaultWANConfig()
	node := &Node{
		registry: &registry{
			memberConf: conf,
		},
	}
	require.Error(t, node.Register())
}

func Test_registry_Register2(t *testing.T) {
	_ = config.GddMemSeed.Write("not exist seed")
	_ = config.GddServiceName.Write("testsvc")
	_, err := NewNode()
	require.Error(t, err)
}

func TestNode_NumNodes(t *testing.T) {
	seed := newSeed("56399")
	_ = config.GddMemSeed.Write(seed.memberNode.Address())
	_ = config.GddServiceName.Write("testsvc")
	_ = config.GddMemPort.Write("56400")
	_ = config.GddMemName.Write("testsvc")
	node, _ := NewNode()
	num := node.NumNodes()
	require.Greater(t, num, 0)
}

func TestNode_Info(t *testing.T) {
	seed := newSeed("56499")
	_ = config.GddMemSeed.Write(seed.memberNode.Address())
	_ = config.GddServiceName.Write("testsvc")
	_ = config.GddMemName.Write("testnode01")
	_ = config.GddMemPort.Write("56099")
	_ = config.GddPort.Write("6060")

	node, _ := NewNode()
	require.NotNil(t, node)
}

func TestNode_String(t *testing.T) {
	seed := newSeed("56599")
	_ = config.GddMemSeed.Write(seed.memberNode.Address())
	_ = config.GddServiceName.Write("testsvc")
	_ = config.GddMemName.Write("testnode01")
	_ = config.GddMemPort.Write("56699")
	_ = config.GddPort.Write("6060")

	node, _ := NewNode()
	require.NotEmpty(t, node.String())
}

func Test_registry_Discover(t *testing.T) {
	seed := newSeed("56899")
	seed.memberNode.Meta = nil
	_ = config.GddMemSeed.Write(seed.memberNode.Address())
	_ = config.GddServiceName.Write("testsvc")
	_ = config.GddMemName.Write("testnode01")
	_ = config.GddMemPort.Write("56799")
	_ = config.GddPort.Write("6060")

	node, _ := NewNode()
	_, err := node.Discover("testsvc")
	require.Error(t, err)
}

func Test_registry_Discover2(t *testing.T) {
	_ = config.GddMemSeed.Write("")
	_ = config.GddServiceName.Write("testsvc01")
	_ = config.GddMemName.Write("testnode01")
	_ = config.GddMemPort.Write("56999")
	_ = config.GddPort.Write("6060")

	node, _ := NewNode()

	_ = config.GddMemSeed.Write(node.memberNode.Address())
	_ = config.GddServiceName.Write("testsvc02")
	_ = config.GddMemName.Write("testnode02")
	_ = config.GddMemPort.Write("57099")
	_ = config.GddPort.Write("6061")

	node, _ = NewNode()
	nodes, _ := node.Discover("testsvc01")
	require.NotEmpty(t, nodes)
}
