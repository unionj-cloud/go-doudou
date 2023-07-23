package codegen

import (
	"bytes"
	"fmt"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/templates"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/constants"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/fileutils"
	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/protobuf/v3"
	"github.com/unionj-cloud/go-doudou/v2/version"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var mainTmplGrpc = templates.EditableHeaderTmpl + `package main

import (
	` + gRPCImportBlock + `
)

func main() {
	conf := config.LoadFromEnv()
	svc := {{.ServiceAlias}}.New{{.SvcName}}(conf)
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
}
`

var appendMainTmplGrpc = `
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
`

var gRPCImportBlock = `
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
`

// GenMainGrpc generates main function for grpc service
func GenMainGrpc(dir string, ic astutils.InterfaceCollector, grpcSvc v3.Service) {
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
	if _, err = Stat(mainfile); os.IsNotExist(err) {
		if f, err = Create(mainfile); err != nil {
			panic(err)
		}
		defer f.Close()
		if tpl, err = template.New("main.go.tmpl").Parse(mainTmplGrpc); err != nil {
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
		}{
			ServicePackage: servicePkg,
			ConfigPackage:  cfgPkg,
			PbPackage:      pbPkg,
			SvcName:        svcName,
			ServiceAlias:   alias,
			Version:        version.Release,
			GrpcSvcName:    grpcSvc.Name,
		}); err != nil {
			panic(err)
		}
		source = strings.TrimSpace(buf.String())
		astutils.FixImport([]byte(source), mainfile)
	} else {
		if f, err = os.OpenFile(mainfile, os.O_APPEND, os.ModePerm); err != nil {
			panic(err)
		}
		defer f.Close()
		var original []byte
		if original, err = ioutil.ReadAll(f); err != nil {
			panic(err)
		}
		registerCall := fmt.Sprintf("pb.Register%sServer", grpcSvc.Name)
		if !strings.Contains(string(original), registerCall) {
			if tpl, err = template.New(appendMainTmplGrpc).Parse(appendMainTmplGrpc); err != nil {
				panic(err)
			}
			var buf bytes.Buffer
			if err = tpl.Execute(&buf, struct {
				GrpcSvcName string
			}{
				GrpcSvcName: grpcSvc.Name,
			}); err != nil {
				panic(err)
			}
			lines, err := fileutils.LinesFromReader(bytes.NewReader(original))
			if err != nil {
				panic(err)
			}
			fileContent := ""
			reg := regexp.MustCompile(fmt.Sprintf(`\.New%s\(`, svcName))
			for _, line := range lines {
				fileContent += line
				fileContent += constants.LineBreak
				if reg.MatchString(line) {
					fileContent += buf.String()
					fileContent += constants.LineBreak
				}
			}
			// fix import block
			var importBuf bytes.Buffer
			if tpl, err = template.New(gRPCImportBlock).Parse(gRPCImportBlock); err != nil {
				panic(err)
			}
			if err = tpl.Execute(&importBuf, struct {
				ServicePackage string
				ConfigPackage  string
				PbPackage      string
				SvcName        string
				ServiceAlias   string
				Version        string
				GrpcSvcName    string
			}{
				ServicePackage: servicePkg,
				ConfigPackage:  cfgPkg,
				PbPackage:      pbPkg,
				SvcName:        svcName,
				ServiceAlias:   alias,
				Version:        version.Release,
				GrpcSvcName:    grpcSvc.Name,
			}); err != nil {
				panic(err)
			}
			original = astutils.AppendImportStatements([]byte(fileContent), importBuf.Bytes())
			astutils.FixImport(original, mainfile)
		}
	}
}
