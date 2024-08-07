package memberlist

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/go-msgpack/codec"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry/constants"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/memberlist"
	logger "github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
)

type Service struct {
	Name          string                 `json:"name"`
	Host          string                 `json:"host"`
	Port          int                    `json:"port"`
	RouteRootPath string                 `json:"routeRootPath"`
	Type          constants.ServiceType  `json:"type"`
	Data          map[string]interface{} `json:"data,omitempty"`
}

func (receiver *Service) BaseUrl() string {
	if receiver == nil {
		return ""
	}
	switch receiver.Type {
	case constants.REST_TYPE:
		return fmt.Sprintf("http://%s:%d%s", receiver.Host, receiver.Port, receiver.RouteRootPath)
	case constants.GRPC_TYPE:
		return fmt.Sprintf("%s:%d", receiver.Host, receiver.Port)
	}
	return ""
}

type NodeMeta struct {
	Services   []Service  `json:"serviceInfo"`
	RegisterAt *time.Time `json:"registerAt"`
	GoVer      string     `json:"goVer"`
	GddVer     string     `json:"gddVer"`
	BuildUser  string     `json:"buildUser"`
	BuildTime  string     `json:"buildTime"`
	Weight     int        `json:"weight"`
}

type delegate struct {
	meta  NodeMeta
	lock  sync.Mutex
	queue *memberlist.TransmitLimitedQueue
}

func (d *delegate) AddService(service Service) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.meta.Services = append(d.meta.Services, service)
}

// NodeMeta return user custom node meta data
func (d *delegate) NodeMeta(limit int) []byte {
	d.lock.Lock()
	defer d.lock.Unlock()

	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.MsgpackHandle{})
	if err := enc.Encode(d.meta); err != nil {
		logger.Panic().Err(err).Msg("[go-doudou] Failed to encode node meta data")
	}
	raw := buf.Bytes()

	if len(raw) > limit {
		logger.Panic().Msgf("[go-doudou] Node meta data '%v' exceeds length limit of %d bytes", d.meta, limit)
	}
	return raw
}

// NotifyMsg callback function when received user data message from remote node
func (d *delegate) NotifyMsg(msg []byte) {
	d.lock.Lock()
	defer d.lock.Unlock()
	// TODO
}

// GetBroadcasts get a number of user data broadcasts
func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	d.lock.Lock()
	defer d.lock.Unlock()

	msgs := d.queue.GetBroadcasts(overhead, limit)
	return msgs
}

// LocalState also sends user data, but by tcp connection when pushPull-ing state with other node
func (d *delegate) LocalState(join bool) []byte {
	return nil
}

// MergeRemoteState gets user data from remote node by tcp connection when pushPull-ing state with other node
func (d *delegate) MergeRemoteState(s []byte, join bool) {
}
