package client

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGenGoClient(t *testing.T) {
	dir := "testdata/testclient"
	defer os.RemoveAll(dir)
	assert.NotPanics(t, func() {
		GenGoClient(dir, "../testdata/petstore3.json", true, "", "client")
	})
}

func TestGenGoClient2(t *testing.T) {
	dir := "testdata/testclient2"
	defer os.RemoveAll(dir)
	assert.NotPanics(t, func() {
		GenGoClient(dir, "../testdata/swagger.json", true, "", "client")
	})
}
