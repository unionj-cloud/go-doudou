package registry

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/memberlist"
)

type eventDelegate struct {
	local *Node
}

// NotifyJoin callback function when node joined
func (e eventDelegate) NotifyJoin(node *memberlist.Node) {
	var (
		mm  mergedMeta
		err error
	)
	if mm, err = newMeta(node); err != nil {
		logrus.Errorln(fmt.Sprintf("%+v", err))
		return
	}
	logrus.Infof("Node %s joined, supplying %s service", node.String(), mm.Meta.Service)
}

// NotifyLeave callback function when node leave
func (e eventDelegate) NotifyLeave(node *memberlist.Node) {
	var (
		mm  mergedMeta
		err error
	)
	if mm, err = newMeta(node); err != nil {
		logrus.Errorln(fmt.Sprintf("%+v", err))
		return
	}
	logrus.Infof("Node %s left, supplying %s service", node.FullAddress(), mm.Meta.Service)
}

// NotifyUpdate callback function when node updated
func (e eventDelegate) NotifyUpdate(node *memberlist.Node) {
	var (
		mm  mergedMeta
		err error
	)
	if mm, err = newMeta(node); err != nil {
		logrus.Errorln(fmt.Sprintf("%+v", err))
		return
	}
	logrus.Infof("Node %s updated, supplying %s service", node.FullAddress(), mm.Meta.Service)
}
