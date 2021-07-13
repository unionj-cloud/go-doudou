package codegen

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/astutils"
	"io/ioutil"
	"os"
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
			want: "pageusers",
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
	ic := astutils.BuildInterfaceCollector(dir+"/svc.go", astutils.ExprString)
	GenHttpHandler(dir, ic)
	expect := `package httpsrv

import (
	"net/http"

	"github.com/unionj-cloud/go-doudou/svc/config"
	ddhttp "github.com/unionj-cloud/go-doudou/svc/http"
)

type TestfileshttphandlerHandler interface {
	PageUsers(w http.ResponseWriter, r *http.Request)
}

func Routes(handler TestfileshttphandlerHandler) []ddhttp.Route {
	rootPath := config.GddRouteRootPath.Load()
	return []ddhttp.Route{
		{
			"PageUsers",
			"POST",
			rootPath + "/testfileshttphandler/pageusers",
			handler.PageUsers,
		},
	}
}
`
	file := dir + "/transport/httpsrv/handler.go"
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
