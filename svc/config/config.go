package config

import (
	"github.com/sirupsen/logrus"
	"os"
)

type envVariable string

const (
	GddBanner     envVariable = "GDD_BANNER"
	GddBannerText envVariable = "GDD_BANNERTEXT"
	// GddLogLevel please reference logrus.ParseLevel
	GddLogLevel      envVariable = "GDD_LOGLEVEL"
	GddLogPath       envVariable = "GDD_LOGPATH"
	GddGraceTimeout  envVariable = "GDD_GRACETIMEOUT"
	GddWriteTimeout  envVariable = "GDD_WRITETIMEOUT"
	GddReadTimeout   envVariable = "GDD_READTIMEOUT"
	GddIdleTimeout   envVariable = "GDD_IDLETIMEOUT"
	GddOutput        envVariable = "GDD_OUTPUT"
	GddRouteRootPath envVariable = "GDD_ROUTE_ROOT_PATH"

	GddServiceName envVariable = "GDD_SERVICE_NAME"
	// GddNodeName if not provided, Memberlist will use hostname instead
	GddNodeName envVariable = "GDD_NODE_NAME"

	GddHost envVariable = "GDD_HOST"
	// GddPort http port
	GddPort    envVariable = "GDD_PORT"
	GddMemPort envVariable = "GDD_MEM_PORT"
	GddBaseUrl envVariable = "GDD_BASE_URL"
	GddSeed    envVariable = "GDD_SEED"
	// GddDepServices dependent service list
	GddDepServices envVariable = "GDD_DEP_SERVICES"
	// Accept 'mono' for monolith mode or 'micro' for microservice mode
	GddMode envVariable = "GDD_MODE"
	// GddManage if true, it will add built-in apis with /go-doudou path prefix for online api document and service status monitor etc.
	GddManage envVariable = "GDD_MANAGE_ENABLE"
	// GddManageUser manage api endpoint http basic auth user
	GddManageUser envVariable = "GDD_MANAGE_USER"
	// GddManagePass manage api endpoint http basic auth password
	GddManagePass envVariable = "GDD_MANAGE_PASS"
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
