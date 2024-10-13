package serversets

import (
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/unionj-cloud/go-doudou/v2/framework/config"
	"github.com/unionj-cloud/toolkit/cast"

	"github.com/go-zookeeper/zk"
)

var json = sonic.ConfigDefault

// An Endpoint is a service (host and port) registered on Zookeeper
// to be discovered by clients/watchers.
type Endpoint struct {
	*ServerSet
	PingRate   time.Duration // default/initial is 1 second
	CloseEvent chan struct{}

	done chan struct{}
	wg   sync.WaitGroup

	host string
	port int

	key   string
	ping  func() error
	alive bool
}

// RegisterEndpoint registers a host and port as alive. It creates the appropriate
// Zookeeper nodes and watchers will be notified this server/endpoint is available.
func (ss *ServerSet) RegisterEndpoint(host string, port int, ping func() error) (*Endpoint, error) {
	endpoint := &Endpoint{
		ServerSet:  ss,
		PingRate:   time.Second,
		CloseEvent: make(chan struct{}, 1),
		done:       make(chan struct{}),
		host:       host,
		port:       port,
		ping:       ping,
		alive:      true,
	}

	if ping != nil {
		endpoint.alive = endpoint.ping() == nil
	}

	connection, sessionEvents, err := ss.connectToZookeeper()
	if err != nil {
		return nil, err
	}

	err = endpoint.update(connection)
	if err != nil {
		return nil, err
	}

	// spawn goroutine to deal with connection/session issues.
	endpoint.wg.Add(1)
	go func() {
		defer endpoint.wg.Done()
		for {
			select {
			case event := <-sessionEvents:
				if event.Type == zk.EventSession && event.State == zk.StateExpired {
					connection.Close()
					connection = nil
				}
			case <-endpoint.done:
				connection.Close()
				return
			}

			if connection == nil {
				connection, sessionEvents, err = ss.connectToZookeeper()
				if err != nil {
					panic(fmt.Errorf("unable to reconnect to zookeeper after session expired: %v", err))
				}

				err = endpoint.update(connection)
				if err != nil {
					panic(fmt.Errorf("unable to reregister endpoint after session expired: %v", err))
				}
			}
		}
	}()

	if ping != nil {
		endpoint.wg.Add(1)
		go func() {
			defer endpoint.wg.Done()
			for {
				select {
				case <-time.After(endpoint.PingRate):
				case <-endpoint.done:
					return
				}

				alive := endpoint.ping() == nil
				if alive != endpoint.alive {
					endpoint.alive = alive
					err := endpoint.update(connection)

					if err != nil {
						panic(fmt.Errorf("unable to reregister after ping change: %v", err))
					}
				}
			}
		}()
	}

	return endpoint, nil
}

func (ss *ServerSet) RegisterEndpointWithMeta(host string, port int, ping func() error, meta map[string]interface{}) (*Endpoint, error) {
	endpoint := &Endpoint{
		ServerSet:  ss,
		PingRate:   time.Second,
		CloseEvent: make(chan struct{}, 1),
		done:       make(chan struct{}),
		host:       host,
		port:       port,
		ping:       ping,
		alive:      true,
	}

	if ping != nil {
		endpoint.alive = endpoint.ping() == nil
	}

	connection, sessionEvents, err := ss.connectToZookeeper()
	if err != nil {
		return nil, err
	}

	err = endpoint.updateWithMeta(connection, meta)
	if err != nil {
		return nil, err
	}

	// spawn goroutine to deal with connection/session issues.
	endpoint.wg.Add(1)
	go func() {
		defer endpoint.wg.Done()
		for {
			select {
			case event := <-sessionEvents:
				if event.Type == zk.EventSession && event.State == zk.StateExpired {
					connection.Close()
					connection = nil
				}
			case <-endpoint.done:
				connection.Close()
				return
			}

			if connection == nil {
				connection, sessionEvents, err = ss.connectToZookeeper()
				if err != nil {
					panic(fmt.Errorf("unable to reconnect to zookeeper after session expired: %v", err))
				}

				err = endpoint.updateWithMeta(connection, meta)
				if err != nil {
					panic(fmt.Errorf("unable to reregister endpoint after session expired: %v", err))
				}
			}
		}
	}()

	if ping != nil {
		endpoint.wg.Add(1)
		go func() {
			defer endpoint.wg.Done()
			for {
				select {
				case <-time.After(endpoint.PingRate):
				case <-endpoint.done:
					return
				}

				alive := endpoint.ping() == nil
				if alive != endpoint.alive {
					endpoint.alive = alive
					err := endpoint.updateWithMeta(connection, meta)

					if err != nil {
						panic(fmt.Errorf("unable to reregister after ping change: %v", err))
					}
				}
			}
		}()
	}

	return endpoint, nil
}

