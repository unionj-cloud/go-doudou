package grpc_resolver_nacos

import (
	"math/rand"
	"sort"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
)

const Name = "nacos_weight_balancer"

var logger = grpclog.Component("nacos_weight_balancer")

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &wPickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
}

type wPickerBuilder struct{}

func (*wPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	logger.Infof("nacos_weight_balancer Picker: Build called with info: %v", info)
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	scs := make([]conn, 0, len(info.ReadySCs))
	for sc, v := range info.ReadySCs {
		//fmt.Println(v)
		weight := v.Address.BalancerAttributes.Value(WeightAttributeKey{}).(WeightAddrInfo).Weight
		scs = append(scs, conn{sc: sc, Weight: weight})
	}
	return &wPicker{
		subConns: scs,
	}
}

type wPicker struct {
	subConns conns
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
	sc     balancer.SubConn
	Weight int
}
type conns []conn

func (a conns) Len() int {
	return len(a)
}

func (a conns) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a conns) Less(i, j int) bool {
	return a[i].Weight < a[j].Weight
}

// Chooser from naming_client package in nacos-sdk-go
type Chooser struct {
	data   []conn
	totals []int
	max    int
}

// NewChooser initializes a new Chooser for picking from the provided Choices.
func newChooser(cs conns) Chooser {
	sort.Sort(conns(cs))
	totals := make([]int, len(cs))
	runningTotal := 0
	for i, c := range cs {
		runningTotal += int(c.Weight)
		totals[i] = runningTotal
	}
	return Chooser{data: cs, totals: totals, max: runningTotal}
}

func (chs Chooser) pick() conn {
	r := rand.Intn(chs.max) + 1
	i := sort.SearchInts(chs.totals, r)
	return chs.data[i]
}
