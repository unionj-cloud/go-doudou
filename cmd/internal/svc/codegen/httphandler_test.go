package codegen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_pattern(t *testing.T) {
	type args struct {
		method string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				method: "GetBooks",
			},
			want: "books",
		},
		{
			name: "2",
			args: args{
				method: "PageUsers",
			},
			want: "page/users",
		},
		{
			name: "3",
			args: args{
				method: "PostSelect_Books",
			},
			want: "select.books",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pattern(tt.args.method); got != tt.want {
				t.Errorf("pattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pattern1(t *testing.T) {
	type args struct {
		method string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				method: "GetHome_Html",
			},
			want: "home.html",
		},
		{
			name: "",
			args: args{
				method: "GetHome_html",
			},
			want: "home.html",
		},
		{
			name: "",
			args: args{
				method: "Get",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, pattern(tt.args.method), "pattern(%v)", tt.args.method)
		})
	}
}
