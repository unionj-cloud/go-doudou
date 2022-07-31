package codegen

import (
	"fmt"
	"path/filepath"
	"testing"
)

func Test_notGenerated(t *testing.T) {
	result := notGenerated(filepath.Join(testDir, "integrationtest"), filepath.Join(testDir, "testcode.postman_collection.json"))
	fmt.Println(result)
}
