package codegen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/executils"
)

func TestInitProj(t *testing.T) {
	dir := filepath.Join("testdata", "init")
	os.MkdirAll(dir, os.ModePerm)
	defer os.RemoveAll(dir)
	conf := InitProjConfig{
		Dir:      dir,
		ModName:  "testinit",
		Runner:   executils.CmdRunner{},
		GenSvcGo: true,
	}
	assert.NotPanics(t, func() {
		InitProj(conf)
	})
}
