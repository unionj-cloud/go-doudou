package templates

var MainTmpl = EditableHeaderTmpl + `package main

import (
	{{- if ne .ProjectType "rest" }}
	"github.com/unionj-cloud/go-doudou/v2/framework/grpcx"
	{{- end }}
	"github.com/unionj-cloud/go-doudou/v2/framework/plugin"
	"github.com/unionj-cloud/go-doudou/v2/framework/rest"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	_ "{{.PluginPackage}}"
)

func main() {
	srv := rest.NewRestServer()
	{{- if ne .ProjectType "rest" }}
	grpcServer := grpcx.NewEmptyGrpcServer()
	{{- end }}
	plugins := plugin.GetServicePlugins()
	for _, key := range plugins.Keys() {
		value, _ := plugins.Get(key)
		{{- if eq .ProjectType "rest" }}
		value.Initialize(srv, nil, nil)
		{{- else }}
		value.Initialize(srv, grpcServer, nil)
		{{- end }}
	}
	defer func() {
		if r := recover(); r != nil {
			zlogger.Info().Msgf("Recovered. Error: %v\n", r)
		}
		for _, key := range plugins.Keys() {
			value, _ := plugins.Get(key)
			value.Close()
		}
	}()
	{{- if ne .ProjectType "rest" }}
	go func() {
		grpcServer.Run()
	}()
	{{- end }}
	srv.Run()
}
`
