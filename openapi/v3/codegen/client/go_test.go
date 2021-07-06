package client

import (
	"github.com/unionj-cloud/go-doudou/pathutils"
	"path"
	"testing"
)

func Test_genGoVo(t *testing.T) {
	testdir := pathutils.Abs("../testfiles")
	api := loadApi(path.Join(testdir, "petstore3.json"))
	genGoVo(api.Components.Schemas, path.Join(testdir, "client", "vo.go"))
}
