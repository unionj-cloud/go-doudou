package maputils

import (
	"reflect"
)

type ChangeType int

const (
	ADDED ChangeType = iota
	MODIFIED
	DELETED
)

type Change struct {
	OldValue   interface{}
	NewValue   interface{}
	ChangeType ChangeType
}

func Diff(new, old map[string]interface{}) map[string]Change {
	mp := map[string]bool{}
	for k, _ := range old {
		mp[k] = true
	}
	changes := make(map[string]Change)

	if new != nil {
		for key, value := range new {
			//key state insert or update
			//insert
			if !mp[key] {
				changes[key] = createAddChange(value)
			} else {
				//update
				oldValue := old[key]
				if !reflect.DeepEqual(oldValue, value) {
					changes[key] = createModifyChange(oldValue, value)
				}
			}
			delete(mp, key)
		}
	}

	// remove del keys
	for key := range mp {
		//get old value and del
		oldValue := old[key]
		changes[key] = createDeletedChange(oldValue)
	}

	return changes
}

func createModifyChange(oldValue interface{}, newValue interface{}) Change {
	return Change{
		OldValue:   oldValue,
		NewValue:   newValue,
		ChangeType: MODIFIED,
	}
}

func createAddChange(newValue interface{}) Change {
	return Change{
		NewValue:   newValue,
		ChangeType: ADDED,
	}
}

func createDeletedChange(oldValue interface{}) Change {
	return Change{
		OldValue:   oldValue,
		ChangeType: DELETED,
	}
}
