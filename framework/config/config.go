package config

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/unionj-cloud/go-doudou/v2/toolkit/constants"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/fileutils"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"github.com/wubin1989/nacos-sdk-go/v2/common/constant"
	"github.com/wubin1989/nacos-sdk-go/v2/vo"
	_ "go.uber.org/automaxprocs"

	"github.com/unionj-cloud/go-doudou/v2/framework/configmgr"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/dotenv"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/envconfig"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/yaml"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
)

var GddConfig = &gddConfig{}

func LoadConfigFromLocal() {
	env := os.Getenv("GDD_ENV")
	if "" == env {
		env = "dev"
	}
	yaml.Load(env)
	dotenv.Load(env)
}

func LoadConfigFromRemote() {
	configType := GddConfigRemoteType.LoadOrDefault(DefaultGddConfigRemoteType)
	switch configType {
	case "":
		return
	case NacosConfigType:
		nacosConfigFormat := GddNacosConfigFormat.LoadOrDefault(string(DefaultGddNacosConfigFormat))
		nacosConfigGroup := GddNacosConfigGroup.LoadOrDefault(DefaultGddNacosConfigGroup)
		nacosConfigDataid := GddNacosConfigDataid.LoadOrDefault(DefaultGddNacosConfigDataid)
		if stringutils.IsEmpty(nacosConfigDataid) {
			panic(errors.New("[go-doudou] nacos config dataId is required"))
		}
		err := configmgr.LoadFromNacos(GetNacosClientParam(), nacosConfigDataid, nacosConfigFormat, nacosConfigGroup)
		if err != nil {
			panic(errors.Wrap(err, "[go-doudou] fail to load config from Nacos"))
		}
	case ApolloConfigType:
		if stringutils.IsEmpty(GddServiceName.Load()) {
			panic(errors.New("[go-doudou] service name is required"))
		}
		apolloCluster := GddApolloCluster.LoadOrDefault(DefaultGddApolloCluster)
		apolloAddr := GddApolloAddr.LoadOrDefault(DefaultGddApolloAddr)
		if stringutils.IsEmpty(apolloAddr) {
			panic(errors.New("[go-doudou] apollo config service address is required"))
		}
		apolloNamespace := GddApolloNamespace.LoadOrDefault(DefaultGddApolloNamespace)
		apolloBackup := cast.ToBoolOrDefault(GddApolloBackupEnable.Load(), DefaultGddApolloBackupEnable)
		apolloBackupPath := GddApolloBackupPath.LoadOrDefault(DefaultGddApolloBackupPath)
		apolloSecret := GddApolloSecret.LoadOrDefault(DefaultGddApolloSecret)
		apolloMustStart := cast.ToBoolOrDefault(GddApolloMuststart.Load(), DefaultGddApolloMuststart)
		apolloLogEnable := cast.ToBoolOrDefault(GddApolloLogEnable.Load(), DefaultGddApolloLogEnable)
		if apolloLogEnable {
			agollo.SetLogger(logrus.StandardLogger())
		}
		c := &config.AppConfig{
			AppID:            GddServiceName.Load(),
			Cluster:          apolloCluster,
			IP:               apolloAddr,
			NamespaceName:    apolloNamespace,
			IsBackupConfig:   apolloBackup,
			Secret:           apolloSecret,
			BackupConfigPath: apolloBackupPath,
			MustStart:        apolloMustStart,
		}
		configmgr.LoadFromApollo(c)
	default:
		panic(fmt.Errorf("[go-doudou] unknown config type: %s\n", configType))
	}
}

func CheckDev() bool {
	return stringutils.IsEmpty(os.Getenv("GDD_ENV")) || os.Getenv("GDD_ENV") == "dev"
}

var G_LogFile *os.File
var lock sync.Mutex

func Shutdown() {
	lock.Lock()
	defer lock.Unlock()
	if G_LogFile != nil {
		G_LogFile.Close()
		G_LogFile = nil
	}
}

