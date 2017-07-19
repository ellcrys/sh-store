package servers

import (
	"testing"

	"github.com/ellcrys/elldb/servers/common"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValidation(t *testing.T) {
	Convey("Validation", t, func() {
		Convey(".validateMapping", func() {
			Convey("Should return error if empty mapping is provided", func() {
				errs := validateMapping(nil)
				So(errs, ShouldHaveLength, 1)
				So(errs[0], ShouldResemble, common.Error{
					Status:  "",
					Code:    "",
					Field:   "mapping",
					Message: "requires at least one custom field",
					Links:   nil,
				})
			})

			Convey("Should return error if custom field does not have a value (object field)", func() {
				errs := validateMapping(map[string]string{"name": ""})
				So(errs, ShouldHaveLength, 1)
				So(errs[0], ShouldResemble, common.Error{
					Status:  "",
					Code:    "",
					Field:   "mapping.name",
					Message: "field name is required",
					Links:   nil,
				})
			})

			Convey("Should return error if custom field points to an unknown value (object field)", func() {
				errs := validateMapping(map[string]string{"name": "unknown"})
				So(errs, ShouldHaveLength, 1)
				So(errs[0], ShouldResemble, common.Error{
					Status:  "",
					Code:    "",
					Field:   "mapping.name",
					Message: "field value 'unknown' is not a valid object field",
					Links:   nil,
				})
			})
		})
	})
}
