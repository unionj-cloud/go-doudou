package cmd

import (
	"github.com/stretchr/testify/assert"
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
		ExecuteCommandC(rootCmd, []string{"svc", "deploy"}...)
	})
}
