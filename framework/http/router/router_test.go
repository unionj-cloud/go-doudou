package router

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
)

type readWriter struct {
	net.Conn
	r bytes.Buffer
	w bytes.Buffer
}

var httpMethods = []string{
	fasthttp.MethodGet,
	fasthttp.MethodHead,
	fasthttp.MethodPost,
	fasthttp.MethodPut,
	fasthttp.MethodPatch,
	fasthttp.MethodDelete,
	fasthttp.MethodConnect,
	fasthttp.MethodOptions,
	fasthttp.MethodTrace,
	MethodWild,
	"CUSTOM",
}

func randomHTTPMethod() string {
	method := httpMethods[rand.Intn(len(httpMethods)-1)]

	for method == MethodWild {
		method = httpMethods[rand.Intn(len(httpMethods)-1)]
	}

	return method
}

func buildLocation(host, path string) string {
	return fmt.Sprintf("http://%s%s", host, path)
}

var zeroTCPAddr = &net.TCPAddr{
	IP: net.IPv4zero,
}

func (rw *readWriter) Close() error {
	return nil
}

func (rw *readWriter) Read(b []byte) (int, error) {
	return rw.r.Read(b)
}

func (rw *readWriter) Write(b []byte) (int, error) {
	return rw.w.Write(b)
}

func (rw *readWriter) RemoteAddr() net.Addr {
	return zeroTCPAddr
}

func (rw *readWriter) LocalAddr() net.Addr {
	return zeroTCPAddr
}

func (rw *readWriter) SetReadDeadline(t time.Time) error {
	return nil
}

func (rw *readWriter) SetWriteDeadline(t time.Time) error {
	return nil
}

type assertFn func(rw *readWriter)

func assertWithTestServer(t *testing.T, uri string, handler fasthttp.RequestHandler, fn assertFn) {
	s := &fasthttp.Server{
		Handler: handler,
	}

	rw := &readWriter{}
	ch := make(chan error)

	rw.r.WriteString(uri)
	go func() {
		ch <- s.ServeConn(rw)
	}()
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("return error %s", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("timeout")
	}

	fn(rw)
}

func catchPanic(testFunc func()) (recv interface{}) {
	defer func() {
		recv = recover()
	}()

	testFunc()
	return
}

