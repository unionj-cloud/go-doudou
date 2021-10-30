package test

import (
	"context"
	"encoding/json"

	"github.com/go-resty/resty/v2"
	_querystring "github.com/google/go-querystring/query"
	"github.com/pkg/errors"
	ddhttp "github.com/unionj-cloud/go-doudou/svc/http"
)

type CustomerClient struct {
	provider ddhttp.IServiceProvider
	client   *resty.Client
}

func (receiver *CustomerClient) SetProvider(provider ddhttp.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *CustomerClient) SetClient(client *resty.Client) {
	receiver.client = client
}
func (receiver *CustomerClient) GetCustomerValidateToken(ctx context.Context,
	queryParams struct {
		// required
		Token string `json:"token" url:"token"`
	}) (ret bool, err error) {
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

	_resp, _err := _req.Get(_server + "/customer/validateToken")
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

	return svcClient
}
