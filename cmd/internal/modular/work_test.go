package modular

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRunner struct {
	goVersion string
}

func (r mockRunner) Output(command string, args ...string) ([]byte, error) {
	return []byte("go version go" + r.goVersion + " darwin/arm64"), nil
}

func (r mockRunner) Run(command string, args ...string) error {
	return nil
}

func (r mockRunner) Start(command string, args ...string) (*exec.Cmd, error) {
	return nil, nil
}

func TestWorkInit_UsesGo125FloorForGoWork(t *testing.T) {
	workDir := t.TempDir()
	receiver := NewWork(WorkConfig{WorkDir: workDir}, mockRunner{goVersion: "1.17.8"})

	assert.NotPanics(t, func() {
		receiver.Init()
	})

	workfile, err := os.ReadFile(filepath.Join(workDir, "go.work"))
	assert.NoError(t, err)
	assert.Contains(t, string(workfile), "go 1.25.0")
	assert.Contains(t, string(workfile), "toolchain go1.25.0")
}
