package grpc_resolver_zk

import (
	"context"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/cast"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
	"net/url"
	"strings"
)

const schemeName = "zk"

func init() {
	resolver.Register(&ZkResolver{})
}

type ZkResolver struct {
	cancelFunc context.CancelFunc
	watcher    Watcher
}

func (r *ZkResolver) Build(url resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	dsn := strings.Join([]string{schemeName + ":/", url.URL.Host, url.URL.Path}, "/")
	config, err := parseURL(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Wrong URL")
	}

	ctx, cancel := context.WithCancel(context.Background())
	pipe := make(chan []serviceInfo)
	go watchZkService(ctx, config.Watcher, pipe)
	go populateEndpoints(ctx, cc, pipe)

	return &ZkResolver{cancelFunc: cancel, watcher: config.Watcher}, nil
}

func (r *ZkResolver) Scheme() string {
	return schemeName
}

func (r *ZkResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (r *ZkResolver) Close() {
	r.watcher.Close()
	r.cancelFunc()
}

type serviceInfo struct {
	Address string
	Weight  int
}

func watchZkService(ctx context.Context, watcher Watcher, out chan<- []serviceInfo) {
	res := make(chan []serviceInfo)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-watcher.Event():
				addrs := convertToAddress(watcher.Endpoints())
				res <- addrs
			case <-quit:
				return
			}
			if watcher.IsClosed() {
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

func convertToAddress(ups []string) (addrs []serviceInfo) {
	for _, up := range ups {
		unescaped, _ := url.PathUnescape(up)
		u, _ := url.Parse(unescaped)
		weight := cast.ToIntOrDefault(u.Query().Get("weight"), 1)
		addrs = append(addrs, serviceInfo{Address: u.Host, Weight: weight})
	}
	return
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
			clientConn.UpdateState(resolver.State{Addresses: conns})
		case <-ctx.Done():
			grpclog.Info("[Zk resolver] Watch has been finished")
			return
		}
	}
}
