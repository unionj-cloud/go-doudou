package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestHttpCmd(t *testing.T) {
	dir := filepath.Join(testDir, "testsvc")
	_ = os.Chdir(dir)
	// go-doudou svc http --handler -c go -o
	_, _, err := ExecuteCommandC(rootCmd, []string{"svc", "http", "--handler", "-c", "go", "-o"}...)
	if err != nil {
		t.Fatal(err)
	}
	expect := `package httpsrv

import (
	"context"
	"encoding/json"
	"net/http"
	service "testsvc"
	"testsvc/vo"

	"github.com/pkg/errors"
)

type TestsvcHandlerImpl struct {
	testsvc service.Testsvc
}

func (receiver *TestsvcHandlerImpl) PageUsers(_writer http.ResponseWriter, _req *http.Request) {
	var (
		ctx   context.Context
		query vo.PageQuery
		code  int
		data  vo.PageRet
		err   error
	)
	ctx = _req.Context()
	if err := json.NewDecoder(_req.Body).Decode(&query); err != nil {
		http.Error(_writer, err.Error(), http.StatusBadRequest)
		return
	}
	defer _req.Body.Close()
	code, data, err = receiver.testsvc.PageUsers(
		ctx,
		query,
	)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			http.Error(_writer, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(_writer, err.Error(), http.StatusInternalServerError)
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

func NewTestsvcHandler(testsvc service.Testsvc) TestsvcHandler {
	return &TestsvcHandlerImpl{
		testsvc,
	}
}
`
	handlerimplfile := filepath.Join(dir, "transport", "httpsrv", "handlerimpl.go")
	f, err := os.Open(handlerimplfile)
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != expect {
		t.Errorf("want %s, got %s\n", expect, string(content))
	}

	expect = `package client

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
		Code int        ` + "`" + `json:"code"` + "`" + `
		Data vo.PageRet ` + "`" + `json:"data"` + "`" + `
		Err  string     ` + "`" + `json:"err"` + "`" + `
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
`

	clientfile := filepath.Join(dir, "client", "client.go")
	f, err = os.Open(clientfile)
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != expect {
		t.Errorf("want %s, got %s\n", expect, string(content))
	}
}
