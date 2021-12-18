package registry

import "github.com/unionj-cloud/memberlist"

// IServiceProvider defines service provider interface for server discovery
type IServiceProvider interface {
	SelectServer() (string, error)
	AddNode(node *memberlist.Node)
	UpdateWeight(node *memberlist.Node)
	RemoveNode(node *memberlist.Node)
}
