package router

import (
	"bufio"
	"strings"
	"testing"

	"github.com/valyala/fasthttp"
)

func TestGroup(t *testing.T) {
	r1 := New()
	r2 := r1.Group("/boo")
	r3 := r1.Group("/goo")
	r4 := r1.Group("/moo")
	r5 := r4.Group("/foo")
	r6 := r5.Group("/foo")

	hit := false

	r1.POST("/foo", func(ctx *fasthttp.RequestCtx) {
		hit = true
		ctx.SetStatusCode(fasthttp.StatusOK)
	})
	r2.POST("/bar", func(ctx *fasthttp.RequestCtx) {
		hit = true
		ctx.SetStatusCode(fasthttp.StatusOK)
	})
	r3.POST("/bar", func(ctx *fasthttp.RequestCtx) {
		hit = true
		ctx.SetStatusCode(fasthttp.StatusOK)
	})
	r4.POST("/bar", func(ctx *fasthttp.RequestCtx) {
		hit = true
		ctx.SetStatusCode(fasthttp.StatusOK)
	})
	r5.POST("/bar", func(ctx *fasthttp.RequestCtx) {
		hit = true
		ctx.SetStatusCode(fasthttp.StatusOK)
	})
	r6.POST("/bar", func(ctx *fasthttp.RequestCtx) {
		hit = true
		ctx.SetStatusCode(fasthttp.StatusOK)
	})
	r6.ServeFiles("/static/{filepath:*}", "./")
	r6.ServeFilesCustom("/custom/static/{filepath:*}", &fasthttp.FS{Root: "./"})

	uris := []string{
		"POST /foo HTTP/1.1\r\n\r\n",
		// testing router group - r2 (grouped from r1)
		"POST /boo/bar HTTP/1.1\r\n\r\n",
		// testing multiple router group - r3 (grouped from r1)
		"POST /goo/bar HTTP/1.1\r\n\r\n",
		// testing multiple router group - r4 (grouped from r1)
		"POST /moo/bar HTTP/1.1\r\n\r\n",
		// testing sub-router group - r5 (grouped from r4)
		"POST /moo/foo/bar HTTP/1.1\r\n\r\n",
		// testing multiple sub-router group - r6 (grouped from r5)
		"POST /moo/foo/foo/bar HTTP/1.1\r\n\r\n",
		// testing multiple sub-router group - r6 (grouped from r5) to serve files
		"GET /moo/foo/foo/static/router.go HTTP/1.1\r\n\r\n",
		// testing multiple sub-router group - r6 (grouped from r5) to serve files with custom settings
		"GET /moo/foo/foo/custom/static/router.go HTTP/1.1\r\n\r\n",
	}

	for _, uri := range uris {
		hit = false

		assertWithTestServer(t, uri, r1.Handler, func(rw *readWriter) {
			br := bufio.NewReader(&rw.w)
			var resp fasthttp.Response
			if err := resp.Read(br); err != nil {
				t.Fatalf("Unexpected error when reading response: %s", err)
			}
			if !(resp.Header.StatusCode() == fasthttp.StatusOK) {
				t.Fatalf("Status code %d, want %d", resp.Header.StatusCode(), fasthttp.StatusOK)
			}
			if !strings.Contains(uri, "static") && !hit {
				t.Fatalf("Regular routing failed with router chaining. %s", uri)
			}
		})
	}

	assertWithTestServer(t, "POST /qax HTTP/1.1\r\n\r\n", r1.Handler, func(rw *readWriter) {
		br := bufio.NewReader(&rw.w)
		var resp fasthttp.Response
		if err := resp.Read(br); err != nil {
			t.Fatalf("Unexpected error when reading response: %s", err)
		}
		if !(resp.Header.StatusCode() == fasthttp.StatusNotFound) {
			t.Errorf("NotFound behavior failed with router chaining.")
			t.FailNow()
		}
	})
}

func TestGroup_shortcutsAndHandle(t *testing.T) {
	r := New()
	g := r.Group("/v1")

	shortcuts := []func(path string, handler fasthttp.RequestHandler){
		g.GET,
		g.HEAD,
		g.POST,
		g.PUT,
		g.PATCH,
		g.DELETE,
		g.CONNECT,
		g.OPTIONS,
		g.TRACE,
		g.ANY,
	}

	for _, fn := range shortcuts {
		fn("/bar", func(_ *fasthttp.RequestCtx) {})
	}

	methods := httpMethods[:len(httpMethods)-1] // Avoid customs methods
	for _, method := range methods {
		h, _ := r.Lookup(method, "/v1/bar", nil)
		if h == nil {
			t.Errorf("Bad shorcurt")
		}
	}

	g2 := g.Group("/foo")

	for _, method := range httpMethods {
		g2.Handle(method, "/bar", func(_ *fasthttp.RequestCtx) {})

		h, _ := r.Lookup(method, "/v1/foo/bar", nil)
		if h == nil {
			t.Errorf("Bad shorcurt")
		}
	}
}
