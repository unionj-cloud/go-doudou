package executils

import (
	"fmt"
	"os"
	"testing"
)

func TestCmdRunner_Run(t *testing.T) {
	var runner CmdRunner
	cs := []string{"-test.run=TestHelperProcess", "--"}
	runner.Run(os.Args[0], cs...)
}

func TestCmdRunner_Start(t *testing.T) {
	var runner CmdRunner
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cmd, _ := runner.Start(os.Args[0], cs...)
	if err := cmd.Wait(); err != nil {
		panic(err)
	}
}

func TestHelperProcess(*testing.T) {
	fmt.Println("testing helper process")
}
