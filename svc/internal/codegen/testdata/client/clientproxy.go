package client

import (
	"context"
	"mime/multipart"
	"os"
	service "testdata"
	"testdata/vo"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/slok/goresilience"
	"github.com/slok/goresilience/circuitbreaker"
	rerrors "github.com/slok/goresilience/errors"
	"github.com/slok/goresilience/metrics"
	"github.com/slok/goresilience/retry"
	"github.com/slok/goresilience/timeout"
	v3 "github.com/unionj-cloud/go-doudou/openapi/v3"
	"github.com/unionj-cloud/go-doudou/svc/config"
)

type UsersvcClientProxy struct {
	client service.Usersvc
	logger *logrus.Logger
	runner goresilience.Runner
}

func (receiver *UsersvcClientProxy) PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, msg error) {
	if _err := receiver.runner.Run(ctx, func(ctx context.Context) error {
		code, data, msg = receiver.client.PageUsers(
			ctx,
			query,
		)
		if msg != nil {
			return errors.Wrap(msg, "call PageUsers fail")
		}
		return nil
	}); _err != nil {
		// you can implement your fallback logic here
		if errors.Is(_err, rerrors.ErrCircuitOpen) {
			receiver.logger.Error(_err)
		}
		msg = errors.Wrap(_err, "call PageUsers fail")
	}
	return
}
func (receiver *UsersvcClientProxy) GetUser(ctx context.Context, userId string, photo string) (code int, data string, msg error) {
	if _err := receiver.runner.Run(ctx, func(ctx context.Context) error {
		code, data, msg = receiver.client.GetUser(
			ctx,
			userId,
			photo,
		)
		if msg != nil {
			return errors.Wrap(msg, "call GetUser fail")
		}
		return nil
	}); _err != nil {
		// you can implement your fallback logic here
		if errors.Is(_err, rerrors.ErrCircuitOpen) {
			receiver.logger.Error(_err)
		}
		msg = errors.Wrap(_err, "call GetUser fail")
	}
	return
}
func (receiver *UsersvcClientProxy) SignUp(ctx context.Context, username string, password int, actived bool, score []int) (code int, data string, msg error) {
	if _err := receiver.runner.Run(ctx, func(ctx context.Context) error {
		code, data, msg = receiver.client.SignUp(
			ctx,
			username,
			password,
			actived,
			score,
		)
		if msg != nil {
			return errors.Wrap(msg, "call SignUp fail")
		}
		return nil
	}); _err != nil {
		// you can implement your fallback logic here
		if errors.Is(_err, rerrors.ErrCircuitOpen) {
			receiver.logger.Error(_err)
		}
		msg = errors.Wrap(_err, "call SignUp fail")
	}
	return
}
func (receiver *UsersvcClientProxy) UploadAvatar(pc context.Context, pf []*v3.FileModel, ps string, pf2 *v3.FileModel, pf3 *multipart.FileHeader, pf4 []*multipart.FileHeader) (ri int, rs string, re error) {
	if _err := receiver.runner.Run(pc, func(ctx context.Context) error {
		ri, rs, re = receiver.client.UploadAvatar(
			pc,
			pf,
			ps,
			pf2,
			pf3,
			pf4,
		)
		if re != nil {
			return errors.Wrap(re, "call UploadAvatar fail")
		}
		return nil
	}); _err != nil {
		// you can implement your fallback logic here
		if errors.Is(_err, rerrors.ErrCircuitOpen) {
			receiver.logger.Error(_err)
		}
		re = errors.Wrap(_err, "call UploadAvatar fail")
	}
	return
}

type ProxyOption func(*UsersvcClientProxy)

func WithRunner(runner goresilience.Runner) ProxyOption {
	return func(proxy *UsersvcClientProxy) {
		proxy.runner = runner
	}
}

func WithLogger(logger *logrus.Logger) ProxyOption {
	return func(proxy *UsersvcClientProxy) {
		proxy.logger = logger
	}
}

func NewUsersvcClientProxy(client service.Usersvc, opts ...ProxyOption) *UsersvcClientProxy {
	cp := &UsersvcClientProxy{
		client: client,
		logger: logrus.StandardLogger(),
	}

	for _, opt := range opts {
		opt(cp)
	}

	if cp.runner == nil {
		var mid []goresilience.Middleware

		if config.GddManage.Load() == "true" {
			mid = append(mid, metrics.NewMiddleware("testdata_client", metrics.NewPrometheusRecorder(prometheus.DefaultRegisterer)))
		}

		mid = append(mid, circuitbreaker.NewMiddleware(circuitbreaker.Config{
			ErrorPercentThresholdToOpen:        50,
			MinimumRequestToOpen:               6,
			SuccessfulRequiredOnHalfOpen:       1,
			WaitDurationInOpenState:            5 * time.Second,
			MetricsSlidingWindowBucketQuantity: 10,
			MetricsBucketDuration:              1 * time.Second,
		}),
			timeout.NewMiddleware(timeout.Config{
				Timeout: 3 * time.Minute,
			}),
			retry.NewMiddleware(retry.Config{
				Times: 3,
			}))

		cp.runner = goresilience.RunnerChain(mid...)
	}

	return cp
}

func (receiver *UsersvcClientProxy) DownloadAvatar(ctx context.Context, userId string) (rf *os.File, re error) {
	if _err := receiver.runner.Run(ctx, func(ctx context.Context) error {
		rf, re = receiver.client.DownloadAvatar(
			ctx,
			userId,
		)
		if re != nil {
			return errors.Wrap(re, "call DownloadAvatar fail")
		}
		return nil
	}); _err != nil {
		// you can implement your fallback logic here
		if errors.Is(_err, rerrors.ErrCircuitOpen) {
			receiver.logger.Error(_err)
		}
		re = errors.Wrap(_err, "call DownloadAvatar fail")
	}
	return
}
