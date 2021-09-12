package registry

import (
	"encoding/json"
	"github.com/unionj-cloud/memberlist"
	"sync"
	"testing"
)

func Test_eventDelegate_NotifyJoin(t *testing.T) {
	mm := mergedMeta{
		Meta: nodeMeta{
			Service:       "test",
			RouteRootPath: "",
			Port:          6060,
			RegisterAt:    nil,
		},
		Data: nil,
	}
	meta, _ := json.Marshal(mm)
	type fields struct {
		local *Node
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
			name: "",
			fields: fields{
				local: &Node{
					mmeta:      mergedMeta{},
					memberNode: nil,
					registry: &registry{
						lock:       sync.Mutex{},
						memberLock: sync.RWMutex{},
					},
					remote: false,
				},
			},
			args: args{
				node: &memberlist.Node{
					Meta: meta,
				},
			},
		},
		{
			name: "",
			fields: fields{
				local: &Node{
					mmeta:      mergedMeta{},
					memberNode: nil,
					registry: &registry{
						lock:       sync.Mutex{},
						memberLock: sync.RWMutex{},
					},
					remote: false,
				},
			},
			args: args{
				node: &memberlist.Node{
					Meta: []byte("{name:"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := eventDelegate{
				local: tt.fields.local,
			}
			e.NotifyJoin(tt.args.node)
		})
	}
}

func Test_eventDelegate_NotifyLeave(t *testing.T) {
	mm := mergedMeta{
		Meta: nodeMeta{
			Service:       "test",
			RouteRootPath: "",
			Port:          6060,
			RegisterAt:    nil,
		},
		Data: nil,
	}
	meta, _ := json.Marshal(mm)
	type fields struct {
		local *Node
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
			name: "",
			fields: fields{
				local: &Node{
					mmeta:      mergedMeta{},
					memberNode: nil,
					registry: &registry{
						lock:       sync.Mutex{},
						memberLock: sync.RWMutex{},
					},
					remote: false,
				},
			},
			args: args{
				node: &memberlist.Node{
					Meta: meta,
				},
			},
		},
		{
			name: "",
			fields: fields{
				local: &Node{
					mmeta:      mergedMeta{},
					memberNode: nil,
					registry: &registry{
						lock:       sync.Mutex{},
						memberLock: sync.RWMutex{},
					},
					remote: false,
				},
			},
			args: args{
				node: &memberlist.Node{
					Meta: []byte("{name:"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := eventDelegate{
				local: tt.fields.local,
			}
			e.NotifyLeave(tt.args.node)
		})
	}
}

func Test_eventDelegate_NotifyUpdate(t *testing.T) {
	mm := mergedMeta{
		Meta: nodeMeta{
			Service:       "test",
			RouteRootPath: "",
			Port:          6060,
			RegisterAt:    nil,
		},
		Data: nil,
	}
	meta, _ := json.Marshal(mm)
	type fields struct {
		local *Node
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
			name: "",
			fields: fields{
				local: &Node{
					mmeta:      mergedMeta{},
					memberNode: nil,
					registry: &registry{
						lock:       sync.Mutex{},
						memberLock: sync.RWMutex{},
					},
					remote: false,
				},
			},
			args: args{
				node: &memberlist.Node{
					Meta: meta,
				},
			},
		},
		{
			name: "",
			fields: fields{
				local: &Node{
					mmeta:      mergedMeta{},
					memberNode: nil,
					registry: &registry{
						lock:       sync.Mutex{},
						memberLock: sync.RWMutex{},
					},
					remote: false,
				},
			},
			args: args{
				node: &memberlist.Node{
					Meta: []byte("{name:"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := eventDelegate{
				local: tt.fields.local,
			}
			e.NotifyUpdate(tt.args.node)
		})
	}
}
