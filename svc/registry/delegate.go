package registry

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

type delegate struct {
	local *Node
}

// NodeMeta return user custom node meta data
func (d *delegate) NodeMeta(limit int) []byte {
	raw, _ := json.Marshal(d.local.mmeta)
	if len(raw) > limit {
		panic(fmt.Errorf("Node meta data '%v' exceeds length limit of %d bytes", d.local.mmeta, limit))
	}
	return raw
}

// NotifyMsg callback function when received user data message from remote node
func (d *delegate) NotifyMsg(msg []byte) {
	d.local.lock.Lock()
	defer d.local.lock.Unlock()
	// TODO
	logrus.Info(string(msg))
}

// GetBroadcasts get a number of user data broadcasts
func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	d.local.lock.Lock()
	defer d.local.lock.Unlock()

	msgs := d.local.broadcasts.GetBroadcasts(overhead, limit)
	return msgs
}

// LocalState also sends user data, but by tcp connection when pushPull-ing state with other node
func (d *delegate) LocalState(join bool) []byte {
	// TODO
	//d.local.lock.Lock()
	//defer d.local.lock.Unlock()
	return nil
}

// MergeRemoteState gets user data from remote node by tcp connection when pushPull-ing state with other node
func (d *delegate) MergeRemoteState(s []byte, join bool) {
	// TODO
	//d.local.lock.Lock()
	//defer d.local.lock.Unlock()
}
