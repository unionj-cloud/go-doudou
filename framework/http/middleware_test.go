package ddhttp_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/agcache/memory"
	apolloConfig "github.com/apolloconfig/agollo/v4/env/config"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/slok/goresilience"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/unionj-cloud/go-doudou/framework/configmgr"
	"github.com/unionj-cloud/go-doudou/framework/configmgr/mock"
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	httpMock "github.com/unionj-cloud/go-doudou/framework/http/mock"
	ddmodel "github.com/unionj-cloud/go-doudou/framework/http/model"
	"github.com/unionj-cloud/go-doudou/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/framework/registry"
	"github.com/unionj-cloud/go-doudou/toolkit/maputils"
	"github.com/wubin1989/nacos-sdk-go/clients/cache"
	"github.com/wubin1989/nacos-sdk-go/clients/config_client"
	"github.com/wubin1989/nacos-sdk-go/vo"
	"net/http"
	"os"
	"testing"
	"time"
)

type IMocksvcHandler interface {
	GetUser(w http.ResponseWriter, r *http.Request)
	SaveUser(w http.ResponseWriter, r *http.Request)
	SignUp(w http.ResponseWriter, r *http.Request)
	GetPanic(w http.ResponseWriter, r *http.Request)
}

type MocksvcHandler struct{}

func (m *MocksvcHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("go-doudou"))
}

func (m *MocksvcHandler) SaveUser(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Data string `json:"data"`
	}{
		Data: "OK",
	}
	resp, _ := json.Marshal(data)
	w.Write(resp)
}

func (m *MocksvcHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Data string `json:"data"`
	}{
		Data: "OK",
	}
	resp, _ := json.Marshal(data)
	w.Write(resp)
}

func (m *MocksvcHandler) GetPanic(w http.ResponseWriter, r *http.Request) {
	panic(context.Canceled)
}

func Routes(handler IMocksvcHandler) []ddmodel.Route {
	return []ddmodel.Route{
		{
			"GetUser",
			"GET",
			"/user",
			handler.GetUser,
		},
		{
			"SaveUser",
			"POST",
			"/save/user",
			handler.SaveUser,
		},
		{
			"SignUp",
			"POST",
			"/sign/up",
			handler.SignUp,
		},
		{
			"GetPanic",
			"GET",
			"/panic",
			handler.GetPanic,
		},
	}
}

func NewMocksvcHandler() IMocksvcHandler {
	return &MocksvcHandler{}
}

type UserVo struct {
	Username string
	Password string
}

type IMockClient interface {
	GetUser(ctx context.Context, _headers map[string]string) (_resp *resty.Response, data string, err error)
	SaveUser(ctx context.Context, _headers map[string]string, payload UserVo) (_resp *resty.Response, data string, err error)
	SignUp(ctx context.Context, _headers map[string]string, username, password string) (_resp *resty.Response, data string, err error)
	GetPanic(ctx context.Context, _headers map[string]string) (_resp *resty.Response, data string, err error)
}

type MockClient struct {
	provider registry.IServiceProvider
	client   *resty.Client
	rootPath string
}

func (receiver *MockClient) SetRootPath(rootPath string) {
	receiver.rootPath = rootPath
}

func (receiver *MockClient) SetProvider(provider registry.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *MockClient) SetClient(client *resty.Client) {
	receiver.client = client
}

func (receiver *MockClient) GetUser(ctx context.Context, _headers map[string]string) (_resp *resty.Response, data string, err error) {
	var _err error
	_req := receiver.client.R()
	_req.SetContext(ctx)
	_path := "/user"
	_resp, _err = _req.Get(_path)
	if _err != nil {
		err = errors.Wrap(_err, "error")
		return
	}
	if _resp.IsError() {
		err = errors.New(_resp.String())
		return
	}
	return _resp, string(_resp.Body()), nil
}

func (receiver *MockClient) GetPanic(ctx context.Context, _headers map[string]string) (_resp *resty.Response, data string, err error) {
	var _err error
	_req := receiver.client.R()
	_req.SetContext(ctx)
	_path := "/panic"
	_resp, _err = _req.Get(_path)
	if _err != nil {
		err = errors.Wrap(_err, "error")
		return
	}
	if _resp.IsError() {
		err = errors.New(_resp.String())
		return
	}
	return _resp, string(_resp.Body()), nil
}

