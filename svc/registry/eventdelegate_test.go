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
						memberConf: nil,
						broadcasts: nil,
						memberlist: nil,
						lock:       sync.Mutex{},
						memberLock: sync.RWMutex{},
						members:    nil,
					},
					remote: false,
				},
			},
			args: args{
				node: &memberlist.Node{
					Name:  "",
					Addr:  "",
					Port:  0,
					Meta:  meta,
					State: 0,
					PMin:  0,
					PMax:  0,
					PCur:  0,
					DMin:  0,
					DMax:  0,
					DCur:  0,
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
						memberConf: nil,
						broadcasts: nil,
						memberlist: nil,
						lock:       sync.Mutex{},
						memberLock: sync.RWMutex{},
						members:    nil,
					},
					remote: false,
				},
			},
			args: args{
				node: &memberlist.Node{
					Name:  "",
					Addr:  "",
					Port:  0,
					Meta:  meta,
					State: 0,
					PMin:  0,
					PMax:  0,
					PCur:  0,
					DMin:  0,
					DMax:  0,
					DCur:  0,
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
						memberConf: nil,
						broadcasts: nil,
						memberlist: nil,
						lock:       sync.Mutex{},
						memberLock: sync.RWMutex{},
						members:    nil,
					},
					remote: false,
				},
			},
			args: args{
				node: &memberlist.Node{
					Name:  "",
					Addr:  "",
					Port:  0,
					Meta:  meta,
					State: 0,
					PMin:  0,
					PMax:  0,
					PCur:  0,
					DMin:  0,
					DMax:  0,
					DCur:  0,
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
