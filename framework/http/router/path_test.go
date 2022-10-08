package router

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/savsgio/gotils/strings"
	"github.com/valyala/fasthttp"
)

type cleanPathTest struct {
	path, result string
}

var cleanTests = []cleanPathTest{
	// Already clean
	{"/", "/"},
	{"/abc", "/abc"},
	{"/a/b/c", "/a/b/c"},
	{"/abc/", "/abc/"},
	{"/a/b/c/", "/a/b/c/"},

	// missing root
	{"", "/"},
	{"a/", "/a/"},
	{"abc", "/abc"},
	{"abc/def", "/abc/def"},
	{"a/b/c", "/a/b/c"},

	// Remove doubled slash
	{"//", "/"},
	{"/abc//", "/abc/"},
	{"/abc/def//", "/abc/def/"},
	{"/a/b/c//", "/a/b/c/"},
	{"/abc//def//ghi", "/abc/def/ghi"},
	{"//abc", "/abc"},
	{"///abc", "/abc"},
	{"//abc//", "/abc/"},

	// Remove . elements
	{".", "/"},
	{"./", "/"},
	{"/abc/./def", "/abc/def"},
	{"/./abc/def", "/abc/def"},
	{"/abc/.", "/abc/"},

	// Remove .. elements
	{"..", "/"},
	{"../", "/"},
	{"../../", "/"},
	{"../..", "/"},
	{"../../abc", "/abc"},
	{"/abc/def/ghi/../jkl", "/abc/def/jkl"},
	{"/abc/def/../ghi/../jkl", "/abc/jkl"},
	{"/abc/def/..", "/abc/"},
	{"/abc/def/../..", "/"},
	{"/abc/def/../../..", "/"},
	{"/abc/def/../../..", "/"},
	{"/abc/def/../../../ghi/jkl/../../../mno", "/mno"},

	// Combinations
	{"abc/./../def", "/def"},
	{"abc//./../def", "/def"},
	{"abc/../../././../def", "/def"},
}

func Test_cleanPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.SkipNow()
	}

	req := new(fasthttp.Request)
	uri := req.URI()

	for _, test := range cleanTests {
		uri.SetPath(test.path)
		if s := cleanPath(string(uri.Path())); s != test.result {
			t.Errorf("cleanPath(%q) = %q, want %q", test.path, s, test.result)
		}

		uri.SetPath(test.result)
		if s := cleanPath(string(uri.Path())); s != test.result {
			t.Errorf("cleanPath(%q) = %q, want %q", test.result, s, test.result)
		}
	}
}

func TestGetOptionalPath(t *testing.T) {
	handler := func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusOK)
	}

	expected := []struct {
		path    string
		tsr     bool
		handler fasthttp.RequestHandler
	}{
		{"/show/{name}", false, handler},
		{"/show/{name}/", true, nil},
		{"/show/{name}/{surname}", false, handler},
		{"/show/{name}/{surname}/", true, nil},
		{"/show/{name}/{surname}/at", false, handler},
		{"/show/{name}/{surname}/at/", true, nil},
		{"/show/{name}/{surname}/at/{address}", false, handler},
		{"/show/{name}/{surname}/at/{address}/", true, nil},
		{"/show/{name}/{surname}/at/{address}/{id}", false, handler},
		{"/show/{name}/{surname}/at/{address}/{id}/", true, nil},
		{"/show/{name}/{surname}/at/{address}/{id}/{phone:.*}", false, handler},
		{"/show/{name}/{surname}/at/{address}/{id}/{phone:.*}/", true, nil},
	}

	r := New()
	r.GET("/show/{name}/{surname?}/at/{address?}/{id}/{phone?:.*}", handler)

	for _, e := range expected {
		ctx := new(fasthttp.RequestCtx)

		h, tsr := r.Lookup("GET", e.path, ctx)

		if tsr != e.tsr {
			t.Errorf("TSR (path: %s) == %v, want %v", e.path, tsr, e.tsr)
		}

		if reflect.ValueOf(h).Pointer() != reflect.ValueOf(e.handler).Pointer() {
			t.Errorf("Handler (path: %s) == %p, want %p", e.path, h, e.handler)
		}
	}

	tests := []struct {
		path          string
		optionalPaths []string
	}{
		{"/hello", nil},
		{"/{name}", nil},
		{"/{name?:[a-zA-Z]{5}}", []string{"/", "/{name:[a-zA-Z]{5}}"}},
		{"/{filepath:^(?!api).*}", nil},
		{"/static/{filepath?:^(?!api).*}", []string{"/static", "/static/{filepath:^(?!api).*}"}},
		{"/show/{name?}", []string{"/show", "/show/{name}"}},
	}

	for _, test := range tests {
		optionalPaths := getOptionalPaths(test.path)

		if len(optionalPaths) != len(test.optionalPaths) {
			t.Errorf("getOptionalPaths() len == %d, want %d", len(optionalPaths), len(test.optionalPaths))
		}

		for _, wantPath := range test.optionalPaths {
			if !strings.Include(optionalPaths, wantPath) {
				t.Errorf("The optional path is not returned for '%s': %s", test.path, wantPath)
			}
		}
	}
}
