package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/v2/cmd"
)

func TestShutdownCmd(t *testing.T) {
	dir := filepath.Join(testDir, "testsvc")
	err := os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	assert.Panics(t, func() {
		ExecuteCommandC(cmd.GetRootCmd(), []string{"svc", "shutdown"}...)
	})
}
