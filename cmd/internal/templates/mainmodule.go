package templates

var MainModuleTmpl = EditableHeaderTmpl + `package main

import (
	"github.com/unionj-cloud/go-doudou/v2/framework/grpcx"
	"github.com/unionj-cloud/go-doudou/v2/framework/plugin"
	"github.com/unionj-cloud/go-doudou/v2/framework/rest"
	_ "{{.PluginPackage}}"
)

func main() {
	srv := rest.NewRestServer()
	grpcServer := grpcx.NewEmptyGrpcServer()
	for _, v := range plugin.GetServicePlugins() {
		v.Initialize(srv, grpcServer, nil)
	}
	defer func() {
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