func (receiver *MockClient) SignUp(ctx context.Context, _headers map[string]string, username, password string) (_resp *resty.Response, data string, err error) {
	var _err error
	_req := receiver.client.R()
	_req.SetContext(ctx)
	formData := make(map[string]string)
	formData["username"] = fmt.Sprintf("%v", username)
	formData["password"] = fmt.Sprintf("%v", password)
	_path := "/sign/up"
	_req.SetMultipartFormData(formData)
	_resp, _err = _req.Post(_path)
	if _err != nil {
		err = errors.Wrap(_err, "error")
		return
	}
	if _resp.IsError() {
		err = errors.New(_resp.String())
		return
	}
	var _result struct {
		Data string `json:"data"`
	}
	if _err = json.Unmarshal(_resp.Body(), &_result); _err != nil {
		err = errors.Wrap(_err, "error")
		return
	}
	return _resp, _result.Data, nil
}

func (receiver *MockClient) SaveUser(ctx context.Context, _headers map[string]string, payload UserVo) (_resp *resty.Response, data string, err error) {
	var _err error
	_req := receiver.client.R()
	_req.SetContext(ctx)
	_req.SetBody(payload)
	_path := "/save/user"
	_resp, _err = _req.Post(_path)
	if _err != nil {
		err = errors.Wrap(_err, "error")
		return
	}
	if _resp.IsError() {
		err = errors.New(_resp.String())
		return
	}
	var _result struct {
		Data string `json:"data"`
	}
	if _err = json.Unmarshal(_resp.Body(), &_result); _err != nil {
		err = errors.Wrap(_err, "error")
		return
	}
	return _resp, _result.Data, nil
}

func NewMockClient(opts ...ddhttp.DdClientOption) *MockClient {
	defaultProvider := ddhttp.NewServiceProvider("DDHTTP")
	defaultClient := ddhttp.NewClient()

	svcClient := &MockClient{
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

func Test_metrics(t *testing.T) {
	Convey("Should be equal to go-doudou", t, func() {
		go func() {
			srv := ddhttp.NewDefaultHttpSrv()
			srv.AddRoute(Routes(NewMocksvcHandler())...)
			srv.Run()
		}()
		time.Sleep(10 * time.Millisecond)
		os.Setenv("DDHTTP", "http://localhost:8088/v1")
		client := NewMockClient()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, ret, err := client.GetUser(ctx, nil)
		So(err, ShouldBeNil)
		So(ret, ShouldEqual, "go-doudou")
	})
}

func Test_NacosConfigType(t *testing.T) {
	Convey("Should be equal to go-doudou with nacos config", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		dataId := ".env"
		configClient := mock.NewMockIConfigClient(ctrl)
		configClient.
			EXPECT().
			GetConfig(vo.ConfigParam{
				DataId: dataId,
				Group:  config.DefaultGddNacosConfigGroup,
			}).
			AnyTimes().
			Return("GDD_SERVICE_NAME=configmgr\n\nGDD_READ_TIMEOUT=60s\nGDD_WRITE_TIMEOUT=60s\nGDD_IDLE_TIMEOUT=120s", nil)

		configClient.
			EXPECT().
			ListenConfig(gomock.Any()).
			AnyTimes().
			Return(nil)

		configmgr.NewConfigClient = func(param vo.NacosClientParam) (iClient config_client.IConfigClient, err error) {
			return configClient, nil
		}

		configmgr.NacosClient = configmgr.NewNacosConfigMgr([]string{dataId},
			config.DefaultGddNacosConfigGroup, configmgr.DotenvConfigFormat, config.DefaultGddNacosNamespaceId, configClient, cache.NewConcurrentMap())

		_ = config.GddConfigRemoteType.Write(config.NacosConfigType)
		config.GddNacosConfigDataid.Write(dataId)
		config.GddPort.Write("6060")
		ddhttp.InitialiseRemoteConfigListener()
		go func() {
			srv := ddhttp.NewDefaultHttpSrv()
			srv.AddRoute(Routes(NewMocksvcHandler())...)
			srv.Run()
		}()
		time.Sleep(10 * time.Millisecond)
		os.Setenv("DDHTTP", "http://localhost:6060/v1")
		client := NewMockClient()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, ret, err := client.GetUser(ctx, nil)
		So(err, ShouldBeNil)
		So(ret, ShouldEqual, "go-doudou")
	})
}

func Test_UnknownRemoteConfigType(t *testing.T) {
	Convey("Should be equal to go-doudou with unknown remote config type", t, func() {
		_ = config.GddConfigRemoteType.Write("Unknown")
		config.GddPort.Write("6061")
		ddhttp.InitialiseRemoteConfigListener()
		go func() {
			srv := ddhttp.NewDefaultHttpSrv()
			srv.AddRoute(Routes(NewMocksvcHandler())...)
			srv.Run()
		}()
		time.Sleep(10 * time.Millisecond)
		os.Setenv("DDHTTP", "http://localhost:6061/v1")
		client := NewMockClient()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, ret, err := client.GetUser(ctx, nil)
		So(err, ShouldBeNil)
		So(ret, ShouldEqual, "go-doudou")
	})
}

func Test_ApolloConfigType(t *testing.T) {
	Convey("Should be equal to go-doudou with apollo config", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		configClient := mock.NewMockClient(ctrl)
		factory := &memory.DefaultCacheFactory{}
		cache := factory.Create()
		cache.Set("gdd.retry.count", "3", 0)
		cache.Set("gdd.weight", "5", 0)
		configClient.
			EXPECT().
			GetConfigCache(config.DefaultGddApolloNamespace).
			AnyTimes().
			Return(cache)

		configClient.
			EXPECT().
			AddChangeListener(gomock.Any()).
			AnyTimes().
			Return()

		configmgr.StartWithConfig = func(loadAppConfig func() (*apolloConfig.AppConfig, error)) (agollo.Client, error) {
			_, _ = loadAppConfig()
			return configClient, nil
		}

		configmgr.ApolloClient = configClient

		_ = config.GddConfigRemoteType.Write(config.ApolloConfigType)
		config.GddPort.Write("6062")
		ddhttp.InitialiseRemoteConfigListener()
		go func() {
			srv := ddhttp.NewDefaultHttpSrv()
			srv.AddRoute(Routes(NewMocksvcHandler())...)
			srv.Run()
		}()
		time.Sleep(10 * time.Millisecond)
		os.Setenv("DDHTTP", "http://localhost:6062/v1")
		client := NewMockClient()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, ret, err := client.GetUser(ctx, nil)
		So(err, ShouldBeNil)
		So(ret, ShouldEqual, "go-doudou")
	})
}

