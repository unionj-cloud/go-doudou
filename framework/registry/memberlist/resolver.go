package memberlist

import (
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/memberlist"
	"google.golang.org/grpc/attributes"
	gresolver "google.golang.org/grpc/resolver"
	"net/url"
	"strings"
	"sync"
)

const schemeName = "memberlist"

var _ gresolver.Builder = (*builder)(nil)
var _ gresolver.Resolver = (*resolver)(nil)
var _ IMemberlistServiceProvider = (*resolver)(nil)

func init() {
	gresolver.Register(&builder{})
}

type builder struct {
}

func parseURL(u string) (string, error) {
	rawURL, err := url.Parse(u)
	if err != nil {
		return "", errors.Wrap(err, "Wrong memberlist URL")
	}
	if rawURL.Scheme != schemeName || len(rawURL.Host) == 0 {
		return "", errors.Wrap(err, "Wrong memberlist URL")
	}
	return rawURL.Host, nil
}

func (r *builder) Scheme() string {
	return schemeName
}

func (b *builder) Build(target gresolver.Target, cc gresolver.ClientConn, opts gresolver.BuildOptions) (gresolver.Resolver, error) {
	dsn := strings.Join([]string{schemeName + ":/", target.URL.Host, target.URL.Path}, "/")
	name, err := parseURL(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Wrong URL")
	}
	r := &resolver{
		base: base{
			name:    name,
			nodeMap: make(map[string]*server),
		},
		cc: cc,
	}
	RegisterServiceProvider(r)
	return r, nil
}

type resolver struct {
	base base
	cc   gresolver.ClientConn
	lock sync.RWMutex
}

func (m *resolver) AddNode(node *memberlist.Node) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.base.AddNode(node)
	m.UpdateCC()
}

func (m *resolver) UpdateWeight(node *memberlist.Node) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.base.UpdateWeight(node)
	m.UpdateCC()
}

func (m *resolver) RemoveNode(node *memberlist.Node) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.base.RemoveNode(node)
	m.UpdateCC()
}

func (m *resolver) UpdateCC() {
	conns := make([]gresolver.Address, 0, len(m.base.nodes))
	for _, item := range m.base.nodes {
		add := gresolver.Address{Addr: item.baseUrl,
			BalancerAttributes: attributes.New(WeightAttributeKey{}, WeightAddrInfo{Weight: item.weight})}
		conns = append(conns, add)
	}
	m.cc.UpdateState(gresolver.State{Addresses: conns})
}

func (m *resolver) ResolveNow(gresolver.ResolveNowOptions) {}

func (m *resolver) Close() {
}
