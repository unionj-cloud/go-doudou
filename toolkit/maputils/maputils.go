package maputils

import (
	"github.com/goccy/go-reflect"
	"github.com/samber/lo"
	"strconv"

	"github.com/unionj-cloud/go-doudou/v2/toolkit/sliceutils"
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

var (
	MaxDepth = 32
)

// Merge recursively merges the src and dst maps. Key conflicts are resolved by
// preferring src, or recursively descending, if both src and dst are maps.
// borrow code from https://github.com/peterbourgon/mergemap
func Merge(dst, src map[string]interface{}) map[string]interface{} {
	return merge(dst, src, 0, false)
}

func MergeOverwriteSlice(dst, src map[string]interface{}) map[string]interface{} {
	return merge(dst, src, 0, true)
}

// overwrite means overwrite slice value
func merge(dst, src map[string]interface{}, depth int, overwrite bool) map[string]interface{} {
	if depth > MaxDepth {
		panic("too deep!")
	}
	for key, srcVal := range src {
		if dstVal, ok := dst[key]; ok {
			srcMap, srcMapOk := mapify(srcVal)
			dstMap, dstMapOk := mapify(dstVal)
			if srcMapOk && dstMapOk {
				srcVal = merge(dstMap, srcMap, depth+1, overwrite)
				goto REWRITE
			}
			if overwrite {
				goto REWRITE
			}
			srcSlice, srcSliceOk := sliceutils.TakeSliceArg(srcVal)
			dstSlice, dstSliceOk := sliceutils.TakeSliceArg(dstVal)
			if srcSliceOk && dstSliceOk {
				merged := make([]interface{}, 0)
				kv := make(map[interface{}]struct{})
				for _, item := range dstSlice {
					if !reflect.ValueOf(item).Type().Comparable() {
						merged = append(merged, item)
						continue
					}
					if _, exists := kv[item]; !exists {
						merged = append(merged, item)
						kv[item] = struct{}{}
					}
				}
				for _, item := range srcSlice {
					if !reflect.ValueOf(item).Type().Comparable() {
						merged = append(merged, item)
						continue
					}
					if _, exists := kv[item]; !exists {
						merged = append(merged, item)
						kv[item] = struct{}{}
					}
				}
				srcVal = merged
			}
		}
	REWRITE:
		dst[key] = srcVal
	}
	return dst
}

func mapify(i interface{}) (map[string]interface{}, bool) {
	value := reflect.ValueOf(i)
	if value.Kind() == reflect.Map {
		m := map[string]interface{}{}
		for _, k := range value.MapKeys() {
			m[k.String()] = value.MapIndex(k).Interface()
		}
		return m, true
	}
	return map[string]interface{}{}, false
}

func ConvertInt642String(data map[string]interface{}) {
	lo.ForEach(lo.Entries(data), func(item lo.Entry[string, interface{}], index int) {
		if item.Value != nil {
			data[item.Key] = convertInt642String(item.Value)
		}
	})
}

func convertInt642String(data interface{}) interface{} {
	switch v := data.(type) {
	case int64:
		return strconv.FormatInt(v, 10)
	case *int64:
		if v != nil {
			return lo.ToPtr(strconv.FormatInt(*v, 10))
		} else {
			return (*string)(nil)
		}
	case []int64:
		return lo.Map[int64, string](v, func(item int64, index int) string {
			return strconv.FormatInt(item, 10)
		})
	case []*int64:
		return lo.Map[*int64, *string](v, func(item *int64, index int) *string {
			if item != nil {
				return lo.ToPtr(strconv.FormatInt(*item, 10))
			}
			return (*string)(nil)
		})
	case []interface{}:
		return lo.Map[interface{}, interface{}](v, func(item interface{}, index int) interface{} {
			return convertInt642String(item)
		})
	}
	return data
}
