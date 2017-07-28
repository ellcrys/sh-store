package servers

import (
	"testing"

	"github.com/ellcrys/elldb/servers/common"
	"github.com/ellcrys/elldb/servers/db"
	"github.com/fatih/structs"
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

func TestValidations(t *testing.T) {
	setup(t, func(rpc, rpc2 *RPC) {
		Convey("Validation", t, func() {
			Convey(".validateObjects", func() {
				Convey("Should return error if no object is provided", func() {
					errs, err := rpc.validateObjects(nil, nil)
					So(err, ShouldBeNil)
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, common.Error{
						Status:  "",
						Code:    "",
						Field:   "",
						Message: "no object provided. At least one object is required",
						Links:   nil,
					})
				})

				Convey("Should return error if objects length exceeds max", func() {
					curMaxObjPerPut := MaxObjectPerPut
					MaxObjectPerPut = 2
					var objs = []map[string]interface{}{}
					objs = append(objs, structs.New(db.Object{}).Map(), structs.New(db.Object{}).Map(), structs.New(db.Object{}).Map())
					errs, err := rpc.validateObjects(objs, nil)
					So(err, ShouldBeNil)
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, common.Error{
						Status:  "",
						Code:    "",
						Field:   "",
						Message: "too many objects. Only a maximum of 2 can be created at once",
						Links:   nil,
					})
					MaxObjectPerPut = curMaxObjPerPut
				})

				Convey("Should return error if object required fields are missing (validation)", func() {
					var objs = []map[string]interface{}{}
					objs = append(objs, map[string]interface{}{})
					errs, err := rpc.validateObjects(objs, nil)
					So(err, ShouldBeNil)
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: owner_id is required", Links: nil})
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: creator_id is required", Links: nil})
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: key is required", Links: nil})
				})

				Convey("Should return error if object required fields have invalid type (validation)", func() {
					var objs = []map[string]interface{}{}
					objs = append(objs, map[string]interface{}{
						"owner_id":   123,
						"creator_id": 123,
						"key":        123,
						"value":      123,
						"ref1":       123,
					})
					errs, err := rpc.validateObjects(objs, nil)
					So(err, ShouldBeNil)
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: owner_id must be a string", Links: nil})
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: creator_id must be a string", Links: nil})
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: key must be a string", Links: nil})
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: value must be a string", Links: nil})
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: ref1 must be a string", Links: nil})
				})

				Convey("Should return error if objects do not share the same owner_id", func() {
					var objs = []map[string]interface{}{}
					objs = append(objs, map[string]interface{}{"owner_id": "xyz", "creator_id": "abc", "key": "some_key"})
					objs = append(objs, map[string]interface{}{"owner_id": "xyz2", "creator_id": "abc", "key": "some_key"})
					errs, err := rpc.validateObjects(objs, nil)
					So(err, ShouldBeNil)
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 1: all objects must share the same 'owner_id' (owner_id) value", Links: nil})
				})

				Convey("Should return errors that use custom field names if mapping is provided", func() {
					var objs = []map[string]interface{}{}
					objs = append(objs, map[string]interface{}{})
					errs, err := rpc.validateObjects(objs, map[string]string{
						"user_id": "owner_id",
						"app_id":  "creator_id",
						"name":    "key",
					})
					So(err, ShouldBeNil)
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: user_id is required", Links: nil})
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: app_id is required", Links: nil})
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: name is required", Links: nil})
				})

				Convey("Should return error if objects' owner does not exists", func() {
					var objs = []map[string]interface{}{}
					objs = append(objs, map[string]interface{}{"owner_id": "unknown", "creator_id": "abc", "key": "mykey"})
					errs, err := rpc.validateObjects(objs, nil)
					So(err, ShouldBeNil)
					So(errs, ShouldContain, common.Error{
						Status:  "",
						Code:    "invalid_parameter",
						Field:   "",
						Message: "owner of object(s) does not exist",
						Links:   nil,
					})
				})
			})
		})
	})
}
