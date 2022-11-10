package config

import "github.com/unionj-cloud/go-doudou/v2/framework/configmgr"

const FrameworkName = "go-doudou"

const (
	// Default configs for framework component
	DefaultGddBanner             = true
	DefaultGddBannerText         = FrameworkName
	DefaultGddLogLevel           = "info"
	DefaultGddLogFormat          = "text"
	DefaultGddLogReqEnable       = false
	DefaultGddLogCaller          = false
	DefaultGddLogDiscard         = false
	DefaultGddGraceTimeout       = "15s"
	DefaultGddWriteTimeout       = "15s"
	DefaultGddReadTimeout        = "15s"
	DefaultGddIdleTimeout        = "60s"
	DefaultGddServiceName        = ""
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
)
