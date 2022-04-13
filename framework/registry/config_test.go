package registry

import (
	"github.com/unionj-cloud/go-doudou/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/framework/memberlist"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	config.GddMemDeadTimeout.Write("60")
	config.GddMemSyncInterval.Write("5")
	config.GddMemReclaimTimeout.Write("60")
	config.GddMemGossipInterval.Write("200")
	config.GddMemProbeInterval.Write("1")
	config.GddMemProbeTimeout.Write("3")
	config.GddMemRetransmitMult.Write("8")
	defer os.Clearenv()
	m.Run()
}

func Test_setGddMemDeadTimeout(t *testing.T) {
	type args struct {
		conf *memberlist.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				conf: newConf(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setGddMemDeadTimeout(tt.args.conf)
		})
	}
}

func Test_setGddMemGossipInterval(t *testing.T) {
	type args struct {
		conf *memberlist.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				conf: newConf(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setGddMemGossipInterval(tt.args.conf)
		})
	}
}

func Test_setGddMemGossipNodes(t *testing.T) {
	type args struct {
		conf *memberlist.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				conf: newConf(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setGddMemGossipNodes(tt.args.conf)
		})
	}
}

func Test_setGddMemIndirectChecks(t *testing.T) {
	type args struct {
		conf *memberlist.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				conf: newConf(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setGddMemIndirectChecks(tt.args.conf)
		})
	}
}

func Test_setGddMemProbeInterval(t *testing.T) {
	type args struct {
		conf *memberlist.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				conf: newConf(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setGddMemProbeInterval(tt.args.conf)
		})
	}
}

func Test_setGddMemProbeTimeout(t *testing.T) {
	type args struct {
		conf *memberlist.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				conf: newConf(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setGddMemProbeTimeout(tt.args.conf)
		})
	}
}

func Test_setGddMemReclaimTimeout(t *testing.T) {
	type args struct {
		conf *memberlist.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				conf: newConf(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setGddMemReclaimTimeout(tt.args.conf)
		})
	}
}

func Test_setGddMemRetransmitMult(t *testing.T) {
	type args struct {
		conf *memberlist.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				conf: newConf(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setGddMemRetransmitMult(tt.args.conf)
		})
	}
}

func Test_setGddMemSuspicionMult(t *testing.T) {
	type args struct {
		conf *memberlist.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				conf: newConf(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setGddMemSuspicionMult(tt.args.conf)
		})
	}
}

func Test_setGddMemSyncInterval(t *testing.T) {
	type args struct {
		conf *memberlist.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				conf: newConf(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setGddMemSyncInterval(tt.args.conf)
		})
	}
}

func Test_setGddMemSyncInterval_unset(t *testing.T) {
	config.GddMemSyncInterval.Write("")
	type args struct {
		conf *memberlist.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				conf: newConf(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setGddMemSyncInterval(tt.args.conf)
		})
	}
}

func Test_setGddMemReclaimTimeout_unset(t *testing.T) {
	config.GddMemReclaimTimeout.Write("")
	type args struct {
		conf *memberlist.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				conf: newConf(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setGddMemReclaimTimeout(tt.args.conf)
		})
	}
}

func Test_setGddMemGossipInterval_unset(t *testing.T) {
	config.GddMemGossipInterval.Write("")
	type args struct {
		conf *memberlist.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				conf: newConf(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setGddMemGossipInterval(tt.args.conf)
		})
	}
}

func Test_setGddMemProbeInterval_unset(t *testing.T) {
	config.GddMemProbeInterval.Write("")
	type args struct {
		conf *memberlist.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				conf: newConf(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setGddMemProbeInterval(tt.args.conf)
		})
	}
}

func Test_setGddMemProbeTimeout_unset(t *testing.T) {
	config.GddMemProbeTimeout.Write("")
	type args struct {
		conf *memberlist.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				conf: newConf(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setGddMemProbeTimeout(tt.args.conf)
		})
	}
}

func Test_setGddMemRetransmitMult_unset(t *testing.T) {
	config.GddMemRetransmitMult.Write("")
	type args struct {
		conf *memberlist.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				conf: newConf(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setGddMemRetransmitMult(tt.args.conf)
		})
	}
}
