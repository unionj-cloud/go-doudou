package registry

import (
	"github.com/unionj-cloud/go-doudou/framework/memberlist"
	"github.com/wubin1989/nacos-sdk-go/clients/naming_client"
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

type INacosServiceProvider interface {
	SetClusters(clusters []string)
	SetGroupName(groupName string)
	SetNamingClient(namingClient naming_client.INamingClient)
}
