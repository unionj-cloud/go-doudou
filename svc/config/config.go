package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func InitEnv() {
	env := os.Getenv("GDD_ENV")
	if "" == env {
		env = "dev"
	}
	wd, _ := os.Getwd()
	_ = godotenv.Load(filepath.Join(wd, ".env."+env+".local"))
	if "test" != env {
		_ = godotenv.Load(filepath.Join(wd, ".env.local"))
	}
	_ = godotenv.Load(filepath.Join(wd, ".env."+env))
	_ = godotenv.Load() // The Original .env
}

var (
	// BuildUser stores user who built the program
	BuildUser string
	// BuildTime stores time at which the program is built
	BuildTime string
	// GddVer stores go-doudou version when program is being built
	GddVer string
)

const FrameworkName = "Go-doudou"

type envVariable string

const (
	// GddBanner indicates banner enabled or not
	GddBanner envVariable = "GDD_BANNER"
	// GddBannerText sets text content of banner
	GddBannerText envVariable = "GDD_BANNER_TEXT"
	// GddLogLevel please reference logrus.ParseLevel
	GddLogLevel envVariable = "GDD_LOG_LEVEL"
	// GddLogPath sets log path
	GddLogPath envVariable = "GDD_LOG_PATH"
	// GddGraceTimeout sets graceful shutdown timeout
	GddGraceTimeout envVariable = "GDD_GRACE_TIMEOUT"
	// GddWriteTimeout sets http connection write timeout
	GddWriteTimeout envVariable = "GDD_WRITE_TIMEOUT"
	// GddReadTimeout sets http connection read timeout
	GddReadTimeout envVariable = "GDD_READ_TIMEOUT"
	// GddIdleTimeout sets http connection idle timeout
	GddIdleTimeout envVariable = "GDD_IDLE_TIMEOUT"
	// GddOutput sets file output path
	GddOutput envVariable = "GDD_OUTPUT"
	// GddRouteRootPath sets root path for all routes
	GddRouteRootPath envVariable = "GDD_ROUTE_ROOT_PATH"
	// GddServiceName sets service name
	GddServiceName envVariable = "GDD_SERVICE_NAME"
	// GddHost sets bind host for http server
	GddHost envVariable = "GDD_HOST"
	// GddPort sets bind port for http server
	GddPort envVariable = "GDD_PORT"
	// GddMode accepts 'mono' for monolith mode or 'micro' for microservice mode
	GddMode envVariable = "GDD_MODE"
	// GddManage if true, it will add built-in apis with /go-doudou path prefix for online api document and service status monitor etc.
	GddManage envVariable = "GDD_MANAGE_ENABLE"
	// GddManageUser manage api endpoint http basic auth user
	GddManageUser envVariable = "GDD_MANAGE_USER"
	// GddManagePass manage api endpoint http basic auth password
	GddManagePass envVariable = "GDD_MANAGE_PASS"
	// GddMemSeed sets cluster seeds for joining
	GddMemSeed envVariable = "GDD_MEM_SEED"
	// GddMemName unique name of this node in cluster. if empty or not set, hostname will be used instead
	GddMemName envVariable = "GDD_MEM_NAME"
	// GddMemHost specify AdvertiseAddr attribute of memberlist config struct.
	// if GddMemHost starts with dot such as .seed-svc-headless.default.svc.cluster.local,
	// it will be prefixed by hostname such as seed-2.seed-svc-headless.default.svc.cluster.local
	// for supporting k8s stateful service
	GddMemHost envVariable = "GDD_MEM_HOST"
	// GddMemPort if empty or not set, an available port will be chosen randomly. recommend specifying a port
	GddMemPort envVariable = "GDD_MEM_PORT"
	// GddMemDeadTimeout dead node will be removed from node map if not received refute messages from it in GddMemDeadTimeout second
	// expose GossipToTheDeadTime property of memberlist.Config
	// default value is 30 in second
	GddMemDeadTimeout envVariable = "GDD_MEM_DEAD_TIMEOUT"
	// GddMemSyncInterval local node will synchronize states from other random node every GddMemSyncInterval second
	// expose PushPullInterval property of memberlist.Config
	// default value is 5 in second
	GddMemSyncInterval envVariable = "GDD_MEM_SYNC_INTERVAL"
	// GddMemReclaimTimeout dead node will be replaced with new node with the same name but different full address in GddMemReclaimTimeout second
	// expose DeadNodeReclaimTime property of memberlist.Config
	// default value is 3 in second
	GddMemReclaimTimeout envVariable = "GDD_MEM_RECLAIM_TIMEOUT"
	// GddMemProbeInterval probe interval
	// expose ProbeInterval property of memberlist.Config
	// default value is 1 in second
	GddMemProbeInterval envVariable = "GDD_MEM_PROBE_INTERVAL"
	// GddMemProbeTimeout probe timeout
	// expose ProbeTimeout property of memberlist.Config
	// default value is 3 in second
	GddMemProbeTimeout envVariable = "GDD_MEM_PROBE_TIMEOUT"
	// GddMemSuspicionMult is the multiplier for determining the time an inaccessible node is considered suspect before declaring it dead.
	// expose SuspicionMult property of memberlist.Config
	// default value is 6
	GddMemSuspicionMult envVariable = "GDD_MEM_SUSPICION_MULT"
	// GddMemGossipNodes how many remote nodes you want to gossip messages
	// expose GossipNodes property of memberlist.Config
	// default value is 4
	GddMemGossipNodes envVariable = "GDD_MEM_GOSSIP_NODES"
	// GddMemGossipInterval gossip interval
	// expose GossipInterval property of memberlist.Config
	// default value is 500 in millisecond
	GddMemGossipInterval envVariable = "GDD_MEM_GOSSIP_INTERVAL"
	// GddMemTCPTimeout tcp timeout
	// expose TCPTimeout property of memberlist.Config
	// default value is 30 in second
	GddMemTCPTimeout envVariable = "GDD_MEM_TCP_TIMEOUT"
	// GddMemWeight node weight
	GddMemWeight envVariable = "GDD_MEM_WEIGHT"
	// GddMemWeightInterval node weight will be calculated every GddMemWeightInterval
	GddMemWeightInterval  envVariable = "GDD_MEM_WEIGHT_INTERVAL"
	GddRetryCount         envVariable = "GDD_RETRY_COUNT"
	GddEnableTracing      envVariable = "GDD_ENABLE_TRACING"
	GddTracingMetricsRoot envVariable = "GDD_TRACING_METRICS_ROOT"
)

// Load loads value from environment variable
func (receiver envVariable) Load() string {
	return os.Getenv(string(receiver))
}

// String return string representation for receiver
func (receiver envVariable) String() string {
	return string(receiver)
}

// Write sets the environment variable to value
func (receiver envVariable) Write(value string) error {
	return os.Setenv(string(receiver), value)
}

// Switch represents a switch
type Switch bool

// Decode decodes string on to true
func (s *Switch) Decode(value string) error {
	if value == "on" {
		*s = true
	}
	return nil
}

// LogLevel alias for logrus.Level
type LogLevel logrus.Level

// Decode decodes value to LogLevel
func (ll *LogLevel) Decode(value string) error {
	switch value {
	case "panic":
		*ll = LogLevel(logrus.PanicLevel)
	case "fatal":
		*ll = LogLevel(logrus.FatalLevel)
	case "error":
		*ll = LogLevel(logrus.ErrorLevel)
	case "warn":
		*ll = LogLevel(logrus.WarnLevel)
	case "debug":
		*ll = LogLevel(logrus.DebugLevel)
	case "trace":
		*ll = LogLevel(logrus.TraceLevel)
	default:
		*ll = LogLevel(logrus.InfoLevel)
	}
	return nil
}
