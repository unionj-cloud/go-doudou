package main

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/toolkit/lab/vo"
	"reflect"
)

var TypeRegistry = make(map[string]reflect.Type)

func init() {
	myTypes := []interface{}{vo.PageQuery{}, UsersvcImpl{}}
	for _, v := range myTypes {
		TypeRegistry[fmt.Sprintf("%T", v)] = reflect.TypeOf(v)
	}
	a := TypeRegistry
	a = a
}
