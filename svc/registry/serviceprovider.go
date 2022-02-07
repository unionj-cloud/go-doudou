package registry

import (
	"github.com/unionj-cloud/go-doudou/memberlist"
)

// IServiceProvider defines service provider interface for server discovery
type IServiceProvider interface {
	SelectServer() string
	AddNode(node *memberlist.Node)
	UpdateWeight(node *memberlist.Node)
	RemoveNode(node *memberlist.Node)
}
