package servers

import (
	"context"
	"testing"

	"github.com/ellcrys/elldb/servers/proto_rpc"
	"github.com/ellcrys/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMapping(t *testing.T) {
	setup(t, func(rpc, rpc2 *RPC) {
		Convey("Mapping", t, func() {
			c1 := &proto_rpc.CreateIdentityMsg{FirstName: "john", LastName: "Doe", Email: util.RandString(5) + "@example.com", Password: "something", Developer: true}
			resp, err := rpc.CreateIdentity(context.Background(), c1)
			So(err, ShouldBeNil)
			var identity map[string]interface{}
			util.FromJSON(resp.Object, &identity)

			c2 := &proto_rpc.CreateIdentityMsg{FirstName: "john2", LastName: "Doe2", Email: util.RandString(5) + "@example.com", Password: "something", Developer: true}
			resp, err = rpc.CreateIdentity(context.Background(), c2)
			So(err, ShouldBeNil)
			var identity2 map[string]interface{}
			util.FromJSON(resp.Object, &identity2)

			ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
			b := &proto_rpc.CreateBucketMsg{Name: util.RandString(5)}
			bucket, err := rpc.CreateBucket(ctx, b)
			So(err, ShouldBeNil)
			So(bucket.Name, ShouldEqual, b.Name)
			So(bucket.ID, ShouldHaveLength, 36)

			Convey(".CreateMapping / .GetMapping / .GetAllMapping", func() {
				c1 := &proto_rpc.CreateIdentityMsg{FirstName: "john", LastName: "Doe", Email: util.RandString(5) + "@example.com", Password: "something"}
				resp, err := rpc.CreateIdentity(context.Background(), c1)
				So(err, ShouldBeNil)
				var identity map[string]interface{}
				util.FromJSON(resp.Object, &identity)

				Convey("Should return error if bucket is not provided", func() {
					_, err := rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "400",
						"message": "bucket is required",
					})
				})

				Convey("Should return error if bucket does not exists", func() {
					_, err := rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{Bucket: "unknown"})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "404",
						"field":   "bucket",
						"message": "bucket not found",
					})
				})

				Convey("Should return error if mapping value is malformed json", func() {
					var mapping = `{ "key": "column"`
					_, err := rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{Name: "some_name", Bucket: b.Name, Mapping: []byte(mapping)})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "400",
						"message": "malformed mapping",
					})
				})

				Convey("Should return error if mapping value has unexpected type", func() {
					var mapping = `{ "key": 123 }`
					_, err := rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{Name: "some_name", Bucket: b.Name, Mapping: []byte(mapping)})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "400",
						"message": "malformed mapping",
					})
				})

				Convey("Should successfully create a mapping", func() {
					ctx := context.WithValue(ctx, CtxIdentity, identity["id"])
					b := &proto_rpc.CreateBucketMsg{Name: util.RandString(5)}
					_, err = rpc.CreateBucket(ctx, b)
					So(err, ShouldBeNil)

					name := util.RandString(5)
					mapping := `{ "user_id": "owner_id" }`
					resp, err := rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{Name: name, Bucket: b.Name, Mapping: []byte(mapping)})
					So(err, ShouldBeNil)
					So(resp.Name, ShouldEqual, name)

					Convey("Should return error if mapping name has been used", func() {
						_, err = rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{Name: name, Bucket: b.Name, Mapping: []byte(mapping)})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"status":  "400",
							"message": "mapping with name '" + name + "' already exist",
						})
					})

					Convey(".GetMapping", func() {

						Convey("Should return error if mapping name is missing", func() {
							_, err := rpc.GetMapping(ctx, &proto_rpc.GetMappingMsg{})
							So(err, ShouldNotBeNil)
							m, err := util.JSONToMap(err.Error())
							So(err, ShouldBeNil)
							errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
							So(errs, ShouldHaveLength, 1)
							So(errs[0], ShouldResemble, map[string]interface{}{
								"message": "name is required",
								"status":  "400",
								"field":   "id",
							})
						})

						Convey("Should return error if mapping does not exists", func() {
							_, err := rpc.GetMapping(ctx, &proto_rpc.GetMappingMsg{Name: "unknown"})
							So(err, ShouldNotBeNil)
							m, err := util.JSONToMap(err.Error())
							So(err, ShouldBeNil)
							errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
							So(errs, ShouldHaveLength, 1)
							So(errs[0], ShouldResemble, map[string]interface{}{
								"status":  "404",
								"message": "mapping not found",
							})
						})

						Convey("Should successfully get an existing mapping", func() {
							resp, err := rpc.GetMapping(ctx, &proto_rpc.GetMappingMsg{Name: name})
							So(err, ShouldBeNil)
							m, err := util.JSONToMap(string(resp.Mapping))
							So(err, ShouldBeNil)
							So(m["name"], ShouldEqual, name)
							So(m["bucket"], ShouldEqual, b.Name)
							mapping, err := util.JSONToMap(m["mapping"].(string))
							So(err, ShouldBeNil)
							So(mapping, ShouldResemble, map[string]interface{}{
								"user_id": "owner_id",
							})
						})
					})

					Convey(".GetAllMapping", func() {
						name := util.RandString(5)
						mapping := `{ "user_id": "owner_id" }`
						_, err := rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{Name: name, Bucket: b.Name, Mapping: []byte(mapping)})
						So(err, ShouldBeNil)

						Convey("Should get all mapping. Limit result to 1 mapping", func() {
							resp, err := rpc.GetAllMapping(ctx, &proto_rpc.GetAllMappingMsg{Limit: 1})
							So(err, ShouldBeNil)
							var m []map[string]interface{}
							err = util.FromJSON(resp.Mappings, &m)
							So(err, ShouldBeNil)
							So(m, ShouldHaveLength, 1)
						})

						Convey("Should get all mapping. Limit result to 2 mapping", func() {
							resp, err := rpc.GetAllMapping(ctx, &proto_rpc.GetAllMappingMsg{Limit: 2})
							So(err, ShouldBeNil)
							var m []map[string]interface{}
							err = util.FromJSON(resp.Mappings, &m)
							So(err, ShouldBeNil)
							So(m, ShouldHaveLength, 2)
						})

						Convey("Should use default max limit if limit is not set", func() {
							MaxGetAllMappingLimit = 1
							resp, err := rpc.GetAllMapping(ctx, &proto_rpc.GetAllMappingMsg{})
							So(err, ShouldBeNil)
							var m []map[string]interface{}
							err = util.FromJSON(resp.Mappings, &m)
							So(err, ShouldBeNil)
							So(m, ShouldHaveLength, 1)
							MaxGetAllMappingLimit = 50
						})

						Convey("Should use default max limit if limit is greater than max limit", func() {
							MaxGetAllMappingLimit = 1
							resp, err := rpc.GetAllMapping(ctx, &proto_rpc.GetAllMappingMsg{Limit: 2})
							So(err, ShouldBeNil)
							var m []map[string]interface{}
							err = util.FromJSON(resp.Mappings, &m)
							So(err, ShouldBeNil)
							So(m, ShouldHaveLength, 1)
							MaxGetAllMappingLimit = 50
						})
					})
				})
			})
		})
	})
}
