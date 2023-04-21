package grpc_resolver_zk

import (
	"net/url"

	"github.com/pkg/errors"
)

// A Watcher represents how a serverset.Watch is used so it can be stubbed out for tests.
type Watcher interface {
	Endpoints() []string
	Event() <-chan struct{}
	IsClosed() bool
	Close()
}

type ZkConfig struct {
	Label       string //用做寻找配置的标签
	ServiceName string //标记服务名称
	Watcher     Watcher
}

var ZkConfigs = make(map[string]*ZkConfig)

func AddZkConfig(config ZkConfig) {
	//ZkConfigs.Store(config.Label, &config)
	ZkConfigs[config.Label] = &config
}

func DelZkConfig(label string) {
	delete(ZkConfigs, label)
}

func parseURL(u string) (*ZkConfig, error) {
	rawURL, err := url.Parse(u)
	if err != nil {
		return &ZkConfig{}, errors.Wrap(err, "Wrong zookeeper URL")
	}
	if rawURL.Scheme != schemeName || len(rawURL.Host) == 0 {
		return &ZkConfig{}, errors.Wrap(err, "Wrong zookeeper URL")
	}
	config, ok := ZkConfigs[rawURL.Host]
	if !ok {
		return &ZkConfig{}, errors.Wrap(err, "The zookeeper config is not exist")
	}
	return config, nil
}
