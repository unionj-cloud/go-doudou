package test

import (
	"context"
	"encoding/json"

	"github.com/go-resty/resty/v2"
	_querystring "github.com/google/go-querystring/query"
	"github.com/pkg/errors"
	ddhttp "github.com/unionj-cloud/go-doudou/svc/http"
)

type UserClient struct {
	provider ddhttp.IServiceProvider
	client   *resty.Client
}

func (receiver *UserClient) SetProvider(provider ddhttp.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *UserClient) SetClient(client *resty.Client) {
	receiver.client = client
}

// Logs user into the system
func (receiver *UserClient) GetUserLogin(ctx context.Context,
	queryParams struct {
		Username string `json:"username,omitempty" url:"username"`
		Password string `json:"password,omitempty" url:"password"`
	}) (ret string, err error) {
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

	_resp, _err := _req.Get(_server + "/user/login")
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

// Get user by user name
func (receiver *UserClient) GetUserUsername(ctx context.Context) (ret User, err error) {
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

	_resp, _err := _req.Get(_server + "/user/{username}")
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

// Creates list of users with given input array
// Creates list of users with given input array
func (receiver *UserClient) PostUserCreateWithList(ctx context.Context,
	bodyJSON []User) (ret User, err error) {
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

	_resp, _err := _req.Post(_server + "/user/createWithList")
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

func NewUser(opts ...ddhttp.DdClientOption) *UserClient {
	defaultProvider := ddhttp.NewServiceProvider("USER")
	defaultClient := ddhttp.NewClient()

	svcClient := &UserClient{
		provider: defaultProvider,
		client:   defaultClient,
	}

	for _, opt := range opts {
		opt(svcClient)
	}

	return svcClient
}
