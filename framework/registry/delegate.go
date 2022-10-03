package registry

import (
	"bytes"
	"github.com/hashicorp/go-msgpack/codec"
	"github.com/unionj-cloud/go-doudou/framework/memberlist"
	logger "github.com/unionj-cloud/go-doudou/toolkit/zlogger"
	"sync"
)

type delegate struct {
	mmeta mergedMeta
	lock  sync.Mutex
	queue *memberlist.TransmitLimitedQueue
}

// NodeMeta return user custom node meta data
func (d *delegate) NodeMeta(limit int) []byte {
	d.lock.Lock()
	defer d.lock.Unlock()

	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.MsgpackHandle{})
	if err := enc.Encode(d.mmeta); err != nil {
		logger.Panic().Err(err).Msg("[go-doudou] Failed to encode node meta data")
	}
	raw := buf.Bytes()

	if len(raw) > limit {
		logger.Panic().Msgf("[go-doudou] Node meta data '%v' exceeds length limit of %d bytes", d.mmeta, limit)
	}
	return raw
}

// NotifyMsg callback function when received user data message from remote node
func (d *delegate) NotifyMsg(msg []byte) {
	d.lock.Lock()
	defer d.lock.Unlock()
	// TODO
}

// GetBroadcasts get a number of user data broadcasts
func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	d.lock.Lock()
	defer d.lock.Unlock()

	msgs := d.queue.GetBroadcasts(overhead, limit)
	return msgs
}

// LocalState also sends user data, but by tcp connection when pushPull-ing state with other node
func (d *delegate) LocalState(join bool) []byte {
	return nil
}

// MergeRemoteState gets user data from remote node by tcp connection when pushPull-ing state with other node
func (d *delegate) MergeRemoteState(s []byte, join bool) {
}
