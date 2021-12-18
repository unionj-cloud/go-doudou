package registry

import (
	"github.com/unionj-cloud/memberlist"
)

type eventDelegate struct {
	ServiceProviders []IServiceProvider
}

func (e *eventDelegate) NotifyWeight(node *memberlist.Node) {
}

// NotifyJoin callback function when node joined
func (e *eventDelegate) NotifyJoin(node *memberlist.Node) {
	for _, sp := range e.ServiceProviders {
		sp.AddNode(node)
	}
}

// NotifyLeave callback function when node leave
func (e *eventDelegate) NotifyLeave(node *memberlist.Node) {
	for _, sp := range e.ServiceProviders {
		sp.RemoveNode(node)
	}
}

// NotifyUpdate callback function when node updated
func (e *eventDelegate) NotifyUpdate(node *memberlist.Node) {
}
