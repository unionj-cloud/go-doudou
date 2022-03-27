package cmd_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/cmd"
	"os"
	"testing"
)

func TestDeployCmd(t *testing.T) {
	dir := testDir + "/testsvc"
	err := os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	assert.Panics(t, func() {
		ExecuteCommandC(cmd.GetRootCmd(), []string{"svc", "deploy"}...)
	})
}
