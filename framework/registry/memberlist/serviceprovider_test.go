package memberlist

import (
	"sync"

	"github.com/unionj-cloud/toolkit/memberlist"
)

type mockServiceProvider struct {
	name      string
	servers   []*memberlist.Node
	serverMap map[string]*memberlist.Node
	lock      sync.RWMutex
}

func (m *mockServiceProvider) SelectServer() string {
	return "mock server"
}

func (m *mockServiceProvider) AddNode(node *memberlist.Node) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.servers = append(m.servers, node)
	m.serverMap[node.Name] = node
}

func (m *mockServiceProvider) UpdateWeight(node *memberlist.Node) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if s, exists := m.serverMap[node.Name]; exists {
		s.Weight = node.Weight
	}
}

func (m *mockServiceProvider) RemoveNode(node *memberlist.Node) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, exists := m.serverMap[node.Name]; exists {
		var idx int
		for i, n := range m.servers {
			if n.Name == node.Name {
				idx = i
			}
		}
		m.servers = append(m.servers[:idx], m.servers[idx+1:]...)
		delete(m.serverMap, node.Name)
	}
}

func newMockServiceProvider(name string) *mockServiceProvider {
	return &mockServiceProvider{
		servers:   make([]*memberlist.Node, 0),
		serverMap: make(map[string]*memberlist.Node),
		name:      name,
	}
}