func TestCallbackOnChange(t *testing.T) {
	Convey("Environment variable GDD_MANAGE_USER should be changed", t, func() {
		listener := ddhttp.NewHttpConfigListener()
		ddhttp.CallbackOnChange(listener)(&configmgr.NacosChangeEvent{
			Namespace: "",
			Group:     "",
			DataId:    "",
			Changes: map[string]maputils.Change{
				"gdd.manage.user": {
					OldValue:   "admin",
					NewValue:   "go-doudou",
					ChangeType: maputils.MODIFIED,
				},
			},
		})
		So(config.GddManageUser.Load(), ShouldEqual, "")
		ddhttp.CallbackOnChange(listener)(&configmgr.NacosChangeEvent{
			Namespace: "",
			Group:     "",
			DataId:    "",
			Changes: map[string]maputils.Change{
				"gdd.manage.user": {
					OldValue:   "admin",
					NewValue:   "go-doudou",
					ChangeType: maputils.MODIFIED,
				},
			},
		})
		So(config.GddManageUser.Load(), ShouldEqual, "go-doudou")
	})
}

func Test_log_get_text(t *testing.T) {
	Convey("Should be equal to go-doudou", t, func() {
		config.GddLogReqEnable.Write("true")
		config.GddPort.Write("6063")
		go func() {
			srv := ddhttp.NewDefaultHttpSrv()
			srv.AddRoute(Routes(NewMocksvcHandler())...)
			srv.Run()
		}()
		time.Sleep(10 * time.Millisecond)
		os.Setenv("DDHTTP", "http://localhost:6063/v1")
		client := NewMockClient()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, ret, err := client.GetUser(ctx, nil)
		So(err, ShouldBeNil)
		So(ret, ShouldEqual, "go-doudou")
	})
}

func Test_log_post_json(t *testing.T) {
	Convey("Should be equal to OK", t, func() {
		config.GddLogReqEnable.Write("true")
		config.GddPort.Write("6064")
		go func() {
			srv := ddhttp.NewDefaultHttpSrv()
			srv.AddRoute(Routes(NewMocksvcHandler())...)
			srv.Run()
		}()
		time.Sleep(10 * time.Millisecond)
		os.Setenv("DDHTTP", "http://localhost:6064/v1")
		client := NewMockClient()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, ret, err := client.SaveUser(ctx, nil, UserVo{
			Username: "go-doudou",
			Password: "go-doudou",
		})
		So(err, ShouldBeNil)
		So(ret, ShouldEqual, "OK")
	})
}

