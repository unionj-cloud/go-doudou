package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	"github.com/unionj-cloud/go-doudou/framework/registry"
)

type StoreClient struct {
	provider registry.IServiceProvider
	client   *resty.Client
	rootPath string
}

func (receiver *StoreClient) SetRootPath(rootPath string) {
	receiver.rootPath = rootPath
}

func (receiver *StoreClient) SetProvider(provider registry.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *StoreClient) SetClient(client *resty.Client) {
	receiver.client = client
}

// GetStoreOrderOrderId Find purchase order by ID
// For valid response try integer IDs with value <= 5 or > 10. Other values will generated exceptions
func (receiver *StoreClient) GetStoreOrderOrderId(ctx context.Context, _headers map[string]string,
	// ID of order that needs to be fetched
	// required
	orderId int64) (ret Order, _resp *resty.Response, err error) {
	var _err error

	_req := receiver.client.R()
	_req.SetContext(ctx)
	if len(_headers) > 0 {
		_req.SetHeaders(_headers)
	}
	_req.SetPathParam("orderId", fmt.Sprintf("%v", orderId))

	_resp, _err = _req.Get("/store/order/{orderId}")
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

// GetStoreInventory Returns pet inventories by status
// Returns a map of status codes to quantities
func (receiver *StoreClient) GetStoreInventory(ctx context.Context, _headers map[string]string) (ret struct {
}, _resp *resty.Response, err error) {
	var _err error

	_req := receiver.client.R()
	_req.SetContext(ctx)
	if len(_headers) > 0 {
		_req.SetHeaders(_headers)
	}

	_resp, _err = _req.Get("/store/inventory")
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

// PostStoreOrder Place an order for a pet
// Place a new order in the store
func (receiver *StoreClient) PostStoreOrder(ctx context.Context, _headers map[string]string,
	bodyJSON *Order) (ret Order, _resp *resty.Response, err error) {
	var _err error

	_req := receiver.client.R()
	_req.SetContext(ctx)
	if len(_headers) > 0 {
		_req.SetHeaders(_headers)
	}
	_req.SetBody(bodyJSON)

	_resp, _err = _req.Post("/store/order")
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

func NewStore(opts ...ddhttp.DdClientOption) *StoreClient {
	defaultProvider := ddhttp.NewServiceProvider("STORE")
	defaultClient := ddhttp.NewClient()

	svcClient := &StoreClient{
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
