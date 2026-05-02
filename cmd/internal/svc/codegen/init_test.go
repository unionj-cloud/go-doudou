package codegen

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitProj(t *testing.T) {
	dir := filepath.Join("testdata", "init")
	os.MkdirAll(dir, os.ModePerm)
	defer os.RemoveAll(dir)
	conf := InitProjConfig{
		Dir:     dir,
		ModName: "testinit",
		Runner:  mockRunner{goVersion: "1.25.3"},
	}
	assert.NotPanics(t, func() {
		InitProj(conf)
	})
}

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

func TestInitProj_UsesGo125FloorForGeneratedFiles(t *testing.T) {
	dir := t.TempDir()
	conf := InitProjConfig{
		Dir:     dir,
		ModName: "testinit",
		Runner:  mockRunner{goVersion: "1.17.8"},
	}

	assert.NotPanics(t, func() {
		InitProj(conf)
	})

	modfile, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	assert.NoError(t, err)
	assert.Contains(t, string(modfile), "go 1.25.0")

	dockerfile, err := os.ReadFile(filepath.Join(dir, "Dockerfile"))
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(dockerfile), "FROM golang:1.25-alpine AS builder"))
}
