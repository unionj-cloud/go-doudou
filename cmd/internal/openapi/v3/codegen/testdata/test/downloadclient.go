package test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-resty/resty/v2"
	_querystring "github.com/google/go-querystring/query"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry"
	"github.com/unionj-cloud/go-doudou/v2/framework/restclient"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/fileutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
)

type DownloadClient struct {
	provider registry.IServiceProvider
	client   *resty.Client
	rootPath string
}

func (receiver *DownloadClient) SetRootPath(rootPath string) {
	receiver.rootPath = rootPath
}

func (receiver *DownloadClient) SetProvider(provider registry.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *DownloadClient) SetClient(client *resty.Client) {
	receiver.client = client
}

// GetDownloadAvatar GetDownloadAvatar demonstrate how to define download file api
// there must be *os.File parameter among output parameters
func (receiver *DownloadClient) GetDownloadAvatar(ctx context.Context, _headers map[string]string,
	queryParams struct {
		// required
		UserId string `json:"userId,omitempty" url:"userId"`
	}) (_downloadFile *os.File, _resp *resty.Response, err error) {
	var _err error

	_req := receiver.client.R()
	_req.SetContext(ctx)
	if len(_headers) > 0 {
		_req.SetHeaders(_headers)
	}
	_queryParams, _ := _querystring.Values(queryParams)
	_req.SetQueryParamsFromValues(_queryParams)
	_req.SetDoNotParseResponse(true)

	_resp, _err = _req.Get("/download/avatar")
	if _err != nil {
		err = errors.Wrap(_err, "")
		return
	}
	if _resp.IsError() {
		err = errors.New(_resp.String())
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
		err = errors.Wrap(_err, "")
		return
	}
	_outFile, _err := os.Create(_file)
	if _err != nil {
		err = errors.Wrap(_err, "")
		return
	}
	defer _outFile.Close()
	defer _resp.RawBody().Close()
	_, _err = io.Copy(_outFile, _resp.RawBody())
	if _err != nil {
		err = errors.Wrap(_err, "")
		return
	}
	_downloadFile = _outFile
	return
}

func NewDownload(opts ...restclient.RestClientOption) *DownloadClient {
	defaultProvider := restclient.NewServiceProvider("DOWNLOAD")
	defaultClient := restclient.NewClient()

	svcClient := &DownloadClient{
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
