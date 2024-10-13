package codegen

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/unionj-cloud/toolkit/fileutils"
)

func Test_notGenerated(t *testing.T) {
	_ = fileutils.CreateDirectory(filepath.Join(testDir, "integrationtest"))
	result := notGenerated(filepath.Join(testDir, "integrationtest"), filepath.Join(testDir, "testcode.postman_collection.json"))
	fmt.Println(result)
}
