package cmd

import (
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/svc"
	"os"
	"testing"
)

func TestDeployCmd(t *testing.T) {
	dir := testDir + "deploycmd"
	receiver := svc.Svc{
		Dir: dir,
	}
	receiver.Init()
	defer os.RemoveAll(dir)
	err := os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	assert.Panics(t, func() {
		ExecuteCommandC(rootCmd, []string{"svc", "deploy"}...)
	})
}
