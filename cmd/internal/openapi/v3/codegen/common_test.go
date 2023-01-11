package codegen

import (
	"github.com/stretchr/testify/assert"
	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/openapi/v3"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/pathutils"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestPattern2Method(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				pattern: "/shelves/{shelf}/books/{book}",
			},
			want: "Shelves_ShelfBooks_Book",
		},
		{
			name: "",
			args: args{
				pattern: "/goodFood/{bigApple}/books/{myBird}",
			},
			want: "Goodfood_BigappleBooks_Mybird",
		},
		{
			name: "",
			args: args{
				pattern: "/api/v1/query_range",
			},
			want: "ApiV1Query_range",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Pattern2Method(tt.args.pattern); got != tt.want {
				t.Errorf("Pattern2Method() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_genGoVo(t *testing.T) {
	testdir := pathutils.Abs("testdata")
	api := v3.LoadAPI(path.Join(testdir, "petstore3.json"))
	generator := &OpenAPICodeGenerator{
		Schemas:       api.Components.Schemas,
		RequestBodies: api.Components.RequestBodies,
		Responses:     api.Components.Responses,
		ApiInfo:       api.Info,
	}
	generator.GenGoDto(api.Components.Schemas, filepath.Join(testdir, "test", "vo.go"), "test", "")
}

func Test_genGoVo_clean(t *testing.T) {
	testdir := pathutils.Abs("testdata")
	api := v3.LoadAPI(path.Join(testdir, "test5.json"))
	generator := &OpenAPICodeGenerator{
		Schemas:       api.Components.Schemas,
		RequestBodies: api.Components.RequestBodies,
		Responses:     api.Components.Responses,
		ApiInfo:       api.Info,
	}
	generator.GenGoDto(api.Components.Schemas, filepath.Join(testdir, "test", "vo.go"), "test", "")
}

func Test_genGoVo_Omit(t *testing.T) {
	testdir := pathutils.Abs("testdata")
	api := v3.LoadAPI(path.Join(testdir, "petstore3.json"))
	generator := &OpenAPICodeGenerator{
		Schemas:       api.Components.Schemas,
		RequestBodies: api.Components.RequestBodies,
		Responses:     api.Components.Responses,
		ApiInfo:       api.Info,
		Omitempty:     true,
	}
	generator.GenGoDto(api.Components.Schemas, filepath.Join(testdir, "test", "vo.go"), "test", "")
}

func Test_genGoHttp(t *testing.T) {
	testdir := pathutils.Abs("testdata")
	api := v3.LoadAPI(path.Join(testdir, "petstore3.json"))
	generator := &OpenAPICodeGenerator{
		Schemas:       api.Components.Schemas,
		RequestBodies: api.Components.RequestBodies,
		Responses:     api.Components.Responses,
		ApiInfo:       api.Info,
	}
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
	operationConverter := &ClientOperationConverter{
		Generator: generator,
	}
	for svcname, paths := range svcmap {
		generator.GenGoHTTP(paths, svcname, filepath.Join(testdir, "test"), "", "test", operationConverter)
	}
}

func Test_genGoHttp1(t *testing.T) {
	testdir := pathutils.Abs("testdata")
	api := v3.LoadAPI(path.Join(testdir, "test1.json"))
	generator := &OpenAPICodeGenerator{
		Schemas:       api.Components.Schemas,
		RequestBodies: api.Components.RequestBodies,
		Responses:     api.Components.Responses,
		ApiInfo:       api.Info,
	}
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
	operationConverter := &ClientOperationConverter{
		Generator: generator,
	}
	for svcname, paths := range svcmap {
		generator.GenGoHTTP(paths, svcname, filepath.Join(testdir, "test"), "", "test", operationConverter)
	}
}

func Test_genGoHttp2(t *testing.T) {
	testdir := pathutils.Abs("testdata")
	api := v3.LoadAPI(path.Join(testdir, "test2.json"))
	generator := &OpenAPICodeGenerator{
		Schemas:       api.Components.Schemas,
		RequestBodies: api.Components.RequestBodies,
		Responses:     api.Components.Responses,
		ApiInfo:       api.Info,
	}
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
	operationConverter := &ClientOperationConverter{
		Generator: generator,
	}
	for svcname, paths := range svcmap {
		generator.GenGoHTTP(paths, svcname, filepath.Join(testdir, "test"), "", "test", operationConverter)
	}
}

func Test_genGoHttp3(t *testing.T) {
	testdir := pathutils.Abs("testdata")
	api := v3.LoadAPI(path.Join(testdir, "test3.json"))
	generator := &OpenAPICodeGenerator{
		Schemas:       api.Components.Schemas,
		RequestBodies: api.Components.RequestBodies,
		Responses:     api.Components.Responses,
		ApiInfo:       api.Info,
	}
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
	operationConverter := &ClientOperationConverter{
		Generator: generator,
	}
	for svcname, paths := range svcmap {
		generator.GenGoHTTP(paths, svcname, filepath.Join(testdir, "test"), "", "test", operationConverter)
	}
}

func Test_genGoHttp4(t *testing.T) {
	testdir := pathutils.Abs("testdata")
	api := v3.LoadAPI("https://petstore3.swagger.io/api/v3/openapi.json")
	generator := &OpenAPICodeGenerator{
		Schemas:       api.Components.Schemas,
		RequestBodies: api.Components.RequestBodies,
		Responses:     api.Components.Responses,
		ApiInfo:       api.Info,
	}
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

	operationConverter := &ClientOperationConverter{
		Generator: generator,
	}
	for svcname, paths := range svcmap {
		generator.GenGoHTTP(paths, svcname, filepath.Join(testdir, "test"), "", "test", operationConverter)
	}
}

func Test_genGoHttp5(t *testing.T) {
	testdir := pathutils.Abs("testdata")
	api := v3.LoadAPI(path.Join(testdir, "test5.json"))
	generator := &OpenAPICodeGenerator{
		Schemas:       api.Components.Schemas,
		RequestBodies: api.Components.RequestBodies,
		Responses:     api.Components.Responses,
		ApiInfo:       api.Info,
	}
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

	operationConverter := &ClientOperationConverter{
		Generator: generator,
	}
	for svcname, paths := range svcmap {
		generator.GenGoHTTP(paths, svcname, filepath.Join(testdir, "test"), "", "test", operationConverter)
	}
}

func Test_loadApiPanic(t *testing.T) {
	assert.Panics(t, func() {
		v3.LoadAPI("notexists.json")
	})
}

func Test_loadApiJsonUnmarshalPanic(t *testing.T) {
	assert.Panics(t, func() {
		v3.LoadAPI("testdata/test4.json")
	})
}

func Test_genGoHttp_Omit(t *testing.T) {
	testdir := pathutils.Abs("testdata")
	api := v3.LoadAPI(path.Join(testdir, "petstore3.json"))
	generator := &OpenAPICodeGenerator{
		Schemas:       api.Components.Schemas,
		RequestBodies: api.Components.RequestBodies,
		Responses:     api.Components.Responses,
		ApiInfo:       api.Info,
		Omitempty:     true,
	}
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

	operationConverter := &ClientOperationConverter{
		Generator: generator,
	}
	for svcname, paths := range svcmap {
		generator.GenGoHTTP(paths, svcname, filepath.Join(testdir, "test"), "", "test", operationConverter)
	}
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
	operationConverter := &ClientOperationConverter{
		Generator: &OpenAPICodeGenerator{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				_, _ = operationConverter.ConvertOperation(tt.args.endpoint, tt.args.httpMethod, tt.args.operation, tt.args.gparams)
			})
		})
	}
}

func Test_genGoVo_api(t *testing.T) {
	testdir := pathutils.Abs("testdata")
	api := v3.LoadAPI(path.Join(testdir, "api-docs.json"))
	generator := &OpenAPICodeGenerator{
		Schemas:       api.Components.Schemas,
		RequestBodies: api.Components.RequestBodies,
		Responses:     api.Components.Responses,
		ApiInfo:       api.Info,
		Omitempty:     true,
	}
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
	operationConverter := &ClientOperationConverter{
		Generator: generator,
	}
	for svcname, paths := range svcmap {
		generator.GenGoHTTP(paths, svcname, filepath.Join(testdir, "test"), "", "test", operationConverter)
	}
}

func Test_genGoVoJava(t *testing.T) {
	testdir := pathutils.Abs("testdata")
	api := v3.LoadAPI(path.Join(testdir, "api-docs.json"))
	generator := &OpenAPICodeGenerator{
		Schemas:       api.Components.Schemas,
		RequestBodies: api.Components.RequestBodies,
		Responses:     api.Components.Responses,
		ApiInfo:       api.Info,
		Omitempty:     true,
	}
	generator.GenGoDto(api.Components.Schemas, filepath.Join(testdir, "test", "vo.go"), "test", "")
}
