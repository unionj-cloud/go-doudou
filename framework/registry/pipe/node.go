package pipe

import (
	"context"
	"net"
	"time"

	"github.com/unionj-cloud/go-doudou/v2/toolkit/pipeconn"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGrpcClientConn(dialCtx pipeconn.DialContextFunc) *grpc.ClientConn {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	grpcConn, err := grpc.DialContext(ctx, "pipe",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(c context.Context, _ string) (net.Conn, error) {
			return dialCtx(c)
		}),
	)
	if err != nil {
		zlogger.Panic().Err(err).Msgf("[go-doudou] failed to connect to pipe")
	}
	return grpcConn
}
