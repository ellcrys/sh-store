package servers

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/ellcrys/elldb/core/servers/db"
	"github.com/ellcrys/elldb/core/servers/proto_rpc"
	"github.com/ellcrys/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestContract(t *testing.T) {
	setup(t, func(rpc, rpc2 *RPC) {
		Convey("Contract", t, func() {
			var account = db.Account{FirstName: "john", LastName: "Doe", Email: util.RandString(5) + "@example.com", Password: "something"}
			err := rpc.db.Create(&account).Error
			So(err, ShouldBeNil)

			Convey(".CreateContract", func() {
				Convey("Should return error if validation is failed", func() {
					_, err := rpc.CreateContract(context.Background(), &proto_rpc.CreateContractMsg{})
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

				Convey("Should successfully create a contract", func() {
					ctx := context.WithValue(context.Background(), CtxAccount, &account)
					name := "my_contract_abc"
					contract, err := rpc.CreateContract(ctx, &proto_rpc.CreateContractMsg{Name: name})
					So(err, ShouldBeNil)
					So(contract, ShouldNotBeNil)
					So(contract.Name, ShouldEqual, name)

					Convey("Should return error if contract name has been used", func() {
						_, err := rpc.CreateContract(ctx, &proto_rpc.CreateContractMsg{Name: name})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"status": "400",
							"source": "/name",
							"detail": "name is not available",
						})
					})
				})
			})
		})
	})
}
