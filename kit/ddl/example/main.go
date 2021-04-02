package main

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/kit/astutils"
)

func init() {

}

func main() {
	fmt.Println(astutils.GetMod())
	fmt.Println(astutils.GetImportPath("/Users/wubin1989/workspace/cloud/go-doudou/kit/ddl/example/domain"))
}
