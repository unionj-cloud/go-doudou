package client

import (
	"context"
	"mime/multipart"
	"os"
	"testdata/vo"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/slok/goresilience"
	"github.com/slok/goresilience/circuitbreaker"
	rerrors "github.com/slok/goresilience/errors"
	"github.com/slok/goresilience/metrics"
	"github.com/slok/goresilience/retry"
	"github.com/slok/goresilience/timeout"
	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/openapi/v3"
)

type UsersvcClientProxy struct {
	client *UsersvcClient
	logger *logrus.Logger
	runner goresilience.Runner
}

func (receiver *UsersvcClientProxy) PageUsers(ctx context.Context, _headers map[string]string, query vo.PageQuery) (_resp *resty.Response, code int, data vo.PageRet, msg error) {
	if _err := receiver.runner.Run(ctx, func(ctx context.Context) error {
		_resp, code, data, msg = receiver.client.PageUsers(
			ctx,
			_headers,
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
func (receiver *UsersvcClientProxy) GetUser(ctx context.Context, _headers map[string]string, userId string, photo string) (_resp *resty.Response, code int, data string, msg error) {
	if _err := receiver.runner.Run(ctx, func(ctx context.Context) error {
		_resp, code, data, msg = receiver.client.GetUser(
			ctx,
			_headers,
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
func (receiver *UsersvcClientProxy) SignUp(ctx context.Context, _headers map[string]string, username string, password int, actived bool, score []int) (_resp *resty.Response, code int, data string, msg error) {
	if _err := receiver.runner.Run(ctx, func(ctx context.Context) error {
		_resp, code, data, msg = receiver.client.SignUp(
			ctx,
			_headers,
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
func (receiver *UsersvcClientProxy) UploadAvatar(ctx context.Context, _headers map[string]string, pf []v3.FileModel, ps string, pf2 v3.FileModel, pf3 *multipart.FileHeader, pf4 []*multipart.FileHeader) (_resp *resty.Response, ri int, ri2 interface{}, re error) {
	if _err := receiver.runner.Run(ctx, func(ctx context.Context) error {
		_resp, ri, ri2, re = receiver.client.UploadAvatar(
			ctx,
			_headers,
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
func (receiver *UsersvcClientProxy) DownloadAvatar(ctx context.Context, _headers map[string]string, userId interface{}, userAttrs ...string) (_resp *resty.Response, rf *os.File, re error) {
	if _err := receiver.runner.Run(ctx, func(ctx context.Context) error {
		_resp, rf, re = receiver.client.DownloadAvatar(
			ctx,
			_headers,
			userId,
			userAttrs...,
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

func NewUsersvcClientProxy(client *UsersvcClient, opts ...ProxyOption) *UsersvcClientProxy {
	cp := &UsersvcClientProxy{
		client: client,
		logger: logrus.StandardLogger(),
	}

	for _, opt := range opts {
		opt(cp)
	}

	if cp.runner == nil {
		var mid []goresilience.Middleware
		mid = append(mid, metrics.NewMiddleware("testdata_client", metrics.NewPrometheusRecorder(prometheus.DefaultRegisterer)))
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
