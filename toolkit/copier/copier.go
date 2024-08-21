package copier

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
	"github.com/pkg/errors"
)

var json = sonic.ConfigDefault

// DeepCopy src to target with json marshal and unmarshal
func DeepCopy(src, target interface{}) error {
	if src == nil || target == nil {
		return nil
	}
	if reflect.ValueOf(target).Kind() != reflect.Ptr {
		return errors.New("Target should be a pointer")
	}
	switch value := src.(type) {
	case *map[string]interface{}:
		if value == nil {
			return nil
		}
		return MapToStruct(*value, target)
	case map[string]interface{}:
		return MapToStruct(value, target)
	default:
		b, err := json.Marshal(src)
		if err != nil {
			return errors.WithStack(err)
		}
		dec := decoder.NewStreamDecoder(bytes.NewReader(b))
		dec.UseInt64()
		return dec.Decode(target)
	}
}

func MapToStruct(m map[string]any, structObj any) error {
	for k, v := range m {
		err := setStructField(structObj, k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func setStructField(structObj any, fieldName string, fieldValue any) error {
	structVal := reflect.ValueOf(structObj).Elem()

	fName := getFieldNameByJsonTag(structObj, fieldName)
	if fName == "" {
		return nil
	}

	fieldVal := structVal.FieldByName(fName)

	if !fieldVal.IsValid() {
		return fmt.Errorf("No such field: %s in obj", fieldName)
	}

	if !fieldVal.CanSet() {
		return fmt.Errorf("Cannot set %s field value", fieldName)
	}

	val := reflect.ValueOf(fieldValue)

	if fieldVal.Type() != val.Type() {

		if val.CanConvert(fieldVal.Type()) {
			fieldVal.Set(val.Convert(fieldVal.Type()))
			return nil
		}

		if m, ok := fieldValue.(map[string]any); ok {

			if fieldVal.Kind() == reflect.Struct {
				return MapToStruct(m, fieldVal.Addr().Interface())
			}

			if fieldVal.Kind() == reflect.Ptr && fieldVal.Type().Elem().Kind() == reflect.Struct {
				if fieldVal.IsNil() {
					fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
				}

				return MapToStruct(m, fieldVal.Interface())
			}

		}

		return fmt.Errorf("Map value type don't match struct field type")
	}

	fieldVal.Set(val)

	return nil
}

func getFieldNameByJsonTag(structObj any, jsonTag string) string {
	s := reflect.TypeOf(structObj).Elem()

	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		tag := field.Tag
		name, _, _ := strings.Cut(tag.Get("json"), ",")
		if name == jsonTag {
			return field.Name
		}
	}

	return ""
}