func TestRouter(t *testing.T) {
	router := New()

	routed := false
	router.Handle(fasthttp.MethodGet, "/user/{name}", func(ctx *fasthttp.RequestCtx) {
		routed = true
		want := "gopher"

		param, ok := ctx.UserValue("name").(string)

		if !ok {
			t.Fatalf("wrong wildcard values: param value is nil")
		}

		if param != want {
			t.Fatalf("wrong wildcard values: want %s, got %s", want, param)
		}
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.SetRequestURI("/user/gopher")

	router.Handler(ctx)

	if !routed {
		t.Fatal("routing failed")
	}
}

func TestRouterAPI(t *testing.T) {
	var handled, get, head, post, put, patch, delete, connect, options, trace, any bool

	httpHandler := func(ctx *fasthttp.RequestCtx) {
		handled = true
	}

	router := New()
	router.GET("/GET", func(ctx *fasthttp.RequestCtx) {
		get = true
	})
	router.HEAD("/HEAD", func(ctx *fasthttp.RequestCtx) {
		head = true
	})
	router.POST("/POST", func(ctx *fasthttp.RequestCtx) {
		post = true
	})
	router.PUT("/PUT", func(ctx *fasthttp.RequestCtx) {
		put = true
	})
	router.PATCH("/PATCH", func(ctx *fasthttp.RequestCtx) {
		patch = true
	})
	router.DELETE("/DELETE", func(ctx *fasthttp.RequestCtx) {
		delete = true
	})
	router.CONNECT("/CONNECT", func(ctx *fasthttp.RequestCtx) {
		connect = true
	})
	router.OPTIONS("/OPTIONS", func(ctx *fasthttp.RequestCtx) {
		options = true
	})
	router.TRACE("/TRACE", func(ctx *fasthttp.RequestCtx) {
		trace = true
	})
	router.ANY("/ANY", func(ctx *fasthttp.RequestCtx) {
		any = true
	})
	router.Handle(fasthttp.MethodGet, "/Handler", httpHandler)

	ctx := new(fasthttp.RequestCtx)

	var request = func(method, path string) {
		ctx.Request.Header.SetMethod(method)
		ctx.Request.SetRequestURI(path)
		router.Handler(ctx)
	}

	request(fasthttp.MethodGet, "/GET")
	if !get {
		t.Error("routing GET failed")
	}

	request(fasthttp.MethodHead, "/HEAD")
	if !head {
		t.Error("routing HEAD failed")
	}

	request(fasthttp.MethodPost, "/POST")
	if !post {
		t.Error("routing POST failed")
	}

	request(fasthttp.MethodPut, "/PUT")
	if !put {
		t.Error("routing PUT failed")
	}

	request(fasthttp.MethodPatch, "/PATCH")
	if !patch {
		t.Error("routing PATCH failed")
	}

	request(fasthttp.MethodDelete, "/DELETE")
	if !delete {
		t.Error("routing DELETE failed")
	}

	request(fasthttp.MethodConnect, "/CONNECT")
	if !connect {
		t.Error("routing CONNECT failed")
	}

	request(fasthttp.MethodOptions, "/OPTIONS")
	if !options {
		t.Error("routing OPTIONS failed")
	}

	request(fasthttp.MethodTrace, "/TRACE")
	if !trace {
		t.Error("routing TRACE failed")
	}

	request(fasthttp.MethodGet, "/Handler")
	if !handled {
		t.Error("routing Handler failed")
	}

	for _, method := range httpMethods {
		request(method, "/ANY")
		if !any {
			t.Errorf("routing ANY failed - Method: %s", method)
		}

		any = false
	}
}

func TestRouterInvalidInput(t *testing.T) {
	router := New()

	handle := func(_ *fasthttp.RequestCtx) {}

	recv := catchPanic(func() {
		router.Handle("", "/", handle)
	})
	if recv == nil {
		t.Fatal("registering empty method did not panic")
	}

	recv = catchPanic(func() {
		router.GET("", handle)
	})
	if recv == nil {
		t.Fatal("registering empty path did not panic")
	}

	recv = catchPanic(func() {
		router.GET("noSlashRoot", handle)
	})
	if recv == nil {
		t.Fatal("registering path not beginning with '/' did not panic")
	}

	recv = catchPanic(func() {
		router.GET("/", nil)
	})
	if recv == nil {
		t.Fatal("registering nil handler did not panic")
	}
}

func TestRouterRegexUserValues(t *testing.T) {
	mux := New()
	mux.GET("/metrics", func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusOK)
	})

	v4 := mux.Group("/v4")
	id := v4.Group("/{id:^[1-9]\\d*}")
	id.GET("/click", func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusOK)
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodGet)
	ctx.Request.SetRequestURI("/v4/123/click")
	mux.Handler(ctx)

	v1 := ctx.UserValue("id")
	if v1 != "123" {
		t.Fatalf(`expected "123" in user value, got %q`, v1)
	}

	ctx.Request.Reset()
	ctx.Request.Header.SetMethod(fasthttp.MethodGet)
	ctx.Request.SetRequestURI("/metrics")
	mux.Handler(ctx)

	if v1 != "123" {
		t.Fatalf(`expected "123" in user value after second call, got %q`, v1)
	}
}

func TestRouterChaining(t *testing.T) {
	router1 := New()
	router2 := New()
	router1.NotFound = router2.Handler

	fooHit := false
	router1.POST("/foo", func(ctx *fasthttp.RequestCtx) {
		fooHit = true
		ctx.SetStatusCode(fasthttp.StatusOK)
	})

	barHit := false
	router2.POST("/bar", func(ctx *fasthttp.RequestCtx) {
		barHit = true
		ctx.SetStatusCode(fasthttp.StatusOK)
	})

	ctx := new(fasthttp.RequestCtx)

	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetRequestURI("/foo")
	router1.Handler(ctx)

	if !(ctx.Response.StatusCode() == fasthttp.StatusOK && fooHit) {
		t.Errorf("Regular routing failed with router chaining.")
		t.FailNow()
	}

	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetRequestURI("/bar")
	router1.Handler(ctx)

	if !(ctx.Response.StatusCode() == fasthttp.StatusOK && barHit) {
		t.Errorf("Chained routing failed with router chaining.")
		t.FailNow()
	}

	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetRequestURI("/qax")
	router1.Handler(ctx)

	if !(ctx.Response.StatusCode() == fasthttp.StatusNotFound) {
		t.Errorf("NotFound behavior failed with router chaining.")
		t.FailNow()
	}
}

