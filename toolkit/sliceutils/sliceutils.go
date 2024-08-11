package sliceutils

import (
	"github.com/goccy/go-reflect"

	"github.com/pkg/errors"
)

// StringSlice2InterfaceSlice converts string slice to interface slice
func StringSlice2InterfaceSlice(strSlice []string) []interface{} {
	ret := make([]interface{}, len(strSlice))
	for i, v := range strSlice {
		ret[i] = v
	}
	return ret
}

// InterfaceSlice2StringSlice converts interface slice to string slice
func InterfaceSlice2StringSlice(strSlice []interface{}) []string {
	ret := make([]string, len(strSlice))
	for i, v := range strSlice {
		ret[i] = v.(string)
	}
	return ret
}

// Contains asserts src contains test
func Contains(src []interface{}, test interface{}) bool {
	for _, item := range src {
		if item == test {
			return true
		}
	}
	return false
}

// ContainsDeep asserts src contains test using reflect.DeepEqual
func ContainsDeep(src []interface{}, test interface{}) bool {
	for _, item := range src {
		if reflect.DeepEqual(item, test) {
			return true
		}
	}
	return false
}

// StringContains asserts src contains test
func StringContains(src []string, test string) bool {
	for _, item := range src {
		if item == test {
			return true
		}
	}
	return false
}

// StringFilter filters string slice by callback function fn
// If fn returns true, the item will be appended to result
func StringFilter(src []string, fn func(item string) bool) []string {
	var ret []string
	for _, item := range src {
		if fn(item) {
			ret = append(ret, item)
		}
	}
	return ret
}

// IndexOf returns index of element in string slice data
func IndexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

// IndexOfAny returns index of element in slice data
func IndexOfAny(target interface{}, anySlice interface{}) (int, error) {
	if reflect.TypeOf(anySlice).Kind() != reflect.Slice {
		return -1, errors.New("not slice")
	}
	data := reflect.ValueOf(anySlice)
	for i := 0; i < data.Len(); i++ {
		elem := data.Index(i)
		if elem.Interface() == target {
			return i, nil
		}
	}
	return -1, nil //not found.
}

// IsEmpty assert src is an empty slice
func IsEmpty(src interface{}) bool {
	if slice, ok := TakeSliceArg(src); ok {
		return slice == nil || len(slice) == 0
	}
	panic("not slice")
}

// TakeSliceArg https://ahmet.im/blog/golang-take-slices-of-any-type-as-input-parameter/
func TakeSliceArg(arg interface{}) (out []interface{}, ok bool) {
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

// ConvertAny2Interface converts interface src to interface slice
func ConvertAny2Interface(src interface{}) ([]interface{}, error) {
	data := reflect.ValueOf(src)
	if data.Type().Kind() == reflect.Ptr {
		data = data.Elem()
	}
	if data.Type().Kind() != reflect.Slice {
		return nil, errors.New("Src not slice")
	}
	ret := make([]interface{}, data.Len())
	for i := 0; i < data.Len(); i++ {
		ret[i] = data.Index(i).Interface()
	}
	return ret, nil
}
