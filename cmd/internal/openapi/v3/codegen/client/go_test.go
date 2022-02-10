package client

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	v3 "github.com/unionj-cloud/go-doudou/toolkit/openapi/v3"
	"github.com/unionj-cloud/go-doudou/toolkit/pathutils"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func Test_genGoVo(t *testing.T) {
	testdir := pathutils.Abs("../testdata")
	api := loadAPI(path.Join(testdir, "petstore3.json"))
	genGoVo(api.Components.Schemas, filepath.Join(testdir, "test", "vo.go"), "test")
}

func Test_genGoVo_clean(t *testing.T) {
	testdir := pathutils.Abs("../testdata")
	api := loadAPI(path.Join(testdir, "test5.json"))
	genGoVo(api.Components.Schemas, filepath.Join(testdir, "test", "vo.go"), "test")
}

func Test_genGoVo_Omit(t *testing.T) {
	testdir := pathutils.Abs("../testdata")
	api := loadAPI(path.Join(testdir, "petstore3.json"))
	omitempty = true
	genGoVo(api.Components.Schemas, filepath.Join(testdir, "test", "vo.go"), "test")
}

func Test_genGoHttp(t *testing.T) {
	testdir := pathutils.Abs("../testdata")
	api := loadAPI(path.Join(testdir, "petstore3.json"))
	schemas = api.Components.Schemas
	requestBodies = api.Components.RequestBodies
	svcmap := make(map[string]map[string]v3.Path)
	for endpoint, path := range api.Paths {
		svcname := strings.Split(strings.Trim(endpoint, "/"), "/")[0]
		if value, exists := svcmap[svcname]; exists {
			value[endpoint] = path
		} else {
			svcmap[svcname] = make(map[string]v3.Path)
			svcmap[svcname][endpoint] = path
		}
	}

	for svcname, paths := range svcmap {
		genGoHTTP(paths, svcname, filepath.Join(testdir, "test"), "", "test")
	}
}

func Test_genGoHttp1(t *testing.T) {
	testdir := pathutils.Abs("../testdata")
	api := loadAPI(path.Join(testdir, "test1.json"))
	schemas = api.Components.Schemas
	requestBodies = api.Components.RequestBodies
	svcmap := make(map[string]map[string]v3.Path)
	for endpoint, path := range api.Paths {
		svcname := strings.Split(strings.Trim(endpoint, "/"), "/")[0]
		if value, exists := svcmap[svcname]; exists {
			value[endpoint] = path
		} else {
			svcmap[svcname] = make(map[string]v3.Path)
			svcmap[svcname][endpoint] = path
		}
	}

	for svcname, paths := range svcmap {
		genGoHTTP(paths, svcname, filepath.Join(testdir, "test"), "", "test")
	}
}

func Test_genGoHttp2(t *testing.T) {
	testdir := pathutils.Abs("../testdata")
	api := loadAPI(path.Join(testdir, "test2.json"))
	schemas = api.Components.Schemas
	requestBodies = api.Components.RequestBodies
	svcmap := make(map[string]map[string]v3.Path)
	for endpoint, path := range api.Paths {
		svcname := strings.Split(strings.Trim(endpoint, "/"), "/")[0]
		if value, exists := svcmap[svcname]; exists {
			value[endpoint] = path
		} else {
			svcmap[svcname] = make(map[string]v3.Path)
			svcmap[svcname][endpoint] = path
		}
	}

	for svcname, paths := range svcmap {
		genGoHTTP(paths, svcname, filepath.Join(testdir, "test"), "", "test")
	}
}

