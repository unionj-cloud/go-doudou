package memberlist

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/memberlist"
	"testing"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Test_eventDelegate_NotifyJoin(t *testing.T) {
	mm := NodeMeta{
		Services: []Service{
			{
				Name:          "test",
				RouteRootPath: "/api",
				Port:          6060,
			},
		},
		Weight: 8,
	}
	meta, _ := json.Marshal(mm)
	type fields struct {
	}
	type args struct {
		node *memberlist.Node
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "",
			fields: fields{},
			args: args{
				node: &memberlist.Node{
					Meta: meta,
				},
			},
		},
		{
			name:   "",
			fields: fields{},
			args: args{
				node: &memberlist.Node{
					Meta: []byte("{name:"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := eventDelegate{}
			e.NotifyJoin(tt.args.node)
		})
	}
}

func Test_eventDelegate_NotifyLeave(t *testing.T) {
	mm := NodeMeta{
		Services: []Service{
			{
				Name:          "test",
				RouteRootPath: "/api",
				Port:          6060,
			},
		},
		Weight: 8,
	}
	meta, _ := json.Marshal(mm)
	type fields struct {
	}
	type args struct {
		node *memberlist.Node
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "",
			fields: fields{},
			args: args{
				node: &memberlist.Node{
					Meta: meta,
				},
			},
		},
		{
			name:   "",
			fields: fields{},
			args: args{
				node: &memberlist.Node{
					Meta: []byte("{name:"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := eventDelegate{}
			e.NotifyLeave(tt.args.node)
		})
	}
}

func Test_eventDelegate_NotifyUpdate(t *testing.T) {
	mm := NodeMeta{
		Services: []Service{
			{
				Name:          "test",
				RouteRootPath: "/api",
				Port:          6060,
			},
		},
		Weight: 8,
	}
	meta, _ := json.Marshal(mm)
	type fields struct {
	}
	type args struct {
		node *memberlist.Node
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "",
			fields: fields{},
			args: args{
				node: &memberlist.Node{
					Meta: meta,
				},
			},
		},
		{
			name:   "",
			fields: fields{},
			args: args{
				node: &memberlist.Node{
					Meta: []byte("{name:"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := eventDelegate{}
			e.NotifyUpdate(tt.args.node)
		})
	}
}

func Test_eventDelegate_NotifySuspectSateChange(t *testing.T) {
	type fields struct {
		ServiceProviders []IMemberlistServiceProvider
	}
	type args struct {
		node *memberlist.Node
	}
	sp := newMockServiceProvider("TEST")
	providers := []IMemberlistServiceProvider{
		sp,
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "",
			fields: fields{
				ServiceProviders: providers,
			},
			args: args{
				node: &memberlist.Node{
					Name:  "test01",
					Addr:  "192.168.1.103",
					Port:  56199,
					State: memberlist.StateAlive,
				},
			},
			want: 1,
		},
		{
			name: "",
			fields: fields{
				ServiceProviders: providers,
			},
			args: args{
				node: &memberlist.Node{
					Name:  "test01",
					Addr:  "192.168.1.103",
					Port:  56199,
					State: memberlist.StateSuspect,
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &eventDelegate{
				ServiceProviders: tt.fields.ServiceProviders,
			}
			e.NotifySuspectSateChange(tt.args.node)
			if got := len(sp.servers); got != tt.want {
				t.Errorf("expected: %d, actual: %d", got, tt.want)
			}
		})
	}
}

func Test_eventDelegate_NotifyWeight(t *testing.T) {
	type fields struct {
		ServiceProviders []IMemberlistServiceProvider
	}
	type args struct {
		node *memberlist.Node
	}
	sp := newMockServiceProvider("TEST")
	providers := []IMemberlistServiceProvider{
		sp,
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "",
			fields: fields{
				ServiceProviders: providers,
			},
			args: args{
				node: &memberlist.Node{
					Name:   "test01",
					Addr:   "192.168.1.103",
					Port:   56199,
					State:  memberlist.StateSuspect,
					Weight: 8,
				},
			},
			want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &eventDelegate{
				ServiceProviders: tt.fields.ServiceProviders,
			}
			e.NotifyJoin(tt.args.node)
			e.NotifyWeight(tt.args.node)
			if got := sp.serverMap["test01"].Weight; got != tt.want {
				t.Errorf("expected: %d, actual: %d", got, tt.want)
			}
		})
	}
}

func Test_eventDelegate_NotifyLeave1(t *testing.T) {
	type fields struct {
		ServiceProviders []IMemberlistServiceProvider
	}
	type args struct {
		node *memberlist.Node
	}
	sp := newMockServiceProvider("TEST")
	providers := []IMemberlistServiceProvider{
		sp,
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "",
			fields: fields{
				ServiceProviders: providers,
			},
			args: args{
				node: &memberlist.Node{
					Name:   "test01",
					Addr:   "192.168.1.103",
					Port:   56199,
					State:  memberlist.StateSuspect,
					Weight: 8,
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &eventDelegate{
				ServiceProviders: tt.fields.ServiceProviders,
			}
			e.NotifyJoin(tt.args.node)
			e.NotifyLeave(tt.args.node)
			if got := len(sp.servers); got != tt.want {
				t.Errorf("expected: %d, actual: %d", got, tt.want)
			}
		})
	}
}

func Test_eventDelegate_NotifyUpdate1(t *testing.T) {
	type fields struct {
		ServiceProviders []IMemberlistServiceProvider
	}
	type args struct {
		node *memberlist.Node
	}
	sp := newMockServiceProvider("TEST")
	providers := []IMemberlistServiceProvider{
		sp,
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "",
			fields: fields{
				ServiceProviders: providers,
			},
			args: args{
				node: &memberlist.Node{
					Name:   "test01",
					Addr:   "192.168.1.103",
					Port:   56199,
					State:  memberlist.StateSuspect,
					Weight: 8,
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &eventDelegate{
				ServiceProviders: tt.fields.ServiceProviders,
			}
			e.NotifyUpdate(tt.args.node)
			if got := len(sp.servers); got != tt.want {
				t.Errorf("expected: %d, actual: %d", got, tt.want)
			}
		})
	}
}
