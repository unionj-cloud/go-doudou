package sliceutils

import "reflect"

func StringSlice2InterfaceSlice(strSlice []string) []interface{} {
	ret := make([]interface{}, len(strSlice))
	for i, v := range strSlice {
		ret[i] = v
	}
	return ret
}

func InterfaceSlice2StringSlice(strSlice []interface{}) []string {
	ret := make([]string, len(strSlice))
	for i, v := range strSlice {
		ret[i] = v.(string)
	}
	return ret
}

func Contains(src []interface{}, test interface{}) bool {
	for _, item := range src {
		if item == test {
			return true
		}
	}
	return false
}

func StringContains(src []string, test string) bool {
	for _, item := range src {
		if item == test {
			return true
		}
	}
	return false
}

func IndexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func IsEmpty(src interface{}) bool {
	if slice, ok := takeSliceArg(src); ok {
		return slice == nil || len(slice) == 0
	}
	panic("not slice")
}

// https://ahmet.im/blog/golang-take-slices-of-any-type-as-input-parameter/
func takeSliceArg(arg interface{}) (out []interface{}, ok bool) {
	slice, success := takeArg(arg, reflect.Slice)
	if !success {
		ok = false
		return
	}
	c := slice.Len()
	out = make([]interface{}, c)
	for i := 0; i < c; i++ {
		out[i] = slice.Index(i).Interface()
	}
	return out, true
}

func takeArg(arg interface{}, kind reflect.Kind) (val reflect.Value, ok bool) {
	val = reflect.ValueOf(arg)
	if val.Kind() == kind {
		ok = true
	}
	return
}
