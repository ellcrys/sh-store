package common

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/ellcrys/jsonapi"
	"github.com/ellcrys/util"
	"google.golang.org/grpc"
)

// EasyHandle takes a custom handle that returns any value type and a status code.
// An error response value whose message is a JSON string is converted to ErrorsPayload and
// returned as JSON. The status code provided in ErrorsPayload is used
// as the response status code if provided, otherwise the status returned by the handle is used.
// Strings are returned as text.
// None errors are coerced and returned as JSON.
func EasyHandle(method string, handler func(w http.ResponseWriter, r *http.Request) (interface{}, int)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// custom method handling
		if r.Method != method {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "404",
				"code":    "unsupported_method",
				"message": "method not supported for this endpoint",
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

			var errs ErrorsPayload
			if err := util.FromJSON([]byte(desc), &errs); err != nil {
				fmt.Fprint(w, []byte("failed to decode error to JSON"))
				return
			}

			if errs.Status != 0 {
				status = errs.Status
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(errs.ErrorsPayload)
			return

		case string:
			fmt.Fprint(w, b)
		case map[string]string, map[string]interface{}:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(b)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)

			if err := jsonapi.MarshalPayload(w, b); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