// Close blocks until the client connection to Zookeeper is closed.
// If already called, will simply return, even if in the process of closing.
func (ep *Endpoint) Close() {
	select {
	case <-ep.done:
		return
	default:
	}

	close(ep.done)
	ep.wg.Wait()
	ep.CloseEvent <- struct{}{}

	return
}

func (ep *Endpoint) update(connection *zk.Conn) error {
	// don't create/remove the node if we're dead
	if !ep.alive {
		if ep.key != "" {
			err := connection.Delete(ep.key, 0)
			ep.key = ""
			return err
		}

		return nil
	}

	entityData, _ := json.Marshal(newEntity(ep.host, ep.port))
	entityMap := make(map[string]interface{})
	json.Unmarshal(entityData, &entityMap)

	var err error
	ep.key, err = ep.ServerSet.registerEndpoint(connection, entityMap)

	return err
}

func (ep *Endpoint) updateWithMeta(connection *zk.Conn, meta map[string]interface{}) error {
	// don't create/remove the node if we're dead
	if !ep.alive {
		if ep.key != "" {
			err := connection.Delete(ep.key, 0)
			ep.key = ""
			return err
		}

		return nil
	}

	entityData := newEntity(ep.host, ep.port)
	meta["serviceEndpoint"] = entityData.ServiceEndpoint
	meta["status"] = entityData.Status

	var err error
	ep.key, err = ep.ServerSet.registerEndpoint(connection, meta)

	return err
}

func (ss *ServerSet) registerEndpoint(connection *zk.Conn, meta map[string]interface{}) (string, error) {
	err := ss.createFullPath(connection)
	if err != nil {
		return "", err
	}
	flags := zk.FlagEphemeral
	if cast.ToBoolOrDefault(config.GddZkSequence.Load(), config.DefaultGddZkSequence) {
		flags = zk.FlagEphemeral | zk.FlagSequence
	}
	querystring := url.Values{}
	if meta["group"] != nil {
		querystring.Set("group", meta["group"].(string))
	}
	if meta["version"] != nil {
		querystring.Set("version", meta["version"].(string))
	}
	if meta["weight"] != nil {
		querystring.Set("weight", strconv.Itoa(meta["weight"].(int)))
	}
	if meta["rootPath"] != nil {
		querystring.Set("rootPath", meta["rootPath"].(string))
	}
	if meta["serviceEndpoint"] != nil {
		serviceEndpoint := meta["serviceEndpoint"].(endpoint)
		if meta["host"] == nil {
			meta["host"] = serviceEndpoint.Host
		}
		if meta["port"] == nil {
			meta["port"] = serviceEndpoint.Port
		}
	}
	memberPath := url.QueryEscape(fmt.Sprintf("%s://%s:%s/%s?%s", meta["scheme"], meta["host"], meta["port"], meta["service"], querystring.Encode()))
	data, _ := json.Marshal(meta)
	return connection.Create(
		ss.directoryPath()+"/"+memberPath,
		data,
		int32(flags),
		zk.WorldACL(zk.PermAll))
}
