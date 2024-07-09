package cache

import (
	"strings"
	"time"
	"unsafe"

	"github.com/dgraph-io/ristretto"
	go_cache "github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	gocache "github.com/unionj-cloud/go-doudou/v2/toolkit/gocache/lib/cache"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/gocache/lib/metrics"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/gocache/lib/store"
	go_cache_store "github.com/unionj-cloud/go-doudou/v2/toolkit/gocache/store/go_cache"
	redis_store "github.com/unionj-cloud/go-doudou/v2/toolkit/gocache/store/redis"
	ristretto_store "github.com/unionj-cloud/go-doudou/v2/toolkit/gocache/store/ristretto"

	"github.com/unionj-cloud/go-doudou/v2/framework/config"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/sliceutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
)

var CacheManager gocache.CacheInterface[any]

const (
	// CacheStoreRistretto TODO
	// There is a bug for CacheStoreRistretto, do not use it.
	// Use CacheStoreGoCache or CacheStoreRedis or both.
	CacheStoreRistretto = "ristretto"
	CacheStoreGoCache   = "go-cache"
	CacheStoreRedis     = "redis"
)

func init() {
	conf := config.Config{
		Cache: struct {
			TTL    int
			Stores string
			Redis  struct {
				Addr           string
				Username       string
				Password       string
				RouteByLatency bool `default:"true"`
				RouteRandomly  bool
				DB             int
				Sentinel       struct {
					Master   string
					Nodes    string
					Password string
				}
			}
			Ristretto struct {
				NumCounters int64 `default:"1000"`
				MaxCost     int64 `default:"100"`
				BufferItems int64 `default:"64"`
			}
			GoCache struct {
				Expiration      time.Duration `default:"5m"`
				CleanupInterval time.Duration `default:"10m"`
			}
		}{
			TTL:    cast.ToIntOrDefault(config.GddCacheTTL.Load(), config.DefaultGddCacheTTL),
			Stores: config.GddCacheStores.LoadOrDefault(config.DefaultGddCacheStores),
			Ristretto: struct {
				NumCounters int64 `default:"1000"`
				MaxCost     int64 `default:"100"`
				BufferItems int64 `default:"64"`
			}{
				NumCounters: cast.ToInt64OrDefault(config.GddCacheRistrettoNumCounters.Load(), config.DefaultGddCacheRistrettoNumCounters),
				MaxCost:     cast.ToInt64OrDefault(config.GddCacheRistrettoMaxCost.Load(), config.DefaultGddCacheRistrettoMaxCost),
				BufferItems: cast.ToInt64OrDefault(config.GddCacheRistrettoBufferItems.Load(), config.DefaultGddCacheRistrettoBufferItems),
			},
			Redis: struct {
				Addr           string
				Username       string
				Password       string
				RouteByLatency bool `default:"true"`
				RouteRandomly  bool
				DB             int
				Sentinel       struct {
					Master   string
					Nodes    string
					Password string
				}
			}{
				Addr:           config.GddCacheRedisAddr.LoadOrDefault(config.DefaultGddCacheRedisAddr),
				Username:       config.GddCacheRedisUser.LoadOrDefault(config.DefaultGddCacheRedisUser),
				Password:       config.GddCacheRedisPass.LoadOrDefault(config.DefaultGddCacheRedisPass),
				RouteByLatency: cast.ToBoolOrDefault(config.GddCacheRedisRouteByLatency.Load(), config.DefaultGddCacheRedisRouteByLatency),
				RouteRandomly:  cast.ToBoolOrDefault(config.GddCacheRedisRouteRandomly.Load(), config.DefaultGddCacheRedisRouteRandomly),
			},
			GoCache: struct {
				Expiration      time.Duration `default:"5m"`
				CleanupInterval time.Duration `default:"10m"`
			}{
				Expiration:      config.GddCacheGocacheExpiration.LoadDurationOrDefault(config.DefaultGddCacheGocacheExpiration),
				CleanupInterval: config.GddCacheGocacheCleanupInterval.LoadDurationOrDefault(config.DefaultGddCacheGocacheCleanupInterval),
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
			ristrettoStore = ristretto_store.NewRistretto(ristrettoCache, store.WithExpiration(time.Duration(ttl)*time.Second), store.WithSynchronousSet())
		} else {
			ristrettoStore = ristretto_store.NewRistretto(ristrettoCache, store.WithSynchronousSet())
		}
		setterCaches = append(setterCaches, gocache.New[any](ristrettoStore))
	}

	if sliceutils.StringContains(stores, CacheStoreGoCache) {
		gocacheClient := go_cache.New(conf.Cache.GoCache.Expiration, conf.Cache.GoCache.CleanupInterval)
		setterCaches = append(setterCaches, gocache.New[any](go_cache_store.NewGoCache(gocacheClient)))
	}

	if sliceutils.StringContains(stores, CacheStoreRedis) {
		var redisClient redis_store.RedisClientInterface
		if stringutils.IsNotEmpty(conf.Cache.Redis.Sentinel.Nodes) {
			redisClient = redis.NewFailoverClusterClient(&redis.FailoverOptions{
				MasterName:       conf.Cache.Redis.Sentinel.Master,
				SentinelAddrs:    strings.Split(conf.Cache.Redis.Sentinel.Nodes, ","),
				SentinelPassword: conf.Cache.Redis.Sentinel.Password,
				Password:         conf.Cache.Redis.Password,
				DB:               conf.Cache.Redis.DB,
			})
		} else if stringutils.IsNotEmpty(conf.Cache.Redis.Addr) {
			redisAddr := conf.Cache.Redis.Addr
			addrs := strings.Split(redisAddr, ",")
			if len(addrs) > 1 {
				redisClient = redis.NewClusterClient(&redis.ClusterOptions{
					Addrs:          addrs,
					Username:       conf.Cache.Redis.Username,
					Password:       conf.Cache.Redis.Password,
					RouteByLatency: conf.Cache.Redis.RouteByLatency,
					RouteRandomly:  conf.Cache.Redis.RouteRandomly,
				})
			} else {
				redisClient = redis.NewClient(&redis.Options{
					Addr:     addrs[0],
					Username: conf.Cache.Redis.Username,
					Password: conf.Cache.Redis.Password,
					DB:       conf.Cache.Redis.DB,
				})
			}
		}
		var redisStore *redis_store.RedisStore
		if ttl > 0 {
			redisStore = redis_store.NewRedis(redisClient, store.WithExpiration(time.Duration(ttl)*time.Second))
		} else {
			redisStore = redis_store.NewRedis(redisClient)
		}
		setterCaches = append(setterCaches, gocache.New[any](redisStore))
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
