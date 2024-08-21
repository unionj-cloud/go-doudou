package copier

import (
	"bytes"
	"github.com/goccy/go-reflect"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
	"github.com/duke-git/lancet/v2/maputil"
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
		return maputil.MapToStruct(*value, target)
	case map[string]interface{}:
		return maputil.MapToStruct(value, target)
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
