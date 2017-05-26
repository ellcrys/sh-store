package common

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test(t *testing.T) {
	Convey("APIErrors", t, func() {
		Convey(".NewSingleAPIErr", func() {
			Convey("Should create and return a single APIError encoded in standard error", func() {
				e := NewSingleAPIErr(200, "some_code", "some_field", "something happened", map[string]string{"about": "http://google.com"})
				So(e, ShouldHaveSameTypeAs, errors.New(""))
				So(e.Error(), ShouldEqual, `{"Message":"something happened","Errors":{"errors":[{"status":"200","code":"some_code","field":"some_field","message":"something happened","links":{"about":"http://google.com"}}]},"StatusCode":200}`)
			})

		})
		Convey(".MultiAPIError", func() {
			Convey("Should create and return multiple APIError encoded in standard error", func() {
				e := NewMultiAPIErr(200, "something happened", []Error{
					{Status: "200", Code: "some_code", Message: "error 1"},
					{Status: "200", Code: "some_code_2", Message: "error 2"},
				})
				So(e, ShouldHaveSameTypeAs, errors.New(""))
				So(e.Error(), ShouldEqual, `{"Message":"something happened","Errors":{"errors":[{"status":"200","code":"some_code","message":"error 1"},{"status":"200","code":"some_code_2","message":"error 2"}]},"StatusCode":200}`)
			})
		})
	})
}
