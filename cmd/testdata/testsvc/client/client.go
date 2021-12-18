package client

import (
	"context"
	"encoding/json"
	"net/url"
	"testsvc/vo"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/stringutils"
	ddhttp "github.com/unionj-cloud/go-doudou/svc/http"
	"github.com/unionj-cloud/go-doudou/svc/registry"
)

type TestsvcClient struct {
	provider registry.IServiceProvider
	client   *resty.Client
}

func (receiver *TestsvcClient) SetProvider(provider registry.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *TestsvcClient) SetClient(client *resty.Client) {
	receiver.client = client
}
func (receiver *TestsvcClient) PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, err error) {
	var (
		_server string
		_err    error
	)
	if _server, _err = receiver.provider.SelectServer(); _err != nil {
		err = errors.Wrap(_err, "")
		return
	}
	_urlValues := url.Values{}
	_req := receiver.client.R()
	_req.SetContext(ctx)
	_req.SetBody(query)
	_path := "/page/users"
	if _req.Body != nil {
		_req.SetQueryParamsFromValues(_urlValues)
	} else {
		_req.SetFormDataFromValues(_urlValues)
	}
	_resp, _err := _req.Post(_server + _path)
	if _err != nil {
		err = errors.Wrap(_err, "")
		return
	}
	if _resp.IsError() {
		err = errors.New(_resp.String())
		return
	}
	var _result struct {
		Code int        `json:"code"`
		Data vo.PageRet `json:"data"`
		Err  string     `json:"err"`
	}
	if _err = json.Unmarshal(_resp.Body(), &_result); _err != nil {
		err = errors.Wrap(_err, "")
		return
	}
	if stringutils.IsNotEmpty(_result.Err) {
		err = errors.New(_result.Err)
		return
	}
	return _result.Code, _result.Data, nil
}

func NewTestsvc(opts ...ddhttp.DdClientOption) *TestsvcClient {
	defaultProvider := ddhttp.NewServiceProvider("TESTSVC")
	defaultClient := ddhttp.NewClient()

	svcClient := &TestsvcClient{
		provider: defaultProvider,
		client:   defaultClient,
	}

	for _, opt := range opts {
		opt(svcClient)
	}

	return svcClient
}
