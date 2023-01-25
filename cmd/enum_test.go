package cmd_test

import (
	"github.com/unionj-cloud/go-doudou/v2/cmd"
	"path/filepath"
	"testing"
)

func TestEnumCmd(t *testing.T) {
	_, _, err := ExecuteCommandC(cmd.GetRootCmd(), []string{"enum", "-f", filepath.Join("testdata", "testsvc", "vo", "vo.go")}...)
	if err != nil {
		t.Fatal(err)
	}
}
