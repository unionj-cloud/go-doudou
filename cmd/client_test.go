package cmd

import (
	"os"
	"testing"
)

func TestClientCmd(t *testing.T) {
	defer os.RemoveAll("client")
	// go-doudou svc http client --file testdata/testsvc/testsvc_openapi3.json
	_, _, err := ExecuteCommandC(rootCmd, []string{"svc", "http", "client", "--file", "testdata/testsvc/testsvc_openapi3.json"}...)
	if err != nil {
		t.Error(err)
		return
	}
}
