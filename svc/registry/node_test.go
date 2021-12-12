package registry

//
//import (
//	"github.com/stretchr/testify/require"
//	"github.com/unionj-cloud/go-doudou/svc/config"
//	"os"
//	"reflect"
//	"testing"
//)
//
//var seed *Node
//
//func TestMain(m *testing.M) {
//	_ = config.GddMemSeed.Write("")
//	_ = config.GddServiceName.Write("seed")
//	_ = config.GddMemName.Write("seed")
//	_ = config.GddMemPort.Write("56199")
//	var err error
//	seed, err = NewNode()
//	if err != nil {
//		panic(err)
//	}
//	defer seed.memberlist.Shutdown()
//	code := m.Run()
//	os.Exit(code)
//}
//
//func Test_seeds(t *testing.T) {
//	type args struct {
//		seedstr string
//	}
//	tests := []struct {
//		name string
//		args args
//		want []string
//	}{
//		{
//			name: "",
//			args: args{
//				seedstr: "seed-01,seed-02,seed-03",
//			},
//			want: []string{"seed-01:56199", "seed-02:56199", "seed-03:56199"},
//		},
//		{
//			name: "",
//			args: args{
//				seedstr: "",
//			},
//			want: nil,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := seeds(tt.args.seedstr); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("seeds() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_registry_Register1(t *testing.T) {
//	node := &Node{}
//	require.Error(t, node.Join())
//}
//
//func Test_registry_Register2(t *testing.T) {
//	_ = config.GddMemSeed.Write("not exist seed")
//	_ = config.GddServiceName.Write("testsvc")
//	_, err := NewNode()
//	require.Error(t, err)
//}
//
//func TestNode_NumNodes(t *testing.T) {
//	_ = config.GddMemSeed.Write(seed.memberNode.Address())
//	_ = config.GddServiceName.Write("testsvc_numnodes")
//	_ = config.GddMemPort.Write("56400")
//	_ = config.GddMemName.Write("testsvc_numnodes")
//	node, _ := NewNode()
//	defer node.memberlist.Shutdown()
//	num := node.NumNodes()
//	require.Greater(t, num, 0)
//}
//
//func TestNode_Info(t *testing.T) {
//	_ = config.GddMemSeed.Write(seed.memberNode.Address())
//	_ = config.GddServiceName.Write("testsvc_info")
//	_ = config.GddMemName.Write("testnode_info")
//	_ = config.GddMemPort.Write("56099")
//	_ = config.GddPort.Write("6060")
//
//	node, _ := NewNode()
//	defer node.memberlist.Shutdown()
//	require.NotNil(t, node.Info())
//}
//
//func TestNode_String(t *testing.T) {
//	_ = config.GddMemSeed.Write(seed.memberNode.Address())
//	_ = config.GddServiceName.Write("testsvc_string")
//	_ = config.GddMemName.Write("testnode_string")
//	_ = config.GddMemPort.Write("56699")
//	_ = config.GddPort.Write("6060")
//
//	node, _ := NewNode()
//	defer node.memberlist.Shutdown()
//	require.NotEmpty(t, node.String())
//}
//
//func Test_registry_AllNodes(t *testing.T) {
//	_ = config.GddMemSeed.Write(seed.memberNode.Address())
//	_ = config.GddServiceName.Write("testsvc_discover1")
//	_ = config.GddMemName.Write("testnode_discover1")
//	_ = config.GddMemPort.Write("56999")
//	_ = config.GddPort.Write("6060")
//	_ = config.GddMemHost.Write(".testsvc_discover1-svc-headless.default.svc.cluster.local")
//	node, _ := NewNode()
//	defer node.memberlist.Shutdown()
//
//	_ = config.GddMemSeed.Write(seed.memberNode.Address())
//	_ = config.GddServiceName.Write("testsvc_discover2")
//	_ = config.GddMemName.Write("testnode_discover2")
//	_ = config.GddMemPort.Write("57099")
//	_ = config.GddPort.Write("6061")
//
//	node, _ = NewNode(WithData("a message"))
//	defer node.memberlist.Shutdown()
//	nodes, _ := node.AllNodes()
//	require.NotEmpty(t, nodes)
//}
