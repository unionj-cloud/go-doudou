package copier

import (
	"bytes"
	"fmt"
	"github.com/spf13/cast"
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

	if !val.IsValid() {
		return nil
	}

	if fieldVal.Type() != val.Type() {

		if val.CanConvert(fieldVal.Type()) {
			fieldVal.Set(val.Convert(fieldVal.Type()))
			return nil
		}

		if fieldVal.Kind() == reflect.Ptr {
			v := reflect.New(fieldVal.Type().Elem())
			if fieldVal.Type().Elem() != val.Type() {
				if val.CanConvert(fieldVal.Type().Elem()) {
					v.Elem().Set(val.Convert(fieldVal.Type().Elem()))
					fieldVal.Set(v)
					return nil
				}
			} else {
				v.Elem().Set(val)
				fieldVal.Set(v)
				return nil
			}
		}

		if val.Kind() == reflect.Ptr {
			v := val.Elem()
			if fieldVal.Type() != v.Type() {
				if v.CanConvert(fieldVal.Type()) {
					fieldVal.Set(v.Convert(fieldVal.Type()))
					return nil
				}
			} else {
				fieldVal.Set(v)
				return nil
			}
		}

		if val.Kind() == reflect.String {
			switch fieldVal.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fieldVal.SetInt(cast.ToInt64(val.String()))
				return nil
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fieldVal.SetUint(cast.ToUint64(val.String()))
				return nil
			case reflect.Float32, reflect.Float64:
				fieldVal.SetFloat(cast.ToFloat64(val.String()))
				return nil
			case reflect.Ptr:
				v := reflect.New(fieldVal.Type().Elem())
				switch fieldVal.Type().Elem().Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					v.Elem().SetInt(cast.ToInt64(val.String()))
					fieldVal.Set(v)
					return nil
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					v.Elem().SetUint(cast.ToUint64(val.String()))
					fieldVal.Set(v)
					return nil
				case reflect.Float32, reflect.Float64:
					v.Elem().SetFloat(cast.ToFloat64(val.String()))
					fieldVal.Set(v)
					return nil
				}
			}
		} else if val.Kind() == reflect.Ptr && val.Type().Elem().Kind() == reflect.String {
			underlyingV := val.Elem().String()
			switch fieldVal.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fieldVal.SetInt(cast.ToInt64(underlyingV))
				return nil
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fieldVal.SetUint(cast.ToUint64(underlyingV))
				return nil
			case reflect.Float32, reflect.Float64:
				fieldVal.SetFloat(cast.ToFloat64(underlyingV))
				return nil
			case reflect.Ptr:
				v := reflect.New(fieldVal.Type().Elem())
				switch fieldVal.Type().Elem().Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					v.Elem().SetInt(cast.ToInt64(underlyingV))
					fieldVal.Set(v)
					return nil
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					v.Elem().SetUint(cast.ToUint64(underlyingV))
					fieldVal.Set(v)
					return nil
				case reflect.Float32, reflect.Float64:
					v.Elem().SetFloat(cast.ToFloat64(underlyingV))
					fieldVal.Set(v)
					return nil
				}
			}
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

		return fmt.Errorf("map attribute [%s] value type don't match struct field [%s] type", fieldName, fName)
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
