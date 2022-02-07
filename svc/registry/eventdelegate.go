package registry

import (
	"github.com/unionj-cloud/go-doudou/memberlist"
)

type eventDelegate struct {
	ServiceProviders []IServiceProvider
}

func (e *eventDelegate) NotifySuspectSateChange(node *memberlist.Node) {
	for _, sp := range e.ServiceProviders {
		if node.State == memberlist.StateSuspect {
			sp.RemoveNode(node)
		} else if node.State == memberlist.StateAlive {
			sp.AddNode(node)
		}
	}
}

func (e *eventDelegate) NotifyWeight(node *memberlist.Node) {
	for _, sp := range e.ServiceProviders {
		sp.UpdateWeight(node)
	}
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
	for _, sp := range e.ServiceProviders {
		sp.AddNode(node)
	}
}