func init() {
	LoadConfigFromLocal()
	LoadConfigFromRemote()
	err := envconfig.Process("gdd", GddConfig)
	if err != nil {
		zlogger.Panic().Msgf("Error processing environment variables: %v", err)
	}
	zl, _ := zerolog.ParseLevel(GddLogLevel.LoadOrDefault(DefaultGddLogLevel))
	opts := []zlogger.LoggerConfigOption{
		zlogger.WithDev(CheckDev()),
		zlogger.WithCaller(cast.ToBoolOrDefault(GddLogCaller.Load(), DefaultGddLogCaller)),
		zlogger.WithDiscard(cast.ToBoolOrDefault(GddLogDiscard.Load(), DefaultGddLogDiscard)),
		zlogger.WithZeroLogLevel(zl),
	}
	logPath := GddLogPath.LoadOrDefault(DefaultGddLogPath)
	if stringutils.IsNotEmpty(logPath) {
		fileutils.CreateDirectory(logPath)
		logFile := filepath.Join(logPath, "app.log")
		G_LogFile, err = os.OpenFile(
			logFile,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)
		if err != nil {
			panic(err)
		}
		logStyle := GddLogStyle.LoadOrDefault(DefaultGddLogStyle)
		switch logStyle {
		case "console":
			output := zerolog.ConsoleWriter{Out: G_LogFile, TimeFormat: constants.FORMAT, NoColor: true}
			opts = append(opts, zlogger.WithWriter(output))
		default:
			opts = append(opts, zlogger.WithWriter(G_LogFile))
		}
		log.Printf("GDD_LOG_PATH is configured. Please view app log from %s", logFile)
	}
	zlogger.InitEntry(zlogger.NewLoggerConfig(opts...))
}

type envVariable string

func (receiver envVariable) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(receiver.Load())), nil
}

const (
	NacosConfigType  = "nacos"
	ApolloConfigType = "apollo"
)

