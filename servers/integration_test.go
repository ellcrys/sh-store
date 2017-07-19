package servers

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/ellcrys/elldb/servers/common"
	"github.com/ellcrys/elldb/servers/db"
	"github.com/ellcrys/elldb/servers/proto_rpc"
	"github.com/ellcrys/elldb/session"
	"github.com/ellcrys/util"
	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

var err error
var testDB *sql.DB
var dbName string
var conStr = "postgresql://root@localhost:26257?sslmode=disable"
var conStrWithDB string

func init() {
	testDB, err = common.OpenDB(conStr)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %s", err))
	}
	dbName, err = common.CreateRandomDB(testDB)
	if err != nil {
		panic(fmt.Errorf("failed to create test database: %s", err))
	}
	conStrWithDB = "postgresql://root@localhost:26257/" + dbName + "?sslmode=disable"
}

func TestRPC(t *testing.T) {

	dbCon, err := gorm.Open("postgres", conStrWithDB)
	if err != nil {
		t.Fatal("failed to connect to test database")
	}

	err = dbCon.CreateTable(&db.Bucket{}, &db.Object{}, &db.Identity{}, &db.Mapping{}).Error
	if err != nil {
		t.Fatalf("failed to create database tables. %s", err)
	}

	rpc := NewRPC()
	rpc.db = dbCon

	rpc.sessionReg, err = session.NewConsulRegistry()
	if err != nil {
		t.Fatalf("failed to connect consul registry. %s", err)
	}

	rpc.dbSession = session.NewSession(rpc.sessionReg)
	rpc.dbSession.SetDB(rpc.db)

	IntegrationTest(t, rpc)

	defer func() {
		common.DropDB(testDB, dbName)
	}()
}

