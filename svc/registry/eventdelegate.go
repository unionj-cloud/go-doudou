package registry

import (
	"fmt"
	"github.com/hashicorp/memberlist"
	"github.com/sirupsen/logrus"
)

type eventDelegate struct {
	local *Node
}

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
