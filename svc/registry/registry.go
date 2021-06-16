package registry

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/memberlist"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"os"
)

type IRegistry interface {
	Register() error
	Discovery(svc string) ([]*Node, error)
}

type registry struct {
	local *Node
}

func (r *registry) Register() error {
	seed := os.Getenv("SEED")
	if stringutils.IsEmpty(seed) {
		return errors.New("No seed found, register failed")
	}
	_, err := r.local.memberlist.Join([]string{seed})
	if err != nil {
		return errors.Wrap(err, "Failed to join cluster")
	}
	return nil
}

func (r *registry) Discovery(svc string) ([]*Node, error) {
	var nodes []*Node
	for _, member := range r.local.memberlist.Members() {
		logrus.Infof("Member: %s %s\n", member.Name, member.Addr)
		if member.State == memberlist.StateAlive {
			var nmeta NodeMeta
			if err := json.Unmarshal(member.Meta, &nmeta); err != nil {
				return nil, errors.Wrap(err, "")
			}
			if nmeta.Service == svc {
				nodes = append(nodes, &Node{
					meta:       nmeta,
					memberNode: member,
				})
			}
		}
	}
	return nodes, nil
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