func IntegrationTest(t *testing.T, rpc *RPC) {
	Convey("RPC Integration Test", t, func() {

		Convey("identity.go", func() {
			Convey(".CreateIdentity", func() {
				Convey("Should return error if validation is failed", func() {
					c1 := &proto_rpc.CreateIdentityMsg{}
					_, err := rpc.CreateIdentity(context.Background(), c1)
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
						c1 := &proto_rpc.CreateIdentityMsg{FirstName: "john", LastName: "Doe", Email: "invalid@", Password: "something"}
						_, err := rpc.CreateIdentity(context.Background(), c1)
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

				Convey("Should successfully create an identity", func() {
					c1 := &proto_rpc.CreateIdentityMsg{FirstName: "john", LastName: "Doe", Email: "user@example.com", Password: "something"}
					resp, err := rpc.CreateIdentity(context.Background(), c1)
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
						_, err = rpc.CreateIdentity(context.Background(), c1)
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

					Convey("Should create client id and secret if identity is a developer", func() {
						c1 := &proto_rpc.CreateIdentityMsg{FirstName: "john", LastName: "Doe", Email: "user2@example.com", Password: "something", Developer: true}
						resp, err := rpc.CreateIdentity(context.Background(), c1)
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

			Convey(".GetIdentity", func() {
				c1 := &proto_rpc.CreateIdentityMsg{FirstName: "john", LastName: "Doe", Email: "user@example.com", Password: "something"}
				_, err := rpc.CreateIdentity(context.Background(), c1)
				So(err, ShouldBeNil)

				Convey("Should successfully get identity by email", func() {
					resp, err := rpc.GetIdentity(context.Background(), &proto_rpc.GetIdentityMsg{ID: c1.Email})
					So(err, ShouldBeNil)
					var m map[string]interface{}
					err = util.FromJSON(resp.Object, &m)
					So(err, ShouldBeNil)
					So(c1.Email, ShouldEqual, m["email"])

					Convey("Should successfully get identity by ID", func() {
						resp, err := rpc.GetIdentity(context.Background(), &proto_rpc.GetIdentityMsg{ID: m["id"].(string)})
						So(err, ShouldBeNil)
						var m map[string]interface{}
						err = util.FromJSON(resp.Object, &m)
						So(err, ShouldBeNil)
						So(c1.Email, ShouldEqual, m["email"])
					})
				})

				Convey("Should return error if identity is not found", func() {
					_, err := rpc.GetIdentity(context.Background(), &proto_rpc.GetIdentityMsg{ID: "unknown"})
					So(err, ShouldNotBeNil)
					m, _ := util.JSONToMap(err.Error())
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "404",
						"message": "Identity not found",
					})
				})
			})
		})

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
					c1 := &proto_rpc.CreateIdentityMsg{FirstName: "john", LastName: "Doe", Email: "user@example.com", Password: "something"}
					resp, err := rpc.CreateIdentity(context.Background(), c1)
					So(err, ShouldBeNil)
					var identity map[string]interface{}
					util.FromJSON(resp.Object, &identity)

					ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
					b := &proto_rpc.CreateBucketMsg{Name: "mybucket"}
					bucket, err := rpc.CreateBucket(ctx, b)
					So(err, ShouldBeNil)
					So(bucket.Name, ShouldEqual, "mybucket")
					So(bucket.ID, ShouldHaveLength, 36)

					Convey("Should return error if bucket name has been taken", func() {
						ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
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

		Convey("validation.go", func() {
			Convey(".validateObjects", func() {
				Convey("Should return error if no object is provided", func() {
					errs, err := rpc.validateObjects(nil, nil)
					So(err, ShouldBeNil)
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, common.Error{
						Status:  "",
						Code:    "",
						Field:   "",
						Message: "no object provided. At least one object is required",
						Links:   nil,
					})
				})

				Convey("Should return error if objects length exceeds max", func() {
					curMaxObjPerPut := MaxObjectPerPut
					MaxObjectPerPut = 2
					var objs = []map[string]interface{}{}
					objs = append(objs, structs.New(db.Object{}).Map(), structs.New(db.Object{}).Map(), structs.New(db.Object{}).Map())
					errs, err := rpc.validateObjects(objs, nil)
					So(err, ShouldBeNil)
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, common.Error{
						Status:  "",
						Code:    "",
						Field:   "",
						Message: "too many objects. Only a maximum of 2 can be created at once",
						Links:   nil,
					})
					MaxObjectPerPut = curMaxObjPerPut
				})

				Convey("Should return error if object required fields are missing (validation)", func() {
					var objs = []map[string]interface{}{}
					objs = append(objs, map[string]interface{}{})
					errs, err := rpc.validateObjects(objs, nil)
					So(err, ShouldBeNil)
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: owner_id is required", Links: nil})
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: creator_id is required", Links: nil})
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: key is required", Links: nil})
				})

				Convey("Should return error if object required fields have invalid type (validation)", func() {
					var objs = []map[string]interface{}{}
					objs = append(objs, map[string]interface{}{
						"owner_id":   123,
						"creator_id": 123,
						"key":        123,
						"value":      123,
						"ref1":       123,
					})
					errs, err := rpc.validateObjects(objs, nil)
					So(err, ShouldBeNil)
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: owner_id must be a string", Links: nil})
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: creator_id must be a string", Links: nil})
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: key must be a string", Links: nil})
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: value must be a string", Links: nil})
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: ref1 must be a string", Links: nil})
				})

				Convey("Should return error if objects do not share the same owner_id", func() {
					var objs = []map[string]interface{}{}
					objs = append(objs, map[string]interface{}{"owner_id": "xyz", "creator_id": "abc", "key": "some_key"})
					objs = append(objs, map[string]interface{}{"owner_id": "xyz2", "creator_id": "abc", "key": "some_key"})
					errs, err := rpc.validateObjects(objs, nil)
					So(err, ShouldBeNil)
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 1: all objects must share the same 'owner_id' (owner_id) value", Links: nil})
				})

				Convey("Should return errors that use custom field names if mapping is provided", func() {
					var objs = []map[string]interface{}{}
					objs = append(objs, map[string]interface{}{})
					errs, err := rpc.validateObjects(objs, map[string]string{
						"user_id": "owner_id",
						"app_id":  "creator_id",
						"name":    "key",
					})
					So(err, ShouldBeNil)
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: user_id is required", Links: nil})
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: app_id is required", Links: nil})
					So(errs, ShouldContain, common.Error{Status: "", Code: "invalid_parameter", Field: "", Message: "object 0: name is required", Links: nil})
				})

				Convey("Should return error if objects' owner does not exists", func() {
					var objs = []map[string]interface{}{}
					objs = append(objs, map[string]interface{}{"owner_id": "unknown", "creator_id": "abc", "key": "mykey"})
					errs, err := rpc.validateObjects(objs, nil)
					So(err, ShouldBeNil)
					So(errs, ShouldContain, common.Error{
						Status:  "",
						Code:    "invalid_parameter",
						Field:   "",
						Message: "owner of object(s) does not exist",
						Links:   nil,
					})
				})
			})
		})

		Convey("mapping.go", func() {
			Convey(".CreateMapping / .GetMapping / .GetAllMapping", func() {
				c1 := &proto_rpc.CreateIdentityMsg{FirstName: "john", LastName: "Doe", Email: "user@example.com", Password: "something"}
				resp, err := rpc.CreateIdentity(context.Background(), c1)
				So(err, ShouldBeNil)
				var identity map[string]interface{}
				util.FromJSON(resp.Object, &identity)

				Convey("Should return error if mapping value is malformed json", func() {
					var mapping = `{ "key": "column"`
					_, err := rpc.CreateMapping(context.Background(), &proto_rpc.CreateMappingMsg{Name: "some_name", Mapping: []byte(mapping)})
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
					_, err := rpc.CreateMapping(context.Background(), &proto_rpc.CreateMappingMsg{Name: "some_name", Mapping: []byte(mapping)})
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
					ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
					name := util.RandString(5)
					mapping := `{ "user_id": "owner_id" }`
					resp, err := rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{Name: name, Mapping: []byte(mapping)})
					So(err, ShouldBeNil)
					So(resp.Name, ShouldEqual, name)

					Convey("Should return error if mapping name has been used", func() {
						_, err = rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{Name: name, Mapping: []byte(mapping)})
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
						_, err := rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{Name: name, Mapping: []byte(mapping)})
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

		Convey("objects.go", func() {
			c1 := &proto_rpc.CreateIdentityMsg{FirstName: "john", LastName: "Doe", Email: "user@example.com", Password: "something"}
			resp, err := rpc.CreateIdentity(context.Background(), c1)
			So(err, ShouldBeNil)
			var identity map[string]interface{}
			util.FromJSON(resp.Object, &identity)

			c2 := &proto_rpc.CreateIdentityMsg{FirstName: "john2", LastName: "Doe2", Email: "user2@example.com", Password: "something"}
			resp, err = rpc.CreateIdentity(context.Background(), c2)
			So(err, ShouldBeNil)
			var identity2 map[string]interface{}
			util.FromJSON(resp.Object, &identity2)

			ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
			b := &proto_rpc.CreateBucketMsg{Name: "mybucket"}
			bucket, err := rpc.CreateBucket(ctx, b)
			So(err, ShouldBeNil)
			So(bucket.Name, ShouldEqual, "mybucket")
			So(bucket.ID, ShouldHaveLength, 36)

			Convey(".CreateObjects", func() {

				Convey("Should return error if bucket is not provided", func() {
					ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
					_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"field":   "bucket",
						"message": "bucket name is required",
						"status":  "400",
					})
				})

				Convey("Should return error if bucket does not exists", func() {
					ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
					_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: "unknown"})
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

				Convey("With bucket", func() {

					Convey("Should return error if validation fails", func() {
						ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
						objs := util.MustStringify([]map[string]interface{}{
							{"owner_id": 123, "key": "mykey"},
						})
						_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"code":    "invalid_parameter",
							"message": "object 0: owner_id must be a string",
						})
					})

					Convey("Should return error if owner of object does not exist", func() {
						ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
						objs := util.MustStringify([]map[string]interface{}{
							{"owner_id": "unknown", "key": "mykey"},
						})
						_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"code":    "invalid_parameter",
							"message": "owner of object(s) does not exist",
						})
					})

					Convey("Should return permission error if authenticated identity does not have permission to create object for the owner", func() {
						ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
						objs := util.MustStringify([]map[string]interface{}{
							{"owner_id": identity2["id"].(string), "key": "mykey"},
						})
						_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"status":  "401",
							"message": "permission denied: you are not authorized to create objects for the owner",
						})
					})

					Convey("Should successfully create new objects", func() {
						ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
						objs := util.MustStringify([]map[string]interface{}{
							{"owner_id": identity["id"].(string), "key": "mykey", "value": "myval"},
							{"owner_id": identity["id"].(string), "key": "mykey2", "value": "myval2"},
						})
						resp, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs})
						So(err, ShouldBeNil)
						var m []map[string]interface{}
						err = util.FromJSON(resp.Objects, &m)
						So(err, ShouldBeNil)
						So(m, ShouldHaveLength, 2)
						So(m[0]["bucket"], ShouldEqual, b.Name)
						So(m[1]["bucket"], ShouldEqual, b.Name)
						So(m[0]["creator_id"], ShouldEqual, identity["id"])
						So(m[1]["creator_id"], ShouldEqual, identity["id"])
						So(m[0]["owner_id"], ShouldEqual, identity["id"])
						So(m[1]["owner_id"], ShouldEqual, identity["id"])
					})

					Convey("With mapping", func() {
						ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
						mappingName := util.RandString(5)
						mapping := `{ "user_id": "owner_id", "first_name": "key", "last_name": "value" }`
						_, err := rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{Name: mappingName, Mapping: []byte(mapping)})
						So(err, ShouldBeNil)

						Convey("Should return error if mapping does not exists", func() {
							ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
							objs := util.MustStringify([]map[string]interface{}{
								{"owner_id": identity["id"].(string), "key": "mykey", "value": "myval"},
							})
							_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs, Mapping: "unknown"})
							So(err, ShouldNotBeNil)
							m, err := util.JSONToMap(err.Error())
							So(err, ShouldBeNil)
							errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
							So(errs, ShouldHaveLength, 1)
							So(errs[0], ShouldResemble, map[string]interface{}{
								"message": "mapping not found",
								"status":  "400",
							})
						})

						Convey("Should successfully create object using a mapping", func() {
							ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
							objs := util.MustStringify([]map[string]interface{}{
								{"user_id": identity["id"].(string), "first_name": "John", "last_name": "Doe"},
							})
							resp, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs, Mapping: mappingName})
							So(err, ShouldBeNil)
							var m []map[string]interface{}
							err = util.FromJSON(resp.Objects, &m)
							So(err, ShouldBeNil)
							So(m, ShouldHaveLength, 1)
							So(m[0], ShouldContainKey, "user_id")
							So(m[0], ShouldContainKey, "first_name")
							So(m[0], ShouldContainKey, "last_name")
							So(m[0]["user_id"], ShouldEqual, identity["id"])
							So(m[0]["first_name"], ShouldEqual, "John")
							So(m[0]["last_name"], ShouldEqual, "Doe")
						})
					})

					Convey("With local session", func() {
						Convey("Should return error if session does not exists locally and in session registry", func() {
							ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
							ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", "invalid"))
							objs := util.MustStringify([]map[string]interface{}{
								{"owner_id": identity["id"].(string), "key": "mykey", "value": "myval"},
							})
							_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs})
							So(err, ShouldNotBeNil)
							m, err := util.JSONToMap(err.Error())
							So(err, ShouldBeNil)
							errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
							So(errs, ShouldHaveLength, 1)
							So(errs[0], ShouldResemble, map[string]interface{}{
								"status":  "404",
								"message": "session not found",
							})
						})

						Convey("Should return error when using a local, unregistered session not belonging to the authenticated identity/caller", func() {
							sessionID := util.UUID4()
							rpc.dbSession.CreateUnregisteredSession(sessionID, identity2["id"].(string))
							ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
							ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID))
							objs := util.MustStringify([]map[string]interface{}{
								{"owner_id": identity["id"].(string), "key": "mykey", "value": "myval"},
							})
							_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs})
							So(err, ShouldNotBeNil)
							m, err := util.JSONToMap(err.Error())
							So(err, ShouldBeNil)
							errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
							So(errs, ShouldHaveLength, 1)
							So(errs[0], ShouldResemble, map[string]interface{}{
								"status":  "401",
								"message": "permission denied: you don't have permission to perform this operation",
							})
						})

						Convey("Should successfully create object using a local, unregistered session", func() {
							sessionID := util.UUID4()
							rpc.dbSession.CreateUnregisteredSession(sessionID, identity["id"].(string))
							So(err, ShouldBeNil)
							ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
							ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID))
							objs := util.MustStringify([]map[string]interface{}{
								{"owner_id": identity["id"].(string), "key": "mykey", "value": "myval"},
							})
							resp, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs})
							So(err, ShouldBeNil)
							var m []map[string]interface{}
							err = util.FromJSON(resp.Objects, &m)
							So(err, ShouldBeNil)
							So(m, ShouldHaveLength, 1)
							So(m[0]["bucket"], ShouldEqual, b.Name)
							So(m[0]["creator_id"], ShouldEqual, identity["id"])
							So(m[0]["owner_id"], ShouldEqual, identity["id"])
							rpc.dbSession.CommitEnd(sessionID)
						})
					})
				})
			})

			Convey(".GetObjects", func() {

				ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
				objs := util.MustStringify([]map[string]interface{}{
					{"owner_id": identity["id"].(string), "key": "mykey", "value": "myval"},
					{"owner_id": identity["id"].(string), "key": "mykey2", "value": "myval2"},
				})
				_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs})
				So(err, ShouldBeNil)

				Convey("Should return error if bucket name is not provided", func() {
					ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
					_, err := rpc.GetObjects(ctx, &proto_rpc.GetObjectMsg{})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "400",
						"field":   "bucket",
						"message": "bucket name is required",
					})
				})

				Convey("Should return error if authenticated user is not the owner of the queried object", func() {
					_, err := rpc.GetObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: b.Name, Owner: "unknown"})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "401",
						"message": "permission denied: you are not authorized to access objects belonging to the owner",
					})
				})

				Convey("Should return error if query parsing error occurred", func() {
					q := []byte(`{"unknown_field": "value"}`)
					_, err := rpc.GetObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: b.Name, Query: q, Owner: identity["id"].(string)})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"message": "unknown query field: unknown_field",
						"status":  "400",
						"code":    "invalid_parameter",
						"field":   "query",
					})
				})

				Convey("Should successfully fetch objects", func() {
					resp, err := rpc.GetObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: b.Name, Owner: identity["id"].(string)})
					So(err, ShouldBeNil)
					var m []map[string]interface{}
					err = util.FromJSON(resp.Objects, &m)
					So(err, ShouldBeNil)
					So(m, ShouldHaveLength, 2)
				})

				Convey("With mapping", func() {
					ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
					mappingName := util.RandString(5)
					mapping := `{ "user_id": "owner_id", "first_name": "key", "last_name": "value" }`
					_, err := rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{Name: mappingName, Mapping: []byte(mapping)})
					So(err, ShouldBeNil)

					Convey("Should return error if mapping does not exist", func() {
						_, err := rpc.GetObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: b.Name, Owner: identity["id"].(string), Mapping: "unknown"})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"status":  "404",
							"code":    "invalid_parameter",
							"field":   "mapping",
							"message": "mapping not found",
						})
					})

					Convey("Should successfully fetch objects with a mapping applied", func() {
						resp, err := rpc.GetObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: b.Name, Owner: identity["id"].(string), Mapping: mappingName})
						So(err, ShouldBeNil)
						var m []map[string]interface{}
						err = util.FromJSON(resp.Objects, &m)
						So(err, ShouldBeNil)
						So(m, ShouldHaveLength, 2)
						So(m[0], ShouldContainKey, "user_id")
						So(m[0], ShouldContainKey, "first_name")
						So(m[0], ShouldContainKey, "last_name")
						So(m[1], ShouldContainKey, "user_id")
						So(m[1], ShouldContainKey, "first_name")
						So(m[1], ShouldContainKey, "last_name")
					})

				})
			})

			Convey(".CountObjects", func() {
				ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
				objs := util.MustStringify([]map[string]interface{}{
					{"owner_id": identity["id"].(string), "key": "mykey", "value": "myval"},
					{"owner_id": identity["id"].(string), "key": "mykey2", "value": "myval2"},
				})
				_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs})
				So(err, ShouldBeNil)

				Convey("Should return error if authenticated user is not authorized to access objects belonging to the owner", func() {
					_, err := rpc.CountObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: b.Name, Owner: identity2["id"].(string)})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "401",
						"message": "permission denied: you are not authorized to access objects belonging to the owner",
					})
				})

				Convey("Should successfully return expected count", func() {
					resp, err := rpc.CountObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: b.Name, Owner: identity["id"].(string)})
					So(err, ShouldBeNil)
					So(resp.Count, ShouldEqual, 2)
				})
			})
		})

		Convey("session.go", func() {

			c1 := &proto_rpc.CreateIdentityMsg{FirstName: "john", LastName: "Doe", Email: "user@example.com", Password: "something"}
			resp, err := rpc.CreateIdentity(context.Background(), c1)
			So(err, ShouldBeNil)
			var identity map[string]interface{}
			util.FromJSON(resp.Object, &identity)

			Convey(".CreateSession", func() {

				Convey("Should return error if session id is not a UUID4", func() {
					ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
					_, err := rpc.CreateSession(ctx, &proto_rpc.Session{ID: "invalid_id"})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"message": "id is invalid. Expected UUIDv4 value",
						"status":  "400",
					})
				})

				Convey("Should successfully create a session with an explicitly provided ID", func() {
					ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
					id := util.UUID4()
					resp, err := rpc.CreateSession(ctx, &proto_rpc.Session{ID: id})
					So(err, ShouldBeNil)
					So(resp.ID, ShouldEqual, id)
				})

				Convey("Should successfully create a session with an implicitly assigned ID when ID is not provide", func() {
					ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
					resp, err := rpc.CreateSession(ctx, &proto_rpc.Session{})
					So(err, ShouldBeNil)
					So(resp.ID, ShouldNotBeEmpty)
				})
			})

			Convey(".GetSession", func() {
				Convey("Should return error if session id is not provided", func() {
					_, err := rpc.GetSession(context.Background(), &proto_rpc.Session{})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "400",
						"message": "session id is required",
					})
				})

				Convey("Should return permission error if authenticated identity/caller is not the owner of the session", func() {
					sessionID := rpc.dbSession.CreateUnregisteredSession(util.UUID4(), "some_id")
					ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
					_, err := rpc.GetSession(ctx, &proto_rpc.Session{ID: sessionID})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "401",
						"message": "permission denied: you don't have permission to perform this operation",
					})
				})

				Convey("Should return error is session does not exists", func() {
					ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
					_, err := rpc.GetSession(ctx, &proto_rpc.Session{ID: "unknown"})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"message": "session not found",
						"status":  "404",
					})
				})

				Convey("Using a registered session; Should return permission error if authenticated identity/caller is not the owner of the session", func() {
					sid := util.UUID4()
					err := rpc.sessionReg.Add(session.RegItem{SID: sid,
						Address: "localhost",
						Port:    9000,
						Meta: map[string]interface{}{
							"identity": "some_identity",
						},
					})
					So(err, ShouldBeNil)
					ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
					_, err = rpc.GetSession(ctx, &proto_rpc.Session{ID: sid})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "401",
						"message": "permission denied: you don't have permission to perform this operation",
					})
				})

				Convey("With a registered session; Should successfully get a session", func() {
					sid := util.UUID4()
					err := rpc.sessionReg.Add(session.RegItem{SID: sid,
						Address: "localhost",
						Port:    9000,
						Meta: map[string]interface{}{
							"identity": identity["id"].(string),
						},
					})
					So(err, ShouldBeNil)
					ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
					s, err := rpc.GetSession(ctx, &proto_rpc.Session{ID: sid})
					So(err, ShouldBeNil)
					So(s.ID, ShouldEqual, sid)
				})
			})
		})

		Reset(func() {
			common.ClearTable(rpc.db.DB(), []string{"buckets", "objects", "identities", "mappings"}...)
		})
	})
}
