package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testdata/vo"

	"github.com/go-resty/resty/v2"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	"github.com/unionj-cloud/go-doudou/framework/registry"
	"github.com/unionj-cloud/go-doudou/toolkit/fileutils"
	v3 "github.com/unionj-cloud/go-doudou/toolkit/openapi/v3"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
)

type UsersvcClient struct {
	provider registry.IServiceProvider
	client   *resty.Client
	rootPath string
}

func (receiver *UsersvcClient) SetRootPath(rootPath string) {
	receiver.rootPath = rootPath
}

func (receiver *UsersvcClient) SetProvider(provider registry.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *UsersvcClient) SetClient(client *resty.Client) {
	receiver.client = client
}
func (receiver *UsersvcClient) PageUsers(ctx context.Context, _headers map[string]string, query vo.PageQuery) (_resp *resty.Response, code int, data vo.PageRet, msg error) {
	var _err error
	_urlValues := url.Values{}
	_req := receiver.client.R()
	if len(_headers) > 0 {
		_req.SetHeaders(_headers)
	}
	_req.SetContext(ctx)
	_req.SetBody(query)
	_path := "/usersvc/pageusers"
	if _req.Body != nil {
		_req.SetQueryParamsFromValues(_urlValues)
	} else {
		_req.SetFormDataFromValues(_urlValues)
	}
	_resp, _err = _req.Post(_path)
	if _err != nil {
		msg = errors.Wrap(_err, "error")
		return
	}
	if _resp.IsError() {
		msg = errors.New(_resp.String())
		return
	}
	var _result struct {
		Code int        `json:"code"`
		Data vo.PageRet `json:"data"`
	}
	if _err = json.Unmarshal(_resp.Body(), &_result); _err != nil {
		msg = errors.Wrap(_err, "error")
		return
	}
	return _resp, _result.Code, _result.Data, nil
}
func (receiver *UsersvcClient) GetUser(ctx context.Context, _headers map[string]string, userId string, photo string) (_resp *resty.Response, code int, data string, msg error) {
	var _err error
	_urlValues := url.Values{}
	_req := receiver.client.R()
	if len(_headers) > 0 {
		_req.SetHeaders(_headers)
	}
	_req.SetContext(ctx)
	_urlValues.Set("userId", fmt.Sprintf("%v", userId))
	_urlValues.Set("photo", fmt.Sprintf("%v", photo))
	_path := "/usersvc/user"
	_resp, _err = _req.SetQueryParamsFromValues(_urlValues).
		Get(_path)
	if _err != nil {
		msg = errors.Wrap(_err, "error")
		return
	}
	if _resp.IsError() {
		msg = errors.New(_resp.String())
		return
	}
	var _result struct {
		Code int    `json:"code"`
		Data string `json:"data"`
	}
	if _err = json.Unmarshal(_resp.Body(), &_result); _err != nil {
		msg = errors.Wrap(_err, "error")
		return
	}
	return _resp, _result.Code, _result.Data, nil
}
func (receiver *UsersvcClient) SignUp(ctx context.Context, _headers map[string]string, username string, password int, actived bool, score []int) (_resp *resty.Response, code int, data string, msg error) {
	var _err error
	_urlValues := url.Values{}
	_req := receiver.client.R()
	if len(_headers) > 0 {
		_req.SetHeaders(_headers)
	}
	_req.SetContext(ctx)
	_urlValues.Set("username", fmt.Sprintf("%v", username))
	_urlValues.Set("password", fmt.Sprintf("%v", password))
	_urlValues.Set("actived", fmt.Sprintf("%v", actived))
	if len(score) == 0 {
		msg = errors.New("size of parameter score should be greater than zero")
		return
	}
	for _, _item := range score {
		_urlValues.Add("score", fmt.Sprintf("%v", _item))
	}
	_path := "/usersvc/signup"
	if _req.Body != nil {
		_req.SetQueryParamsFromValues(_urlValues)
	} else {
		_req.SetFormDataFromValues(_urlValues)
	}
	_resp, _err = _req.Post(_path)
	if _err != nil {
		msg = errors.Wrap(_err, "error")
		return
	}
	if _resp.IsError() {
		msg = errors.New(_resp.String())
		return
	}
	var _result struct {
		Code int    `json:"code"`
		Data string `json:"data"`
	}
	if _err = json.Unmarshal(_resp.Body(), &_result); _err != nil {
		msg = errors.Wrap(_err, "error")
		return
	}
	return _resp, _result.Code, _result.Data, nil
}
func (receiver *UsersvcClient) UploadAvatar(ctx context.Context, _headers map[string]string, pf []v3.FileModel, ps string, pf2 v3.FileModel, pf3 *multipart.FileHeader, pf4 []*multipart.FileHeader) (_resp *resty.Response, ri int, ri2 interface{}, re error) {
	var _err error
	_urlValues := url.Values{}
	_req := receiver.client.R()
	if len(_headers) > 0 {
		_req.SetHeaders(_headers)
	}
	_req.SetContext(ctx)
	if len(pf) == 0 {
		re = errors.New("at least one file should be uploaded for parameter pf")
		return
	}
	for _, _f := range pf {
		_req.SetFileReader("pf", _f.Filename, _f.Reader)
	}
	_urlValues.Set("ps", fmt.Sprintf("%v", ps))
	_req.SetFileReader("pf2", pf2.Filename, pf2.Reader)
	if _f, _err := pf3.Open(); _err != nil {
		re = errors.Wrap(_err, "error")
		return
	} else {
		_req.SetFileReader("pf3", pf3.Filename, _f)
	}
	for _, _fh := range pf4 {
		_f, _err := _fh.Open()
		if _err != nil {
			re = errors.Wrap(_err, "error")
			return
		}
		_req.SetFileReader("pf4", _fh.Filename, _f)
	}
	_path := "/usersvc/uploadavatar"
	if _req.Body != nil {
		_req.SetQueryParamsFromValues(_urlValues)
	} else {
		_req.SetFormDataFromValues(_urlValues)
	}
	_resp, _err = _req.Post(_path)
	if _err != nil {
		re = errors.Wrap(_err, "error")
		return
	}
	if _resp.IsError() {
		re = errors.New(_resp.String())
		return
	}
	var _result struct {
		Ri  int         `json:"ri"`
		Ri2 interface{} `json:"ri2"`
	}
	if _err = json.Unmarshal(_resp.Body(), &_result); _err != nil {
		re = errors.Wrap(_err, "error")
		return
	}
	return _resp, _result.Ri, _result.Ri2, nil
}
func (receiver *UsersvcClient) DownloadAvatar(ctx context.Context, _headers map[string]string, userId interface{}, data []byte, userAttrs ...string) (_resp *resty.Response, rf *os.File, re error) {
	var _err error
	_urlValues := url.Values{}
	_req := receiver.client.R()
	if len(_headers) > 0 {
		_req.SetHeaders(_headers)
	}
	_req.SetContext(ctx)
	_req.SetBody(userId)
	if len(data) == 0 {
		re = errors.New("size of parameter data should be greater than zero")
		return
	}
	for _, _item := range data {
		_urlValues.Add("data", fmt.Sprintf("%v", _item))
	}
	if userAttrs != nil {
		for _, _item := range userAttrs {
			_urlValues.Add("userAttrs", fmt.Sprintf("%v", _item))
		}
	}
	_req.SetDoNotParseResponse(true)
	_path := "/usersvc/downloadavatar"
	if _req.Body != nil {
		_req.SetQueryParamsFromValues(_urlValues)
	} else {
		_req.SetFormDataFromValues(_urlValues)
	}
	_resp, _err = _req.Post(_path)
	if _err != nil {
		re = errors.Wrap(_err, "error")
		return
	}
	if _resp.IsError() {
		re = errors.New(_resp.String())
		return
	}
	_disp := _resp.Header().Get("Content-Disposition")
	_file := strings.TrimPrefix(_disp, "attachment; filename=")
	_output := os.TempDir()
	if stringutils.IsNotEmpty(_output) {
		_file = _output + string(filepath.Separator) + _file
	}
	_file = filepath.Clean(_file)
	if _err = fileutils.CreateDirectory(filepath.Dir(_file)); _err != nil {
		re = errors.Wrap(_err, "error")
		return
	}
	_outFile, _err := os.Create(_file)
	if _err != nil {
		re = errors.Wrap(_err, "error")
		return
	}
	defer _outFile.Close()
	defer _resp.RawBody().Close()
	_, _err = io.Copy(_outFile, _resp.RawBody())
	if _err != nil {
		re = errors.Wrap(_err, "error")
		return
	}
	rf = _outFile
	return
}

func NewUsersvcClient(opts ...ddhttp.DdClientOption) *UsersvcClient {
	defaultProvider := ddhttp.NewServiceProvider("USERSVC")
	defaultClient := ddhttp.NewClient()

	svcClient := &UsersvcClient{
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