func Test_genGoHttp3(t *testing.T) {
	testdir := pathutils.Abs("../testdata")
	api := loadAPI(path.Join(testdir, "test3.json"))
	schemas = api.Components.Schemas
	requestBodies = api.Components.RequestBodies
	svcmap := make(map[string]map[string]v3.Path)
	for endpoint, path := range api.Paths {
		svcname := strings.Split(strings.Trim(endpoint, "/"), "/")[0]
		if value, exists := svcmap[svcname]; exists {
			value[endpoint] = path
		} else {
			svcmap[svcname] = make(map[string]v3.Path)
			svcmap[svcname][endpoint] = path
		}
	}

	for svcname, paths := range svcmap {
		genGoHTTP(paths, svcname, filepath.Join(testdir, "test"), "", "test")
	}
}

func Test_genGoHttp4(t *testing.T) {
	api := loadAPI("https://petstore3.swagger.io/api/v3/openapi.json")
	schemas = api.Components.Schemas
	requestBodies = api.Components.RequestBodies
	svcmap := make(map[string]map[string]v3.Path)
	for endpoint, path := range api.Paths {
		svcname := strings.Split(strings.Trim(endpoint, "/"), "/")[0]
		if value, exists := svcmap[svcname]; exists {
			value[endpoint] = path
		} else {
			svcmap[svcname] = make(map[string]v3.Path)
			svcmap[svcname][endpoint] = path
		}
	}

	for svcname, paths := range svcmap {
		genGoHTTP(paths, svcname, filepath.Join(pathutils.Abs("../testdata"), "test"), "", "test")
	}
}

func Test_genGoHttp5(t *testing.T) {
	testdir := pathutils.Abs("../testdata")
	api := loadAPI(path.Join(testdir, "test5.json"))
	schemas = api.Components.Schemas
	requestBodies = api.Components.RequestBodies
	svcmap := make(map[string]map[string]v3.Path)
	for endpoint, path := range api.Paths {
		svcname := strings.Split(strings.Trim(endpoint, "/"), "/")[0]
		if value, exists := svcmap[svcname]; exists {
			value[endpoint] = path
		} else {
			svcmap[svcname] = make(map[string]v3.Path)
			svcmap[svcname][endpoint] = path
		}
	}

	for svcname, paths := range svcmap {
		genGoHTTP(paths, svcname, filepath.Join(testdir, "test"), "", "test")
	}
}

func Test_loadApiPanic(t *testing.T) {
	assert.Panics(t, func() {
		loadAPI("notexists.json")
	})
}

func Test_loadApiJsonUnmarshalPanic(t *testing.T) {
	assert.Panics(t, func() {
		loadAPI("../testdata/test4.json")
	})
}

func Test_genGoHttp_Omit(t *testing.T) {
	testdir := pathutils.Abs("../testdata")
	api := loadAPI(path.Join(testdir, "petstore3.json"))
	omitempty = true
	schemas = api.Components.Schemas
	requestBodies = api.Components.RequestBodies
	svcmap := make(map[string]map[string]v3.Path)
	for endpoint, path := range api.Paths {
		svcname := strings.Split(strings.Trim(endpoint, "/"), "/")[0]
		if value, exists := svcmap[svcname]; exists {
			value[endpoint] = path
		} else {
			svcmap[svcname] = make(map[string]v3.Path)
			svcmap[svcname][endpoint] = path
		}
	}

	for svcname, paths := range svcmap {
		genGoHTTP(paths, svcname, filepath.Join(testdir, "test"), "", "test")
	}
}

func Example1() {
	a := []int{1, 2, 3}
	ret, _ := json.Marshal(a)
	fmt.Println(string(ret))

	var _ret []int
	b := `[1,2,3]`
	err := json.Unmarshal([]byte(b), &_ret)
	if err != nil {
		panic(err)
	}
	fmt.Println(_ret)
	// Output:
	// [1,2,3]
	//[1 2 3]
}

func Example2() {
	a := [][]float64{{1.0, 2.4, 3.7}, {1.0, 2.4, 3.7}, {1.0, 2.4, 3.7}}
	ret, _ := json.Marshal(a)
	fmt.Println(string(ret))

	var _ret [][]float64
	b := `[[1,2.4,3.7],[1,2.4,3.7],[1,2.4,3.7]]`
	err := json.Unmarshal([]byte(b), &_ret)
	if err != nil {
		panic(err)
	}
	fmt.Println(_ret)
	// Output:
	// [[1,2.4,3.7],[1,2.4,3.7],[1,2.4,3.7]]
	//[[1 2.4 3.7] [1 2.4 3.7] [1 2.4 3.7]]
}

