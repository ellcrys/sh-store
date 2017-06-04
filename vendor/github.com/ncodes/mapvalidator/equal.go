package mapvalidator

import "fmt"

// EqualStringWithMsg creates a new rule that checks if a field equals a string value.
// If the value type is not a string, error message is returned
func EqualStringWithMsg(targetField string, cmpVal string, msg string) Rule {
	missingFieldMsg := fmt.Sprintf(`field '%s' is required`, targetField)
	return NewRuleMaker(targetField, false, msg, missingFieldMsg, func(r *RuleMaker, fieldName string, fieldValue interface{}, fullMap map[string]interface{}) (errs []error) {
		if v, ok := fieldValue.(string); ok {
			if v != cmpVal {
				errs = append(errs, fmt.Errorf(msg))
			}
		} else {
			errs = append(errs, fmt.Errorf(msg))
		}
		return
	})
}

// EqualIntWithMsg creates a new rule that checks if a field equals an int value.
// If the value type is not a int, error message is returned
func EqualIntWithMsg(targetField string, cmpVal int, msg string) Rule {
	missingFieldMsg := fmt.Sprintf(`field '%s' is required`, targetField)
	return NewRuleMaker(targetField, false, msg, missingFieldMsg, func(r *RuleMaker, fieldName string, fieldValue interface{}, fullMap map[string]interface{}) (errs []error) {
		if v, ok := fieldValue.(int); ok {
			if v != cmpVal {
				errs = append(errs, fmt.Errorf(msg))
			}
		} else {
			errs = append(errs, fmt.Errorf(msg))
		}
		return
	})
}

// EqualInt32WithMsg creates a new rule that checks if a field equals an int32 value.
// If the value type is not a int32, error message is returned
func EqualInt32WithMsg(targetField string, cmpVal int32, msg string) Rule {
	missingFieldMsg := fmt.Sprintf(`field '%s' is required`, targetField)
	return NewRuleMaker(targetField, false, msg, missingFieldMsg, func(r *RuleMaker, fieldName string, fieldValue interface{}, fullMap map[string]interface{}) (errs []error) {
		if v, ok := fieldValue.(int32); ok {
			if v != cmpVal {
				errs = append(errs, fmt.Errorf(msg))
			}
		} else {
			errs = append(errs, fmt.Errorf(msg))
		}
		return
	})
}

// EqualInt64WithMsg creates a new rule that checks if a field equals an int64 value.
// If the value type is not a int64, error message is returned
func EqualInt64WithMsg(targetField string, cmpVal int64, msg string) Rule {
	missingFieldMsg := fmt.Sprintf(`field '%s' is required`, targetField)
	return NewRuleMaker(targetField, false, msg, missingFieldMsg, func(r *RuleMaker, fieldName string, fieldValue interface{}, fullMap map[string]interface{}) (errs []error) {
		if v, ok := fieldValue.(int64); ok {
			if v != cmpVal {
				errs = append(errs, fmt.Errorf(msg))
			}
		} else {
			errs = append(errs, fmt.Errorf(msg))
		}
		return
	})
}

// EqualFloat32WithMsg creates a new rule that checks if a field equals an float32 value.
// If the value type is not a float32, error message is returned
func EqualFloat32WithMsg(targetField string, cmpVal float32, msg string) Rule {
	missingFieldMsg := fmt.Sprintf(`field '%s' is required`, targetField)
	return NewRuleMaker(targetField, false, msg, missingFieldMsg, func(r *RuleMaker, fieldName string, fieldValue interface{}, fullMap map[string]interface{}) (errs []error) {
		if v, ok := fieldValue.(float32); ok {
			if v != cmpVal {
				errs = append(errs, fmt.Errorf(msg))
			}
		} else {
			errs = append(errs, fmt.Errorf(msg))
		}
		return
	})
}

// EqualFloat64WithMsg creates a new rule that checks if a field equals an float32 value.
// If the value type is not a float32, error message is returned
func EqualFloat64WithMsg(targetField string, cmpVal float64, msg string) Rule {
	missingFieldMsg := fmt.Sprintf(`field '%s' is required`, targetField)
	return NewRuleMaker(targetField, false, msg, missingFieldMsg, func(r *RuleMaker, fieldName string, fieldValue interface{}, fullMap map[string]interface{}) (errs []error) {
		if v, ok := fieldValue.(float64); ok {
			if v != cmpVal {
				errs = append(errs, fmt.Errorf(msg))
			}
		} else {
			errs = append(errs, fmt.Errorf(msg))
		}
		return
	})
}
