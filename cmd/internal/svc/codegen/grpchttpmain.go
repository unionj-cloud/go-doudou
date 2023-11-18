package codegen

import (
	"bytes"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/templates"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/protobuf/v3"
	"github.com/unionj-cloud/go-doudou/v2/version"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var mainTmplGrpcHttp = templates.EditableHeaderTmpl + `package main

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpczerolog "github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	"google.golang.org/grpc"
    "github.com/unionj-cloud/go-doudou/v2/framework/grpcx"
	{{.ServiceAlias}} "{{.ServicePackage}}"
    "{{.ConfigPackage}}"
	pb "{{.PbPackage}}"
	"{{.HttpPackage}}"
	"github.com/unionj-cloud/go-doudou/v2/framework/rest"
)

func main() {
	conf := config.LoadFromEnv()
	svc := {{.ServiceAlias}}.New{{.SvcName}}(conf)

	go func() {
		grpcServer := grpcx.NewGrpcServer(
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
				grpc_ctxtags.StreamServerInterceptor(),
				grpc_opentracing.StreamServerInterceptor(),
				grpc_prometheus.StreamServerInterceptor,
				tags.StreamServerInterceptor(tags.WithFieldExtractor(tags.CodeGenRequestFieldExtractor)),
				logging.StreamServerInterceptor(grpczerolog.InterceptorLogger(zlogger.Logger)),
				grpc_recovery.StreamServerInterceptor(),
			)),
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(),
				grpc_opentracing.UnaryServerInterceptor(),
				grpc_prometheus.UnaryServerInterceptor,
				tags.UnaryServerInterceptor(tags.WithFieldExtractor(tags.CodeGenRequestFieldExtractor)),
				logging.UnaryServerInterceptor(grpczerolog.InterceptorLogger(zlogger.Logger)),
				grpc_recovery.UnaryServerInterceptor(),
			)),
		)
		pb.Register{{.GrpcSvcName}}Server(grpcServer, svc)
		grpcServer.Run()
	}()

	handler := httpsrv.New{{.SvcName}}Http2Grpc(svc)
	srv := rest.NewRestServer()
	srv.AddRoute(httpsrv.Routes(handler)...)
	srv.Run()
}
`

// GenMainGrpcHttp generates main function for grpc service
func GenMainGrpcHttp(dir string, ic astutils.InterfaceCollector, grpcSvc v3.Service) {
	var (
		err      error
		mainfile string
		f        *os.File
		tpl      *template.Template
		cmdDir   string
		svcName  string
		alias    string
		source   string
	)
	cmdDir = filepath.Join(dir, "cmd")
	if err = MkdirAll(cmdDir, os.ModePerm); err != nil {
		panic(err)
	}
	svcName = ic.Interfaces[0].Name
	alias = ic.Package.Name
	mainfile = filepath.Join(cmdDir, "main.go")
	servicePkg := astutils.GetPkgPath(dir)
	cfgPkg := astutils.GetPkgPath(filepath.Join(dir, "config"))
	pbPkg := astutils.GetPkgPath(filepath.Join(dir, "transport", "grpc"))
	httpsrvPkg := astutils.GetPkgPath(filepath.Join(dir, "transport", "httpsrv"))
	if _, err = Stat(mainfile); os.IsNotExist(err) {
		if f, err = Create(mainfile); err != nil {
			panic(err)
		}
		defer f.Close()
		if tpl, err = template.New(mainTmplGrpcHttp).Parse(mainTmplGrpcHttp); err != nil {
			panic(err)
		}
		var buf bytes.Buffer
		if err = tpl.Execute(&buf, struct {
			ServicePackage string
			ConfigPackage  string
			PbPackage      string
			SvcName        string
			ServiceAlias   string
			Version        string
			GrpcSvcName    string
			HttpPackage    string
		}{
			ServicePackage: servicePkg,
			ConfigPackage:  cfgPkg,
			PbPackage:      pbPkg,
			SvcName:        svcName,
			ServiceAlias:   alias,
			Version:        version.Release,
			GrpcSvcName:    grpcSvc.Name,
			HttpPackage:    httpsrvPkg,
		}); err != nil {
			panic(err)
		}
		source = strings.TrimSpace(buf.String())
		astutils.FixImport([]byte(source), mainfile)
	}
}
