package common

import (
	"fmt"
)

// UnMapFields takes a mapping, an object or slice of objects that have mapped fields and produces
// a new object or slice of objects where the mapped, custom fields are converted to the
// original columns
func UnMapFields(mapping map[string]string, mappedObj interface{}) error {
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
