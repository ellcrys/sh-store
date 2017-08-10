package servers

import (
	"context"
	"testing"

	"github.com/ellcrys/elldb/core/servers/db"
	"github.com/ellcrys/elldb/core/servers/proto_rpc"
	"github.com/ellcrys/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAccount(t *testing.T) {
	setup(t, func(rpc, rpc2 *RPC) {
		Convey("Account", t, func() {

			var account = db.Account{FirstName: "john", LastName: "Doe", Email: util.RandString(5) + "@example.com", Password: "something"}
			err := rpc.db.Create(&account).Error
			So(err, ShouldBeNil)

			var account2 = db.Account{FirstName: "john2", LastName: "Doe2", Email: util.RandString(5) + "@example.com", Password: "something"}
			err = rpc.db.Create(&account2).Error
			So(err, ShouldBeNil)

			Convey("account.go", func() {

				Convey(".GetAccount", func() {

					Convey("Should successfully get account by email", func() {
						acct, err := rpc.GetAccount(context.Background(), &proto_rpc.GetAccountMsg{ID: account.Email})
						So(err, ShouldBeNil)
						So(acct, ShouldNotBeNil)
						So(account.Email, ShouldEqual, acct.Email)

						Convey("Should successfully get account by ID", func() {
							acct2, err := rpc.GetAccount(context.Background(), &proto_rpc.GetAccountMsg{ID: acct.ID})
							So(err, ShouldBeNil)
							So(acct2, ShouldNotBeNil)
							So(account.Email, ShouldEqual, acct2.Email)
						})
					})

					Convey("Should return error if account is not found", func() {
						_, err := rpc.GetAccount(context.Background(), &proto_rpc.GetAccountMsg{ID: "unknown"})
						So(err, ShouldNotBeNil)
						m, _ := util.JSONToMap(err.Error())
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"status": "404",
							"source": "id",
							"detail": "account not found",
						})
					})
				})
			})
		})
	})
}
