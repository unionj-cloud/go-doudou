package registry

import (
	"fmt"
	"github.com/hashicorp/memberlist"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/sliceutils"
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
	e.local.registry.members = append(e.local.registry.members, node)
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
	index, _ := sliceutils.IndexOfAny(node, e.local.registry.members)
	e.local.registry.members = append(e.local.registry.members[:index], e.local.registry.members[index+1:]...)
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
	index, _ := sliceutils.IndexOfAny(node, e.local.registry.members)
	e.local.registry.members[index] = node
	logrus.Infof("Node %s updated, supplying %s service", node.FullAddress(), mm.Meta.Service)
}
