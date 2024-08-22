package copier

import (
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
	"github.com/pkg/errors"
	"reflect"
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
	b, err := json.MarshalToString(src)
	if err != nil {
		return errors.WithStack(err)
	}
	dec := decoder.NewDecoder(b)
	dec.UseInt64()
	if err = dec.Decode(target); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
