package radix

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/savsgio/gotils/bytes"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
)

func generateHandler() fasthttp.RequestHandler {
	hex := bytes.Rand(make([]byte, 10))

	return func(ctx *fasthttp.RequestCtx) {
		ctx.Write(hex)
	}
}

func testHandlerAndParams(
	t *testing.T, tree *Tree, reqPath string, handler fasthttp.RequestHandler, wantTSR bool, params map[string]interface{},
) {
	for _, ctx := range []*fasthttp.RequestCtx{new(fasthttp.RequestCtx), nil} {

		h, tsr := tree.Get(reqPath, ctx)
		if reflect.ValueOf(handler).Pointer() != reflect.ValueOf(h).Pointer() {
			t.Errorf("Path '%s' handler == %p, want %p", reqPath, h, handler)
		}

		if wantTSR != tsr {
			t.Errorf("Path '%s' tsr == %v, want %v", reqPath, tsr, wantTSR)
		}

		if ctx != nil {
			resultParams := make(map[string]interface{})
			if params == nil {
				params = make(map[string]interface{})
			}

			ctx.VisitUserValues(func(key []byte, value interface{}) {
				resultParams[string(key)] = value
			})

			if !reflect.DeepEqual(resultParams, params) {
				t.Errorf("Path '%s' User values == %v, want %v", reqPath, resultParams, params)
			}
		}
	}
}

func Test_Tree(t *testing.T) {
	type args struct {
		path    string
		reqPath string
		handler fasthttp.RequestHandler
	}

	type want struct {
		tsr    bool
		params map[string]interface{}
	}

	tests := []struct {
		args args
		want want
	}{
		{
			args: args{
				path:    "/users/{name}",
				reqPath: "/users/atreugo",
				handler: generateHandler(),
			},
			want: want{
				params: map[string]interface{}{
					"name": "atreugo",
				},
			},
		},
		{
			args: args{
				path:    "/users",
				reqPath: "/users",
				handler: generateHandler(),
			},
			want: want{
				params: nil,
			},
		},
		{
			args: args{
				path:    "/user/",
				reqPath: "/user",
				handler: generateHandler(),
			},
			want: want{
				tsr:    true,
				params: nil,
			},
		},
		{
			args: args{
				path:    "/",
				reqPath: "/",
				handler: generateHandler(),
			},
			want: want{
				params: nil,
			},
		},
		{
			args: args{
				path:    "/users/{name}/jobs",
				reqPath: "/users/atreugo/jobs",
				handler: generateHandler(),
			},
			want: want{
				params: map[string]interface{}{
					"name": "atreugo",
				},
			},
		},
		{
			args: args{
				path:    "/users/admin",
				reqPath: "/users/admin",
				handler: generateHandler(),
			},
			want: want{
				params: nil,
			},
		},
		{
			args: args{
				path:    "/users/{status}/proc",
				reqPath: "/users/active/proc",
				handler: generateHandler(),
			},
			want: want{
				params: map[string]interface{}{
					"status": "active",
				},
			},
		},
		{
			args: args{
				path:    "/static/{filepath:*}",
				reqPath: "/static/assets/js/main.js",
				handler: generateHandler(),
			},
			want: want{
				params: map[string]interface{}{
					"filepath": "assets/js/main.js",
				},
			},
		},
	}

	tree := New()

	for _, test := range tests {
		tree.Add(test.args.path, test.args.handler)
	}

	for _, test := range tests {
		wantHandler := test.args.handler
		if test.want.tsr {
			wantHandler = nil
		}

		testHandlerAndParams(t, tree, test.args.reqPath, wantHandler, test.want.tsr, test.want.params)
	}

	filepathHandler := generateHandler()

	tree.Add("/{filepath:*}", filepathHandler)

	testHandlerAndParams(t, tree, "/js/main.js", filepathHandler, false, map[string]interface{}{
		"filepath": "js/main.js",
	})
}

func Test_Get(t *testing.T) {
	handler := generateHandler()

	tree := New()
	tree.Add("/api/", handler)

	testHandlerAndParams(t, tree, "/api", nil, true, nil)
	testHandlerAndParams(t, tree, "/api/", handler, false, nil)
	testHandlerAndParams(t, tree, "/notfound", nil, false, nil)

	tree = New()
	tree.Add("/api", handler)

	testHandlerAndParams(t, tree, "/api", handler, false, nil)
	testHandlerAndParams(t, tree, "/api/", nil, true, nil)
	testHandlerAndParams(t, tree, "/notfound", nil, false, nil)
}

