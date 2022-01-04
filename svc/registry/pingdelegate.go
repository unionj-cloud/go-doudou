package registry

import (
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/memberlist"
	"time"
)

type pingDelegate struct {
}

func (p pingDelegate) AckPayload() []byte {
	return nil
}

func (p pingDelegate) NotifyPingComplete(other *memberlist.Node, rtt time.Duration, payload []byte) {
	logrus.Debugf("[go-doudou] ping remote node %s success in %s", other.Name, rtt.String())
}
