package etcd

import (
	"sync"

	"github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

const Name = "etcd_weight_balancer"

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &wPickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
}

type wPickerBuilder struct{}

func (*wPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	zlogger.Debug().Msgf("[go-doudou] etcd_weight_balancer Picker: Build called with info: %v", info)
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	scs := make([]*conn, 0, len(info.ReadySCs))
	for sc, v := range info.ReadySCs {
		weight := 1
		if metadata, ok := v.Address.Metadata.(map[string]interface{}); !ok {
			zlogger.Error().Msg("[go-doudou] etcd endpoint metadata is not map[string]string type")
		} else {
			weight = int(metadata["weight"].(float64))
		}
		scs = append(scs, &conn{sc: sc, weight: weight})
	}
	return &wPicker{
		subConns: scs,
	}
}

type wPicker struct {
	subConns []*conn
	mu       sync.Mutex
}

func (p *wPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	p.mu.Lock()
	sc := newChooser(p.subConns).pick().sc
	p.mu.Unlock()
	return balancer.PickResult{SubConn: sc}, nil
}

type conn struct {
	sc            balancer.SubConn
	weight        int
	currentWeight int
}

// Chooser from naming_client package in nacos-sdk-go
type Chooser struct {
	data []*conn
}

// NewChooser initializes a new Chooser for picking from the provided Choices.
func newChooser(cs []*conn) Chooser {
	return Chooser{data: cs}
}

func (chs Chooser) pick() conn {
	var selected *conn
	total := 0
	for i := 0; i < len(chs.data); i++ {
		s := chs.data[i]
		s.currentWeight += s.weight
		total += s.weight
		if selected == nil || s.currentWeight > selected.currentWeight {
			selected = s
		}
	}
	selected.currentWeight -= total
	return *selected
}
