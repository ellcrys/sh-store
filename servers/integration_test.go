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
				ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
				b := &proto_rpc.CreateBucketMsg{Name: "mybucket"}
				bucket, err := rpc.CreateBucket(ctx, b)
				So(err, ShouldBeNil)
				So(bucket.Name, ShouldEqual, "mybucket")
				So(bucket.ID, ShouldHaveLength, 36)

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

				// Convey("Shou")
			})

			Convey("Should return error if validation is failed", func() {
				ctx := context.WithValue(context.Background(), CtxIdentity, identity["id"])
				_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: "unknown"})
				So(err, ShouldNotBeNil)
			})
		})

		Reset(func() {
			common.ClearTable(rpc.db.DB(), []string{"buckets", "objects", "identities", "mappings"}...)
		})
	})
}
