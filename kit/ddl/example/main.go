package main

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/kit/astutils"
	"github.com/unionj-cloud/go-doudou/kit/reflectutils"
	"reflect"
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

	var a interface{}
	a = 10
	var b interface{}
	b = a

	fmt.Println(reflect.ValueOf(b).Kind() == reflect.Ptr)

	fmt.Printf("%v\n", reflectutils.ValueOf(b))

	var d int
	var e int
	d = 10
	e = 5
	var c interface{}
	c = []*int{&d, &e}

	data := reflect.ValueOf(c)
	for i := 0; i < data.Len(); i++ {
		fmt.Printf("%v\n", reflectutils.ValueOfValue(data.Index(i)))
	}

}
