package codegen

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	"io/ioutil"
	"os"
	"path/filepath"
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
			want: "select/books",
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

func Test_routeName(t *testing.T) {
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
			want: "Books",
		},
		{
			name: "2",
			args: args{
				method: "PageUsers",
			},
			want: "PageUsers",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := routeName(tt.args.method); got != tt.want {
				t.Errorf("routeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenHttpHandler(t *testing.T) {
	dir := testDir + "httphandler"
	InitSvc(dir)
	defer os.RemoveAll(dir)
	ic := astutils.BuildInterfaceCollector(filepath.Join(dir, "svc.go"), astutils.ExprString)
	GenHttpHandler(dir, ic, 1)
	expect := `package httpsrv

import (
	"net/http"

	ddmodel "github.com/unionj-cloud/go-doudou/framework/http/model"
)

type TestdatahttphandlerHandler interface {
	PageUsers(w http.ResponseWriter, r *http.Request)
}

func Routes(handler TestdatahttphandlerHandler) []ddmodel.Route {
	return []ddmodel.Route{
		{
			"PageUsers",
			"POST",
			"/testdatahttphandler/pageusers",
			handler.PageUsers,
		},
	}
}
`
	file := filepath.Join(dir, "transport", "httpsrv", "handler.go")
	f, err := os.Open(file)
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != expect {
		fmt.Printf("want %s, got %s\n", expect, string(content))
	}
	assert.Equal(t, expect, string(content))
}