func TestRouterMutable(t *testing.T) {
	handler1 := func(_ *fasthttp.RequestCtx) {}
	handler2 := func(_ *fasthttp.RequestCtx) {}

	router := New()
	router.Mutable(true)

	if !router.treeMutable {
		t.Errorf("Router.treesMutables is false")
	}

	for _, method := range httpMethods {
		router.Handle(method, "/", handler1)
	}

	for method := range router.trees {
		if !router.trees[method].Mutable {
			t.Errorf("Method %d - Mutable == %v, want %v", method, router.trees[method].Mutable, true)
		}
	}

	routes := []string{
		"/",
		"/api/{version}",
		"/{filepath:*}",
		"/user{user:.*}",
	}

	router = New()

	for _, route := range routes {
		for _, method := range httpMethods {
			router.Handle(method, route, handler1)
		}

		for _, method := range httpMethods {
			err := catchPanic(func() {
				router.Handle(method, route, handler2)
			})

			if err == nil {
				t.Errorf("Mutable 'false' - Method %s - Route %s - Expected panic", method, route)
			}

			h, _ := router.Lookup(method, route, nil)
			if reflect.ValueOf(h).Pointer() != reflect.ValueOf(handler1).Pointer() {
				t.Errorf("Mutable 'false' - Method %s - Route %s - Handler updated", method, route)
			}
		}

		router.Mutable(true)

		for _, method := range httpMethods {
			err := catchPanic(func() {
				router.Handle(method, route, handler2)
			})

			if err != nil {
				t.Errorf("Mutable 'true' - Method %s - Route %s - Unexpected panic: %v", method, route, err)
			}

			h, _ := router.Lookup(method, route, nil)
			if reflect.ValueOf(h).Pointer() != reflect.ValueOf(handler2).Pointer() {
				t.Errorf("Method %s - Route %s - Handler is not updated", method, route)
			}
		}

		router.Mutable(false)
	}

}

func TestRouterOPTIONS(t *testing.T) {
	handlerFunc := func(_ *fasthttp.RequestCtx) {}

	router := New()
	router.POST("/path", handlerFunc)

	ctx := new(fasthttp.RequestCtx)

	var checkHandling = func(path, expectedAllowed string, expectedStatusCode int) {
		ctx.Request.Header.SetMethod(fasthttp.MethodOptions)
		ctx.Request.SetRequestURI(path)
		router.Handler(ctx)

		if !(ctx.Response.StatusCode() == expectedStatusCode) {
			t.Errorf("OPTIONS handling failed: Code=%d, Header=%v", ctx.Response.StatusCode(), ctx.Response.Header.String())
		} else if allow := string(ctx.Response.Header.Peek("Allow")); allow != expectedAllowed {
			t.Error("unexpected Allow header value: " + allow)
		}
	}

	// test not allowed
	// * (server)
	checkHandling("*", "OPTIONS, POST", fasthttp.StatusOK)

	// path
	checkHandling("/path", "OPTIONS, POST", fasthttp.StatusOK)

	ctx.Request.Header.SetMethod(fasthttp.MethodOptions)
	ctx.Request.SetRequestURI("/doesnotexist")
	router.Handler(ctx)
	if !(ctx.Response.StatusCode() == fasthttp.StatusNotFound) {
		t.Errorf("OPTIONS handling failed: Code=%d, Header=%v", ctx.Response.StatusCode(), ctx.Response.Header.String())
	}

	// add another method
	router.GET("/path", handlerFunc)

	// set a global OPTIONS handler
	router.GlobalOPTIONS = func(ctx *fasthttp.RequestCtx) {
		// Adjust status code to 204
		ctx.SetStatusCode(fasthttp.StatusNoContent)
	}

	// test again
	// * (server)
	checkHandling("*", "GET, OPTIONS, POST", fasthttp.StatusNoContent)

	// path
	checkHandling("/path", "GET, OPTIONS, POST", fasthttp.StatusNoContent)

	// custom handler
	var custom bool
	router.OPTIONS("/path", func(ctx *fasthttp.RequestCtx) {
		custom = true
	})

	// test again
	// * (server)
	checkHandling("*", "GET, OPTIONS, POST", fasthttp.StatusNoContent)
	if custom {
		t.Error("custom handler called on *")
	}

	// path
	ctx.Request.Header.SetMethod(fasthttp.MethodOptions)
	ctx.Request.SetRequestURI("/path")
	router.Handler(ctx)
	if !(ctx.Response.StatusCode() == fasthttp.StatusNoContent) {
		t.Errorf("OPTIONS handling failed: Code=%d, Header=%v", ctx.Response.StatusCode(), ctx.Response.Header.String())
	}
	if !custom {
		t.Error("custom handler not called")
	}
}

