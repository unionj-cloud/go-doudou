package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	_querystring "github.com/google/go-querystring/query"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	"github.com/unionj-cloud/go-doudou/framework/registry"
)

type SignClient struct {
	provider registry.IServiceProvider
	client   *resty.Client
}

func (receiver *SignClient) SetProvider(provider registry.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *SignClient) SetClient(client *resty.Client) {
	receiver.client = client
}

// PostSignUp SignUp demonstrate how to define POST and Content-Type as application/x-www-form-urlencoded api
func (receiver *SignClient) PostSignUp(ctx context.Context,
	bodyParams SignUpReq) (ret SignUpResp, _resp *resty.Response, err error) {
	var _err error

	_req := receiver.client.R()
	_req.SetContext(ctx)
	_bodyParams, _ := _querystring.Values(bodyParams)
	_req.SetFormDataFromValues(_bodyParams)

	_resp, _err = _req.Post("/sign/up")
	if _err != nil {
		err = errors.Wrap(_err, "")
		return
	}
	if _resp.IsError() {
		err = errors.New(_resp.String())
		return
	}
	if _err = json.Unmarshal(_resp.Body(), &ret); _err != nil {
		err = errors.Wrap(_err, "")
		return
	}
	return
}

func NewSign(opts ...ddhttp.DdClientOption) *SignClient {
	defaultProvider := ddhttp.NewServiceProvider("SIGN")
	defaultClient := ddhttp.NewClient()

	svcClient := &SignClient{
		provider: defaultProvider,
		client:   defaultClient,
	}

	for _, opt := range opts {
		opt(svcClient)
	}

	svcClient.client.OnBeforeRequest(func(_ *resty.Client, request *resty.Request) error {
		request.URL = svcClient.provider.SelectServer() + request.URL
		return nil
	})

	svcClient.client.SetPreRequestHook(func(_ *resty.Client, request *http.Request) error {
		traceReq, _ := nethttp.TraceRequest(opentracing.GlobalTracer(), request,
			nethttp.OperationName(fmt.Sprintf("HTTP %s: %s", request.Method, request.RequestURI)))
		*request = *traceReq
		return nil
	})

	svcClient.client.OnAfterResponse(func(_ *resty.Client, response *resty.Response) error {
		nethttp.TracerFromRequest(response.Request.RawRequest).Finish()
		return nil
	})

	return svcClient
}
