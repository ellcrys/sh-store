package mapvalidator

import "fmt"

// Custom creates a new rule that calls a custom function to perform
// validation logic not supported by existing rules.
func Custom(targetField string, msg string, customFunc func(fieldValue interface{}, fullMap map[string]interface{}) bool) Rule {
	missingFieldMsg := fmt.Sprintf(`field '%s' is required`, targetField)
	return NewRuleMaker(targetField, false, msg, missingFieldMsg, func(r *RuleMaker, fieldName string, fieldValue interface{}, fullMap map[string]interface{}) (errs []error) {
		if !customFunc(fieldValue, fullMap) {
			errs = append(errs, fmt.Errorf(msg))
		}
		return
	})
}
