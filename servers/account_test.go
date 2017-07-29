package servers

import (
	"context"
	"testing"

	"github.com/ellcrys/elldb/servers/proto_rpc"
	"github.com/ellcrys/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAccount(t *testing.T) {
	setup(t, func(rpc, rpc2 *RPC) {
		Convey("Account", t, func() {
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

			Convey("account.go", func() {
				Convey(".CreateAccount", func() {
					Convey("Should return error if validation is failed", func() {
						c1 := &proto_rpc.CreateAccountMsg{}
						_, err := rpc.CreateAccount(context.Background(), c1)
						So(err, ShouldNotBeNil)
						m, _ := util.JSONToMap(err.Error())
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 6)

						So(errs[0], ShouldResemble, map[string]interface{}{
							"code":    "invalid_parameter",
							"field":   "first_name",
							"message": "First Name is required",
						})

						So(errs[1], ShouldResemble, map[string]interface{}{
							"code":    "invalid_parameter",
							"field":   "last_name",
							"message": "Last Name is required",
						})

						So(errs[2], ShouldResemble, map[string]interface{}{
							"code":    "invalid_parameter",
							"field":   "email",
							"message": "Email is required",
						})

						So(errs[3], ShouldResemble, map[string]interface{}{
							"message": "Email is not valid",
							"code":    "invalid_parameter",
							"field":   "email",
						})

						So(errs[4], ShouldResemble, map[string]interface{}{
							"code":    "invalid_parameter",
							"field":   "password",
							"message": "Password is required",
						})

						So(errs[5], ShouldResemble, map[string]interface{}{
							"field":   "password",
							"message": "Password is too short. Must be at least 6 characters long",
							"code":    "invalid_parameter",
						})

						Convey("Should return error if email is invalid", func() {
							c1 := &proto_rpc.CreateAccountMsg{FirstName: "john", LastName: "Doe", Email: "invalid@", Password: "something"}
							_, err := rpc.CreateAccount(context.Background(), c1)
							So(err, ShouldNotBeNil)
							m, _ := util.JSONToMap(err.Error())
							errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
							So(errs, ShouldHaveLength, 1)
							So(errs[0], ShouldResemble, map[string]interface{}{
								"code":    "invalid_parameter",
								"field":   "email",
								"message": "Email is not valid",
							})
						})
					})

					Convey("Should successfully create an account", func() {
						c1 := &proto_rpc.CreateAccountMsg{FirstName: "john", LastName: "Doe", Email: util.RandString(5) + "@example.com", Password: "something"}
						resp, err := rpc.CreateAccount(context.Background(), c1)
						So(err, ShouldBeNil)
						var m map[string]interface{}
						err = util.FromJSON(resp.Object, &m)
						So(err, ShouldBeNil)
						So(m, ShouldHaveLength, 7)
						So(m["id"], ShouldHaveLength, 36)
						So(m["first_name"], ShouldEqual, c1.FirstName)
						So(m["last_name"], ShouldEqual, c1.LastName)
						So(m["email"], ShouldEqual, c1.Email)
						So(m["password"], ShouldNotBeEmpty)
						So(m["password"], ShouldNotEqual, c1.Password)
						So(m["confirmed"], ShouldEqual, false)

						Convey("Should return error if email has been used", func() {
							_, err = rpc.CreateAccount(context.Background(), c1)
							So(err, ShouldNotBeNil)
							m, _ := util.JSONToMap(err.Error())
							errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
							So(errs, ShouldHaveLength, 1)
							So(errs[0], ShouldResemble, map[string]interface{}{
								"code":    "invalid_parameter",
								"field":   "email",
								"message": "Email is not available",
								"status":  "400",
							})
						})

						Convey("Should create client id and secret if account is a developer", func() {
							c1 := &proto_rpc.CreateAccountMsg{FirstName: "john", LastName: "Doe", Email: "user2@example.com", Password: "something", Developer: true}
							resp, err := rpc.CreateAccount(context.Background(), c1)
							So(err, ShouldBeNil)
							var m map[string]interface{}
							err = util.FromJSON(resp.Object, &m)
							So(err, ShouldBeNil)
							So(m, ShouldHaveLength, 10)
							So(m, ShouldContainKey, "client_id")
							So(m, ShouldContainKey, "client_secret")
						})
					})
				})

				Convey(".GetAccount", func() {
					c1 := &proto_rpc.CreateAccountMsg{FirstName: "john", LastName: "Doe", Email: util.RandString(5) + "@example.com", Password: "something"}
					_, err := rpc.CreateAccount(context.Background(), c1)
					So(err, ShouldBeNil)

					Convey("Should successfully get account by email", func() {
						resp, err := rpc.GetAccount(context.Background(), &proto_rpc.GetAccountMsg{ID: c1.Email})
						So(err, ShouldBeNil)
						var m map[string]interface{}
						err = util.FromJSON(resp.Object, &m)
						So(err, ShouldBeNil)
						So(c1.Email, ShouldEqual, m["email"])

						Convey("Should successfully get account by ID", func() {
							resp, err := rpc.GetAccount(context.Background(), &proto_rpc.GetAccountMsg{ID: m["id"].(string)})
							So(err, ShouldBeNil)
							var m map[string]interface{}
							err = util.FromJSON(resp.Object, &m)
							So(err, ShouldBeNil)
							So(c1.Email, ShouldEqual, m["email"])
						})
					})

					Convey("Should return error if account is not found", func() {
						_, err := rpc.GetAccount(context.Background(), &proto_rpc.GetAccountMsg{ID: "unknown"})
						So(err, ShouldNotBeNil)
						m, _ := util.JSONToMap(err.Error())
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"status":  "404",
							"message": "Account not found",
						})
					})
				})
			})
		})
	})
}