const (
	// GddBanner indicates banner enabled or not
	GddBanner envVariable = "GDD_BANNER"
	// GddBannerText sets text content of banner
	GddBannerText envVariable = "GDD_BANNER_TEXT"
	// GddLogLevel accepts panic, fatal, error, warn, warning, info, debug, trace, disabled. please reference zerolog.ParseLevel
	GddLogLevel envVariable = "GDD_LOG_LEVEL"
	// GddLogFormat text or json
	// Deprecated: move to zerolog
	GddLogFormat envVariable = "GDD_LOG_FORMAT"
	// GddLogReqEnable enables request and response logging
	GddLogReqEnable envVariable = "GDD_LOG_REQ_ENABLE"
	GddLogCaller    envVariable = "GDD_LOG_CALLER"
	GddLogDiscard   envVariable = "GDD_LOG_DISCARD"
	GddLogPath      envVariable = "GDD_LOG_PATH"
	// GddLogStyle is only valid when GDD_LOG_PATH is specified, accepts json or console.
	// Default is json
	GddLogStyle envVariable = "GDD_LOG_STYLE"
	// GddGraceTimeout sets graceful shutdown timeout
	GddGraceTimeout envVariable = "GDD_GRACE_TIMEOUT"
	// GddWriteTimeout sets http connection write timeout
	GddWriteTimeout envVariable = "GDD_WRITE_TIMEOUT"
	// GddReadTimeout sets http connection read timeout
	GddReadTimeout envVariable = "GDD_READ_TIMEOUT"
	// GddIdleTimeout sets http connection idle timeout
	GddIdleTimeout envVariable = "GDD_IDLE_TIMEOUT"
	// GddRouteRootPath sets root path for all routes
	GddRouteRootPath envVariable = "GDD_ROUTE_ROOT_PATH"
	// GddServiceName sets service name
	GddServiceName    envVariable = "GDD_SERVICE_NAME"
	GddServiceGroup   envVariable = "GDD_SERVICE_GROUP"
	GddServiceVersion envVariable = "GDD_SERVICE_VERSION"
	// GddHost sets bind host for http server
	GddHost envVariable = "GDD_HOST"
	// GddPort sets bind port for http server
	GddPort envVariable = "GDD_PORT"
	// GddGrpcPort sets bind port for grpc server
	GddGrpcPort envVariable = "GDD_GRPC_PORT"
	// GddManage if true, it will add built-in apis with /go-doudou path prefix for online api document and service status monitor etc.
	GddManage envVariable = "GDD_MANAGE_ENABLE"
	// GddManageUser manage api endpoint http basic auth user
	GddManageUser envVariable = "GDD_MANAGE_USER"
	// GddManagePass manage api endpoint http basic auth password
	GddManagePass envVariable = "GDD_MANAGE_PASS"

	GddEnableResponseGzip envVariable = "GDD_ENABLE_RESPONSE_GZIP"
	// Deprecated: move to GddFallbackContentType
	GddAppType envVariable = "GDD_APP_TYPE"
	// GddFallbackContentType fallback response content-type header value
	GddFallbackContentType        envVariable = "GDD_FALLBACK_CONTENTTYPE"
	GddRouterSaveMatchedRoutePath envVariable = "GDD_ROUTER_SAVEMATCHEDROUTEPATH"

	// GddConfigRemoteType has two options available: nacos, apollo
	GddConfigRemoteType envVariable = "GDD_CONFIG_REMOTE_TYPE"

	GddRetryCount         envVariable = "GDD_RETRY_COUNT"
	GddTracingMetricsRoot envVariable = "GDD_TRACING_METRICS_ROOT"

	GddServiceDiscoveryMode envVariable = "GDD_SERVICE_DISCOVERY_MODE"

	GddNacosNamespaceId         envVariable = "GDD_NACOS_NAMESPACE_ID"
	GddNacosTimeoutMs           envVariable = "GDD_NACOS_TIMEOUT_MS"
	GddNacosNotLoadCacheAtStart envVariable = "GDD_NACOS_NOTLOADCACHEATSTART"
	GddNacosLogDir              envVariable = "GDD_NACOS_LOG_DIR"
	GddNacosCacheDir            envVariable = "GDD_NACOS_CACHE_DIR"
	GddNacosLogLevel            envVariable = "GDD_NACOS_LOG_LEVEL"
	GddNacosLogDiscard          envVariable = "GDD_NACOS_LOG_DISCARD"
	GddNacosServerAddr          envVariable = "GDD_NACOS_SERVER_ADDR"
	GddNacosRegisterHost        envVariable = "GDD_NACOS_REGISTER_HOST"
	GddNacosClusterName         envVariable = "GDD_NACOS_CLUSTER_NAME"
	GddNacosGroupName           envVariable = "GDD_NACOS_GROUP_NAME"
	// GddNacosConfigFormat has two options available: dotenv, yaml
	GddNacosConfigFormat envVariable = "GDD_NACOS_CONFIG_FORMAT"
	GddNacosConfigGroup  envVariable = "GDD_NACOS_CONFIG_GROUP"
	GddNacosConfigDataid envVariable = "GDD_NACOS_CONFIG_DATAID"

	// GddWeight node weight
	GddWeight envVariable = "GDD_WEIGHT"

	GddApolloCluster      envVariable = "GDD_APOLLO_CLUSTER"
	GddApolloAddr         envVariable = "GDD_APOLLO_ADDR"
	GddApolloNamespace    envVariable = "GDD_APOLLO_NAMESPACE"
	GddApolloBackupEnable envVariable = "GDD_APOLLO_BACKUP_ENABLE"
	GddApolloBackupPath   envVariable = "GDD_APOLLO_BACKUP_PATH"
	GddApolloMuststart    envVariable = "GDD_APOLLO_MUSTSTART"
	GddApolloSecret       envVariable = "GDD_APOLLO_SECRET"
	GddApolloLogEnable    envVariable = "GDD_APOLLO_LOG_ENABLE"

	// GddSqlLogEnable only for doc purpose
	GddSqlLogEnable envVariable = "GDD_SQL_LOG_ENABLE"

	GddStatsFreq envVariable = "GDD_STATS_FREQ"

	GddRegisterHost  envVariable = "GDD_REGISTER_HOST"
	GddEtcdEndpoints envVariable = "GDD_ETCD_ENDPOINTS"
	GddEtcdLease     envVariable = "GDD_ETCD_LEASE"

	// configs for memberlist component
	// GddMemSeed sets cluster seeds for joining
	GddMemSeed envVariable = "GDD_MEM_SEED"
	// GddMemName unique name of this node in cluster. if empty or not set, hostname will be used instead
	GddMemName envVariable = "GDD_MEM_NAME"
	// GddMemHost specify AdvertiseAddr attribute of memberlist config struct.
	// if GddMemHost starts with dot such as .seed-svc-headless.default.svc.cluster.local,
	// it will be prefixed by hostname such as seed-2.seed-svc-headless.default.svc.cluster.local
	// for supporting k8s stateful service
	// if empty or not set, private ip will be used instead.
	GddMemHost envVariable = "GDD_MEM_HOST"
	// GddMemPort if empty or not set, an available port will be chosen randomly. recommend specifying a port
	GddMemPort envVariable = "GDD_MEM_PORT"
	// GddMemDeadTimeout dead node will be removed from node map if not received refute messages from it in GddMemDeadTimeout second
	// expose GossipToTheDeadTime property of memberlist.Config
	GddMemDeadTimeout envVariable = "GDD_MEM_DEAD_TIMEOUT"
	// GddMemSyncInterval local node will synchronize states from other random node every GddMemSyncInterval second
	// expose PushPullInterval property of memberlist.Config
	GddMemSyncInterval envVariable = "GDD_MEM_SYNC_INTERVAL"
	// GddMemReclaimTimeout dead node will be replaced with new node with the same name but different full address in GddMemReclaimTimeout second
	// expose DeadNodeReclaimTime property of memberlist.Config
	GddMemReclaimTimeout envVariable = "GDD_MEM_RECLAIM_TIMEOUT"
	// GddMemProbeInterval probe interval
	// expose ProbeInterval property of memberlist.Config
	GddMemProbeInterval envVariable = "GDD_MEM_PROBE_INTERVAL"
	// GddMemProbeTimeout probe timeout
	// expose ProbeTimeout property of memberlist.Config
	GddMemProbeTimeout envVariable = "GDD_MEM_PROBE_TIMEOUT"
	// GddMemSuspicionMult is the multiplier for determining the time an inaccessible node is considered suspect before declaring it dead.
	// expose SuspicionMult property of memberlist.Config
	GddMemSuspicionMult  envVariable = "GDD_MEM_SUSPICION_MULT"
	GddMemRetransmitMult envVariable = "GDD_MEM_RETRANSMIT_MULT"
	// GddMemGossipNodes how many remote nodes you want to gossip messages
	// expose GossipNodes property of memberlist.Config
	GddMemGossipNodes envVariable = "GDD_MEM_GOSSIP_NODES"
	// GddMemGossipInterval gossip interval
	// expose GossipInterval property of memberlist.Config
	GddMemGossipInterval envVariable = "GDD_MEM_GOSSIP_INTERVAL"
	// GddMemTCPTimeout tcp timeout
	// expose TCPTimeout property of memberlist.Config
	GddMemTCPTimeout envVariable = "GDD_MEM_TCP_TIMEOUT"
	// GddMemWeight node weight
	GddMemWeight envVariable = "GDD_MEM_WEIGHT"
	// GddMemWeightInterval node weight will be calculated every GddMemWeightInterval
	GddMemWeightInterval envVariable = "GDD_MEM_WEIGHT_INTERVAL"
	GddMemIndirectChecks envVariable = "GDD_MEM_INDIRECT_CHECKS"
	GddMemLogDisable     envVariable = "GDD_MEM_LOG_DISABLE"
	// GddMemCIDRsAllowed If not set, allow any connection (default), otherwise specify all networks
	// allowed connecting (you must specify IPv6/IPv4 separately)
	GddMemCIDRsAllowed envVariable = "GDD_MEM_CIDRS_ALLOWED"

	GddDBDisableAutoConfigure   envVariable = "GDD_DB_DISABLEAUTOCONFIGURE"
	GddDBDriver                 envVariable = "GDD_DB_DRIVER"
	GddDBDsn                    envVariable = "GDD_DB_DSN"
	GddDBTablePrefix            envVariable = "GDD_DB_TABLE_PREFIX"
	GddDBSkipDefaultTransaction envVariable = "GDD_DB_SKIPDEFAULTTRANSACTION"
	GddDBMaxIdleConns           envVariable = "GDD_DB_POOL_MAXIDLECONNS"
	GddDBMaxOpenConns           envVariable = "GDD_DB_POOL_MAXOPENCONNS"
	GddDBConnMaxLifetime        envVariable = "GDD_DB_POOL_CONNMAXLIFETIME"
	GddDBConnMaxIdleTime        envVariable = "GDD_DB_POOL_CONNMAXIDLETIME"
	GddDBPrepareStmt            envVariable = "GDD_DB_PREPARESTMT"

	GddDBLogSlowThreshold             envVariable = "GDD_DB_LOG_SLOWTHRESHOLD"
	GddDBLogIgnoreRecordNotFoundError envVariable = "GDD_DB_LOG_IGNORERECORDNOTFOUNDERROR"
	GddDBLogParameterizedQueries      envVariable = "GDD_DB_LOG_PARAMETERIZEDQUERIES"
	GddDBLogLevel                     envVariable = "GDD_DB_LOG_LEVEL"

	GddDBMysqlSkipInitializeWithVersion envVariable = "GDD_DB_MYSQL_SKIPINITIALIZEWITHVERSION"
	GddDBMysqlDefaultStringSize         envVariable = "GDD_DB_MYSQL_DEFAULTSTRINGSIZE"
	//GddDBMysqlDefaultDatetimePrecision      envVariable = "GDD_DB_MYSQL_DEFAULTDATETIMEPRECISION"
	GddDBMysqlDisableWithReturning          envVariable = "GDD_DB_MYSQL_DISABLEWITHRETURNING"
	GddDBMysqlDisableDatetimePrecision      envVariable = "GDD_DB_MYSQL_DISABLEDATETIMEPRECISION"
	GddDBMysqlDontSupportRenameIndex        envVariable = "GDD_DB_MYSQL_DONTSUPPORTRENAMEINDEX"
	GddDBMysqlDontSupportRenameColumn       envVariable = "GDD_DB_MYSQL_DONTSUPPORTRENAMECOLUMN"
	GddDBMysqlDontSupportForShareClause     envVariable = "GDD_DB_MYSQL_DONTSUPPORTFORSHARECLAUSE"
	GddDBMysqlDontSupportNullAsDefaultValue envVariable = "GDD_DB_MYSQL_DONTSUPPORTNULLASDEFAULTVALUE"
	GddDBMysqlDontSupportRenameColumnUnique envVariable = "GDD_DB_MYSQL_DONTSUPPORTRENAMECOLUMNUNIQUE"

	GddDBPostgresPreferSimpleProtocol envVariable = "GDD_DB_POSTGRES_PREFERSIMPLEPROTOCOL"
	GddDBPostgresWithoutReturning     envVariable = "GDD_DB_POSTGRES_WITHOUTRETURNING"

	GddDbPrometheusEnable          envVariable = "GDD_DB_PROMETHEUS_ENABLE"
	GddDbPrometheusRefreshInterval envVariable = "GDD_DB_PROMETHEUS_REFRESHINTERVAL"
	GddDbPrometheusDBName          envVariable = "GDD_DB_PROMETHEUS_DBNAME"

	GddDbCacheEnable envVariable = "GDD_DB_CACHE_ENABLE"
	GddCacheTTL      envVariable = "GDD_CACHE_TTL"
	GddCacheStores   envVariable = "GDD_CACHE_STORES"

	GddCacheRedisAddr           envVariable = "GDD_CACHE_REDIS_ADDR"
	GddCacheRedisUser           envVariable = "GDD_CACHE_REDIS_USER"
	GddCacheRedisPass           envVariable = "GDD_CACHE_REDIS_PASS"
	GddCacheRedisRouteByLatency envVariable = "GDD_CACHE_REDIS_ROUTEBYLATENCY"
	GddCacheRedisRouteRandomly  envVariable = "GDD_CACHE_REDIS_ROUTERANDOMLY"

	GddCacheRistrettoNumCounters envVariable = "GDD_CACHE_RISTRETTO_NUMCOUNTERS"
	GddCacheRistrettoMaxCost     envVariable = "GDD_CACHE_RISTRETTO_MAXCOST"
	GddCacheRistrettoBufferItems envVariable = "GDD_CACHE_RISTRETTO_BUFFERITEMS"

	GddCacheGocacheExpiration      envVariable = "GDD_CACHE_GOCACHE_EXPIRATION"
	GddCacheGocacheCleanupInterval envVariable = "GDD_CACHE_GOCACHE_CLEANUP_INTERVAL"

	GddZkServers          envVariable = "GDD_ZK_SERVERS"
	GddZkSequence         envVariable = "GDD_ZK_SEQUENCE"
	GddZkDirectoryPattern envVariable = "GDD_ZK_DIRECTORY_PATTERN"
)

