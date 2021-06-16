package registry

import (
	"encoding/json"
	"fmt"
)

type registry struct {
	local *Node
}

func (r *registry) NodeMeta(limit int) []byte {
	raw, _ := json.Marshal(r.local.meta)
	if len(raw) > limit {
		panic(fmt.Errorf("Node meta data '%v' exceeds length limit of %d bytes", r.local.meta, limit))
	}
	return raw
}

func (r *registry) NotifyMsg(msg []byte) {
	r.local.lock.Lock()
	defer r.local.lock.Unlock()

	cp := make([]byte, len(msg))
	copy(cp, msg)
	// TODO
}

func (r *registry) GetBroadcasts(overhead, limit int) [][]byte {
	r.local.lock.Lock()
	defer r.local.lock.Unlock()

	msgs := r.local.broadcasts.GetBroadcasts(overhead, limit)
	return msgs
}

func (r *registry) LocalState(join bool) []byte {
	r.local.lock.Lock()
	defer r.local.lock.Unlock()

	// TODO
	return nil
}

func (r *registry) MergeRemoteState(s []byte, join bool) {
	r.local.lock.Lock()
	defer r.local.lock.Unlock()
	// TODO
}
