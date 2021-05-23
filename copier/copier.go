package copier

import (
	"encoding/json"
	"github.com/pkg/errors"
	"reflect"
)

func DeepCopy(src, target interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return errors.Wrap(err, "")
	}
	if reflect.ValueOf(target).Kind() != reflect.Ptr {
		return errors.New("Target should be a pointer")
	}
	return json.Unmarshal(b, target)
}
