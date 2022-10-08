package radix

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

func Test_findWildPath(t *testing.T) {
	type test struct {
		path string
		want wildPath
	}

	tests := []test{
		{
			path: "/api/{param1}/data",
			want: wildPath{
				path:  "{param1}",
				keys:  []string{"param1"},
				start: 5,
				end:   13,
				pType: param,
				regex: nil,
			},
		},
		{
			path: "/api/{param1}_{param2}/data",
			want: wildPath{
				path:  "{param1}_{param2}",
				keys:  []string{"param1", "param2"},
				start: 5,
				end:   22,
				pType: param,
				regex: regexp.MustCompile("(.*)_(.*)"),
			},
		},
		{
			path: "/api/{param1:[a-z]{2}}/data",
			want: wildPath{
				path:  "{param1:[a-z]{2}}",
				keys:  []string{"param1"},
				start: 5,
				end:   22,
				pType: param,
				regex: regexp.MustCompile("([a-z]{2})"),
			},
		},
		{
			path: "/api/{param1:[a-z]{1}:[0-9]{1}}/data",
			want: wildPath{
				path:  "{param1:[a-z]{1}:[0-9]{1}}",
				keys:  []string{"param1"},
				start: 5,
				end:   31,
				pType: param,
				regex: regexp.MustCompile("([a-z]{1}:[0-9]{1})"),
			},
		},
		{
			path: "/api/{param1:[a-z]{3}}_{param2}/data",
			want: wildPath{
				path:  "{param1:[a-z]{3}}_{param2}",
				keys:  []string{"param1", "param2"},
				start: 5,
				end:   31,
				pType: param,
				regex: regexp.MustCompile("([a-z]{3})_(.*)"),
			},
		},
		{
			path: "/api/prefix{param1:[a-z]{3}}_{param2}suffix/data",
			want: wildPath{
				path:  "{param1:[a-z]{3}}_{param2}suffix",
				keys:  []string{"param1", "param2"},
				start: 11,
				end:   43,
				pType: param,
				regex: regexp.MustCompile("([a-z]{3})_(.*)suffix"),
			},
		},
	}

	for _, test := range tests {
		fullPath := test.path

		result := findWildPath(test.path, fullPath)

		if result.path != test.want.path {
			t.Errorf("wildPath.path == %s, want %s", result.path, test.want.path)
		}

		if !reflect.DeepEqual(result.keys, test.want.keys) {
			t.Errorf("wildPath.key == %v, want %v", result.keys, test.want.keys)
		}

		if result.start != test.want.start {
			t.Errorf("wildPath.start == %d, want %d", result.start, test.want.start)
		}

		if result.end != test.want.end {
			t.Errorf("wildPath.end == %d, want %d", result.end, test.want.end)
		}

		resultHasRegex := result.regex != nil
		wantHasRegex := test.want.regex != nil
		if resultHasRegex && wantHasRegex {
			resultRegex := result.regex.String()
			wantRegex := test.want.regex.String()

			if resultRegex != wantRegex {
				t.Errorf("wildPath.regex == %s, want %s", resultRegex, wantRegex)
			}
		} else if resultHasRegex != wantHasRegex {
			t.Errorf("wildPath.regex == %v, want %v", result.regex, test.want.regex)
		}
	}
}

func Test_findWildPathConflict(t *testing.T) {
	type test struct {
		path        string
		wantErr     bool
		wantErrText string
	}

	tests := []test{
		{
			path:        "/api/{version}{fake}/data",
			wantErr:     true,
			wantErrText: "the wildcards must be separated by at least 1 char",
		},
		{
			path:        "/api/{}/data",
			wantErr:     true,
			wantErrText: "wildcards must be named with a non-empty name in path '/api/{}/data'",
		},
		{
			path:        "/api/{version{bad}}/data",
			wantErr:     true,
			wantErrText: "the char '{' is not allowed in the param name",
		},
		{
			path:        "/api/{version{bad}:^[a-z]{2}}/data",
			wantErr:     true,
			wantErrText: "the char '{' is not allowed in the param name",
		},
		{
			path:    "/api/{version:^[a-z]{2}}/data",
			wantErr: false,
		},
		{
			path:        "/api/{version:^[a-z]{2}}{bad}/data",
			wantErr:     true,
			wantErrText: "the wildcards must be separated by at least 1 char",
		},
		{
			path:        "/api/{version{bad1}:^[a-z]{2}:123}{bad2}/data",
			wantErr:     true,
			wantErrText: "the char '{' is not allowed in the param name",
		},
		{
			path:        "/api/{version:^[a-z]{2}:123}/data",
			wantErr:     false,
			wantErrText: "",
		},
		{
			path:        "/api/{param1:^[a-z]{3}}_{param2}/data",
			wantErr:     false,
			wantErrText: "",
		},
	}

	for _, test := range tests {
		fullPath := test.path

		err := catchPanic(func() {
			findWildPath(test.path, fullPath)
		})

		if test.wantErr != (err != nil) {
			t.Errorf("Unexpected panic for path '%s': %v", test.path, err)
		}

		if err != nil && test.wantErrText != fmt.Sprint(err) {
			t.Errorf("Invalid conflict error text for path '%s': %v", test.path, err)
		}
	}
}