func TestRouterNotAllowed(t *testing.T) {
	handlerFunc := func(_ *fasthttp.RequestCtx) {}

	router := New()
	router.POST("/path", handlerFunc)

	ctx := new(fasthttp.RequestCtx)

	var checkHandling = func(path, expectedAllowed string, expectedStatusCode int) {
		ctx.Request.Header.SetMethod(fasthttp.MethodGet)
		ctx.Request.SetRequestURI(path)
		router.Handler(ctx)

		if !(ctx.Response.StatusCode() == expectedStatusCode) {
			t.Errorf("NotAllowed handling failed:: Code=%d, Header=%v", ctx.Response.StatusCode(), ctx.Response.Header.String())
		} else if allow := string(ctx.Response.Header.Peek("Allow")); allow != expectedAllowed {
			t.Error("unexpected Allow header value: " + allow)
		}
	}

	// test not allowed
	checkHandling("/path", "OPTIONS, POST", fasthttp.StatusMethodNotAllowed)

	// add another method
	router.DELETE("/path", handlerFunc)
	router.OPTIONS("/path", handlerFunc) // must be ignored

	// test again
	checkHandling("/path", "DELETE, OPTIONS, POST", fasthttp.StatusMethodNotAllowed)

	// test custom handler
	responseText := "custom method"
	router.MethodNotAllowed = func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusTeapot)
		ctx.Write([]byte(responseText))
	}

	ctx.Response.Reset()
	router.Handler(ctx)

	if got := string(ctx.Response.Body()); !(got == responseText) {
		t.Errorf("unexpected response got %q want %q", got, responseText)
	}
	if ctx.Response.StatusCode() != fasthttp.StatusTeapot {
		t.Errorf("unexpected response code %d want %d", ctx.Response.StatusCode(), fasthttp.StatusTeapot)
	}
	if allow := string(ctx.Response.Header.Peek("Allow")); allow != "DELETE, OPTIONS, POST" {
		t.Error("unexpected Allow header value: " + allow)
	}
}

func testRouterNotFoundByMethod(t *testing.T, method string) {
	handlerFunc := func(_ *fasthttp.RequestCtx) {}
	host := "fast"

	router := New()
	router.Handle(method, "/path", handlerFunc)
	router.Handle(method, "/dir/", handlerFunc)
	router.Handle(method, "/", handlerFunc)
	router.Handle(method, "/{proc}/StaTus", handlerFunc)
	router.Handle(method, "/USERS/{name}/enTRies/", handlerFunc)
	router.Handle(method, "/static/{filepath:*}", handlerFunc)

	// Moved Permanently, request with GET method
	expectedCode := fasthttp.StatusMovedPermanently
	if method == fasthttp.MethodConnect {
		// CONNECT method does not allow redirects, so Not Found (404)
		expectedCode = fasthttp.StatusNotFound
	} else if method != fasthttp.MethodGet {
		// Permanent Redirect, request with same method
		expectedCode = fasthttp.StatusPermanentRedirect
	}

	type testRoute struct {
		route    string
		code     int
		location string
	}

	testRoutes := []testRoute{
		{"", fasthttp.StatusOK, ""},                              // TSR +/ (Not clean by router, this path is cleaned by fasthttp `ctx.Path()`)
		{"/../path", expectedCode, buildLocation(host, "/path")}, // CleanPath (Not clean by router, this path is cleaned by fasthttp `ctx.Path()`)
		{"/nope", fasthttp.StatusNotFound, ""},                   // NotFound
	}

	if method != fasthttp.MethodConnect {
		testRoutes = append(testRoutes, []testRoute{
			{"/path/", expectedCode, buildLocation(host, "/path")},                                   // TSR -/
			{"/dir", expectedCode, buildLocation(host, "/dir/")},                                     // TSR +/
			{"/PATH", expectedCode, buildLocation(host, "/path")},                                    // Fixed Case
			{"/DIR/", expectedCode, buildLocation(host, "/dir/")},                                    // Fixed Case
			{"/PATH/", expectedCode, buildLocation(host, "/path")},                                   // Fixed Case -/
			{"/DIR", expectedCode, buildLocation(host, "/dir/")},                                     // Fixed Case +/
			{"/paTh/?name=foo", expectedCode, buildLocation(host, "/path?name=foo")},                 // Fixed Case With Query Params +/
			{"/paTh?name=foo", expectedCode, buildLocation(host, "/path?name=foo")},                  // Fixed Case With Query Params +/
			{"/sergio/status/", expectedCode, buildLocation(host, "/sergio/StaTus")},                 // Fixed Case With Params -/
			{"/users/atreugo/eNtriEs", expectedCode, buildLocation(host, "/USERS/atreugo/enTRies/")}, // Fixed Case With Params +/
			{"/STatiC/test.go", expectedCode, buildLocation(host, "/static/test.go")},                // Fixed Case Wildcard
		}...)
	}

	reqMethod := method
	if method == MethodWild {
		reqMethod = randomHTTPMethod()
	}

	for _, tr := range testRoutes {
		if runtime.GOOS == "windows" && strings.HasPrefix(tr.route, "/../") {
			// See: https://github.com/valyala/fasthttp/issues/1226
			t.Logf("skipping route '%s %s' on %s, unsupported yet", reqMethod, tr.route, runtime.GOOS)

			continue
		}

		ctx := new(fasthttp.RequestCtx)

		ctx.Request.Header.SetMethod(reqMethod)
		ctx.Request.SetRequestURI(tr.route)
		ctx.Request.SetHost(host)
		router.Handler(ctx)

		statusCode := ctx.Response.StatusCode()
		location := string(ctx.Response.Header.Peek("Location"))
		if !(statusCode == tr.code && (statusCode == fasthttp.StatusNotFound || location == tr.location)) {
			t.Errorf("NotFound handling route %s failed: ReqMethod=%s, Code=%d, Header=%v", method, tr.route, statusCode, location)
		}
	}

	ctx := new(fasthttp.RequestCtx)

	// Test custom not found handler
	var notFound bool
	router.NotFound = func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		notFound = true
	}

	ctx.Request.Header.SetMethod(reqMethod)
	ctx.Request.SetRequestURI("/nope")
	router.Handler(ctx)
	if !(ctx.Response.StatusCode() == fasthttp.StatusNotFound && notFound == true) {
		t.Errorf("Custom NotFound handler failed: Code=%d, Header=%v", ctx.Response.StatusCode(), ctx.Response.Header.String())
	}
	ctx.Response.Reset()
}

