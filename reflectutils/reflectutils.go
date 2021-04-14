package reflectutils

import "reflect"

func ValueOf(data interface{}) reflect.Value {
	if reflect.ValueOf(data).Kind() == reflect.Ptr {
		return reflect.ValueOf(data).Elem()
	} else {
		return reflect.ValueOf(data)
	}
}

func ValueOfValue(data reflect.Value) reflect.Value {
	if data.Kind() == reflect.Ptr {
		return data.Elem()
	} else {
		return data
	}
}
