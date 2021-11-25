package test

import (
	"context"
	"encoding/json"

	"github.com/go-resty/resty/v2"
	_querystring "github.com/google/go-querystring/query"
	"github.com/pkg/errors"
	v3 "github.com/unionj-cloud/go-doudou/openapi/v3"
	ddhttp "github.com/unionj-cloud/go-doudou/svc/http"
)

type PetClient struct {
	provider ddhttp.IServiceProvider
	client   *resty.Client
}

func (receiver *PetClient) SetProvider(provider ddhttp.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *PetClient) SetClient(client *resty.Client) {
	receiver.client = client
}

// PostPet Add a new pet to the store
// Add a new pet to the store
func (receiver *PetClient) PostPet(ctx context.Context,
	bodyJSON Pet) (ret Pet, err error) {
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

	_resp, _err := _req.Post(_server + "/pet")
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
func (receiver *PetClient) PutPet(ctx context.Context,
	bodyJSON Pet) (ret Pet, err error) {
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

	_resp, _err := _req.Put(_server + "/pet")
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
func (receiver *PetClient) GetPetFindByTags(ctx context.Context,
	queryParams struct {
		Tags []string `json:"tags,omitempty" url:"tags"`
	}) (ret []Pet, err error) {
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

	_resp, _err := _req.Get(_server + "/pet/findByTags")
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

// GetPetPetId Find pet by ID
// Returns a single pet
func (receiver *PetClient) GetPetPetId(ctx context.Context) (ret Pet, err error) {
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

	_resp, _err := _req.Get(_server + "/pet/{petId}")
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
func (receiver *PetClient) PostPetPetIdUploadImage(ctx context.Context,
	queryParams struct {
		AdditionalMetadata string `json:"additionalMetadata,omitempty" url:"additionalMetadata"`
	},
	file *v3.FileModel) (ret ApiResponse, err error) {
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
	_req.SetFileReader("file", file.Filename, file.Reader)

	_resp, _err := _req.Post(_server + "/pet/{petId}/uploadImage")
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
func (receiver *PetClient) GetPetFindByStatus(ctx context.Context,
	queryParams struct {
		Status string `json:"status,omitempty" url:"status"`
	}) (ret []Pet, err error) {
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

	_resp, _err := _req.Get(_server + "/pet/findByStatus")
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

	return svcClient
}
