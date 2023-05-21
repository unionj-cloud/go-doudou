package copier

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/errorx"
	"reflect"
)

// DeepCopy src to target with json marshal and unmarshal
func DeepCopy(src, target interface{}) error {
	if src == nil || target == nil {
		return nil
	}
	b, err := json.Marshal(src)
	if err != nil {
		return errorx.Wrap(err)
	}
	if reflect.ValueOf(target).Kind() != reflect.Ptr {
		return errors.New("Target should be a pointer")
	}
	return json.Unmarshal(b, target)
}
