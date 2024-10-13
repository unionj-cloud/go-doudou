package utils

import (
	"github.com/hashicorp/go-sockaddr"
	"github.com/unionj-cloud/go-doudou/v2/framework/config"
	"github.com/unionj-cloud/toolkit/stringutils"
	"github.com/unionj-cloud/toolkit/zlogger"
)

var GetPrivateIP = sockaddr.GetPrivateIP

func GetRegisterHost() string {
	registerHost := config.DefaultGddRegisterHost
	if stringutils.IsNotEmpty(config.GddRegisterHost.Load()) {
		registerHost = config.GddRegisterHost.Load()
	}
	if stringutils.IsEmpty(registerHost) {
		var err error
		registerHost, err = GetPrivateIP()
		if err != nil {
			zlogger.Panic().Err(err).Msg("[go-doudou] failed to get interface addresses")
		}
		if stringutils.IsEmpty(registerHost) {
			zlogger.Panic().Msg("[go-doudou] no private IP address found, and explicit IP not provided")
		}
	}
	return registerHost
}
