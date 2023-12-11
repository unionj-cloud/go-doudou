package cache

import (
	"strings"
	"time"
	"unsafe"

	"github.com/dgraph-io/ristretto"
	gocache "github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/metrics"
	"github.com/eko/gocache/lib/v4/store"
	redis_store "github.com/eko/gocache/store/redis/v4"
	ristretto_store "github.com/eko/gocache/store/ristretto/v4"
	"github.com/redis/go-redis/v9"

	"github.com/unionj-cloud/go-doudou/v2/framework/config"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/sliceutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
)

var CacheManager gocache.CacheInterface[any]

const (
	CacheStoreRistretto = "ristretto"
	CacheStoreRedis     = "redis"
)

func init() {
	conf := config.Config{
		Cache: struct {
			TTL    int
			Stores string
			Redis  struct {
				Addr           string
				RouteByLatency bool "default:\"true\""
				RouteRandomly  bool
			}
			Ristretto struct {
				NumCounters int64 "default:\"1000\""
				MaxCost     int64 "default:\"100\""
				BufferItems int64 "default:\"64\""
			}
		}{
			TTL:    cast.ToIntOrDefault(config.GddCacheTTL.Load(), config.DefaultGddCacheTTL),
			Stores: config.GddCacheStores.LoadOrDefault(config.DefaultGddCacheStores),
			Ristretto: struct {
				NumCounters int64 "default:\"1000\""
				MaxCost     int64 "default:\"100\""
				BufferItems int64 "default:\"64\""
			}{
				NumCounters: cast.ToInt64OrDefault(config.GddCacheRistrettoNumCounters.Load(), config.DefaultGddCacheRistrettoNumCounters),
				MaxCost:     cast.ToInt64OrDefault(config.GddCacheRistrettoMaxCost.Load(), config.DefaultGddCacheRistrettoMaxCost),
				BufferItems: cast.ToInt64OrDefault(config.GddCacheRistrettoBufferItems.Load(), config.DefaultGddCacheRistrettoBufferItems),
			},
			Redis: struct {
				Addr           string
				RouteByLatency bool "default:\"true\""
				RouteRandomly  bool
			}{
				Addr:           config.GddCacheRedisAddr.LoadOrDefault(config.DefaultGddCacheRedisAddr),
				RouteByLatency: cast.ToBoolOrDefault(config.GddCacheRedisRouteByLatency.Load(), config.DefaultGddCacheRedisRouteByLatency),
				RouteRandomly:  cast.ToBoolOrDefault(config.GddCacheRedisRouteRandomly.Load(), config.DefaultGddCacheRedisRouteRandomly),
			},
		},
		Service: struct{ Name string }{
			Name: config.GddServiceName.LoadOrDefault(config.DefaultGddServiceName),
		},
	}

	CacheManager = NewCacheManager(conf)
}

func NewCacheManager(conf config.Config) gocache.CacheInterface[any] {
	storesStr := conf.Cache.Stores
	if stringutils.IsEmpty(storesStr) {
		return nil
	}
	stores := strings.Split(storesStr, ",")

	var setterCaches []gocache.SetterCacheInterface[any]
	ttl := conf.Cache.TTL

	if sliceutils.StringContains(stores, CacheStoreRistretto) {
		ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
			NumCounters: conf.Cache.Ristretto.NumCounters,
			MaxCost:     conf.Cache.Ristretto.MaxCost,
			BufferItems: conf.Cache.Ristretto.BufferItems,
			Cost: func(value interface{}) int64 {
				return int64(unsafe.Sizeof(value))
			},
		})
		if err != nil {
			panic(err)
		}
		var ristrettoStore *ristretto_store.RistrettoStore
		if ttl > 0 {
			ristrettoStore = ristretto_store.NewRistretto(ristrettoCache, store.WithExpiration(time.Duration(ttl)*time.Second))
		} else {
			ristrettoStore = ristretto_store.NewRistretto(ristrettoCache)
		}
		setterCaches = append(setterCaches, gocache.New[any](ristrettoStore))
	}

	if sliceutils.StringContains(stores, CacheStoreRedis) {
		redisAddr := conf.Cache.Redis.Addr
		if stringutils.IsNotEmpty(redisAddr) {
			addrs := strings.Split(redisAddr, ",")
			var redisClient redis_store.RedisClientInterface
			if len(addrs) > 1 {
				redisClient = redis.NewClusterClient(&redis.ClusterOptions{
					Addrs:          addrs,
					RouteByLatency: conf.Cache.Redis.RouteByLatency,
					RouteRandomly:  conf.Cache.Redis.RouteRandomly,
				})
			} else {
				redisClient = redis.NewClient(&redis.Options{Addr: addrs[0]})
			}
			var redisStore *redis_store.RedisStore
			if ttl > 0 {
				redisStore = redis_store.NewRedis(redisClient, store.WithExpiration(time.Duration(ttl)*time.Second))
			} else {
				redisStore = redis_store.NewRedis(redisClient)
			}
			setterCaches = append(setterCaches, gocache.New[any](redisStore))
		}
	}

	var cacheManager gocache.CacheInterface[any]

	// Initialize chained cache
	cacheManager = gocache.NewChain[any](setterCaches...)

	serviceName := conf.Service.Name
	if stringutils.IsNotEmpty(serviceName) {
		// Initializes Prometheus metrics service
		promMetrics := metrics.NewPrometheus(serviceName)

		// Initialize chained cache
		cacheManager = gocache.NewMetric[any](promMetrics, cacheManager)
	}

	return cacheManager
}
