package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	_querystring "github.com/google/go-querystring/query"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	ddhttp "github.com/unionj-cloud/go-doudou/svc/http"
	"github.com/unionj-cloud/go-doudou/svc/registry"
)

type OcrClient struct {
	provider registry.IServiceProvider
	client   *resty.Client
}

func (receiver *OcrClient) SetProvider(provider registry.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *OcrClient) SetClient(client *resty.Client) {
	receiver.client = client
}
func (receiver *OcrClient) PostOcrCharacter(ctx context.Context,
	queryParams struct {
		MinHeight      int     `json:"minHeight,omitempty" url:"minHeight"`
		MinProbability float32 `json:"minProbability,omitempty" url:"minProbability"`
	},
	bodyJSON *os.File) (ret ResultListRecognizeCharacterResultVO, err error) {
	var _err error

	_req := receiver.client.R()
	_req.SetContext(ctx)
	_queryParams, _ := _querystring.Values(queryParams)
	_req.SetQueryParamsFromValues(_queryParams)
	_req.SetBody(bodyJSON)

	_resp, _err := _req.Post("/ocr/character")
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
func (receiver *OcrClient) PostOcrCharacterText(ctx context.Context,
	queryParams struct {
		MinHeight      int     `json:"minHeight,omitempty" url:"minHeight"`
		MinProbability float32 `json:"minProbability,omitempty" url:"minProbability"`
	},
	bodyJSON *os.File) (ret Resultstring, err error) {
	var _err error

	_req := receiver.client.R()
	_req.SetContext(ctx)
	_queryParams, _ := _querystring.Values(queryParams)
	_req.SetQueryParamsFromValues(_queryParams)
	_req.SetBody(bodyJSON)

	_resp, _err := _req.Post("/ocr/character/text")
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
func (receiver *OcrClient) PostOcrPdf(ctx context.Context,
	bodyJSON *os.File) (ret ResultRecognizePdfResultVO, err error) {
	var _err error

	_req := receiver.client.R()
	_req.SetContext(ctx)
	_req.SetBody(bodyJSON)

	_resp, _err := _req.Post("/ocr/pdf")
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
func (receiver *OcrClient) PostOcrPdfText(ctx context.Context,
	bodyJSON *os.File) (ret Resultstring, err error) {
	var _err error

	_req := receiver.client.R()
	_req.SetContext(ctx)
	_req.SetBody(bodyJSON)

	_resp, _err := _req.Post("/ocr/pdf/text")
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

func NewOcr(opts ...ddhttp.DdClientOption) *OcrClient {
	defaultProvider := ddhttp.NewServiceProvider("OCR")
	defaultClient := ddhttp.NewClient()

	svcClient := &OcrClient{
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
