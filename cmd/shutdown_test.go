package cmd

import (
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/svc"
	"os"
	"testing"
)

func TestShutdownCmd(t *testing.T) {
	dir := testDir + "shutdowncmd"
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
		ExecuteCommandC(rootCmd, []string{"svc", "shutdown"}...)
	})
}