// Load loads value from environment variable
func (receiver envVariable) Load() string {
	return os.Getenv(string(receiver))
}

func (receiver envVariable) LoadOrDefault(d string) string {
	val := d
	if stringutils.IsNotEmpty(receiver.Load()) {
		val = receiver.Load()
	}
	return val
}

func (receiver envVariable) LoadDurationOrDefault(d string) time.Duration {
	val, _ := time.ParseDuration(d)
	if stringutils.IsNotEmpty(receiver.Load()) {
		var err error
		val, err = time.ParseDuration(receiver.Load())
		if err != nil {
			panic(err)
		}
	}
	return val
}

// String return string representation for receiver
func (receiver envVariable) String() string {
	return receiver.Load()
}

// Write sets the environment variable to value
func (receiver envVariable) Write(value string) error {
	return os.Setenv(string(receiver), value)
}

func GetNacosClientParam() vo.NacosClientParam {
	namespaceId := DefaultGddNacosNamespaceId
	if stringutils.IsNotEmpty(GddNacosNamespaceId.Load()) {
		namespaceId = GddNacosNamespaceId.Load()
	}
	timeoutMs := DefaultGddNacosTimeoutMs
	if stringutils.IsNotEmpty(GddNacosTimeoutMs.Load()) {
		if t, err := cast.ToIntE(GddNacosTimeoutMs.Load()); err == nil {
			timeoutMs = t
		}
	}
	notLoadCacheAtStart := DefaultGddNacosNotLoadCacheAtStart
	if stringutils.IsNotEmpty(GddNacosNotLoadCacheAtStart.Load()) {
		notLoadCacheAtStart, _ = cast.ToBoolE(GddNacosNotLoadCacheAtStart.Load())
	}
	logDir := DefaultGddNacosLogDir
	if stringutils.IsNotEmpty(GddNacosLogDir.Load()) {
		logDir = GddNacosLogDir.Load()
	}
	cacheDir := DefaultGddNacosCacheDir
	if stringutils.IsNotEmpty(GddNacosCacheDir.Load()) {
		cacheDir = GddNacosCacheDir.Load()
	}
	logLevel := DefaultGddNacosLogLevel
	if stringutils.IsNotEmpty(GddNacosLogLevel.Load()) {
		logLevel = GddNacosLogLevel.Load()
	}
	logDiscard := DefaultGddNacosLogDiscard
	if stringutils.IsNotEmpty(GddNacosLogDiscard.Load()) {
		logDiscard, _ = cast.ToBoolE(GddNacosLogDiscard.Load())
	}
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(namespaceId),
		constant.WithTimeoutMs(uint64(timeoutMs)),
		constant.WithNotLoadCacheAtStart(notLoadCacheAtStart),
		constant.WithLogDir(logDir),
		constant.WithCacheDir(cacheDir),
		constant.WithLogLevel(logLevel),
		constant.WithLogDiscard(logDiscard),
	)
	serverAddrStr := DefaultGddNacosServerAddr
	if stringutils.IsNotEmpty(GddNacosServerAddr.Load()) {
		serverAddrStr = GddNacosServerAddr.Load()
	}
	var serverConfigs []constant.ServerConfig
	addrs := strings.Split(serverAddrStr, ",")
	for _, addr := range addrs {
		u, err := url.Parse(addr)
		if err != nil {
			panic(fmt.Errorf("[go-doudou] failed to create nacos discovery client: %v", err))
		}
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			panic(fmt.Errorf("[go-doudou] failed to create nacos discovery client: %v", err))
		}
		serverConfigs = append(serverConfigs, *constant.NewServerConfig(
			host,
			uint64(cast.ToInt(port)),
			constant.WithScheme(u.Scheme),
			constant.WithContextPath(u.Path),
		))
	}

	return vo.NacosClientParam{
		ClientConfig:  &clientConfig,
		ServerConfigs: serverConfigs,
	}
}

