package cmd_test

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/pathutils"
)

var testDir string

func init() {
	testDir = pathutils.Abs("testdata")
}

func ExecuteCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

// NewMockSvc new Svc instance for unit test purpose
func NewMockSvc(dir string, opts ...svc.SvcOption) svc.ISvc {
	return svc.NewSvc(dir, svc.WithRunner(mockRunner{}))
}

type mockRunner struct {
}

func (r mockRunner) Output(command string, args ...string) ([]byte, error) {
	return []byte("go version go1.17.8 darwin/amd64"), nil
}

func (r mockRunner) Run(command string, args ...string) error {
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, args...)
	c := exec.Command(os.Args[0], cs...)
	c.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		panic(err)
	}
	return nil
}

func (r mockRunner) Start(command string, args ...string) (*exec.Cmd, error) {
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	return cmd, nil
}
