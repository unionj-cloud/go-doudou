package reflectutils

import "github.com/goccy/go-reflect"

// ValueOf returns underlying value of interface data
func ValueOf(data interface{}) reflect.Value {
	return ValueOfValue(reflect.ValueOf(data))
}

// ValueOfValue returns underlying value of reflect.Value data
func ValueOfValue(data reflect.Value) reflect.Value {
	if data.Kind() == reflect.Ptr {
		return data.Elem()
	}
	return data
}
