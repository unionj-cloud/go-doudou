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

func TestGenHttpHandler(t *testing.T) {
	ic := astutils.BuildInterfaceCollector(filepath.Join(testDir, "svc.go"), astutils.ExprString)
	GenHttpHandler(testDir, ic, 1)
	expect := `package httpsrv

import (
	"net/http"

	ddmodel "github.com/unionj-cloud/go-doudou/framework/http/model"
)

type UsersvcHandler interface {
	PageUsers(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	SignUp(w http.ResponseWriter, r *http.Request)
	UploadAvatar(w http.ResponseWriter, r *http.Request)
	DownloadAvatar(w http.ResponseWriter, r *http.Request)
}

func Routes(handler UsersvcHandler) []ddmodel.Route {
	return []ddmodel.Route{
		{
			"PageUsers",
			"POST",
			"/usersvc/pageusers",
			handler.PageUsers,
		},
		{
			"GetUser",
			"GET",
			"/usersvc/user",
			handler.GetUser,
		},
		{
			"SignUp",
			"POST",
			"/usersvc/signup",
			handler.SignUp,
		},
		{
			"UploadAvatar",
			"POST",
			"/usersvc/uploadavatar",
			handler.UploadAvatar,
		},
		{
			"DownloadAvatar",
			"POST",
			"/usersvc/downloadavatar",
			handler.DownloadAvatar,
		},
	}
}

var RouteAnnotationStore = ddmodel.AnnotationStore{
	"PageUsers": {
		{
			Name: "@role",
			Params: []string{
				"user",
			},
		},
	},
	"GetUser": {
		{
			Name: "@role",
			Params: []string{
				"admin",
			},
		},
	},
	"SignUp": {
		{
			Name: "@permission",
			Params: []string{
				"create",
				"update",
			},
		},
		{
			Name: "@role",
			Params: []string{
				"admin",
			},
		},
	},
	"UploadAvatar": {
		{
			Name: "@role",
			Params: []string{
				"user",
			},
		},
	},
}
`
	file := filepath.Join(testDir, "transport", "httpsrv", "handler.go")
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
