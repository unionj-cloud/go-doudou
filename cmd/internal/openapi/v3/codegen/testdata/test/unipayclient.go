package test

import (
	"context"
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

type UnipayClient struct {
	provider registry.IServiceProvider
	client   *resty.Client
}

func (receiver *UnipayClient) SetProvider(provider registry.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *UnipayClient) SetClient(client *resty.Client) {
	receiver.client = client
}
func (receiver *UnipayClient) GetUnipayStartUnionPay(ctx context.Context,
	queryParams struct {
		// required
		TxnAmt string `json:"txnAmt,omitempty" url:"txnAmt"`
		// required
		Token string `json:"token,omitempty" url:"token"`
		// required
		CompanyId string `json:"companyId,omitempty" url:"companyId"`
		// required
		FrontUrl string `json:"frontUrl,omitempty" url:"frontUrl"`
	}) (ret string, _resp *resty.Response, err error) {
	var _err error

	_req := receiver.client.R()
	_req.SetContext(ctx)
	_queryParams, _ := _querystring.Values(queryParams)
	_req.SetQueryParamsFromValues(_queryParams)

	_resp, _err = _req.Get("/unipay/startUnionPay")
	if _err != nil {
		err = errors.Wrap(_err, "")
		return
	}
	if _resp.IsError() {
		err = errors.New(_resp.String())
		return
	}
	ret = _resp.String()
	return
}

func NewUnipay(opts ...ddhttp.DdClientOption) *UnipayClient {
	defaultProvider := ddhttp.NewServiceProvider("UNIPAY")
	defaultClient := ddhttp.NewClient()

	svcClient := &UnipayClient{
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
			nethttp.OperationName(fmt.Sprintf("HTTP %s: %s", request.Method, request.URL.Path)))
		*request = *traceReq
		return nil
	})

	svcClient.client.OnAfterResponse(func(_ *resty.Client, response *resty.Response) error {
		nethttp.TracerFromRequest(response.Request.RawRequest).Finish()
		return nil
	})

	return svcClient
}
