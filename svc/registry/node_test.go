package registry

import (
	"github.com/stretchr/testify/require"
	"github.com/unionj-cloud/go-doudou/svc/config"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	_ = config.GddMemSeed.Write("")
	_ = config.GddServiceName.Write("seed")
	_ = config.GddMemName.Write("seed")
	_ = config.GddMemPort.Write("56199")
	err := NewNode()
	if err != nil {
		panic(err)
	}

	defer mlist.Shutdown()
	code := m.Run()
	os.Exit(code)
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

func Test_registry_Register1(t *testing.T) {
	require.NoError(t, join())
}

func Test_registry_Register2(t *testing.T) {
	_ = config.GddMemSeed.Write("not exist seed")
	_ = config.GddServiceName.Write("testsvc")
	require.Error(t, NewNode())
}
