package codegen

import (
	"github.com/unionj-cloud/go-doudou/v2/toolkit/constants"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/fileutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

var goModFixGrpc = `
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2 v2.0.0-rc.2
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.0.0-rc.2
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
`

// FixModGrpc ...
func FixModGrpc(dir string) {
	var (
		err     error
		modfile string
	)
	modfile = filepath.Join(dir, "go.mod")
	lines, err := fileutils.File2lines(modfile)
	if err != nil {
		panic(err)
	}
	fileContent := ""
	reg := regexp.MustCompile(`require \(`)
	var found bool
	for _, line := range lines {
		fileContent += line
		fileContent += constants.LineBreak
		if reg.MatchString(line) && !found {
			found = true
			fileContent += goModFixGrpc
			fileContent += constants.LineBreak
		}
	}
	ioutil.WriteFile(modfile, []byte(fileContent), os.ModePerm)
}