func TestRouterNotFound(t *testing.T) {
	for _, method := range httpMethods {
		testRouterNotFoundByMethod(t, method)
	}

	router := New()
	handlerFunc := func(_ *fasthttp.RequestCtx) {}
	host := "fast"
	ctx := new(fasthttp.RequestCtx)

	// Test other method than GET (want 308 instead of 301)
	router.PATCH("/path", handlerFunc)

	ctx.Request.Header.SetMethod(fasthttp.MethodPatch)
	ctx.Request.SetRequestURI("/path/?key=val")
	ctx.Request.SetHost(host)
	router.Handler(ctx)
	if !(ctx.Response.StatusCode() == fasthttp.StatusPermanentRedirect && string(ctx.Response.Header.Peek("Location")) == buildLocation(host, "/path?key=val")) {
		t.Errorf("Custom NotFound handler failed: Code=%d, Header=%v", ctx.Response.StatusCode(), ctx.Response.Header.String())
	}
	ctx.Response.Reset()

	// Test special case where no node for the prefix "/" exists
	router = New()
	router.GET("/a", handlerFunc)

	ctx.Request.Header.SetMethod(fasthttp.MethodPatch)
	ctx.Request.SetRequestURI("/")
	router.Handler(ctx)
	if !(ctx.Response.StatusCode() == fasthttp.StatusNotFound) {
		t.Errorf("NotFound handling route / failed: Code=%d", ctx.Response.StatusCode())
	}
}

func TestRouterNotFound_MethodWild(t *testing.T) {
	postFound, anyFound := false, false

	router := New()
	router.ANY("/{path:*}", func(ctx *fasthttp.RequestCtx) { anyFound = true })
	router.POST("/specific", func(ctx *fasthttp.RequestCtx) { postFound = true })

	for i := 0; i < 100; i++ {
		router.Handle(
			randomHTTPMethod(),
			fmt.Sprintf("/%d", rand.Int63()),
			func(ctx *fasthttp.RequestCtx) {},
		)
	}

	ctx := new(fasthttp.RequestCtx)
	var request = func(method, path string) {
		ctx.Request.Header.SetMethod(method)
		ctx.Request.SetRequestURI(path)
		router.Handler(ctx)
	}

	for _, method := range httpMethods {
		request(method, "/specific")

		if method == fasthttp.MethodPost {
			if !postFound {
				t.Errorf("Method '%s': not found", method)
			}
		} else {
			if !anyFound {
				t.Errorf("Method 'ANY' not found with request method %s", method)
			}
		}

		status := ctx.Response.StatusCode()
		if status != fasthttp.StatusOK {
			t.Errorf("Response status code == %d, want %d", status, fasthttp.StatusOK)
		}

		postFound, anyFound = false, false
		ctx.Response.Reset()
	}
}

