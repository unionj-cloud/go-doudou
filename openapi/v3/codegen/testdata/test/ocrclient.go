package test

import (
	"context"
	"encoding/json"
	"os"

	"github.com/go-resty/resty/v2"
	_querystring "github.com/google/go-querystring/query"
	"github.com/pkg/errors"
	ddhttp "github.com/unionj-cloud/go-doudou/svc/http"
)

type OcrClient struct {
	provider ddhttp.IServiceProvider
	client   *resty.Client
}

func (receiver *OcrClient) SetProvider(provider ddhttp.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *OcrClient) SetClient(client *resty.Client) {
	receiver.client = client
}
func (receiver *OcrClient) PostOcrCharacterText(ctx context.Context,
	queryParams struct {
		MinProbability float32 `json:"minProbability,omitempty" url:"minProbability"`
		MinHeight      int     `json:"minHeight,omitempty" url:"minHeight"`
	},
	bodyJSON *os.File) (ret Resultstring, err error) {
	var (
		_server string
		_err    error
	)
	if _server, _err = receiver.provider.SelectServer(); _err != nil {
		err = errors.Wrap(_err, "")
		return
	}

	_req := receiver.client.R()
	_req.SetContext(ctx)
	_queryParams, _ := _querystring.Values(queryParams)
	_req.SetQueryParamsFromValues(_queryParams)
	_req.SetBody(bodyJSON)

	_resp, _err := _req.Post(_server + "/ocr/character/text")
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
	var (
		_server string
		_err    error
	)
	if _server, _err = receiver.provider.SelectServer(); _err != nil {
		err = errors.Wrap(_err, "")
		return
	}

	_req := receiver.client.R()
	_req.SetContext(ctx)
	_req.SetBody(bodyJSON)

	_resp, _err := _req.Post(_server + "/ocr/pdf")
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
	var (
		_server string
		_err    error
	)
	if _server, _err = receiver.provider.SelectServer(); _err != nil {
		err = errors.Wrap(_err, "")
		return
	}

	_req := receiver.client.R()
	_req.SetContext(ctx)
	_req.SetBody(bodyJSON)

	_resp, _err := _req.Post(_server + "/ocr/pdf/text")
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
func (receiver *OcrClient) PostOcrCharacter(ctx context.Context,
	queryParams struct {
		MinHeight      int     `json:"minHeight,omitempty" url:"minHeight"`
		MinProbability float32 `json:"minProbability,omitempty" url:"minProbability"`
	},
	bodyJSON *os.File) (ret ResultListRecognizeCharacterResultVO, err error) {
	var (
		_server string
		_err    error
	)
	if _server, _err = receiver.provider.SelectServer(); _err != nil {
		err = errors.Wrap(_err, "")
		return
	}

	_req := receiver.client.R()
	_req.SetContext(ctx)
	_queryParams, _ := _querystring.Values(queryParams)
	_req.SetQueryParamsFromValues(_queryParams)
	_req.SetBody(bodyJSON)

	_resp, _err := _req.Post(_server + "/ocr/character")
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

	return svcClient
}
