package mock

import (
	"container/list"
	"github.com/apolloconfig/agollo/v4/agcache"
	"github.com/apolloconfig/agollo/v4/storage"
)

//go:generate mockgen -destination mock_apollo_client_interface.go -package mock -source=./apollo_client_interface.go

//Client apollo 客户端接口
type Client interface {
	GetConfig(namespace string) *storage.Config
	GetConfigAndInit(namespace string) *storage.Config
	GetConfigCache(namespace string) agcache.CacheInterface
	GetDefaultConfigCache() agcache.CacheInterface
	GetApolloConfigCache() agcache.CacheInterface
	GetValue(key string) string
	GetStringValue(key string, defaultValue string) string
	GetIntValue(key string, defaultValue int) int
	GetFloatValue(key string, defaultValue float64) float64
	GetBoolValue(key string, defaultValue bool) bool
	GetStringSliceValue(key string, defaultValue []string) []string
	GetIntSliceValue(key string, defaultValue []int) []int
	AddChangeListener(listener storage.ChangeListener)
	RemoveChangeListener(listener storage.ChangeListener)
	GetChangeListeners() *list.List
	UseEventDispatch()
}