func TestRouterPanicHandler(t *testing.T) {
	router := New()
	panicHandled := false

	router.PanicHandler = func(ctx *fasthttp.RequestCtx, p interface{}) {
		panicHandled = true
	}

	router.Handle(fasthttp.MethodPut, "/user/{name}", func(ctx *fasthttp.RequestCtx) {
		panic("oops!")
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(fasthttp.MethodPut)
	ctx.Request.SetRequestURI("/user/gopher")

	defer func() {
		if rcv := recover(); rcv != nil {
			t.Fatal("handling panic failed")
		}
	}()

	router.Handler(ctx)

	if !panicHandled {
		t.Fatal("simulating failed")
	}
}

func testRouterLookupByMethod(t *testing.T, method string) {
	reqMethod := method
	if method == MethodWild {
		reqMethod = randomHTTPMethod()
	}

	routed := false
	wantHandle := func(_ *fasthttp.RequestCtx) {
		routed = true
	}
	wantParams := map[string]string{"name": "gopher"}

	ctx := new(fasthttp.RequestCtx)
	router := New()

	// try empty router first
	handle, tsr := router.Lookup(reqMethod, "/nope", ctx)
	if handle != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", handle)
	}
	if tsr {
		t.Error("Got wrong TSR recommendation!")
	}

	// insert route and try again
	router.Handle(method, "/user/{name}", wantHandle)
	handle, _ = router.Lookup(reqMethod, "/user/gopher", ctx)
	if handle == nil {
		t.Fatal("Got no handle!")
	} else {
		handle(nil)
		if !routed {
			t.Fatal("Routing failed!")
		}
	}

	for expectedKey, expectedVal := range wantParams {
		if ctx.UserValue(expectedKey) != expectedVal {
			t.Errorf("The values %s = %s is not save in context", expectedKey, expectedVal)
		}
	}

	routed = false

	// route without param
	router.Handle(method, "/user", wantHandle)
	handle, _ = router.Lookup(reqMethod, "/user", ctx)
	if handle == nil {
		t.Fatal("Got no handle!")
	} else {
		handle(nil)
		if !routed {
			t.Fatal("Routing failed!")
		}
	}

	for expectedKey, expectedVal := range wantParams {
		if ctx.UserValue(expectedKey) != expectedVal {
			t.Errorf("The values %s = %s is not save in context", expectedKey, expectedVal)
		}
	}

	handle, tsr = router.Lookup(reqMethod, "/user/gopher/", ctx)
	if handle != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", handle)
	}
	if !tsr {
		t.Error("Got no TSR recommendation!")
	}

	handle, tsr = router.Lookup(reqMethod, "/nope", ctx)
	if handle != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", handle)
	}
	if tsr {
		t.Error("Got wrong TSR recommendation!")
	}
}

func TestRouterLookup(t *testing.T) {
	for _, method := range httpMethods {
		testRouterLookupByMethod(t, method)
	}
}

func TestRouterMatchedRoutePath(t *testing.T) {
	route1 := "/user/{name}"
	routed1 := false
	handle1 := func(ctx *fasthttp.RequestCtx) {
		route := ctx.UserValue(MatchedRoutePathParam)
		if route != route1 {
			t.Fatalf("Wrong matched route: want %s, got %s", route1, route)
		}
		routed1 = true
	}

	route2 := "/user/{name}/details"
	routed2 := false
	handle2 := func(ctx *fasthttp.RequestCtx) {
		route := ctx.UserValue(MatchedRoutePathParam)
		if route != route2 {
			t.Fatalf("Wrong matched route: want %s, got %s", route2, route)
		}
		routed2 = true
	}

	route3 := "/"
	routed3 := false
	handle3 := func(ctx *fasthttp.RequestCtx) {
		route := ctx.UserValue(MatchedRoutePathParam)
		if route != route3 {
			t.Fatalf("Wrong matched route: want %s, got %s", route3, route)
		}
		routed3 = true
	}

	router := New()
	router.SaveMatchedRoutePath = true
	router.Handle(fasthttp.MethodGet, route1, handle1)
	router.Handle(fasthttp.MethodGet, route2, handle2)
	router.Handle(fasthttp.MethodGet, route3, handle3)

	ctx := new(fasthttp.RequestCtx)

	ctx.Request.Header.SetMethod(fasthttp.MethodGet)
	ctx.Request.SetRequestURI("/user/gopher")
	router.Handler(ctx)
	if !routed1 || routed2 || routed3 {
		t.Fatal("Routing failed!")
	}

	ctx.Request.Header.SetMethod(fasthttp.MethodGet)
	ctx.Request.SetRequestURI("/user/gopher/details")
	router.Handler(ctx)
	if !routed2 || routed3 {
		t.Fatal("Routing failed!")
	}

	ctx.Request.Header.SetMethod(fasthttp.MethodGet)
	ctx.Request.SetRequestURI("/")
	router.Handler(ctx)
	if !routed3 {
		t.Fatal("Routing failed!")
	}
}

