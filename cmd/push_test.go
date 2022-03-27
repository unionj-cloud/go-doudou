package cmd_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/cmd"
	"os"
	"testing"
)

func TestPushCmd(t *testing.T) {
	dir := testDir + "/pushcmd"
	receiver := NewMockSvc(dir)
	receiver.Init()
	defer os.RemoveAll(dir)
	err := os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	assert.Panics(t, func() {
		ExecuteCommandC(cmd.GetRootCmd(), []string{"svc", "push", "-r", "testprivaterepo"}...)
	})
}
