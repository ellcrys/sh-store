package common

import "github.com/ellcrys/jsonapi"
import "fmt"
import "github.com/ellcrys/util"

var (
	// ServerError returns a jsonapi error indicating a server error
	ServerError = Error(500, "server_error", "", "server error")

	// BodyMalformedError returns a jsonapi error indicating malformed json body error
	BodyMalformedError = Error(400, "invalid_parameter", "", "malformed json body")
)

// ErrorObject represents jsonapi error object
type ErrorObject struct {
	*jsonapi.ErrorObject
}

// ToJSON returns json representation of object
func (e *ErrorObject) ToJSON() string {
	data, _ := util.ToJSON(e)
	return string(data)
}

// ErrorsPayload represents jsonapi errors payload
type ErrorsPayload struct {
	*jsonapi.ErrorsPayload
	Status int `json:"status"`
}

// ToJSON returns json representation of object
func (e *ErrorsPayload) ToJSON() string {
	data, _ := util.ToJSON(e)
	return string(data)
}

// Error creates an error.
// The jsonapi payload is converted to json and passed to fmt.Error
func Error(status int, code, source, detail string) error {
	return fmt.Errorf((&ErrorsPayload{
		Status: status,
		ErrorsPayload: &jsonapi.ErrorsPayload{
			Errors: []*jsonapi.ErrorObject{{
				Status: fmt.Sprintf("%d", status),
				Code:   code,
				Detail: detail,
				Source: source,
			}},
		},
	}).ToJSON())
}

// Errors creates an error out of many jsonapi.ErrorObjects
// The jsonapi payload is converted to json and passed to fmt.Error
func Errors(status int, errs []*jsonapi.ErrorObject) error {
	return fmt.Errorf((&ErrorsPayload{
		Status: status,
		ErrorsPayload: &jsonapi.ErrorsPayload{
			Errors: errs,
		},
	}).ToJSON())
}
