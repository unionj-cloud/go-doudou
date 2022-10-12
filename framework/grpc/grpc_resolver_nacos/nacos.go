package grpc_resolver_nacos

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
)

const schemeName = "nacos"

const interval = 5 * time.Second

func init() {
	resolver.Register(&NacosResolver{})
}

type NacosResolver struct {
	cancelFunc context.CancelFunc
}

func (r *NacosResolver) Build(url resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	dsn := strings.Join([]string{schemeName + ":/", url.URL.Host, url.URL.Path}, "/")
	config, err := parseURL(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Wrong URL")
	}

	ctx, cancel := context.WithCancel(context.Background())
	pipe := make(chan []serviceInfo)
	go watchNacosService(ctx, config, pipe)
	go populateEndpoints(ctx, cc, pipe)

	return &NacosResolver{cancelFunc: cancel}, nil
}

func (r *NacosResolver) Scheme() string {
	return schemeName
}

func (r *NacosResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (r *NacosResolver) Close() {
	r.cancelFunc()
}

type serviceInfo struct {
	Address string
	Weight  int
}

func watchNacosService(ctx context.Context, config *NacosConfig, out chan<- []serviceInfo) {
	res := make(chan []serviceInfo)
	quit := make(chan struct{})
	go func() {
		start := true
		for {
			if !start {
				time.Sleep(interval)
			}
			start = false
			inss, err := GetHealthyInstances(config.ServiceName, config.Clusters, config.GroupName, config.NacosClient)
			if err != nil {
				select {
				case <-quit:
					return
				default:
					grpclog.Errorf("[Nacos resolver] Couldn't fetch endpoints. label={%s}; error={%v}", config.Label, err)
					continue
				}
			}
			grpclog.Infof("[Nacos resolver] %d endpoints fetched for label={%s}",
				len(inss),
				config.Label,
			)
			ee := make([]serviceInfo, 0, len(inss))
			for _, s := range inss {
				address := fmt.Sprintf("%s:%d", s.Ip, s.Port)
				ee = append(ee, serviceInfo{Address: address, Weight: (int)(s.Weight)})
			}
			select {
			case res <- ee:
				continue
			case <-quit:
				return
			}
		}
	}()
	for {
		// If in the below select both channels have values that can be read,
		// Go picks one pseudo-randomly.
		// But when the context is canceled we want to act upon it immediately.
		if ctx.Err() != nil {
			// Close quit so the goroutine returns and doesn't leak.
			// Do NOT close res because that can lead to panics in the goroutine.
			// res will be garbage collected at some point.
			close(quit)
			return
		}
		select {
		case ee := <-res:
			out <- ee
		case <-ctx.Done():
			close(quit)
			return
		}
	}
}

func populateEndpoints(ctx context.Context, clientConn resolver.ClientConn, input <-chan []serviceInfo) {
	for {
		select {
		case cc := <-input:
			connsSet := make(map[serviceInfo]struct{}, len(cc))
			for _, c := range cc {
				connsSet[c] = struct{}{}
			}
			conns := make([]resolver.Address, 0, len(connsSet))
			for c := range connsSet {
				add := resolver.Address{Addr: c.Address,
					BalancerAttributes: attributes.New(WeightAttributeKey{}, WeightAddrInfo{Weight: c.Weight})}
				//fmt.Printf("%v/n", add)
				conns = append(conns, add)
			}
			//fmt.Printf("%v/n", conns)
			sort.Sort(byAddressString(conns)) // Don't replace the same address list in the balancer

			clientConn.UpdateState(resolver.State{Addresses: conns})
		case <-ctx.Done():
			grpclog.Info("[Nacos resolver] Watch has been finished")
			return
		}
	}
}

// byAddressString sorts resolver.Address by Address Field  sorting in increasing order.
type byAddressString []resolver.Address

func (p byAddressString) Len() int           { return len(p) }
func (p byAddressString) Less(i, j int) bool { return p[i].Addr < p[j].Addr }
func (p byAddressString) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
