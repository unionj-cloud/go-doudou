package config

import (
	"github.com/unionj-cloud/go-doudou/v2/framework/configmgr"
	"gorm.io/gorm/logger"
)

const FrameworkName = "go-doudou"

const (
	// Default configs for framework component
	DefaultGddBanner             = true
	DefaultGddBannerText         = FrameworkName
	DefaultGddLogLevel           = "info"
	DefaultGddLogFormat          = "text"
	DefaultGddLogReqEnable       = false
	DefaultGddLogCaller          = true
	DefaultGddLogDiscard         = false
	DefaultGddGraceTimeout       = "15s"
	DefaultGddWriteTimeout       = "15s"
	DefaultGddReadTimeout        = "15s"
	DefaultGddIdleTimeout        = "60s"
	DefaultGddServiceName        = ""
	DefaultGddServiceGroup       = ""
	DefaultGddServiceVersion     = ""
	DefaultGddRouteRootPath      = ""
	DefaultGddHost               = ""
	DefaultGddPort               = 6060
	DefaultGddGrpcPort           = 50051
	DefaultGddRetryCount         = 0
	DefaultGddManage             = true
	DefaultGddManageUser         = "admin"
	DefaultGddManagePass         = "admin"
	DefaultGddTracingMetricsRoot = "tracing"
	DefaultGddWeight             = 1

	DefaultGddServiceDiscoveryMode = ""

	DefaultGddNacosNamespaceId         = "public"
	DefaultGddNacosTimeoutMs           = 10000
	DefaultGddNacosNotLoadCacheAtStart = false
	DefaultGddNacosLogDir              = "/tmp/nacos/log"
	DefaultGddNacosCacheDir            = "/tmp/nacos/cache"
	DefaultGddNacosLogLevel            = "info"
	DefaultGddNacosLogDiscard          = false
	DefaultGddNacosServerAddr          = ""
	DefaultGddNacosRegisterHost        = ""
	DefaultGddNacosClusterName         = "DEFAULT"
	DefaultGddNacosGroupName           = "DEFAULT_GROUP"

	DefaultGddNacosConfigFormat = configmgr.DotenvConfigFormat
	DefaultGddNacosConfigGroup  = "DEFAULT_GROUP"
	DefaultGddNacosConfigDataid = ""

	DefaultGddEnableResponseGzip         = true
	DefaultGddAppType                    = "rest"
	DefaultGddFallbackContentType        = "application/json; charset=UTF-8"
	DefaultGddRouterSaveMatchedRoutePath = true
	DefaultGddConfigRemoteType           = ""

	DefaultGddApolloCluster      = "default"
	DefaultGddApolloAddr         = ""
	DefaultGddApolloNamespace    = "application.properties"
	DefaultGddApolloBackupEnable = true
	DefaultGddApolloBackupPath   = ""
	DefaultGddApolloSecret       = ""
	DefaultGddApolloMuststart    = false
	DefaultGddApolloLogEnable    = false

	// DefaultGddSqlLogEnable only for doc purpose
	DefaultGddSqlLogEnable = false

	DefaultGddStatsFreq = "1s"

	DefaultGddRegisterHost        = ""
	DefaultGddEtcdEndpoints       = ""
	DefaultGddEtcdLease     int64 = 5

	// Default configs for memberlist component
	DefaultGddMemSeed           = ""
	DefaultGddMemPort           = 7946
	DefaultGddMemDeadTimeout    = "60s"
	DefaultGddMemSyncInterval   = "60s"
	DefaultGddMemReclaimTimeout = "3s"
	DefaultGddMemProbeInterval  = "5s"
	DefaultGddMemProbeTimeout   = "3s"
	DefaultGddMemSuspicionMult  = 6
	DefaultGddMemRetransmitMult = 4
	DefaultGddMemGossipNodes    = 4
	DefaultGddMemGossipInterval = "500ms"
	DefaultGddMemTCPTimeout     = "30s"
	DefaultGddMemIndirectChecks = 3
	DefaultGddMemWeight         = 1
	DefaultGddMemWeightInterval = 0
	DefaultGddMemName           = ""
	DefaultGddMemHost           = ""
	DefaultGddMemCIDRsAllowed   = ""
	DefaultGddMemLogDisable     = false

	DefaultGddDBDisableAutoConfigure = false
	DefaultGddDBDriver               = ""
	DefaultGddDBDsn                  = ""
	DefaultGddDBTablePrefix          = ""
	DefaultGddDBMaxIdleConns         = 2
	DefaultGddDBMaxOpenConns         = -1
	DefaultGddDBConnMaxLifetime      = -1
	DefaultGddDBConnMaxIdleTime      = -1

	DefaultGddDBLogSlowThreshold             = "200ms"
	DefaultGddDBLogIgnoreRecordNotFoundError = false
	DefaultGddDBLogParameterizedQueries      = false
	DefaultGddDBLogLevel                     = logger.Warn

	DefaultGddDBMysqlSkipInitializeWithVersion = false
	DefaultGddDBMysqlDefaultStringSize         = 0
	//DefaultGddDBMysqlDefaultDatetimePrecision      envVariable = "GDD_DB_MYSQL_DEFAULTDATETIMEPRECISION"
	DefaultGddDBMysqlDisableWithReturning          = false
	DefaultGddDBMysqlDisableDatetimePrecision      = false
	DefaultGddDBMysqlDontSupportRenameIndex        = false
	DefaultGddDBMysqlDontSupportRenameColumn       = false
	DefaultGddDBMysqlDontSupportForShareClause     = false
	DefaultGddDBMysqlDontSupportNullAsDefaultValue = false
	DefaultGddDBMysqlDontSupportRenameColumnUnique = false

	DefaultGddDBPostgresPreferSimpleProtocol = false
	DefaultGddDBPostgresWithoutReturning     = false

	DefaultGddZkServers          = ""
	DefaultGddZkSequence         = false
	DefaultGddZkDirectoryPattern = "/registry/%s/providers"
)
