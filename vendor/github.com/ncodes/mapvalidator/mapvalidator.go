package mapvalidator

import (
	"fmt"
)

// Rule describes a validation check
type Rule interface {

	// TargetMustExist is used to override the default behaviour
	// where this rule is ignored if the target field is not
	// included in the map.
	TargetMustExist() bool

	// GetTargetMustExistMessage returns the message to return when the target field does not exist
	GetTargetMustExistMessage() string

	// GetTargetField returns the name of the field to validate
	GetTargetField() string

	// Check performs the verification according to the rule
	// and return every thing wrong with the field
	Check(fieldName string, fieldValue interface{}, fullMap map[string]interface{}) []error
}

// Validate performs validation on a map. It uses a list of one or more rules.
// Returns the errors from rules that failed. Rules targetting fields that do not
// exists in the map are ignored
func Validate(m map[string]interface{}, rules ...Rule) (errs []error) {
	for _, rule := range rules {
		fieldVal, hasField := m[rule.GetTargetField()]
		if hasField {
			_errs := rule.Check(rule.GetTargetField(), fieldVal, m)
			if _errs != nil {
				errs = append(errs, _errs...)
			}
		} else {
			if rule.TargetMustExist() {
				msg := rule.GetTargetMustExistMessage()
				if msg == "" {
					msg = "field '" + rule.GetTargetField() + "' is required"
				}
				errs = append(errs, fmt.Errorf(msg))
			}
		}
	}
	return
}
