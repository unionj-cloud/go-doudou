package templates

var MainModuleTmpl = EditableHeaderTmpl + `package main

import (
	"github.com/unionj-cloud/go-doudou/v2/framework/grpcx"
	"github.com/unionj-cloud/go-doudou/v2/framework/plugin"
	"github.com/unionj-cloud/go-doudou/v2/framework/rest"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	_ "{{.PluginPackage}}"
)

func main() {
	srv := rest.NewRestServer()
	grpcServer := grpcx.NewEmptyGrpcServer()
	for _, v := range plugin.GetServicePlugins() {
		v.Initialize(srv, grpcServer, nil)
	}
	defer func() {
		if r := recover(); r != nil {
			zlogger.Info().Msgf("Recovered. Error: %v\n", r)
		}
		for _, v := range plugin.GetServicePlugins() {
			v.Close()
		}
	}()
	go func() {
		grpcServer.Run()
	}()
	srv.Run()
}
`
