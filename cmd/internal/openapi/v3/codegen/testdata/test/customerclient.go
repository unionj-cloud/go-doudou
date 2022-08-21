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

type CustomerClient struct {
	provider registry.IServiceProvider
	client   *resty.Client
	rootPath string
}

func (receiver *CustomerClient) SetRootPath(rootPath string) {
	receiver.rootPath = rootPath
}

func (receiver *CustomerClient) SetProvider(provider registry.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *CustomerClient) SetClient(client *resty.Client) {
	receiver.client = client
}
func (receiver *CustomerClient) GetCustomerValidateToken(ctx context.Context, _headers map[string]string,
	queryParams struct {
		// required
		Token string `json:"token,omitempty" url:"token"`
	}) (ret bool, _resp *resty.Response, err error) {
	var _err error

	_req := receiver.client.R()
	_req.SetContext(ctx)
	if len(_headers) > 0 {
		_req.SetHeaders(_headers)
	}
	_queryParams, _ := _querystring.Values(queryParams)
	_req.SetQueryParamsFromValues(_queryParams)

	_resp, _err = _req.Get("/customer/validateToken")
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

func NewCustomer(opts ...ddhttp.DdClientOption) *CustomerClient {
	defaultProvider := ddhttp.NewServiceProvider("CUSTOMER")
	defaultClient := ddhttp.NewClient()

	svcClient := &CustomerClient{
		provider: defaultProvider,
		client:   defaultClient,
	}

	for _, opt := range opts {
		opt(svcClient)
	}

	svcClient.client.OnBeforeRequest(func(_ *resty.Client, request *resty.Request) error {
		request.URL = svcClient.provider.SelectServer() + svcClient.rootPath + request.URL
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
