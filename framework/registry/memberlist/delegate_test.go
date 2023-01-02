package memberlist

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/memberlist"
	"sync"
	"testing"
)

func Test_delegate_NodeMeta(t *testing.T) {
	type fields struct {
	}
	type args struct {
		limit int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "",
			fields: fields{},
			args: args{
				limit: 1024,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &delegate{}
			assert.NotPanics(t, func() {
				d.NodeMeta(tt.args.limit)
			})
		})
	}
}
func Test_delegate_NodeMeta_panic(t *testing.T) {
	type fields struct {
	}
	type args struct {
		limit int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name:   "",
			fields: fields{},
			args: args{
				limit: 1,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &delegate{
				meta: NodeMeta{
					Services: []Service{
						{
							Name:          "test",
							RouteRootPath: "/api",
							Port:          6060,
						},
					},
					Weight: 8,
				},
				lock:  sync.Mutex{},
				queue: nil,
			}
			assert.Panics(t, func() {
				d.NodeMeta(tt.args.limit)
			})
		})
	}
}

func Test_delegate_NotifyMsg(t *testing.T) {
	type fields struct {
	}
	type args struct {
		msg []byte
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
				msg: []byte("this is a test msg"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &delegate{}
			d.NotifyMsg(tt.args.msg)
		})
	}
}

type testBroadcast struct {
	node   string
	msg    []byte
	notify chan struct{}
}

func (b *testBroadcast) Invalidates(other memberlist.Broadcast) bool {
	// Check if that broadcast is a memberlist type
	mb, ok := other.(*testBroadcast)
	if !ok {
		return false
	}

	// Invalidates any message about the same node
	return b.node == mb.node
}

// memberlist.NamedBroadcast optional interface
func (b *testBroadcast) Name() string {
	return b.node
}

func (b *testBroadcast) Message() []byte {
	return b.msg
}

func (b *testBroadcast) Finished() {
	select {
	case b.notify <- struct{}{}:
	default:
	}
}

func prettyPrintMessages(msgs [][]byte) []string {
	var out []string
	for _, msg := range msgs {
		out = append(out, "'"+string(msg)+"'")
	}
	return out
}

func Test_delegate_GetBroadcasts(t *testing.T) {
	q := &memberlist.TransmitLimitedQueue{RetransmitMult: 3, NumNodes: func() int { return 10 }}

	// 18 bytes per message
	q.QueueBroadcast(&testBroadcast{"test", []byte("1. this is a test."), nil})
	q.QueueBroadcast(&testBroadcast{"foo", []byte("2. this is a test."), nil})
	q.QueueBroadcast(&testBroadcast{"bar", []byte("3. this is a test."), nil})
	q.QueueBroadcast(&testBroadcast{"baz", []byte("4. this is a test."), nil})

	//
	//// 3 byte overhead, should only get 3 messages back
	//partial := q.GetBroadcasts(3, 80)
	//require.Equal(t, 3, len(partial), "missing messages: %v", prettyPrintMessages(partial))

	type fields struct {
	}
	type args struct {
		overhead int
		limit    int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   [][]byte
	}{
		{
			name:   "",
			fields: fields{},
			args: args{
				overhead: 2,
				limit:    80,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &delegate{
				queue: q,
			}
			got := d.GetBroadcasts(tt.args.overhead, tt.args.limit)
			require.Equal(t, 4, len(got), "missing messages: %v", prettyPrintMessages(got))
		})
	}
}

func Test_delegate_LocalState(t *testing.T) {
	d := delegate{}
	d.LocalState(false)
}

func Test_delegate_MergeRemoteState(t *testing.T) {
	d := delegate{}
	d.MergeRemoteState(nil, false)
}