func GetServiceName() string {
	service := GddServiceName.LoadOrDefault(DefaultGddServiceName)
	if stringutils.IsEmpty(service) {
		zlogger.Panic().Msgf("[go-doudou] no value for environment variable %s found", string(GddServiceName))
	}
	return service
}

func GetPort() uint64 {
	httpPort := DefaultGddPort
	if stringutils.IsNotEmpty(GddPort.Load()) {
		if port, err := cast.ToIntE(GddPort.Load()); err == nil {
			httpPort = port
		}
	}
	return uint64(httpPort)
}

func GetGrpcPort() uint64 {
	grpcPort := DefaultGddGrpcPort
	if stringutils.IsNotEmpty(GddGrpcPort.Load()) {
		if port, err := cast.ToIntE(GddGrpcPort.Load()); err == nil {
			grpcPort = port
		}
	}
	return uint64(grpcPort)
}

func ServiceDiscoveryMap() map[string]struct{} {
	modeStr := GddServiceDiscoveryMode.LoadOrDefault(DefaultGddServiceDiscoveryMode)
	if stringutils.IsEmpty(modeStr) {
		return nil
	}
	modes := strings.Split(modeStr, ",")
	modemap := make(map[string]struct{})
	for _, mode := range modes {
		modemap[mode] = struct{}{}
	}
	return modemap
}

