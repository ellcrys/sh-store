package servers

import (
	"context"
	"testing"

	"github.com/ellcrys/elldb/servers/db"
	"github.com/ellcrys/elldb/servers/proto_rpc"
	"github.com/ellcrys/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMapping(t *testing.T) {
	setup(t, func(rpc, rpc2 *RPC) {
		Convey("Mapping", t, func() {
			var account = db.Account{FirstName: "john", LastName: "Doe", Email: util.RandString(5) + "@example.com", Password: "something"}
			err := rpc.db.Create(&account).Error
			So(err, ShouldBeNil)

			var contract = db.Contract{ID: util.UUID4(), Creator: account.ID, Name: util.RandString(5), ClientID: util.RandString(5), ClientSecret: util.RandString(5)}
			err = rpc.db.Create(&contract).Error
			So(err, ShouldBeNil)

			var contract2 = db.Contract{ID: util.UUID4(), Creator: account.ID, Name: util.RandString(5), ClientID: util.RandString(5), ClientSecret: util.RandString(5)}
			err = rpc.db.Create(&contract2).Error
			So(err, ShouldBeNil)

			ctx := context.WithValue(context.Background(), CtxContract, &contract)
			b := &proto_rpc.CreateBucketMsg{Name: util.RandString(5)}
			bucket, err := rpc.CreateBucket(ctx, b)
			So(err, ShouldBeNil)
			So(bucket.Name, ShouldEqual, b.Name)
			So(bucket.ID, ShouldHaveLength, 36)

			Convey(".CreateMapping / .GetMapping / .GetAllMapping", func() {

				Convey("Should return error if bucket is not provided", func() {
					_, err := rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"detail": "bucket is required",
						"status": "400",
					})
				})

				Convey("Should return error if bucket does not exists", func() {
					_, err := rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{Bucket: "unknown"})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status": "404",
						"source": "/bucket",
						"detail": "bucket not found",
					})
				})

				Convey("Should return error if authenticated contract does not own the bucket", func() {
					ctx = context.WithValue(context.Background(), CtxContract, &contract2)
					_, err := rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{Name: "some_name", Bucket: b.Name, Mapping: ""})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"detail": "You do not have permission to perform this action",
						"status": "401",
					})
				})

				Convey("Should return error if mapping value is malformed json", func() {
					var mapping = `{ "key": "column"`
					_, err := rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{Name: "some_name", Bucket: b.Name, Mapping: mapping})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"detail": "malformed mapping",
						"status": "400",
					})
				})

				Convey("Should return error if mapping value has unexpected type", func() {
					var mapping = `{ "key": 123 }`
					_, err := rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{Name: "some_name", Bucket: b.Name, Mapping: mapping})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"detail": "malformed mapping",
						"status": "400",
					})
				})

				Convey("Should successfully create a mapping", func() {

					name := util.RandString(5)
					mapping := `{ "user_id": "owner_id" }`
					newMapping, err := rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{Name: name, Bucket: b.Name, Mapping: mapping})
					So(err, ShouldBeNil)
					So(newMapping, ShouldNotBeNil)
					So(newMapping.Name, ShouldEqual, name)
					So(newMapping.Creator, ShouldEqual, contract.ID)
					So(newMapping.Mapping, ShouldEqual, mapping)
					So(newMapping.Bucket, ShouldEqual, b.Name)

					Convey("Should return error if mapping name has been used", func() {
						_, err = rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{Name: name, Bucket: b.Name, Mapping: mapping})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"detail": "mapping with name '" + name + "' already exist",
							"status": "400",
						})
					})

					Convey(".GetMapping", func() {

						Convey("Should return error if mapping name is missing", func() {
							_, err := rpc.GetMapping(ctx, &proto_rpc.GetMappingMsg{})
							So(err, ShouldNotBeNil)
							m, err := util.JSONToMap(err.Error())
							So(err, ShouldBeNil)
							errs := m["errors"].([]interface{})
							So(errs, ShouldHaveLength, 1)
							So(errs[0], ShouldResemble, map[string]interface{}{
								"detail": "name is required",
								"status": "400",
								"source": "name",
							})
						})

						Convey("Should return error if mapping does not exists", func() {
							_, err := rpc.GetMapping(ctx, &proto_rpc.GetMappingMsg{Name: "unknown"})
							So(err, ShouldNotBeNil)
							m, err := util.JSONToMap(err.Error())
							So(err, ShouldBeNil)
							errs := m["errors"].([]interface{})
							So(errs, ShouldHaveLength, 1)
							So(errs[0], ShouldResemble, map[string]interface{}{
								"status": "404",
								"detail": "mapping not found",
							})
						})

						Convey("Should successfully get an existing mapping", func() {
							resp, err := rpc.GetMapping(ctx, &proto_rpc.GetMappingMsg{Name: name})
							So(err, ShouldBeNil)
							So(resp.ID, ShouldEqual, newMapping.ID)
							So(resp.Mapping, ShouldNotBeEmpty)
							So(resp.Creator, ShouldEqual, contract.ID)
							So(resp.Bucket, ShouldEqual, newMapping.Bucket)
						})
					})

					Convey(".GetAllMapping", func() {
						name := util.RandString(5)
						mapping := `{ "user_id": "owner_id" }`
						_, err := rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{Name: name, Bucket: b.Name, Mapping: mapping})
						So(err, ShouldBeNil)

						Convey("Should get all mapping. Limit result to 1 mapping", func() {
							resp, err := rpc.GetAllMapping(ctx, &proto_rpc.GetAllMappingMsg{Limit: 1})
							So(err, ShouldBeNil)
							So(resp, ShouldNotBeNil)
							So(resp.Mappings, ShouldHaveLength, 1)
						})

						Convey("Should get all mapping. Limit result to 2 mapping", func() {
							resp, err := rpc.GetAllMapping(ctx, &proto_rpc.GetAllMappingMsg{Limit: 2})
							So(err, ShouldBeNil)
							So(resp, ShouldNotBeNil)
							So(resp.Mappings, ShouldHaveLength, 2)
						})

						Convey("Should use default max limit if limit is not set", func() {
							MaxGetAllMappingLimit = 1
							resp, err := rpc.GetAllMapping(ctx, &proto_rpc.GetAllMappingMsg{})
							So(err, ShouldBeNil)
							So(resp, ShouldNotBeNil)
							So(resp.Mappings, ShouldHaveLength, 1)
							MaxGetAllMappingLimit = 50
						})

						Convey("Should use default max limit if limit is greater than max limit", func() {
							MaxGetAllMappingLimit = 1
							resp, err := rpc.GetAllMapping(ctx, &proto_rpc.GetAllMappingMsg{Limit: 2})
							So(err, ShouldBeNil)
							So(resp, ShouldNotBeNil)
							So(resp.Mappings, ShouldHaveLength, 1)
							MaxGetAllMappingLimit = 50
						})
					})
				})
			})
		})
	})
}
