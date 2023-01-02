package memberlist

import (
	"net"
	"time"
)

//go:generate mockgen -destination ./mock/mock_memberlist_interface.go -package mock -source=./memberlist_interface.go

type IMemberlist interface {
	Join(existing []string) (int, error)
	Ping(node string, addr net.Addr) (time.Duration, error)
	LocalNode() *Node
	UpdateNode(timeout time.Duration) error
	SendToAddress(a Address, msg []byte) error
	SendBestEffort(to *Node, msg []byte) error
	SendReliable(to *Node, msg []byte) error
	Members() []*Node
	NumMembers() (alive int)
	Leave(timeout time.Duration) error
	GetHealthScore() int
	ProtocolVersion() uint8
	Shutdown() error
	Config() *Config
	AdvertiseAddr() string
}
