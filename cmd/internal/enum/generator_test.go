package enum

import (
	"path/filepath"
	"testing"
)

func TestGenerator_Generate(t *testing.T) {
	file := filepath.Join("./testdata", "dto.go")
	receiver := Generator{
		File: file,
	}
	receiver.Generate()
}
