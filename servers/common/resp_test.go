package common

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test(t *testing.T) {
	Convey("APIErrors", t, func() {

		Convey(".SingleObjectResp", func() {
			attr := map[string]interface{}{
				"id":   "abc",
				"name": "John",
				"age":  100,
			}
			r := SingleObjectResp("user", attr)
			So(len(r), ShouldEqual, 1)
			So(len(r["data"].(map[string]interface{})), ShouldEqual, 3)
			So(r["data"].(map[string]interface{})["type"], ShouldEqual, "user")
			So(r["data"].(map[string]interface{})["id"], ShouldEqual, "abc")
			So(r["data"].(map[string]interface{})["attributes"], ShouldResemble, attr)
		})

		Convey(".MultiObjectResp", func() {
			attrs := []map[string]interface{}{
				{
					"id":   "abc",
					"name": "John",
					"age":  100,
				},
				{
					"id":   "xyz",
					"name": "Jane",
					"age":  80,
				},
			}
			r := MultiObjectResp("user", attrs)
			So(len(r), ShouldEqual, 1)
			So(len(r["data"].([]map[string]interface{})), ShouldEqual, 2)
			So(r["data"].([]map[string]interface{})[0]["type"], ShouldEqual, "user")
			So(r["data"].([]map[string]interface{})[0]["id"], ShouldEqual, "abc")
			So(r["data"].([]map[string]interface{})[0]["attributes"], ShouldResemble, attrs[0])
		})
	})
}
