package codegen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/toolkit/executils"
)

func TestInitProj(t *testing.T) {
	dir := filepath.Join("testdata", "init")
	os.MkdirAll(dir, os.ModePerm)
	defer os.RemoveAll(dir)
	conf := InitProjConfig{
		Dir:      dir,
		ModName:  "testinit",
		Runner:   executils.CmdRunner{},
	}
	assert.NotPanics(t, func() {
		InitProj(conf)
	})
}
