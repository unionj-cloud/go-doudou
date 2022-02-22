package config

const FrameworkName = "Go-doudou"

const (
	// Default configs for framework component
	DefaultGddBanner             = true
	DefaultGddBannerText         = FrameworkName
	DefaultGddLogLevel           = "info"
	DefaultGddLogFormat          = "text"
	DefaultGddLogReqEnable       = false
	DefaultGddGraceTimeout       = "15s"
	DefaultGddWriteTimeout       = "15s"
	DefaultGddReadTimeout        = "15s"
	DefaultGddIdleTimeout        = "60s"
	DefaultGddServiceName        = ""
	DefaultGddRouteRootPath      = ""
	DefaultGddHost               = ""
	DefaultGddPort               = 6060
	DefaultGddRetryCount         = 0
	DefaultGddManage             = true
	DefaultGddManageUser         = "admin"
	DefaultGddManagePass         = "admin"
	DefaultGddTracingMetricsRoot = FrameworkName

	// Default configs for memberlist component
	DefaultGddMemSeed           = ""
	DefaultGddMemPort           = 7946
	DefaultGddMemDeadTimeout    = "60s"
	DefaultGddMemSyncInterval   = "60s"
	DefaultGddMemReclaimTimeout = "3s"
	DefaultGddMemProbeInterval  = "5s"
	DefaultGddMemProbeTimeout   = "3s"
	DefaultGddMemSuspicionMult  = 6
	DefaultGddMemGossipNodes    = 4
	DefaultGddMemGossipInterval = "500ms"
	DefaultGddMemTCPTimeout     = "30s"
	DefaultGddMemIndirectChecks = 3
	DefaultGddMemWeight         = 0
	DefaultGddMemWeightInterval = 0
	DefaultGddMemName           = ""
	DefaultGddMemHost           = ""
	DefaultGddMemCIDRsAllowed   = ""
	DefaultGddMemLogDisable     = false
)
