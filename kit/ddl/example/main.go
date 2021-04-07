package main

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/kit/astutils"
	"regexp"
)

func init() {

}

func main() {
	fmt.Println(astutils.GetMod())
	fmt.Println(astutils.GetImportPath("/Users/wubin1989/workspace/cloud/go-doudou/kit/ddl/example/domain"))

	re, err := regexp.Compile(`^\(.+\)$`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", re.MatchString("()"))
}
