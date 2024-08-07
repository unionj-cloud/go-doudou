package cmd_test

import (
	"path/filepath"
	"testing"

	"github.com/unionj-cloud/go-doudou/v2/cmd"
)

func TestEnumCmd(t *testing.T) {
	_, _, err := ExecuteCommandC(cmd.GetRootCmd(), []string{"enum", "-f", filepath.Join("testdata", "testsvc", "vo", "vo.go")}...)
	if err != nil {
		t.Fatal(err)
	}
}
