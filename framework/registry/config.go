package registry

import (
	"github.com/unionj-cloud/go-doudou/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/framework/memberlist"
	"github.com/unionj-cloud/go-doudou/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"strconv"
	"time"
)

func setGddMemDeadTimeout(conf *memberlist.Config) {
	deadTimeoutStr := config.GddMemDeadTimeout.Load()
	if stringutils.IsNotEmpty(deadTimeoutStr) {
		if deadTimeout, err := strconv.Atoi(deadTimeoutStr); err == nil {
			conf.GossipToTheDeadTime = time.Duration(deadTimeout) * time.Second
		} else {
			if duration, err := time.ParseDuration(deadTimeoutStr); err == nil {
				conf.GossipToTheDeadTime = duration
			}
		}
	}
}

func setGddMemSyncInterval(conf *memberlist.Config) {
	syncIntervalStr := config.GddMemSyncInterval.Load()
	if stringutils.IsNotEmpty(syncIntervalStr) {
		if syncInterval, err := strconv.Atoi(syncIntervalStr); err == nil {
			conf.PushPullInterval = time.Duration(syncInterval) * time.Second
		} else {
			if duration, err := time.ParseDuration(syncIntervalStr); err == nil {
				conf.PushPullInterval = duration
			}
		}
	}
}

func setGddMemReclaimTimeout(conf *memberlist.Config) {
	reclaimTimeoutStr := config.GddMemReclaimTimeout.Load()
	if stringutils.IsNotEmpty(reclaimTimeoutStr) {
		if reclaimTimeout, err := strconv.Atoi(reclaimTimeoutStr); err == nil {
			conf.DeadNodeReclaimTime = time.Duration(reclaimTimeout) * time.Second
		} else {
			if duration, err := time.ParseDuration(reclaimTimeoutStr); err == nil {
				conf.DeadNodeReclaimTime = duration
			}
		}
	}
}

func setGddMemGossipInterval(conf *memberlist.Config) {
	gossipIntervalStr := config.GddMemGossipInterval.Load()
	if stringutils.IsNotEmpty(gossipIntervalStr) {
		if gossipInterval, err := strconv.Atoi(gossipIntervalStr); err == nil {
			conf.GossipInterval = time.Duration(gossipInterval) * time.Millisecond
		} else {
			if duration, err := time.ParseDuration(gossipIntervalStr); err == nil {
				conf.GossipInterval = duration
			}
		}
	}
}

func setGddMemProbeInterval(conf *memberlist.Config) {
	probeIntervalStr := config.GddMemProbeInterval.Load()
	if stringutils.IsNotEmpty(probeIntervalStr) {
		if probeInterval, err := strconv.Atoi(probeIntervalStr); err == nil {
			conf.ProbeInterval = time.Duration(probeInterval) * time.Second
		} else {
			if duration, err := time.ParseDuration(probeIntervalStr); err == nil {
				conf.ProbeInterval = duration
			}
		}
	}
}

func setGddMemProbeTimeout(conf *memberlist.Config) {
	probeTimeoutStr := config.GddMemProbeTimeout.Load()
	if stringutils.IsNotEmpty(probeTimeoutStr) {
		if probeTimeout, err := strconv.Atoi(probeTimeoutStr); err == nil {
			conf.ProbeTimeout = time.Duration(probeTimeout) * time.Second
		} else {
			if duration, err := time.ParseDuration(probeTimeoutStr); err == nil {
				conf.ProbeTimeout = duration
			}
		}
	}
}

func setGddMemSuspicionMult(conf *memberlist.Config) {
	if eg, err := cast.ToIntE(config.GddMemSuspicionMult.Load()); err == nil {
		conf.SuspicionMult = eg
	}
}

func setGddMemRetransmitMult(conf *memberlist.Config) {
	if eg, err := cast.ToIntE(config.GddMemRetransmitMult.Load()); err == nil {
		conf.RetransmitMult = eg
	}
}

func setGddMemGossipNodes(conf *memberlist.Config) {
	if eg, err := cast.ToIntE(config.GddMemGossipNodes.Load()); err == nil {
		conf.GossipNodes = eg
	}
}

func setGddMemIndirectChecks(conf *memberlist.Config) {
	if eg, err := cast.ToIntE(config.GddMemIndirectChecks.Load()); err == nil {
		conf.IndirectChecks = eg
	}
}