type Config struct {
	Db struct {
		Name   string
		Driver string
		Dsn    string
		Table  struct {
			// for postgresql schema only
			Prefix string
		}
		PrepareStmt              bool
		SkipDefaultTransaction   bool
		DisableNestedTransaction bool
		Log                      struct {
			Level                     string
			SlowThreshold             string `default:"200ms"`
			IgnoreRecordNotFoundError bool
			ParameterizedQueries      bool
		}
		Mysql struct {
			SkipInitializeWithVersion     bool
			DefaultStringSize             int
			DisableWithReturning          bool
			DisableDatetimePrecision      bool
			DontSupportRenameIndex        bool
			DontSupportRenameColumn       bool
			DontSupportForShareClause     bool
			DontSupportNullAsDefaultValue bool
			DontSupportRenameColumnUnique bool
		}
		Postgres struct {
			PreferSimpleProtocol bool
			WithoutReturning     bool
		}
		Pool struct {
			MaxIdleConns    int `default:"2"`
			MaxOpenConns    int `default:"-1"`
			ConnMaxLifetime string
			ConnMaxIdleTime string
		}
		Cache struct {
			Enable          bool
			ManualConfigure bool
		}
		Prometheus struct {
			Enable          bool
			RefreshInterval int `default:"15"`
			DBName          string
		}
	}
	Cache struct {
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
	}
	Service struct {
		Name string
	}
}

type gddConfig struct {
	Banner             bool
	BannerText         string        `default:"go-doudou" split_words:"true"`
	RouteRootPath      string        `default:"/" split_words:"true"`
	WriteTimeout       time.Duration `default:"15s" split_words:"true"`
	ReadTimeout        time.Duration `default:"15s" split_words:"true"`
	IdleTimeout        time.Duration `default:"60s" split_words:"true"`
	GraceTimeout       time.Duration `default:"15s" split_words:"true"`
	Port               string        `default:"6060"`
	Host               string
	EnableResponseGzip bool `default:"true" split_words:"true"`
	LogReqEnable       bool `default:"false" split_words:"true"`
	ManageEnable       bool `default:"true" split_words:"true"`
	Grpc               struct {
		Port string `default:"50051"`
	}
	Config
}
