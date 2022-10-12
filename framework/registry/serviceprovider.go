package registry

import (
	"github.com/unionj-cloud/go-doudou/framework/memberlist"
)

type IServiceProvider interface {
	SelectServer() string
}

// IMemberlistServiceProvider defines service provider interface for server discovery
type IMemberlistServiceProvider interface {
	IServiceProvider
	AddNode(node *memberlist.Node)
	UpdateWeight(node *memberlist.Node)
	RemoveNode(node *memberlist.Node)
}
