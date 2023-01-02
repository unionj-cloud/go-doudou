package memberlist

import (
	"github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	"sync"

	"google.golang.org/grpc/balancer"
	balancerbase "google.golang.org/grpc/balancer/base"
)

const Name = "memberlist_weight_balancer"

func newBuilder() balancer.Builder {
	return balancerbase.NewBalancerBuilder(Name, &wPickerBuilder{}, balancerbase.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
}

type wPickerBuilder struct{}

func (*wPickerBuilder) Build(info balancerbase.PickerBuildInfo) balancer.Picker {
	zlogger.Debug().Msgf("[go-doudou] memberlist_weight_balancer Picker: Build called with info: %v", info)
	if len(info.ReadySCs) == 0 {
		return balancerbase.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	scs := make([]*conn, 0, len(info.ReadySCs))
	for sc, v := range info.ReadySCs {
		weight := v.Address.BalancerAttributes.Value(WeightAttributeKey{}).(WeightAddrInfo).Weight
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

type WeightAttributeKey struct{}

type WeightAddrInfo struct {
	Weight int
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
