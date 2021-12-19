package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testdata/vo"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/fileutils"
	v3 "github.com/unionj-cloud/go-doudou/openapi/v3"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/svc/config"
	ddhttp "github.com/unionj-cloud/go-doudou/svc/http"
	"github.com/unionj-cloud/go-doudou/svc/registry"
)

type UsersvcClient struct {
	provider registry.IServiceProvider
	client   *resty.Client
}

func (receiver *UsersvcClient) SetProvider(provider registry.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *UsersvcClient) SetClient(client *resty.Client) {
	receiver.client = client
}
func (receiver *UsersvcClient) PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, msg error) {
	var _err error
	_urlValues := url.Values{}
	_req := receiver.client.R()
	_req.SetContext(ctx)
	_req.SetBody(query)
	_path := "/usersvc/pageusers"
	if _req.Body != nil {
		_req.SetQueryParamsFromValues(_urlValues)
	} else {
		_req.SetFormDataFromValues(_urlValues)
	}
	_resp, _err := _req.Post(_path)
	if _err != nil {
		msg = errors.Wrap(_err, "")
		return
	}
	if _resp.IsError() {
		msg = errors.New(_resp.String())
		return
	}
	var _result struct {
		Code int        `json:"code"`
		Data vo.PageRet `json:"data"`
		Msg  string     `json:"msg"`
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
func (receiver *UsersvcClient) GetUser(ctx context.Context, userId string, photo string) (code int, data string, msg error) {
	var _err error
	_urlValues := url.Values{}
	_req := receiver.client.R()
	_req.SetContext(ctx)
	_urlValues.Set("userId", fmt.Sprintf("%v", userId))
	_urlValues.Set("photo", fmt.Sprintf("%v", photo))
	_path := "/usersvc/user"
	_resp, _err := _req.SetQueryParamsFromValues(_urlValues).
		Get(_path)
	if _err != nil {
		msg = errors.Wrap(_err, "")
		return
	}
	if _resp.IsError() {
		msg = errors.New(_resp.String())
		return
	}
	var _result struct {
		Code int    `json:"code"`
		Data string `json:"data"`
		Msg  string `json:"msg"`
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
func (receiver *UsersvcClient) SignUp(ctx context.Context, username string, password int, actived bool, score []int) (code int, data string, msg error) {
	var _err error
	_urlValues := url.Values{}
	_req := receiver.client.R()
	_req.SetContext(ctx)
	_urlValues.Set("username", fmt.Sprintf("%v", username))
	_urlValues.Set("password", fmt.Sprintf("%v", password))
	_urlValues.Set("actived", fmt.Sprintf("%v", actived))
	for _, _item := range score {
		_urlValues.Add("score", fmt.Sprintf("%v", _item))
	}
	_path := "/usersvc/signup"
	if _req.Body != nil {
		_req.SetQueryParamsFromValues(_urlValues)
	} else {
		_req.SetFormDataFromValues(_urlValues)
	}
	_resp, _err := _req.Post(_path)
	if _err != nil {
		msg = errors.Wrap(_err, "")
		return
	}
	if _resp.IsError() {
		msg = errors.New(_resp.String())
		return
	}
	var _result struct {
		Code int    `json:"code"`
		Data string `json:"data"`
		Msg  string `json:"msg"`
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
func (receiver *UsersvcClient) UploadAvatar(pc context.Context, pf []*v3.FileModel, ps string, pf2 *v3.FileModel, pf3 *multipart.FileHeader, pf4 []*multipart.FileHeader) (ri int, rs string, re error) {
	var _err error
	_urlValues := url.Values{}
	_req := receiver.client.R()
	_req.SetContext(pc)
	for _, _f := range pf {
		_req.SetFileReader("pf", _f.Filename, _f.Reader)
	}
	_urlValues.Set("ps", fmt.Sprintf("%v", ps))
	_req.SetFileReader("pf2", pf2.Filename, pf2.Reader)
	if _f, _err := pf3.Open(); _err != nil {
		re = errors.Wrap(_err, "")
		return
	} else {
		_req.SetFileReader("pf3", pf3.Filename, _f)
	}
	for _, _fh := range pf4 {
		_f, _err := _fh.Open()
		if _err != nil {
			re = errors.Wrap(_err, "")
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
	_resp, _err := _req.Post(_path)
	if _err != nil {
		re = errors.Wrap(_err, "")
		return
	}
	if _resp.IsError() {
		re = errors.New(_resp.String())
		return
	}
	var _result struct {
		Ri int    `json:"ri"`
		Rs string `json:"rs"`
		Re string `json:"re"`
	}
	if _err = json.Unmarshal(_resp.Body(), &_result); _err != nil {
		re = errors.Wrap(_err, "")
		return
	}
	if stringutils.IsNotEmpty(_result.Re) {
		re = errors.New(_result.Re)
		return
	}
	return _result.Ri, _result.Rs, nil
}
func (receiver *UsersvcClient) DownloadAvatar(ctx context.Context, userId string) (rf *os.File, re error) {
	var _err error
	_urlValues := url.Values{}
	_req := receiver.client.R()
	_req.SetContext(ctx)
	_urlValues.Set("userId", fmt.Sprintf("%v", userId))
	_req.SetDoNotParseResponse(true)
	_path := "/usersvc/downloadavatar"
	if _req.Body != nil {
		_req.SetQueryParamsFromValues(_urlValues)
	} else {
		_req.SetFormDataFromValues(_urlValues)
	}
	_resp, _err := _req.Post(_path)
	if _err != nil {
		re = errors.Wrap(_err, "")
		return
	}
	if _resp.IsError() {
		re = errors.New(_resp.String())
		return
	}
	_disp := _resp.Header().Get("Content-Disposition")
	_file := strings.TrimPrefix(_disp, "attachment; filename=")
	_output := config.GddOutput.Load()
	if stringutils.IsNotEmpty(_output) {
		_file = _output + string(filepath.Separator) + _file
	}
	_file = filepath.Clean(_file)
	if _err = fileutils.CreateDirectory(filepath.Dir(_file)); _err != nil {
		re = errors.Wrap(_err, "")
		return
	}
	_outFile, _err := os.Create(_file)
	if _err != nil {
		re = errors.Wrap(_err, "")
		return
	}
	defer _outFile.Close()
	defer _resp.RawBody().Close()
	_, _err = io.Copy(_outFile, _resp.RawBody())
	if _err != nil {
		re = errors.Wrap(_err, "")
		return
	}
	rf = _outFile
	return
}

func NewUsersvc(opts ...ddhttp.DdClientOption) *UsersvcClient {
	defaultProvider := ddhttp.NewServiceProvider("USERSVC")
	defaultClient := ddhttp.NewClient()

	svcClient := &UsersvcClient{
		provider: defaultProvider,
		client:   defaultClient,
	}

	for _, opt := range opts {
		opt(svcClient)
	}

	svcClient.client.OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
		client.SetHostURL(svcClient.provider.SelectServer())
		return nil
	})

	return svcClient
}
