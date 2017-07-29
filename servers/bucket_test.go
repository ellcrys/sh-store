package servers

import (
	"context"
	"testing"

	"github.com/ellcrys/elldb/servers/proto_rpc"
	"github.com/ellcrys/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBucket(t *testing.T) {
	setup(t, func(rpc, rpc2 *RPC) {
		Convey("Bucket", t, func() {
			c1 := &proto_rpc.CreateAccountMsg{FirstName: "john", LastName: "Doe", Email: util.RandString(5) + "@example.com", Password: "something", Developer: true}
			resp, err := rpc.CreateAccount(context.Background(), c1)
			So(err, ShouldBeNil)
			var account map[string]interface{}
			util.FromJSON(resp.Object, &account)

			c2 := &proto_rpc.CreateAccountMsg{FirstName: "john2", LastName: "Doe2", Email: util.RandString(5) + "@example.com", Password: "something", Developer: true}
			resp, err = rpc.CreateAccount(context.Background(), c2)
			So(err, ShouldBeNil)
			var account2 map[string]interface{}
			util.FromJSON(resp.Object, &account2)

			ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
			b := &proto_rpc.CreateBucketMsg{Name: util.RandString(5)}
			bucket, err := rpc.CreateBucket(ctx, b)
			So(err, ShouldBeNil)
			So(bucket.Name, ShouldEqual, b.Name)
			So(bucket.ID, ShouldHaveLength, 36)

			Convey("bucket.go", func() {
				Convey(".CreateBucket", func() {
					Convey("Should return error if validation is failed", func() {
						_, err := rpc.CreateBucket(context.Background(), &proto_rpc.CreateBucketMsg{})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"message": "Name is required",
							"code":    "invalid_parameter",
							"field":   "name",
						})
					})

					Convey("Should successfully create a bucket", func() {
						c1 := &proto_rpc.CreateAccountMsg{FirstName: "john", LastName: "Doe", Email: "user@example.com", Password: "something"}
						resp, err := rpc.CreateAccount(context.Background(), c1)
						So(err, ShouldBeNil)
						var account map[string]interface{}
						util.FromJSON(resp.Object, &account)

						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						b := &proto_rpc.CreateBucketMsg{Name: "mybucket"}
						bucket, err := rpc.CreateBucket(ctx, b)
						So(err, ShouldBeNil)
						So(bucket.Name, ShouldEqual, "mybucket")
						So(bucket.ID, ShouldHaveLength, 36)

						Convey("Should return error if bucket name has been taken", func() {
							ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
							b := &proto_rpc.CreateBucketMsg{Name: "mybucket"}
							_, err = rpc.CreateBucket(ctx, b)
							So(err, ShouldNotBeNil)

							m, err := util.JSONToMap(err.Error())
							So(err, ShouldBeNil)
							errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
							So(errs, ShouldHaveLength, 1)
							So(errs[0], ShouldResemble, map[string]interface{}{
								"status":  "400",
								"field":   "name",
								"message": "Bucket name has been used",
							})
						})
					})
				})
			})

		})
	})
}
