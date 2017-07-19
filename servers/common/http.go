package common

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/ellcrys/util"
	"google.golang.org/grpc"
)

// EasyHandle takes a custom handle that returns any value type and a status code.
// An error response value whose message is a JSON string is considered to be a decoded common.APIError
// and will be encoded back to APIError and returned as JSON. The status code provided in APIError is used
// as the response status code. Other errors are returned as text.
// None errors are coerced and returned as JSON.
func EasyHandle(method string, handler func(w http.ResponseWriter, r *http.Request) (interface{}, int)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// custom method handling
		if r.Method != method {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(Error{
				Status:  "404",
				Code:    "unsupported_method",
				Message: "method not supported for this endpoint",
			})
			return
		}

		// call handler
		body, status := handler(w, r)

		// send response
		switch b := body.(type) {
		case error:
			desc := grpc.ErrorDesc(b)

			// Error does not contain a json string
			if !govalidator.IsJSON(desc) {
				w.WriteHeader(status)
				fmt.Fprintf(w, "%s", desc)
				return
			}

			// json must be an APIError
			var apiError APIError
			if err := util.FromJSON([]byte(desc), &apiError); err != nil {
				fmt.Fprint(w, []byte("failed to decode error to JSON"))
				return
			}

			if apiError.StatusCode != 0 {
				status = apiError.StatusCode
				apiError.StatusCode = 0
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(apiError.Errors)
			return

		case string:
			fmt.Fprint(w, b)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(b)
		}
	}
}
