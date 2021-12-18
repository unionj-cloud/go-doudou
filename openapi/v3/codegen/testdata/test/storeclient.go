package test

import (
	"context"
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	ddhttp "github.com/unionj-cloud/go-doudou/svc/http"
)

type StoreClient struct {
	provider ddhttp.IServiceProvider
	client   *resty.Client
}

func (receiver *StoreClient) SetProvider(provider ddhttp.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *StoreClient) SetClient(client *resty.Client) {
	receiver.client = client
}

// GetStoreOrderOrderId Find purchase order by ID
// For valid response try integer IDs with value <= 5 or > 10. Other values will generated exceptions
func (receiver *StoreClient) GetStoreOrderOrderId(ctx context.Context) (ret Order, err error) {
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

	_resp, _err := _req.Get(_server + "/store/order/{orderId}")
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
func (receiver *StoreClient) PostStoreOrder(ctx context.Context,
	bodyJSON Order) (ret Order, err error) {
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

	_resp, _err := _req.Post(_server + "/store/order")
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
func (receiver *StoreClient) GetStoreInventory(ctx context.Context) (ret struct {
}, err error) {
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

	_resp, _err := _req.Get(_server + "/store/inventory")
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

	return svcClient
}
