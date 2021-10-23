package executils

import (
	"fmt"
	"os"
	"testing"
)

func ExampleCmdRunner_Run() {
	var runner CmdRunner
	cs := []string{"-test.run=TestHelperProcess", "--"}
	runner.Run(os.Args[0], cs...)
	// Output:
	// testing helper process
}

func ExampleCmdRunner_Start() {
	var runner CmdRunner
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cmd, _ := runner.Start(os.Args[0], cs...)
	if err := cmd.Wait(); err != nil {
		panic(err)
	}
	// Output:
	// testing helper process
}

func TestHelperProcess(*testing.T) {
	fmt.Println("testing helper process")
}
