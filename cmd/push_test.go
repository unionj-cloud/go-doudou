package cmd

import (
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/svc"
	"os"
	"testing"
)

func TestPushCmd(t *testing.T) {
	dir := testDir + "pushcmd"
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
		ExecuteCommandC(rootCmd, []string{"svc", "push", "-r", "testprivaterepo"}...)
	})
}
