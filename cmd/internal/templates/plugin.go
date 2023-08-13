package templates

var PluginTmpl = EditableHeaderTmpl + `package plugin

import (
	{{.ServiceAlias}} "{{.ServicePackage}}"
	"{{.ConfigPackage}}"
	pb "{{.TransportGrpcPackage}}"
	"{{.TransportHttpPackage}}"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpczerolog "github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/unionj-cloud/go-doudou/v2/framework/grpcx"
	"github.com/unionj-cloud/go-doudou/v2/framework/plugin"
	"github.com/unionj-cloud/go-doudou/v2/framework/rest"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/pipeconn"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	"google.golang.org/grpc"
	"os"
)

var _ plugin.ServicePlugin = (*{{.SvcName}}Plugin)(nil)

type {{.SvcName}}Plugin struct {
	grpcConns []*grpc.ClientConn
}

func (receiver *{{.SvcName}}Plugin) Close() {
	for _, item := range receiver.grpcConns {
		item.Close()
	}
}

func (receiver *{{.SvcName}}Plugin) GoDoudouServicePlugin() {

}

func (receiver *{{.SvcName}}Plugin) GetName() string {
	name := os.Getenv("GDD_SERVICE_NAME")
	if stringutils.IsEmpty(name) {
		name = "cloud.unionj.{{.SvcName}}"
	}
	return name
}

func (receiver *{{.SvcName}}Plugin) Initialize(restServer *rest.RestServer, grpcServer *grpcx.GrpcServer, dialCtx pipeconn.DialContextFunc) {
	conf := config.LoadFromEnv()
	svc := {{.ServiceAlias}}.New{{.SvcName}}(conf)
	restServer.AddRoute(httpsrv.Routes(httpsrv.New{{.SvcName}}Handler(svc))...)
	if grpcServer.Server == nil {
		grpcServer.Server = grpc.NewServer(
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
	}
	pb.Register{{.SvcName}}ServiceServer(grpcServer, svc)
}

func init() {
	plugin.RegisterServicePlugin(&{{.SvcName}}Plugin{})
}
`