func Test_log_post_formdata(t *testing.T) {
	Convey("Should be equal to OK", t, func() {
		config.GddLogReqEnable.Write("true")
		config.GddPort.Write("6065")
		go func() {
			srv := ddhttp.NewDefaultHttpSrv()
			srv.AddRoute(Routes(NewMocksvcHandler())...)
			srv.Run()
		}()
		time.Sleep(10 * time.Millisecond)
		os.Setenv("DDHTTP", "http://localhost:6065/v1")
		client := NewMockClient()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, ret, err := client.SignUp(ctx, nil, "go-doudou", "go-doudou")
		So(err, ShouldBeNil)
		So(ret, ShouldEqual, "OK")
	})
}

func Test_basicauth_401(t *testing.T) {
	Convey("Should return 401", t, func() {
		config.GddPort.Write("6066")
		config.GddManagePass.Write("admin")
		go func() {
			srv := ddhttp.NewDefaultHttpSrv()
			srv.AddRoute(Routes(NewMocksvcHandler())...)
			srv.Run()
		}()
		time.Sleep(10 * time.Millisecond)
		resp, err := http.Get("http://localhost:6066/go-doudou/config")
		So(err, ShouldBeNil)
		So(resp.StatusCode, ShouldEqual, 401)
	})
}

func Test_basicauth_200(t *testing.T) {
	Convey("Should return 200", t, func() {
		config.GddPort.Write("6066")
		config.GddManageUser.Write("admin")
		config.GddManagePass.Write("admin")
		go func() {
			srv := ddhttp.NewDefaultHttpSrv()
			srv.AddRoute(Routes(NewMocksvcHandler())...)
			srv.Run()
		}()
		time.Sleep(10 * time.Millisecond)
		resp, err := http.Get("http://admin:admin@localhost:6066/go-doudou/config")
		So(err, ShouldBeNil)
		So(resp.StatusCode, ShouldEqual, 200)
	})
}

func Test_recovery(t *testing.T) {
	Convey("Should recovery from panic", t, func() {
		config.GddPort.Write("6067")
		go func() {
			srv := ddhttp.NewDefaultHttpSrv()
			srv.AddRoute(Routes(NewMocksvcHandler())...)
			srv.Run()
		}()
		time.Sleep(10 * time.Millisecond)
		os.Setenv("DDHTTP", "http://localhost:6067/v1")
		client := NewMockClient()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, _, err := client.GetPanic(ctx, nil)
		So(err, ShouldNotBeNil)
	})
}

func Test_bulkhead(t *testing.T) {
	Convey("Should work with bulkhead", t, func() {
		config.GddPort.Write("6068")
		go func() {
			srv := ddhttp.NewDefaultHttpSrv()
			srv.AddRoute(Routes(NewMocksvcHandler())...)
			srv.AddMiddleware(ddhttp.BulkHead(4, 500*time.Millisecond))
			srv.Run()
		}()
		time.Sleep(10 * time.Millisecond)
		os.Setenv("DDHTTP", "http://localhost:6068/v1")
		client := NewMockClient()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, ret, err := client.SignUp(ctx, nil, "go-doudou", "go-doudou")
		So(err, ShouldBeNil)
		So(ret, ShouldEqual, "OK")
	})
}

func Test_bulkhead_fail(t *testing.T) {
	Convey("Should fail with bulkhead", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		runner := httpMock.NewMockRunner(ctrl)
		runner.
			EXPECT().
			Run(gomock.Any(), gomock.Any()).
			AnyTimes().
			Return(errors.New("mock runner test error"))

		ddhttp.RunnerChain = func(middlewares ...goresilience.Middleware) goresilience.Runner {
			return runner
		}
		config.GddPort.Write("6069")
		go func() {
			srv := ddhttp.NewDefaultHttpSrv()
			srv.AddRoute(Routes(NewMocksvcHandler())...)
			srv.AddMiddleware(ddhttp.BulkHead(4, 500*time.Millisecond))
			srv.Run()
		}()
		time.Sleep(10 * time.Millisecond)
		os.Setenv("DDHTTP", "http://localhost:6069/v1")
		client := NewMockClient()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, _, err := client.SignUp(ctx, nil, "go-doudou", "go-doudou")
		So(err.Error(), ShouldEqual, "too many requests")
	})
}
