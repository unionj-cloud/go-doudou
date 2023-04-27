package grpc_resolver_zk

import (
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/cast"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"net/url"
	"strings"
)

const schemeName = "zk"

func init() {
	resolver.Register(&ZkResolver{})
}

type ZkResolver struct {
	*ZkConfig
}

func (r *ZkResolver) Build(url resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	dsn := strings.Join([]string{schemeName + ":/", url.URL.Host, url.URL.Path}, "/")
	config, err := parseURL(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Wrong URL")
	}
	zkResolver := &ZkResolver{
		ZkConfig: config,
	}
	go zkResolver.watchZkService(cc)
	return zkResolver, nil
}

func (r *ZkResolver) Scheme() string {
	return schemeName
}

func (r *ZkResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (r *ZkResolver) Close() {
	r.Watcher.Close()
}

type serviceInfo struct {
	Address string
	Weight  int
}

func (r *ZkResolver) updateState(clientConn resolver.ClientConn) {
	services := r.convertToAddress(r.Watcher.Endpoints())
	connsSet := make(map[serviceInfo]struct{}, len(services))
	for _, c := range services {
		connsSet[c] = struct{}{}
	}
	addrs := make([]resolver.Address, 0, len(connsSet))
	for c := range connsSet {
		addr := resolver.Address{Addr: c.Address,
			BalancerAttributes: attributes.New(WeightAttributeKey{}, WeightAddrInfo{Weight: c.Weight})}
		addrs = append(addrs, addr)
	}
	clientConn.UpdateState(resolver.State{Addresses: addrs})
}

func (r *ZkResolver) watchZkService(clientConn resolver.ClientConn) {
	r.updateState(clientConn)
	for {
		select {
		case _, ok := <-r.Watcher.Event():
			if !ok {
				return
			}
			r.updateState(clientConn)
		}

		if r.Watcher.IsClosed() {
			return
		}
	}
}

func (r *ZkResolver) convertToAddress(ups []string) (addrs []serviceInfo) {
	for _, up := range ups {
		unescaped, _ := url.QueryUnescape(up)
		u, _ := url.Parse(unescaped)
		weight := cast.ToIntOrDefault(u.Query().Get("weight"), 1)
		group := u.Query().Get("group")
		version := u.Query().Get("version")
		if group != r.Group || version != r.Version {
			continue
		}
		addrs = append(addrs, serviceInfo{Address: u.Host, Weight: weight})
	}
	return
}
