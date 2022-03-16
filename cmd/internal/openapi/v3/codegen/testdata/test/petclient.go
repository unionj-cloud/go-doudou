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

type PetClient struct {
	provider registry.IServiceProvider
	client   *resty.Client
	rootPath string
}

func (receiver *PetClient) SetRootPath(rootPath string) {
	receiver.rootPath = rootPath
}

func (receiver *PetClient) SetProvider(provider registry.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *PetClient) SetClient(client *resty.Client) {
	receiver.client = client
}

// GetPetPetId Find pet by ID
// Returns a single pet
func (receiver *PetClient) GetPetPetId(ctx context.Context, _headers map[string]string,
	// ID of pet to return
	// required
	petId int64) (ret Pet, _resp *resty.Response, err error) {
	var _err error

	_req := receiver.client.R()
	_req.SetContext(ctx)
	if len(_headers) > 0 {
		_req.SetHeaders(_headers)
	}
	_req.SetPathParam("petId", fmt.Sprintf("%v", petId))

	_resp, _err = _req.Get("/pet/{petId}")
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

// GetPetFindByTags Finds Pets by tags
// Multiple tags can be provided with comma separated strings. Use tag1, tag2, tag3 for testing.
func (receiver *PetClient) GetPetFindByTags(ctx context.Context, _headers map[string]string,
	queryParams *struct {
		Tags *[]string `json:"tags,omitempty" url:"tags"`
	}) (ret []Pet, _resp *resty.Response, err error) {
	var _err error

	_req := receiver.client.R()
	_req.SetContext(ctx)
	if len(_headers) > 0 {
		_req.SetHeaders(_headers)
	}
	_queryParams, _ := _querystring.Values(queryParams)
	_req.SetQueryParamsFromValues(_queryParams)

	_resp, _err = _req.Get("/pet/findByTags")
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

// PostPetPetIdUploadImage uploads an image
func (receiver *PetClient) PostPetPetIdUploadImage(ctx context.Context, _headers map[string]string,
	queryParams *struct {
		AdditionalMetadata *string `json:"additionalMetadata,omitempty" url:"additionalMetadata"`
	},
	// ID of pet to update
	// required
	petId int64,
	file *v3.FileModel) (ret ApiResponse, _resp *resty.Response, err error) {
	var _err error

	_req := receiver.client.R()
	_req.SetContext(ctx)
	if len(_headers) > 0 {
		_req.SetHeaders(_headers)
	}
	_queryParams, _ := _querystring.Values(queryParams)
	_req.SetQueryParamsFromValues(_queryParams)
	_req.SetPathParam("petId", fmt.Sprintf("%v", petId))
	if file != nil {
		_req.SetFileReader("file", file.Filename, file.Reader)
	}

	_resp, _err = _req.Post("/pet/{petId}/uploadImage")
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

// PostPet Add a new pet to the store
// Add a new pet to the store
func (receiver *PetClient) PostPet(ctx context.Context, _headers map[string]string,
	bodyJSON Pet) (ret Pet, _resp *resty.Response, err error) {
	var _err error

	_req := receiver.client.R()
	_req.SetContext(ctx)
	if len(_headers) > 0 {
		_req.SetHeaders(_headers)
	}
	_req.SetBody(bodyJSON)

	_resp, _err = _req.Post("/pet")
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

// PutPet Update an existing pet
// Update an existing pet by Id
func (receiver *PetClient) PutPet(ctx context.Context, _headers map[string]string,
	bodyJSON *Pet) (ret Pet, _resp *resty.Response, err error) {
	var _err error

	_req := receiver.client.R()
	_req.SetContext(ctx)
	if len(_headers) > 0 {
		_req.SetHeaders(_headers)
	}
	_req.SetBody(bodyJSON)

	_resp, _err = _req.Put("/pet")
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

// GetPetFindByStatus Finds Pets by status
// Multiple status values can be provided with comma separated strings
func (receiver *PetClient) GetPetFindByStatus(ctx context.Context, _headers map[string]string,
	queryParams *struct {
		Status *string `json:"status,omitempty" url:"status"`
	}) (ret []Pet, _resp *resty.Response, err error) {
	var _err error

	_req := receiver.client.R()
	_req.SetContext(ctx)
	if len(_headers) > 0 {
		_req.SetHeaders(_headers)
	}
	_queryParams, _ := _querystring.Values(queryParams)
	_req.SetQueryParamsFromValues(_queryParams)

	_resp, _err = _req.Get("/pet/findByStatus")
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

func NewPet(opts ...ddhttp.DdClientOption) *PetClient {
	defaultProvider := ddhttp.NewServiceProvider("PET")
	defaultClient := ddhttp.NewClient()

	svcClient := &PetClient{
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
