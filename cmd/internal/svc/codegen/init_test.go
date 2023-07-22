package codegen

import (
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/executils"
	"os"
	"path/filepath"
	"testing"
)

func TestInitProj(t *testing.T) {
	dir := filepath.Join("testdata", "init")
	os.MkdirAll(dir, os.ModePerm)
	defer os.RemoveAll(dir)
	assert.NotPanics(t, func() {
		InitProj(dir, "testinit", executils.CmdRunner{}, true)
	})
}
