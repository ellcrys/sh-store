package servers

import (
	"context"
	"testing"

	"github.com/ellcrys/elldb/core/servers/db"
	"github.com/ellcrys/elldb/core/servers/proto_rpc"
	"github.com/ellcrys/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBucket(t *testing.T) {
	setup(t, func(rpc, rpc2 *RPC) {
		Convey("Bucket", t, func() {
			var account = db.Account{FirstName: "john", LastName: "Doe", Email: util.RandString(5) + "@example.com", Password: "something"}
			err := rpc.db.Create(&account).Error
			So(err, ShouldBeNil)

			var contract = db.Contract{Creator: account.ID, Name: util.RandString(5), ClientID: util.RandString(5), ClientSecret: util.RandString(5)}
			err = rpc.db.Create(&contract).Error
			So(err, ShouldBeNil)

			Convey(".CreateBucket", func() {
				Convey("Should return error if validation is failed", func() {
					_, err := rpc.CreateBucket(context.Background(), &proto_rpc.CreateBucketMsg{})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"detail": "name is required",
						"code":   "invalid_parameter",
						"source": "/name",
					})
				})

				Convey("Should successfully create a bucket", func() {

					ctx := context.WithValue(context.Background(), CtxContract, &contract)
					b := &proto_rpc.CreateBucketMsg{Name: "mybucket"}
					bucket, err := rpc.CreateBucket(ctx, b)
					So(err, ShouldBeNil)
					So(bucket.Name, ShouldEqual, "mybucket")
					So(bucket.ID, ShouldHaveLength, 36)
					So(bucket.Creator, ShouldEqual, contract.ID)
					So(bucket.Immutable, ShouldBeFalse)

					Convey("Should return error if bucket name has been taken", func() {
						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						b := &proto_rpc.CreateBucketMsg{Name: "mybucket"}
						_, err = rpc.CreateBucket(ctx, b)
						So(err, ShouldNotBeNil)

						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"status": "400",
							"source": "/name",
							"detail": "name has been used",
						})
					})
				})
			})
		})
	})
}
