package codegen

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/fileutils"
	"path/filepath"
	"testing"
)

func Test_notGenerated(t *testing.T) {
	_ = fileutils.CreateDirectory(filepath.Join(testDir, "integrationtest"))
	result := notGenerated(filepath.Join(testDir, "integrationtest"), filepath.Join(testDir, "testcode.postman_collection.json"))
	fmt.Println(result)
}
