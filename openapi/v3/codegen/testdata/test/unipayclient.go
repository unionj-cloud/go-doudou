package test

import (
	"context"

	"github.com/go-resty/resty/v2"
	_querystring "github.com/google/go-querystring/query"
	"github.com/pkg/errors"
	ddhttp "github.com/unionj-cloud/go-doudou/svc/http"
)

type UnipayClient struct {
	provider ddhttp.IServiceProvider
	client   *resty.Client
}

func (receiver *UnipayClient) SetProvider(provider ddhttp.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *UnipayClient) SetClient(client *resty.Client) {
	receiver.client = client
}
func (receiver *UnipayClient) GetUnipayStartUnionPay(ctx context.Context,
	queryParams struct {
		// required
		FrontUrl string `json:"frontUrl,omitempty" url:"frontUrl"`
		// required
		TxnAmt string `json:"txnAmt,omitempty" url:"txnAmt"`
		// required
		Token string `json:"token,omitempty" url:"token"`
		// required
		CompanyId string `json:"companyId,omitempty" url:"companyId"`
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

	_resp, _err := _req.Get(_server + "/unipay/startUnionPay")
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

	return svcClient
}
