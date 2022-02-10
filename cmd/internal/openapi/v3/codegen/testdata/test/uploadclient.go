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
	v3 "github.com/unionj-cloud/go-doudou/toolkit/openapi/v3"
)

type UploadClient struct {
	provider registry.IServiceProvider
	client   *resty.Client
}

func (receiver *UploadClient) SetProvider(provider registry.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *UploadClient) SetClient(client *resty.Client) {
	receiver.client = client
}

// PostUploadAvatar UploadAvatar demonstrate how to define upload files api
// there must be one []v3.FileModel or v3.FileModel parameter among input parameters
// remember to close the readers by Close method of v3.FileModel if you don't need them anymore when you finished your own business logic
func (receiver *UploadClient) PostUploadAvatar(ctx context.Context,
	bodyParams struct {
		// required
		Ps string `json:"ps,omitempty" url:"ps"`
	},
	pf []v3.FileModel) (ret UploadAvatarResp, _resp *resty.Response, err error) {
	var _err error

	_req := receiver.client.R()
	_req.SetContext(ctx)
	_bodyParams, _ := _querystring.Values(bodyParams)
	_req.SetFormDataFromValues(_bodyParams)
	if len(pf) == 0 {
		err = errors.New("at least one file should be uploaded for parameter pf")
		return
	}
	for _, _f := range pf {
		_req.SetFileReader("pf", _f.Filename, _f.Reader)
	}

	_resp, _err = _req.Post("/upload/avatar")
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

// PostUploadAvatar2 UploadAvatar2 demonstrate how to define upload files api
// remember to close the readers by Close method of v3.FileModel if you don't need them anymore when you finished your own business logic
func (receiver *UploadClient) PostUploadAvatar2(ctx context.Context,
	bodyParams struct {
		// required
		Ps string `json:"ps,omitempty" url:"ps"`
	},
	pf []v3.FileModel,
	pf2 *v3.FileModel,
	pf3 *v3.FileModel) (ret UploadAvatar2Resp, _resp *resty.Response, err error) {
	var _err error

	_req := receiver.client.R()
	_req.SetContext(ctx)
	_bodyParams, _ := _querystring.Values(bodyParams)
	_req.SetFormDataFromValues(_bodyParams)
	if len(pf) == 0 {
		err = errors.New("at least one file should be uploaded for parameter pf")
		return
	}
	for _, _f := range pf {
		_req.SetFileReader("pf", _f.Filename, _f.Reader)
	}
	if pf2 != nil {
		_req.SetFileReader("pf2", pf2.Filename, pf2.Reader)
	}
	if pf3 != nil {
		_req.SetFileReader("pf3", pf3.Filename, pf3.Reader)
	}

	_resp, _err = _req.Post("/upload/avatar/2")
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

func NewUpload(opts ...ddhttp.DdClientOption) *UploadClient {
	defaultProvider := ddhttp.NewServiceProvider("UPLOAD")
	defaultClient := ddhttp.NewClient()

	svcClient := &UploadClient{
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
