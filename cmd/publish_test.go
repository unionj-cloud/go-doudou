package cmd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/svc"
	"github.com/unionj-cloud/go-doudou/test"
	"os"
	"testing"
)

var esHost string
var esPort int

func TestMain(m *testing.M) {
	var terminator func()
	terminator, esHost, esPort = test.PrepareTestEnvironment()
	code := m.Run()
	terminator()
	os.Exit(code)
}

func TestPublishCmd(t *testing.T) {
	dir := testDir + "publishcmd"
	receiver := svc.Svc{
		Dir: dir,
	}
	receiver.Init()
	defer os.RemoveAll(dir)
	err := os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	ExecuteCommandC(rootCmd, []string{"svc", "http", "--handler", "-c", "go", "-o", "--doc"}...)
	assert.NotPanics(t, func() {
		ExecuteCommandC(rootCmd, []string{"svc", "publish", "--esaddr", fmt.Sprintf("http://%s:%d", esHost, esPort),
			"--esindex", "doc"}...)
	})
}
