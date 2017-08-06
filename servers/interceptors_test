package servers

import (
	"context"
	"testing"

	"github.com/ellcrys/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInterceptor(t *testing.T) {
	setup(t, func(rpc, rpc2 *RPC) {
		Convey("Interceptor", t, func() {
			Convey("Should return error if account does not exists", func() {
				_, err := rpc.processAppTokenClaims(context.Background(), map[string]interface{}{
					"id": "unknown",
				})
				m, err := util.JSONToMap(err.Error())
				So(err, ShouldBeNil)
				errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
				So(errs, ShouldHaveLength, 1)
				So(errs[0], ShouldResemble, map[string]interface{}{
					"status":  "401",
					"code":    "auth_error",
					"message": "permission denied. Invalid token",
				})
			})
		})
	})
}
