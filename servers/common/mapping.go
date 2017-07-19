package common

import (
	"fmt"
)

// UnMapFields takes a mapping, an object or slice of objects that have mapped fields and produces
// a new object or slice of objects where the mapped, custom fields are converted to the
// original columns
func UnMapFields(mapping map[string]string, mappedObj interface{}) error {
	if len(mapping) == 0 {
		return nil
	}
	switch obj := mappedObj.(type) {
	case map[string]interface{}:
		for field, value := range obj {
			if mappedColumn, ok := mapping[field]; ok {
				delete(obj, field)
				obj[mappedColumn] = value
			}
			if _value, isMap := value.(map[string]interface{}); isMap {
				UnMapFields(mapping, _value)
			}
		}
	case []map[string]interface{}:
		for _, m := range obj {
			UnMapFields(mapping, m)
		}
	default:
		return fmt.Errorf("invalid mappedObj type. Expects a map or slice of map")
	}
	return nil
}

// MapFields applies a mapping to a unmapped object
func MapFields(mapping map[string]string, unmappedObj interface{}) error {
	if len(mapping) == 0 {
		return nil
	}
	switch obj := unmappedObj.(type) {
	case map[string]interface{}:
		for field, value := range obj {
			for mField, mValue := range mapping {
				if field == mValue {
					delete(obj, field)
					obj[mField] = value
					break
				}
			}
			if _value, isMap := value.(map[string]interface{}); isMap {
				MapFields(mapping, _value)
			}
		}
	case []map[string]interface{}:
		for _, m := range obj {
			MapFields(mapping, m)
		}
	default:
		return fmt.Errorf("invalid mappedObj type. Expected a map or slice of map")
	}
	return nil
}
