package client

import (
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
		genGoHttp(paths, svcname, filepath.Join(testdir, "client"))
	}
}
