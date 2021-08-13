package config

import (
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func init() {
	wd, _ := os.Getwd()
	err := godotenv.Load(filepath.Join(wd, ".env"))
	if err != nil {
		err = godotenv.Load(filepath.Join(wd, "../.env"))
		if err != nil {
			logrus.Warnln(errors.Wrap(err, "Error loading .env file"))
		}
	}
}

type envVariable string

const (
	GddBanner     envVariable = "GDD_BANNER"
	GddBannerText envVariable = "GDD_BANNER_TEXT"
	// GddLogLevel please reference logrus.ParseLevel
	GddLogLevel      envVariable = "GDD_LOG_LEVEL"
	GddLogPath       envVariable = "GDD_LOG_PATH"
	GddGraceTimeout  envVariable = "GDD_GRACE_TIMEOUT"
	GddWriteTimeout  envVariable = "GDD_WRITE_TIMEOUT"
	GddReadTimeout   envVariable = "GDD_READ_TIMEOUT"
	GddIdleTimeout   envVariable = "GDD_IDLE_TIMEOUT"
	GddOutput        envVariable = "GDD_OUTPUT"
	GddRouteRootPath envVariable = "GDD_ROUTE_ROOT_PATH"
	GddServiceName   envVariable = "GDD_SERVICE_NAME"
	GddHost          envVariable = "GDD_HOST"
	// GddPort http port
	GddPort    envVariable = "GDD_PORT"
	GddBaseUrl envVariable = "GDD_BASE_URL"
	// Accept 'mono' for monolith mode or 'micro' for microservice mode
	GddMode envVariable = "GDD_MODE"
	// GddManage if true, it will add built-in apis with /go-doudou path prefix for online api document and service status monitor etc.
	GddManage envVariable = "GDD_MANAGE_ENABLE"
	// GddManageUser manage api endpoint http basic auth user
	GddManageUser envVariable = "GDD_MANAGE_USER"
	// GddManagePass manage api endpoint http basic auth password
	GddManagePass envVariable = "GDD_MANAGE_PASS"

	// GddMemNodeName if not provided, hostname will be used instead
	GddMemNodeName envVariable = "GDD_MEM_NODE_NAME"
	GddMemSeed     envVariable = "GDD_MEM_SEED"
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
)

func (receiver envVariable) Load() string {
	return os.Getenv(string(receiver))
}

func (receiver envVariable) String() string {
	return string(receiver)
}

type Switch bool

func (s *Switch) Decode(value string) error {
	if value == "on" {
		*s = true
	}
	return nil
}

type LogLevel logrus.Level

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
