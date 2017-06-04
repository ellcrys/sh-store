package mapvalidator

// RuleMakerCheckFunc is the function to pass to the rule maker which will be called in
// the new rule's check method.
type RuleMakerCheckFunc func(r *RuleMaker, fieldName string, fieldValue interface{}, fullMap map[string]interface{}) []error

// RuleMaker describes a rule for a RuleMaker target field. It satisfies the Rule interface
type RuleMaker struct {
	targetField            string
	message                string
	targetMustExistMessage string
	checkFunc              RuleMakerCheckFunc
	targetRequired         bool
}

// NewRuleMaker creates a new required rule
func NewRuleMaker(targetField string, targetRequired bool, message, targetMustExistMessage string, checkFunc RuleMakerCheckFunc) *RuleMaker {
	return &RuleMaker{
		targetField:            targetField,
		targetRequired:         targetRequired,
		message:                message,
		targetMustExistMessage: targetMustExistMessage,
		checkFunc:              checkFunc,
	}
}

// GetTargetField returns the name of the field to perform validation against
func (r *RuleMaker) GetTargetField() string {
	return r.targetField
}

// TargetMustExist forces an error to returned if target field does not exist in map
func (r *RuleMaker) TargetMustExist() bool {
	return r.targetRequired
}

// GetTargetMustExistMessage returns the error message to return when target field does not exists in map
func (r *RuleMaker) GetTargetMustExistMessage() string {
	return r.targetMustExistMessage
}

// Check performs all validations against the field value and returns every validation error failed
func (r *RuleMaker) Check(fieldName string, fieldValue interface{}, fullMap map[string]interface{}) []error {
	return r.checkFunc(r, fieldName, fieldValue, fullMap)
}