func TestRouterServeFiles(t *testing.T) {
	r := New()

	recv := catchPanic(func() {
		r.ServeFiles("/noFilepath", os.TempDir())
	})
	if recv == nil {
		t.Fatal("registering path not ending with '{filepath:*}' did not panic")
	}
	body := []byte("fake ico")
	ioutil.WriteFile(os.TempDir()+"/favicon.ico", body, 0644)

	r.ServeFiles("/{filepath:*}", os.TempDir())

	assertWithTestServer(t, "GET /favicon.ico HTTP/1.1\r\n\r\n", r.Handler, func(rw *readWriter) {
		br := bufio.NewReader(&rw.w)
		var resp fasthttp.Response
		if err := resp.Read(br); err != nil {
			t.Fatalf("Unexpected error when reading response: %s", err)
		}
		if resp.Header.StatusCode() != 200 {
			t.Fatalf("Unexpected status code %d. Expected %d", resp.Header.StatusCode(), 200)
		}
		if !bytes.Equal(resp.Body(), body) {
			t.Fatalf("Unexpected body %q. Expected %q", resp.Body(), string(body))
		}
	})
}

func TestRouterServeFilesCustom(t *testing.T) {
	r := New()

	root := os.TempDir()

	fs := &fasthttp.FS{
		Root: root,
	}

	recv := catchPanic(func() {
		r.ServeFilesCustom("/noFilepath", fs)
	})
	if recv == nil {
		t.Fatal("registering path not ending with '{filepath:*}' did not panic")
	}
	body := []byte("fake ico")
	ioutil.WriteFile(root+"/favicon.ico", body, 0644)

	r.ServeFilesCustom("/{filepath:*}", fs)

	assertWithTestServer(t, "GET /favicon.ico HTTP/1.1\r\n\r\n", r.Handler, func(rw *readWriter) {
		br := bufio.NewReader(&rw.w)
		var resp fasthttp.Response
		if err := resp.Read(br); err != nil {
			t.Fatalf("Unexpected error when reading response: %s", err)
		}
		if resp.Header.StatusCode() != 200 {
			t.Fatalf("Unexpected status code %d. Expected %d", resp.Header.StatusCode(), 200)
		}
		if !bytes.Equal(resp.Body(), body) {
			t.Fatalf("Unexpected body %q. Expected %q", resp.Body(), string(body))
		}
	})
}

func TestRouterList(t *testing.T) {
	expected := map[string][]string{
		"GET":    {"/bar"},
		"PATCH":  {"/foo"},
		"POST":   {"/v1/users/{name}/{surname?}"},
		"DELETE": {"/v1/users/{id?}"},
	}

	r := New()
	r.GET("/bar", func(ctx *fasthttp.RequestCtx) {})
	r.PATCH("/foo", func(ctx *fasthttp.RequestCtx) {})

	v1 := r.Group("/v1")
	v1.POST("/users/{name}/{surname?}", func(ctx *fasthttp.RequestCtx) {})
	v1.DELETE("/users/{id?}", func(ctx *fasthttp.RequestCtx) {})

	result := r.List()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Router.List() == %v, want %v", result, expected)
	}

}

