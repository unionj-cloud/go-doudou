package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	"github.com/unionj-cloud/go-doudou/framework/registry"
)

type PageClient struct {
	provider registry.IServiceProvider
	client   *resty.Client
}

func (receiver *PageClient) SetProvider(provider registry.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *PageClient) SetClient(client *resty.Client) {
	receiver.client = client
}

// PostPageUsers2 PageUsers2 demonstrate how to define POST and Content-Type as application/json api
func (receiver *PageClient) PostPageUsers2(ctx context.Context, _headers map[string]string,
	// comments above input and output struct type parameters in vo package will display on online document
	// not comments here
	bodyJSON *PageQuery) (ret PageUsers2Resp, _resp *resty.Response, err error) {
	var _err error

	_req := receiver.client.R()
	_req.SetContext(ctx)
	if len(_headers) > 0 {
		_req.SetHeaders(_headers)
	}
	_req.SetBody(bodyJSON)

	_resp, _err = _req.Post("/page/users/2")
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

// PostPageUsers PageUsers demonstrate how to define POST and Content-Type as application/json api
func (receiver *PageClient) PostPageUsers(ctx context.Context, _headers map[string]string,
	// comments above input and output struct type parameters in vo package will display on online document
	// not comments here
	bodyJSON PageQuery) (ret PageUsersResp, _resp *resty.Response, err error) {
	var _err error

	_req := receiver.client.R()
	_req.SetContext(ctx)
	if len(_headers) > 0 {
		_req.SetHeaders(_headers)
	}
	_req.SetBody(bodyJSON)

	_resp, _err = _req.Post("/page/users")
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

func NewPage(opts ...ddhttp.DdClientOption) *PageClient {
	defaultProvider := ddhttp.NewServiceProvider("PAGE")
	defaultClient := ddhttp.NewClient()

	svcClient := &PageClient{
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
