package cmd_test

import (
	"os"
	"testing"

	"github.com/unionj-cloud/go-doudou/v2/cmd"
)

func TestClientCmd(t *testing.T) {
	defer os.RemoveAll("client")
	// go-doudou svc http client --file testdata/testsvc/testsvc_openapi3.json
	_, _, err := ExecuteCommandC(cmd.GetRootCmd(), []string{"svc", "http", "client", "--file", "testdata/testsvc/testsvc_openapi3.json"}...)
	if err != nil {
		t.Error(err)
		return
	}
}
