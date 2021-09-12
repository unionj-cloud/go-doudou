package registry

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

type delegate struct {
	local *Node
}

func (d *delegate) NodeMeta(limit int) []byte {
	raw, _ := json.Marshal(d.local.mmeta)
	if len(raw) > limit {
		panic(fmt.Errorf("Node meta data '%v' exceeds length limit of %d bytes", d.local.mmeta, limit))
	}
	return raw
}

func (d *delegate) NotifyMsg(msg []byte) {
	d.local.lock.Lock()
	defer d.local.lock.Unlock()
	// TODO
	logrus.Info(string(msg))
}

func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	d.local.lock.Lock()
	defer d.local.lock.Unlock()

	msgs := d.local.broadcasts.GetBroadcasts(overhead, limit)
	return msgs
}

func (d *delegate) LocalState(join bool) []byte {
	// TODO
	//d.local.lock.Lock()
	//defer d.local.lock.Unlock()
	return nil
}

func (d *delegate) MergeRemoteState(s []byte, join bool) {
	// TODO
	//d.local.lock.Lock()
	//defer d.local.lock.Unlock()
}
