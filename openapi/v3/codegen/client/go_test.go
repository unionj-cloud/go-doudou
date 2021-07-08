package client

import (
	"encoding/json"
	"fmt"
	v3 "github.com/unionj-cloud/go-doudou/openapi/v3"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func Test_genGoVo(t *testing.T) {
	testdir := pathutils.Abs("../testfiles")
	api := loadApi(path.Join(testdir, "petstore3.json"))
	genGoVo(api.Components.Schemas, filepath.Join(testdir, "client", "vo.go"))
}

func Test_genGoVo_Omit(t *testing.T) {
	testdir := pathutils.Abs("../testfiles")
	api := loadApi(path.Join(testdir, "petstore3.json"))
	omitempty = true
	genGoVo(api.Components.Schemas, filepath.Join(testdir, "client", "vo.go"))
}

func Test_genGoHttp(t *testing.T) {
	testdir := pathutils.Abs("../testfiles")
	api := loadApi(path.Join(testdir, "petstore3.json"))
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
		genGoHttp(paths, svcname, filepath.Join(testdir, "client"), "")
	}
}

func Test_genGoHttp1(t *testing.T) {
	testdir := pathutils.Abs("../testfiles")
	api := loadApi(path.Join(testdir, "test1.json"))
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
		genGoHttp(paths, svcname, filepath.Join(testdir, "client"), "")
	}
}

func Test_genGoHttp2(t *testing.T) {
	testdir := pathutils.Abs("../testfiles")
	api := loadApi(path.Join(testdir, "test2.json"))
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
		genGoHttp(paths, svcname, filepath.Join(testdir, "client"), "")
	}
}

func Test_genGoHttp_Omit(t *testing.T) {
	testdir := pathutils.Abs("../testfiles")
	api := loadApi(path.Join(testdir, "petstore3.json"))
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
		genGoHttp(paths, svcname, filepath.Join(testdir, "client"), "")
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
