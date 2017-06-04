package mapvalidator

import "fmt"
import "reflect"

// Required creates a new rule that requires a field
func Required(targetField string) Rule {
	msg := fmt.Sprintf(`field '%s' is required`, targetField)
	return RequiredWithMsg(targetField, msg)
}

// RequiredWithMsg creates a new rule that requires a field with a custom error message
func RequiredWithMsg(targetField, msg string) Rule {
	return NewRuleMaker(targetField, true, msg, msg, func(r *RuleMaker, fieldName string, fieldValue interface{}, fullMap map[string]interface{}) (errs []error) {
		switch _value := fieldValue.(type) {
		case string:
			if len(_value) == 0 {
				errs = append(errs, fmt.Errorf(r.message))
			}
		case uint, uint8, uint16, uint32, uint64:
			rv := reflect.ValueOf(_value)
			if rv.Uint() == 0 {
				errs = append(errs, fmt.Errorf(r.message))
			}
		case int, int8, int16, int32, int64:
			rv := reflect.ValueOf(_value)
			if rv.Int() == 0 {
				errs = append(errs, fmt.Errorf(r.message))
			}
		case float32, float64:
			rv := reflect.ValueOf(_value)
			if rv.Float() == 0 {
				errs = append(errs, fmt.Errorf(r.message))
			}
		case bool:
			if _value == false {
				errs = append(errs, fmt.Errorf(r.message))
			}
		default:
			rv := reflect.ValueOf(_value)
			if rv.Kind() == reflect.Map || rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
				if rv.Len() == 0 {
					if rv.Len() == 0 {
						errs = append(errs, fmt.Errorf(r.message))
					}
				}
			}
		}
		return
	})
}
