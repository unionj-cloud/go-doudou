package executils

import (
	"os"
	"os/exec"
)

//go:generate mockgen -destination ../../mock/mock_executils_runner.go -package mock -source=./executils.go

// Runner is mainly for executing shell command
type Runner interface {
	Run(string, ...string) error
	Start(string, ...string) (*exec.Cmd, error)
	Output(string, ...string) ([]byte, error)
}

// CmdRunner implements Runner interface
type CmdRunner struct{}

func (r CmdRunner) Output(command string, args ...string) ([]byte, error) {
	return exec.Command(command, args...).Output()
}

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