func TestRouterSamePrefixParamRoute(t *testing.T) {
	var id1, id2, id3, pageSize, page, iid string
	var routed1, routed2, routed3 bool

	r := New()
	v1 := r.Group("/v1")
	v1.GET("/foo/{id}/{pageSize}/{page}", func(ctx *fasthttp.RequestCtx) {
		id1 = ctx.UserValue("id").(string)
		pageSize = ctx.UserValue("pageSize").(string)
		page = ctx.UserValue("page").(string)
		routed1 = true
	})
	v1.GET("/foo/{id}/{iid}", func(ctx *fasthttp.RequestCtx) {
		id2 = ctx.UserValue("id").(string)
		iid = ctx.UserValue("iid").(string)
		routed2 = true
	})
	v1.GET("/foo/{id}", func(ctx *fasthttp.RequestCtx) {
		id3 = ctx.UserValue("id").(string)
		routed3 = true
	})

	req := new(fasthttp.RequestCtx)
	req.Request.SetRequestURI("/v1/foo/1/20/4")
	r.Handler(req)
	req = new(fasthttp.RequestCtx)
	req.Request.SetRequestURI("/v1/foo/2/3")
	r.Handler(req)
	req = new(fasthttp.RequestCtx)
	req.Request.SetRequestURI("/v1/foo/v3")
	r.Handler(req)

	if !routed1 {
		t.Error("/foo/{id}/{pageSize}/{page} not routed.")
	}
	if !routed2 {
		t.Error("/foo/{id}/{iid} not routed")
	}

	if !routed3 {
		t.Error("/foo/{id} not routed")
	}

	if id1 != "1" {
		t.Errorf("/foo/{id}/{pageSize}/{page} id expect: 1 got %s", id1)
	}

	if pageSize != "20" {
		t.Errorf("/foo/{id}/{pageSize}/{page} pageSize expect: 20 got %s", pageSize)
	}

	if page != "4" {
		t.Errorf("/foo/{id}/{pageSize}/{page} page expect: 4 got %s", page)
	}

	if id2 != "2" {
		t.Errorf("/foo/{id}/{iid} id expect: 2 got %s", id2)
	}

	if iid != "3" {
		t.Errorf("/foo/{id}/{iid} iid expect: 3 got %s", iid)
	}

	if id3 != "v3" {
		t.Errorf("/foo/{id} id expect: v3 got %s", id3)
	}
}

func BenchmarkAllowed(b *testing.B) {
	handlerFunc := func(_ *fasthttp.RequestCtx) {}

	router := New()
	router.POST("/path", handlerFunc)
	router.GET("/path", handlerFunc)

	b.Run("Global", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = router.allowed("*", fasthttp.MethodOptions)
		}
	})
	b.Run("Path", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = router.allowed("/path", fasthttp.MethodOptions)
		}
	})
}

func BenchmarkRouterGet(b *testing.B) {
	r := New()
	r.GET("/hello", func(ctx *fasthttp.RequestCtx) {})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetRequestURI("/hello")

	for i := 0; i < b.N; i++ {
		r.Handler(ctx)
	}
}

func BenchmarkRouterParams(b *testing.B) {
	r := New()
	r.GET("/{id}", func(ctx *fasthttp.RequestCtx) {})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetRequestURI("/hello")

	for i := 0; i < b.N; i++ {
		r.Handler(ctx)
	}
}

func BenchmarkRouterANY(b *testing.B) {
	r := New()
	r.GET("/data", func(ctx *fasthttp.RequestCtx) {})
	r.ANY("/", func(ctx *fasthttp.RequestCtx) {})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetRequestURI("/")

	for i := 0; i < b.N; i++ {
		r.Handler(ctx)
	}
}

func BenchmarkRouterGet_ANY(b *testing.B) {
	resp := []byte("Bench GET")
	respANY := []byte("Bench GET (ANY)")

	r := New()
	r.GET("/", func(ctx *fasthttp.RequestCtx) {
		ctx.Success("text/plain", resp)
	})
	r.ANY("/", func(ctx *fasthttp.RequestCtx) {
		ctx.Success("text/plain", respANY)
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod("UNICORN")
	ctx.Request.SetRequestURI("/")

	for i := 0; i < b.N; i++ {
		r.Handler(ctx)
	}
}

func BenchmarkRouterNotFound(b *testing.B) {
	r := New()
	r.GET("/bench", func(ctx *fasthttp.RequestCtx) {})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetRequestURI("/notfound")

	for i := 0; i < b.N; i++ {
		r.Handler(ctx)
	}
}

func BenchmarkRouterFindCaseInsensitive(b *testing.B) {
	r := New()
	r.GET("/bench", func(ctx *fasthttp.RequestCtx) {})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetRequestURI("/BenCh/.")

	for i := 0; i < b.N; i++ {
		r.Handler(ctx)
	}
}

func BenchmarkRouterRedirectTrailingSlash(b *testing.B) {
	r := New()
	r.GET("/bench/", func(ctx *fasthttp.RequestCtx) {})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetRequestURI("/bench")

	for i := 0; i < b.N; i++ {
		r.Handler(ctx)
	}
}

func Benchmark_Get(b *testing.B) {
	handler := func(ctx *fasthttp.RequestCtx) {}

	r := New()

	r.GET("/", handler)
	r.GET("/plaintext", handler)
	r.GET("/json", handler)
	r.GET("/fortune", handler)
	r.GET("/fortune-quick", handler)
	r.GET("/db", handler)
	r.GET("/queries", handler)
	r.GET("/update", handler)

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetRequestURI("/update")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.Handler(ctx)
	}
}
