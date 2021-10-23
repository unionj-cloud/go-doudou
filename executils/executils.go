package executils

import (
	"os"
	"os/exec"
)

// Runner is mainly for executing shell command
type Runner interface {
	Run(string, ...string) error
	Start(string, ...string) (*exec.Cmd, error)
}

// CmdRunner implements Runner interface
type CmdRunner struct{}

// Run executes commands
func (r CmdRunner) Run(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Start starts the specified command but does not wait for it to complete.
func (r CmdRunner) Start(command string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, cmd.Start()
}
