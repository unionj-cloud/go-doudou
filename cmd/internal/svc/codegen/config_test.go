package codegen

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenConfig(t *testing.T) {
	dir := testDir + "config"
	InitSvc(dir)
	defer os.RemoveAll(dir)
	GenConfig(dir)
}

func TestGenConfig1(t *testing.T) {
	GenConfig(filepath.Join(testDir))
}
