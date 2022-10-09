package router

import (
	"runtime"
	"testing"

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
