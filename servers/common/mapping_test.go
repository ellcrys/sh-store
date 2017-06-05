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

		Convey("Should return nil if mapping is empty", func() {
			err := UnMapFields(map[string]string{}, "invalid_type")
			So(err, ShouldBeNil)
		})

		Convey("Should return error if mapped object type is invalid", func() {
			err := UnMapFields(map[string]string{
				"custom_key": "key_to_map",
			}, "invalid_type")
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

	Convey(".MapFields", t, func() {

		Convey("Should successfully apply mapping to a single map", func() {
			mapping := map[string]string{
				"my_key": "key",
				"my_val": "val",
			}

			unmapped := map[string]interface{}{
				"key": map[string]interface{}{
					"key": map[string]interface{}{
						"key": "123",
					},
				},
				"val": "12",
			}

			expected := map[string]interface{}{
				"my_key": map[string]interface{}{
					"my_key": map[string]interface{}{
						"my_key": "123",
					},
				},
				"my_val": "12",
			}

			err := MapFields(mapping, unmapped)
			So(err, ShouldBeNil)
			So(unmapped, ShouldResemble, expected)
		})

		Convey("Should successfully apply mapping to a slice of maps", func() {
			mapping := map[string]string{
				"my_key": "key",
				"my_val": "val",
			}

			unmapped := []map[string]interface{}{
				{
					"key": map[string]interface{}{
						"key": map[string]interface{}{
							"key": "123",
						},
					},
					"val": "12",
				},
				{
					"key": map[string]interface{}{
						"key": map[string]interface{}{
							"key": "123",
						},
					},
					"val": "12",
				},
			}

			expected := map[string]interface{}{
				"my_key": map[string]interface{}{
					"my_key": map[string]interface{}{
						"my_key": "123",
					},
				},
				"my_val": "12",
			}

			err := MapFields(mapping, unmapped)
			So(err, ShouldBeNil)
			So(unmapped[0], ShouldResemble, expected)
			So(unmapped[1], ShouldResemble, expected)
		})
	})

}
