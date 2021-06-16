package codegen

import (
	"os"
)

func ExampleInitSvc() {
	dir := testDir + "initsvc"
	InitSvc(dir)
	defer os.RemoveAll(dir)
	// Output:
	// 1.15
	// testfilesinitsvc

}
