package common

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMapping(t *testing.T) {

	mapping := map[string]string{
		"name": "ref1",
		"age":  "ref2",
	}

	Convey(".UnMapFields", t, func() {

		Convey("Should return error if mapped object type is invalid", func() {
			err := UnMapFields(map[string]string{}, "invalid_type")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "invalid mappedObj type. Expects a map or slice of map")
		})

		Convey("Should successfully unmap the mapped fields", func() {
			mappedObj := map[string]interface{}{
				"name": "ken",
				"age":  20,
				"inner": map[string]interface{}{
					"age": 50,
					"name": map[string]interface{}{
						"name": "pete",
					},
				},
			}

			expected := map[string]interface{}{
				"ref1": "ken",
				"ref2": 20,
				"inner": map[string]interface{}{
					"ref2": 50,
					"ref1": map[string]interface{}{
						"ref1": "pete",
					},
				},
			}
			UnMapFields(mapping, mappedObj)
			So(mappedObj, ShouldResemble, expected)

			mappedObjs := []map[string]interface{}{
				map[string]interface{}{
					"name": "ken",
					"age":  20,
					"inner": map[string]interface{}{
						"age": 50,
						"name": map[string]interface{}{
							"name": "pete",
						},
					},
				},
			}

			UnMapFields(mapping, mappedObjs)
			for _, m := range mappedObjs {
				So(m, ShouldResemble, expected)
			}
		})
	})

}
