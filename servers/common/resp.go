package common

import (
	"fmt"

	"github.com/ellcrys/util"
)

var (
	// ServerError describe a server error
	ServerError = NewSingleAPIErr(500, "", "server_error", "server error", nil)
	// BodyMalformedError describes an error about a malformed json body
	BodyMalformedError = NewSingleAPIErr(400, "", "invalid_body", "malformed json body", nil)
)

// Error defines a single error
type Error struct {
	Status  string            `json:"status,omitempty" structs:"status,omitempty" mapstructure:"status,omitempty"`
	Code    string            `json:"code,omitempty" structs:"code,omitempty" mapstructure:"code,omitempty"`
	Field   string            `json:"field,omitempty" structs:"field,omitempty" mapstructure:"field,omitempty"`
	Message string            `json:"message,omitempty" structs:"message,omitempty" mapstructure:"message,omitempty"`
	Links   map[string]string `json:"links,omitempty" structs:"links,omitempty" mapstructure:"links,omitempty"`
}

// Errors defines a collection of errors
type Errors struct {
	Errors []Error `json:"errors,omitempty" structs:"errors,omitempty" mapstructure:"errors,omitempty"`
}

// APIError describes an error from a server endpoint
type APIError struct {
	Message    string
	Errors     Errors
	StatusCode int
}

// NewAPIError creates an APIError that is encoded to JSON and returned as a standard error.
func NewAPIError(status int, msg string, errors Errors) error {
	errJSON := util.MustStringify(APIError{Message: msg, Errors: errors, StatusCode: status})
	return fmt.Errorf("%s", errJSON)
}

func (e *APIError) Error() string {
	return e.Message
}

// NewSingleAPIErr returns a response with a single error
func NewSingleAPIErr(status int, code, field, message string, links map[string]string) error {
	errors := Errors{
		Errors: []Error{{
			Status:  fmt.Sprintf("%d", status),
			Code:    code,
			Field:   field,
			Message: message,
			Links:   links,
		}},
	}
	return NewAPIError(status, message, errors)
}

// NewMultiAPIErr returns a response with multiple errors
func NewMultiAPIErr(status int, message string, errs []Error) error {
	errors := Errors{
		Errors: errs,
	}
	return NewAPIError(status, message, errors)
}

// SingleObjectResp creates standard response format for only a single object
func SingleObjectResp(objectType string, attr map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"type":       objectType,
			"id":         attr["id"],
			"attributes": attr,
		},
	}
}

// MultiObjectResp creates standard response format for multiple objects
func MultiObjectResp(objectType string, attrs []map[string]interface{}) map[string]interface{} {
	var data = make([]map[string]interface{}, 0)
	for _, attr := range attrs {
		data = append(data, map[string]interface{}{
			"type":       objectType,
			"id":         attr["id"],
			"attributes": attr,
		})
	}
	return map[string]interface{}{
		"data": data,
	}
}
