package cmd

import (
	"github.com/unionj-cloud/go-doudou/svc"
	"io/ioutil"
	"os"
	"testing"
)

func TestHttpCmd(t *testing.T) {
	dir := testDir + "httpcmd"
	receiver := svc.Svc{
		Dir: dir,
	}
	receiver.Init()
	defer os.RemoveAll(dir)
	err := os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	// go-doudou svc http --handler -c go -o
	_, _, err = ExecuteCommandC(rootCmd, []string{"svc", "http", "--handler", "-c", "go", "-o"}...)
	if err != nil {
		t.Fatal(err)
	}
	expect := `package httpsrv

import (
	"context"
	"encoding/json"
	"net/http"
	service "testfileshttpcmd"
	"testfileshttpcmd/vo"
)

type TestfileshttpcmdHandlerImpl struct {
	testfileshttpcmd service.Testfileshttpcmd
}

func (receiver *TestfileshttpcmdHandlerImpl) PageUsers(_writer http.ResponseWriter, _req *http.Request) {
	var (
		ctx   context.Context
		query vo.PageQuery
		code  int
		data  vo.PageRet
		msg   error
	)
	ctx = _req.Context()
	if err := json.NewDecoder(_req.Body).Decode(&query); err != nil {
		http.Error(_writer, err.Error(), http.StatusBadRequest)
		return
	}
	defer _req.Body.Close()
	code, data, msg = receiver.testfileshttpcmd.PageUsers(
		ctx,
		query,
	)
	if msg != nil {
		if msg == context.Canceled {
			http.Error(_writer, msg.Error(), http.StatusBadRequest)
		} else {
			http.Error(_writer, msg.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewEncoder(_writer).Encode(struct {
		Code int        ` + "`" + `json:"code,omitempty"` + "`" + `
		Data vo.PageRet ` + "`" + `json:"data,omitempty"` + "`" + `
	}{
		Code: code,
		Data: data,
	}); err != nil {
		http.Error(_writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NewTestfileshttpcmdHandler(testfileshttpcmd service.Testfileshttpcmd) TestfileshttpcmdHandler {
	return &TestfileshttpcmdHandlerImpl{
		testfileshttpcmd,
	}
}
`
	handlerimplfile := dir + "/transport/httpsrv/handlerimpl.go"
	f, err := os.Open(handlerimplfile)
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != expect {
		t.Errorf("want %s, go %s\n", expect, string(content))
	}

	expect = `package client

import (
	"context"
	"encoding/json"
	"net/url"
	service "testfileshttpcmd"
	"testfileshttpcmd/vo"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/stringutils"
	ddhttp "github.com/unionj-cloud/go-doudou/svc/http"
)

type TestfileshttpcmdClient struct {
	provider ddhttp.IServiceProvider
	client   *resty.Client
}

func (receiver *TestfileshttpcmdClient) PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, msg error) {
	var (
		_server string
		_err    error
	)
	if _server, _err = receiver.provider.SelectServer(); _err != nil {
		msg = errors.Wrap(_err, "")
		return
	}
	_urlValues := url.Values{}
	_req := receiver.client.R()
	_req.SetContext(ctx)
	_req.SetBody(query)
	if _req.Body != nil {
		_req.SetQueryParamsFromValues(_urlValues)
	} else {
		_req.SetFormDataFromValues(_urlValues)
	}
	_resp, _err := _req.Post(_server + "/testfileshttpcmd/pageusers")
	if _err != nil {
		msg = errors.Wrap(_err, "")
		return
	}
	if _resp.IsError() {
		msg = errors.New(_resp.String())
		return
	}
	var _result struct {
		Code int        ` + "`" + `json:"code"` + "`" + `
		Data vo.PageRet ` + "`" + `json:"data"` + "`" + `
		Msg  string     ` + "`" + `json:"msg"` + "`" + `
	}
	if _err = json.Unmarshal(_resp.Body(), &_result); _err != nil {
		msg = errors.Wrap(_err, "")
		return
	}
	if stringutils.IsNotEmpty(_result.Msg) {
		msg = errors.New(_result.Msg)
		return
	}
	return _result.Code, _result.Data, nil
}

type TestfileshttpcmdClientOption func(*TestfileshttpcmdClient)

func WithProvider(provider ddhttp.IServiceProvider) TestfileshttpcmdClientOption {
	return func(a *TestfileshttpcmdClient) {
		a.provider = provider
	}
}

func WithClient(client *resty.Client) TestfileshttpcmdClientOption {
	return func(a *TestfileshttpcmdClient) {
		a.client = client
	}
}

func NewTestfileshttpcmd(opts ...TestfileshttpcmdClientOption) service.Testfileshttpcmd {
	defaultProvider := ddhttp.NewServiceProvider("Testfileshttpcmd")
	defaultClient := ddhttp.NewClient()

	svcClient := &TestfileshttpcmdClient{
		provider: defaultProvider,
		client:   defaultClient,
	}

	for _, opt := range opts {
		opt(svcClient)
	}

	return svcClient
}
`

	clientfile := dir + "/client/client.go"
	f, err = os.Open(clientfile)
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != expect {
		t.Errorf("want %s, go %s\n", expect, string(content))
	}
}
