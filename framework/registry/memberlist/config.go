package memberlist

import (
	"github.com/unionj-cloud/go-doudou/v2/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/memberlist"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"strconv"
	"time"
)

func setGddMemDeadTimeout(conf *memberlist.Config) {
	var isSet bool
	deadTimeoutStr := config.GddMemDeadTimeout.Load()
	if stringutils.IsNotEmpty(deadTimeoutStr) {
		if deadTimeout, err := strconv.Atoi(deadTimeoutStr); err == nil {
			conf.GossipToTheDeadTime = time.Duration(deadTimeout) * time.Second
			isSet = true
		} else {
			if duration, err := time.ParseDuration(deadTimeoutStr); err == nil {
				conf.GossipToTheDeadTime = duration
				isSet = true
			}
		}
	}
	if !isSet {
		conf.GossipToTheDeadTime, _ = time.ParseDuration(config.DefaultGddMemDeadTimeout)
	}
}

func setGddMemSyncInterval(conf *memberlist.Config) {
	var isSet bool
	syncIntervalStr := config.GddMemSyncInterval.Load()
	if stringutils.IsNotEmpty(syncIntervalStr) {
		if syncInterval, err := strconv.Atoi(syncIntervalStr); err == nil {
			conf.PushPullInterval = time.Duration(syncInterval) * time.Second
			isSet = true
		} else {
			if duration, err := time.ParseDuration(syncIntervalStr); err == nil {
				conf.PushPullInterval = duration
				isSet = true
			}
		}
	}
	if !isSet {
		conf.PushPullInterval, _ = time.ParseDuration(config.DefaultGddMemSyncInterval)
	}
}

func setGddMemReclaimTimeout(conf *memberlist.Config) {
	var isSet bool
	reclaimTimeoutStr := config.GddMemReclaimTimeout.Load()
	if stringutils.IsNotEmpty(reclaimTimeoutStr) {
		if reclaimTimeout, err := strconv.Atoi(reclaimTimeoutStr); err == nil {
			conf.DeadNodeReclaimTime = time.Duration(reclaimTimeout) * time.Second
			isSet = true
		} else {
			if duration, err := time.ParseDuration(reclaimTimeoutStr); err == nil {
				conf.DeadNodeReclaimTime = duration
				isSet = true
			}
		}
	}
	if !isSet {
		conf.DeadNodeReclaimTime, _ = time.ParseDuration(config.DefaultGddMemReclaimTimeout)
	}
}

func setGddMemGossipInterval(conf *memberlist.Config) {
	var isSet bool
	gossipIntervalStr := config.GddMemGossipInterval.Load()
	if stringutils.IsNotEmpty(gossipIntervalStr) {
		if gossipInterval, err := strconv.Atoi(gossipIntervalStr); err == nil {
			conf.GossipInterval = time.Duration(gossipInterval) * time.Millisecond
			isSet = true
		} else {
			if duration, err := time.ParseDuration(gossipIntervalStr); err == nil {
				conf.GossipInterval = duration
				isSet = true
			}
		}
	}
	if !isSet {
		conf.GossipInterval, _ = time.ParseDuration(config.DefaultGddMemGossipInterval)
	}
}

func setGddMemProbeInterval(conf *memberlist.Config) {
	var isSet bool
	probeIntervalStr := config.GddMemProbeInterval.Load()
	if stringutils.IsNotEmpty(probeIntervalStr) {
		if probeInterval, err := strconv.Atoi(probeIntervalStr); err == nil {
			conf.ProbeInterval = time.Duration(probeInterval) * time.Second
			isSet = true
		} else {
			if duration, err := time.ParseDuration(probeIntervalStr); err == nil {
				conf.ProbeInterval = duration
				isSet = true
			}
		}
	}
	if !isSet {
		conf.ProbeInterval, _ = time.ParseDuration(config.DefaultGddMemProbeInterval)
	}
}

func setGddMemProbeTimeout(conf *memberlist.Config) {
	var set bool
	probeTimeoutStr := config.GddMemProbeTimeout.Load()
	if stringutils.IsNotEmpty(probeTimeoutStr) {
		if probeTimeout, err := strconv.Atoi(probeTimeoutStr); err == nil {
			conf.ProbeTimeout = time.Duration(probeTimeout) * time.Second
			set = true
		} else {
			if duration, err := time.ParseDuration(probeTimeoutStr); err == nil {
				conf.ProbeTimeout = duration
				set = true
			}
		}
	}
	if !set {
		conf.ProbeTimeout, _ = time.ParseDuration(config.DefaultGddMemProbeTimeout)
	}
}

func setGddMemSuspicionMult(conf *memberlist.Config) {
	var set bool
	if eg, err := cast.ToIntE(config.GddMemSuspicionMult.Load()); err == nil {
		conf.SuspicionMult = eg
		set = true
	}
	if !set {
		conf.SuspicionMult = config.DefaultGddMemSuspicionMult
	}
}

func setGddMemRetransmitMult(conf *memberlist.Config) {
	var set bool
	if eg, err := cast.ToIntE(config.GddMemRetransmitMult.Load()); err == nil {
		conf.RetransmitMult = eg
		set = true
	}
	if !set {
		conf.RetransmitMult = config.DefaultGddMemRetransmitMult
	}
}

func setGddMemGossipNodes(conf *memberlist.Config) {
	var set bool
	if eg, err := cast.ToIntE(config.GddMemGossipNodes.Load()); err == nil {
		conf.GossipNodes = eg
		set = true
	}
	if !set {
		conf.GossipNodes = config.DefaultGddMemGossipNodes
	}
}

func setGddMemIndirectChecks(conf *memberlist.Config) {
	var set bool
	if eg, err := cast.ToIntE(config.GddMemIndirectChecks.Load()); err == nil {
		conf.IndirectChecks = eg
		set = true
	}
	if !set {
		conf.IndirectChecks = config.DefaultGddMemIndirectChecks
	}
}