func Example3() {
	a := 15
	ret, _ := json.Marshal(a)
	fmt.Println(string(ret))

	var _ret int
	b := `15`
	err := json.Unmarshal([]byte(b), &_ret)
	if err != nil {
		panic(err)
	}
	fmt.Println(_ret)
	// Output:
	// 15
	//15

}

func Example4() {
	a := "a normal string"
	ret, _ := json.Marshal(a)
	fmt.Println(string(ret))

	var _ret string
	b := `"a normal string"`
	err := json.Unmarshal([]byte(b), &_ret)
	if err != nil {
		panic(err)
	}
	fmt.Println(_ret)
	// Output:
	// "a normal string"
	//a normal string
}

func Example5() {
	a := []string{"a normal string"}
	ret, _ := json.Marshal(a)
	fmt.Println(string(ret))

	var _ret []string
	b := `["a normal string"]`
	err := json.Unmarshal([]byte(b), &_ret)
	if err != nil {
		panic(err)
	}
	fmt.Println(_ret)
	// Output:
	// ["a normal string"]
	//[a normal string]
}

func Test_toMethod(t *testing.T) {
	type args struct {
		endpoint string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				endpoint: "/apps/a32(34/~name/{id}",
			},
			want: "AppsA3234NameId",
		},
		{
			name: "",
			args: args{
				endpoint: "/678/9apps/a32(34/~name/{id}",
			},
			want: "AppsA3234NameId6789",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toMethod(tt.args.endpoint); got != tt.want {
				t.Errorf("toMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpMethod(t *testing.T) {
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
				method: "CreateUser",
			},
			want: "POST",
		},
		{
			name: "",
			args: args{
				method: "GetUserInfo",
			},
			want: "GET",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := httpMethod(tt.args.method); got != tt.want {
				t.Errorf("httpMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenGoClient(t *testing.T) {
	dir := "../testdata/testclient"
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(dir)
	assert.NotPanics(t, func() {
		GenGoClient(dir, "../testdata/petstore3.json", true, "", "client")
	})
}

func Test_operation2Method(t *testing.T) {
	type args struct {
		endpoint   string
		httpMethod string
		operation  *v3.Operation
		gparams    []v3.Parameter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				endpoint:   "/test/operation2Mthod/{pid}",
				httpMethod: "GET",
				operation: &v3.Operation{
					Tags:        []string{"test"},
					Summary:     "This is only for test",
					Description: "Description for test",
					OperationID: "TestOperation2Mthod",
					Parameters:  []v3.Parameter{},
					RequestBody: &v3.RequestBody{
						Description: "This is a description",
						Content: &v3.Content{
							FormURL: &v3.MediaType{
								Schema: &v3.Schema{
									Type: "object",
									Properties: map[string]*v3.Schema{
										"id":    v3.Int64,
										"name":  v3.String,
										"score": v3.Float64,
										"isBoy": v3.Bool,
									},
								},
							},
						},
						Required: true,
					},
					Responses: &v3.Responses{
						Resp200: &v3.Response{
							Description: "this is a response",
							Content: &v3.Content{
								JSON: &v3.MediaType{
									Schema: &v3.Schema{
										Type: "object",
										Properties: map[string]*v3.Schema{
											"code": v3.Int,
											"data": v3.String,
											"err":  v3.String,
										},
									},
								},
							},
						},
					},
				},
				gparams: []v3.Parameter{
					{
						Name:        "companyId",
						In:          v3.InQuery,
						Description: "company ID",
						Required:    true,
						Schema:      v3.Int64,
					},
					{
						Name:        "name",
						In:          v3.InQuery,
						Description: "user name",
						Required:    true,
						Deprecated:  false,
						Schema:      v3.String,
					},
					{
						Name:        "pid",
						In:          v3.InPath,
						Description: "project Id",
						Required:    false,
						Deprecated:  false,
						Schema:      v3.Int64,
					},
					{
						Name:        "token",
						In:          v3.InHeader,
						Description: "user token",
						Required:    true,
						Deprecated:  false,
						Schema:      v3.String,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				endpoint:   "/test/operation2Mthod/{pid}",
				httpMethod: "GET",
				operation: &v3.Operation{
					Tags:        []string{"test"},
					Summary:     "This is only for test",
					Description: "Description for test",
					OperationID: "TestOperation2Mthod",
					Parameters:  []v3.Parameter{},
					RequestBody: &v3.RequestBody{
						Description: "This is a description",
						Content: &v3.Content{
							FormData: &v3.MediaType{
								Schema: &v3.Schema{
									Type: "object",
									Properties: map[string]*v3.Schema{
										"id":      v3.Int64,
										"name":    v3.String,
										"score":   v3.Float64,
										"isBoy":   v3.Bool,
										"photoes": v3.FileArray,
										"doc":     v3.File,
									},
								},
							},
						},
						Required: true,
					},
					Responses: &v3.Responses{
						Resp200: &v3.Response{
							Description: "this is a response",
							Content: &v3.Content{
								JSON: &v3.MediaType{
									Schema: &v3.Schema{
										Type: "object",
										Properties: map[string]*v3.Schema{
											"code": v3.Int,
											"data": v3.String,
											"err":  v3.String,
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				endpoint:   "/test/operation2Mthod/{pid}",
				httpMethod: "GET",
				operation: &v3.Operation{
					Tags:        []string{"test"},
					Summary:     "This is only for test",
					Description: "Description for test",
					OperationID: "TestOperation2Mthod",
					Parameters:  []v3.Parameter{},
					RequestBody: &v3.RequestBody{
						Description: "This is a description",
						Content: &v3.Content{
							FormData: &v3.MediaType{
								Schema: &v3.Schema{
									Type: "object",
									Properties: map[string]*v3.Schema{
										"id":      v3.Int64,
										"name":    v3.String,
										"score":   v3.Float64,
										"isBoy":   v3.Bool,
										"photoes": v3.FileArray,
										"doc":     v3.File,
									},
								},
							},
						},
						Required: true,
					},
					Responses: &v3.Responses{
						Resp200: &v3.Response{
							Description: "this is a response",
							Content: &v3.Content{
								Stream: &v3.MediaType{
									Schema: v3.File,
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				endpoint:   "/test/operation2Mthod/{pid}",
				httpMethod: "GET",
				operation: &v3.Operation{
					Tags:        []string{"test"},
					Summary:     "This is only for test",
					Description: "Description for test",
					OperationID: "TestOperation2Mthod",
					Parameters:  []v3.Parameter{},
					RequestBody: &v3.RequestBody{
						Description: "This is a description",
						Content: &v3.Content{
							TextPlain: &v3.MediaType{
								Schema: &v3.Schema{
									Type: "object",
									Properties: map[string]*v3.Schema{
										"id":    v3.Int64,
										"name":  v3.String,
										"score": v3.Float64,
										"isBoy": v3.Bool,
									},
								},
							},
						},
						Required: true,
					},
					Responses: &v3.Responses{
						Resp200: &v3.Response{
							Description: "this is a response",
							Content: &v3.Content{
								TextPlain: &v3.MediaType{
									Schema: &v3.Schema{
										Type: "object",
										Properties: map[string]*v3.Schema{
											"code": v3.Int,
											"data": v3.String,
											"err":  v3.String,
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				_, _ = operation2Method(tt.args.endpoint, tt.args.httpMethod, tt.args.operation, tt.args.gparams)
			})
		})
	}
}
