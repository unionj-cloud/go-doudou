package registry

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/unionj-cloud/go-doudou/svc/config"
	"github.com/unionj-cloud/memberlist"
	"log"
	"os"
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := seeds(tt.args.seedstr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("seeds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testConfigNet(tb testing.TB, bport int) *memberlist.Config {
	tb.Helper()

	lanConfig := memberlist.DefaultLANConfig()
	hostname, _ := os.Hostname()
	lanConfig.Name = hostname + fmt.Sprint(time.Now().UnixNano())
	lanConfig.BindPort = bport
	lanConfig.AdvertisePort = bport
	lanConfig.RequireNodeNames = true
	lanConfig.Logger = log.New(os.Stderr, lanConfig.Name, log.LstdFlags)
	return lanConfig
}

func Test_registry_Register(t *testing.T) {
	port, err := getFreePort()
	if err != nil {
		panic(err)
	}
	c1 := testConfigNet(t, port)
	m1, err := memberlist.Create(c1)
	require.NoError(t, err)
	defer m1.Shutdown()

	_ = config.GddMemSeed.Write(m1.LocalNode().Address())
	_ = config.GddServiceName.Write("testsvc")

	_, err = NewNode()
	require.NoError(t, err)
}
