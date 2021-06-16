package copier

import (
	"encoding/json"
	"github.com/pkg/errors"
	"reflect"
)

func DeepCopy(src, target interface{}) error {
	b, _ := json.Marshal(src)
	if reflect.ValueOf(target).Kind() != reflect.Ptr {
		return errors.New("Target should be a pointer")
	}
	return json.Unmarshal(b, target)
}
