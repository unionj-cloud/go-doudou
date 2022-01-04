package registry

import (
	"encoding/json"
	"fmt"
	"github.com/unionj-cloud/memberlist"
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
	raw, _ := json.Marshal(d.mmeta)
	if len(raw) > limit {
		panic(fmt.Errorf("[go-doudou] Node meta data '%v' exceeds length limit of %d bytes", d.mmeta, limit))
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
	d.lock.Lock()
	defer d.lock.Unlock()
	// TODO
	return nil
}

// MergeRemoteState gets user data from remote node by tcp connection when pushPull-ing state with other node
func (d *delegate) MergeRemoteState(s []byte, join bool) {
	d.lock.Lock()
	defer d.lock.Unlock()
	// TODO
}
