package mapvalidator

import "fmt"
import "github.com/asaskevich/govalidator"
import "reflect"

// TypeChecker represents a function for testing a value type
type TypeChecker func(value interface{}) bool

// TypeString checks whether a value is a string type
func TypeString(val interface{}) bool {
	if _, ok := val.(string); ok {
		return true
	}
	return false
}

// TypeNumber checks whether a value is an int or a float
func TypeNumber(val interface{}) bool {
	switch val.(type) {
	case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32:
		return true
	default:
		return false
	}
}

// TypeDigit check whether a value is a digit. String with
// only numeric chars will pass.
func TypeDigit(val interface{}) bool {
	switch v := val.(type) {
	case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32:
		return true
	case string:
		return govalidator.IsNumeric(v)
	default:
		return false
	}
}

// TypeMap checks whether a value is a map
func TypeMap(val interface{}) bool {
	rv := reflect.ValueOf(val)
	return rv.Kind() == reflect.Map
}

// TypeWithMsg creates a new rule that checks if the target field value
// has a specific type
func TypeWithMsg(targetField string, typeChecker TypeChecker, msg string) Rule {
	missingFieldMsg := fmt.Sprintf(`field '%s' is required`, targetField)
	return NewRuleMaker(targetField, false, msg, missingFieldMsg, func(r *RuleMaker, fieldName string, fieldValue interface{}, fullMap map[string]interface{}) (errs []error) {
		if typeChecker != nil {
			if !typeChecker(fieldValue) {
				errs = append(errs, fmt.Errorf(msg))
			}
		}
		return
	})
}
