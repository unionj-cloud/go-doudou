package grpcx

import (
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/unionj-cloud/go-doudou/v2/framework"
	"github.com/unionj-cloud/go-doudou/v2/framework/config"
	register "github.com/unionj-cloud/go-doudou/v2/framework/registry"
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
	if !config.CheckDev() {
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
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GddConfig.Grpc.Port))
	if err != nil {
		logger.Panic().Msgf("failed to listen: %v", err)
	}
	srv.ServeWithPipe(ln, pipe)
	defer func() {
		logger.Info().Msgf("Grpc server is gracefully shutting down in %s", config.GddConfig.GraceTimeout)
		// Make sure to set a deadline on exiting the process
		// after upg.Exit() is closed. No new upgrades can be
		// performed if the parent doesn't exit.
		time.AfterFunc(config.GddConfig.GraceTimeout, func() {
			logger.Error().Msg("Graceful shutdown timed out")
			os.Exit(1)
		})
		register.ShutdownGrpc()
		if err := timeutils.CallWithCtx(context.Background(), func() struct{} {
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

func (srv *GrpcServer) Serve(ln net.Listener) {
	srv.ServeWithPipe(ln, nil)
}

func (srv *GrpcServer) ServeWithPipe(ln net.Listener, pipe net.Listener) {
	if srv.Server == nil {
		return
	}
	framework.PrintBanner()
	framework.PrintLock.Lock()
	register.NewGrpc(srv.data)
	reflection.Register(srv)
	srv.printServices()
	go func() {
		if err := srv.Server.Serve(ln); err != nil {
			logger.Error().Msgf("failed to serve: %v", err)
		}
	}()
	if pipe != nil {
		go func() {
			if err := srv.Server.Serve(pipe); err != nil {
				logger.Error().Msgf("failed to serve: %v", err)
			}
		}()
	}
	logger.Info().Msgf("Grpc server is listening at %v", ln.Addr())
	logger.Info().Msgf("Grpc server started in %s", time.Since(startAt))
	framework.PrintLock.Unlock()
}
