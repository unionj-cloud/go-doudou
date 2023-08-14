package grpcx

import (
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/unionj-cloud/go-doudou/v2/framework"
	"github.com/unionj-cloud/go-doudou/v2/framework/internal/banner"
	"github.com/unionj-cloud/go-doudou/v2/framework/internal/config"
	register "github.com/unionj-cloud/go-doudou/v2/framework/registry"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/timeutils"
	logger "github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

var startAt time.Time
var reflectionRegisterOnce sync.Once

func init() {
	startAt = time.Now()
}

type GrpcServer struct {
	*grpc.Server
	data map[string]interface{}
}

func NewGrpcServer(opt ...grpc.ServerOption) *GrpcServer {
	return &GrpcServer{
		Server: grpc.NewServer(opt...),
	}
}

func NewEmptyGrpcServer() *GrpcServer {
	return &GrpcServer{}
}

func NewGrpcServerWithData(data map[string]interface{}, opt ...grpc.ServerOption) *GrpcServer {
	server := GrpcServer{
		data: data,
	}
	server.Server = grpc.NewServer(opt...)
	return &server
}

func (srv *GrpcServer) printServices() {
	if !framework.CheckDev() {
		return
	}
	logger.Info().Msg("================ Registered Services ================")
	data := [][]string{}
	for k, v := range srv.GetServiceInfo() {
		for i, method := range v.Methods {
			if i == 0 {
				data = append(data, []string{k, method.Name})
			} else {
				data = append(data, []string{"", method.Name})
			}
		}
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"SERVICE", "RPC"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
	rows := strings.Split(strings.TrimSpace(tableString.String()), "\n")
	for _, row := range rows {
		logger.Info().Msg(row)
	}
	logger.Info().Msg("===================================================")
}

// Run runs grpc server
func (srv *GrpcServer) Run() {
	srv.RunWithPipe(nil)
}

// RunWithPipe runs grpc server
func (srv *GrpcServer) RunWithPipe(pipe net.Listener) {
	if srv.Server == nil {
		return
	}
	banner.Print()
	config.PrintLock.Lock()
	register.NewGrpc(srv.data)
	port := config.DefaultGddGrpcPort
	if p, err := cast.ToIntE(config.GddGrpcPort.Load()); err == nil {
		port = p
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Panic().Msgf("failed to listen: %v", err)
	}
	reflection.Register(srv)
	srv.printServices()
	go func() {
		if err := srv.Serve(lis); err != nil {
			logger.Error().Msgf("failed to serve: %v", err)
		}
	}()
	if pipe != nil {
		go func() {
			if err := srv.Serve(pipe); err != nil {
				logger.Error().Msgf("failed to serve: %v", err)
			}
		}()
	}
	logger.Info().Msgf("Grpc server is listening at %v", lis.Addr())
	logger.Info().Msgf("Grpc server started in %s", time.Since(startAt))
	config.PrintLock.Unlock()
	defer func() {
		register.ShutdownGrpc()

		grace, err := time.ParseDuration(config.GddGraceTimeout.Load())
		if err != nil {
			logger.Debug().Msgf("Parse %s %s as time.Duration failed: %s, use default %s instead.\n", string(config.GddGraceTimeout),
				config.GddGraceTimeout.Load(), err.Error(), config.DefaultGddGraceTimeout)
			grace, _ = time.ParseDuration(config.DefaultGddGraceTimeout)
		}
		logger.Info().Msgf("Grpc server is gracefully shutting down in %s", grace)

		ctx, cancel := context.WithTimeout(context.Background(), grace)
		defer cancel()
		if err := timeutils.CallWithCtx(ctx, func() struct{} {
			srv.GracefulStop()
			return struct{}{}
		}); err != nil {
			logger.Error().Err(err).Msg("")
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
}