func Test_AddWithParam(t *testing.T) {
	handler := generateHandler()

	tree := New()
	tree.Add("/test", handler)
	tree.Add("/api/prefix{version:V[0-9]}_{name:[a-z]+}_sufix/files", handler)
	tree.Add("/api/prefix{version:V[0-9]}_{name:[a-z]+}_sufix/data", handler)
	tree.Add("/api/prefix/files", handler)
	tree.Add("/prefix{name:[a-z]+}suffix/data", handler)
	tree.Add("/prefix{name:[a-z]+}/data", handler)
	tree.Add("/api/{file}.json", handler)

	testHandlerAndParams(t, tree, "/api/prefixV1_atreugo_sufix/files", handler, false, map[string]interface{}{
		"version": "V1", "name": "atreugo",
	})
	testHandlerAndParams(t, tree, "/api/prefixV1_atreugo_sufix/data", handler, false, map[string]interface{}{
		"version": "V1", "name": "atreugo",
	})
	testHandlerAndParams(t, tree, "/prefixatreugosuffix/data", handler, false, map[string]interface{}{
		"name": "atreugo",
	})
	testHandlerAndParams(t, tree, "/prefixatreugo/data", handler, false, map[string]interface{}{
		"name": "atreugo",
	})
	testHandlerAndParams(t, tree, "/api/name.json", handler, false, map[string]interface{}{
		"file": "name",
	})

	// Not found
	testHandlerAndParams(t, tree, "/api/prefixV1_1111_sufix/fake", nil, false, nil)
}

func Test_TreeRootWildcard(t *testing.T) {
	handler := generateHandler()

	tree := New()
	tree.Add("/{filepath:*}", handler)

	testHandlerAndParams(t, tree, "/", handler, false, map[string]interface{}{
		"filepath": "",
	})

	tree.Add("/hello/{a}/{b}/{c}", handler)
	testHandlerAndParams(t, tree, "/hello/a", handler, false, map[string]interface{}{
		"filepath": "hello/a",
	})
}

func Test_TreeNilHandler(t *testing.T) {
	const panicMsg = "nil handler"

	tree := New()

	err := catchPanic(func() {
		tree.Add("/", nil)
	})

	if err == nil {
		t.Fatal("Expected panic")
	}

	if err != nil && panicMsg != fmt.Sprint(err) {
		t.Errorf("Invalid conflict error text (%v)", err)
	}
}

func Test_TreeMutable(t *testing.T) {
	routes := []string{
		"/",
		"/api/{version}",
		"/{filepath:*}",
		"/user{user:a-Z+}",
	}

	handler := generateHandler()
	tree := New()

	for _, route := range routes {
		tree.Add(route, handler)

		err := catchPanic(func() {
			tree.Add(route, handler)
		})

		if err == nil {
			t.Errorf("Route '%s' - Expected panic", route)
		}
	}

	tree.Mutable = true

	for _, route := range routes {
		err := catchPanic(func() {
			tree.Add(route, handler)
		})

		if err != nil {
			t.Errorf("Route '%s' - Unexpected panic: %v", route, err)
		}
	}
}

func Benchmark_Get(b *testing.B) {
	handler := func(ctx *fasthttp.RequestCtx) {}

	tree := New()

	// for i := 0; i < 3000; i++ {
	// 	tree.Add(
	// 		fmt.Sprintf("/%s", bytes.Rand(make([]byte, 15))), handler,
	// 	)
	// }

	tree.Add("/", handler)
	tree.Add("/plaintext", handler)
	tree.Add("/json", handler)
	tree.Add("/fortune", handler)
	tree.Add("/fortune-quick", handler)
	tree.Add("/db", handler)
	tree.Add("/queries", handler)
	tree.Add("/update", handler)

	ctx := new(fasthttp.RequestCtx)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Get("/update", ctx)
	}
}

func Benchmark_GetWithRegex(b *testing.B) {
	handler := func(ctx *fasthttp.RequestCtx) {}

	tree := New()
	ctx := new(fasthttp.RequestCtx)

	tree.Add("/api/{version:v[0-9]}/data", handler)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Get("/api/v1/data", ctx)
	}
}

func Benchmark_GetWithParams(b *testing.B) {
	handler := func(ctx *fasthttp.RequestCtx) {}

	tree := New()
	ctx := new(fasthttp.RequestCtx)

	tree.Add("/api/{version}/data", handler)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Get("/api/v1/data", ctx)
	}
}

func Benchmark_FindCaseInsensitivePath(b *testing.B) {
	handler := func(ctx *fasthttp.RequestCtx) {}

	tree := New()
	buf := bytebufferpool.Get()

	tree.Add("/endpoint", handler)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.FindCaseInsensitivePath("/ENdpOiNT", false, buf)
		buf.Reset()
	}
}
