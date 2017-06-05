package servers

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"google.golang.org/grpc/metadata"

	"github.com/asaskevich/govalidator"
	"github.com/ellcrys/util"
	"github.com/jinzhu/gorm"
	"github.com/ncodes/patchain/cockroach"
	"github.com/ncodes/patchain/cockroach/tables"
	"github.com/ncodes/patchain/object"
	"github.com/ncodes/safehold/servers/common"
	"github.com/ncodes/safehold/servers/proto_rpc"
	"github.com/ncodes/safehold/session"
	"github.com/op/go-logging"
	. "github.com/smartystreets/goconvey/convey"
)

var testDB *sql.DB

var dbName = "test_" + strings.ToLower(util.RandString(5))
var conStr = "postgresql://root@localhost:26257?sslmode=disable"
var conStrWithDB = "postgresql://root@localhost:26257/" + dbName + "?sslmode=disable"

func init() {
	var err error
	testDB, err = sql.Open("postgres", conStr)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %s", err))
	}
	logging.SetLevel(logging.CRITICAL, logRPC.Module)
}

func createDb(t *testing.T) error {
	_, err := testDB.Query(fmt.Sprintf("CREATE DATABASE %s;", dbName))
	return err
}

func dropDB(t *testing.T) error {
	_, err := testDB.Query(fmt.Sprintf("DROP DATABASE %s;", dbName))
	return err
}

func clearTable(db *gorm.DB, tables ...string) error {
	_, err := db.CommonDB().Exec("TRUNCATE " + strings.Join(tables, ","))
	if err != nil {
		return err
	}
	return nil
}

func TestRPC(t *testing.T) {

	if err := createDb(t); err != nil {
		t.Fatalf("failed to create test database. %s", err)
	}
	defer dropDB(t)

	cdb := cockroach.NewDB()
	cdb.ConnectionString = conStrWithDB
	cdb.NoLogging()
	if err := cdb.Connect(0, 10); err != nil {
		t.Fatalf("failed to connect to database. %s", err)
	}

	if err := cdb.CreateTables(); err != nil {
		t.Fatalf("failed to create tables. %s", err)
	}

	rpcServer := NewRPC(cdb)

	consulReg, err := session.NewConsulRegistry()
	if err != nil {
		t.Fatalf("failed to connect consul registry. %s", err)
	}

	rpcServer.sessionReg = consulReg
	rpcServer.dbSession = session.NewSession(consulReg)
	rpcServer.dbSession.SetDB(rpcServer.db)
	rpcServer.object = object.NewObject(rpcServer.db)

	Convey("TestRPC", t, func() {

		Convey(".validateIdentity", func() {
			errs := validateIdentity(&proto_rpc.CreateIdentityMsg{})
			So(errs, ShouldContain, common.Error{Code: "invalid_parameter", Message: "Email is required", Field: "email"})
			So(errs, ShouldContain, common.Error{Code: "invalid_parameter", Message: "Email is not valid", Field: "email"})
			So(errs, ShouldContain, common.Error{Code: "invalid_parameter", Message: "Password is required", Field: "password"})
			So(errs, ShouldContain, common.Error{Code: "invalid_parameter", Message: "Password is too short. Must be at least 6 characters long", Field: "password"})
			errs = validateIdentity(&proto_rpc.CreateIdentityMsg{Email: "email@gmail.com", Password: "some_password"})
			So(errs, ShouldBeEmpty)
		})

		Convey(".CreateIdentity", func() {
			Convey("Should return error if system identity does not exist", func() {
				newIdentity := &proto_rpc.CreateIdentityMsg{
					Email:    fmt.Sprintf("%s@email.com", util.RandString(5)),
					Password: "some_password",
				}
				_, err := rpcServer.CreateIdentity(context.Background(), newIdentity)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, `{"Message":"system identity not found","Errors":{"errors":[{"status":"400","message":"system identity not found"}]},"StatusCode":400}`)
			})

			Convey("Should successfully create an identity", func() {
				err := rpcServer.createSystemResources()
				So(err, ShouldBeNil)
				newIdentity := &proto_rpc.CreateIdentityMsg{
					Email:    fmt.Sprintf("%s@email.com", util.RandString(5)),
					Password: "some_password",
				}
				resp, err := rpcServer.CreateIdentity(context.Background(), newIdentity)
				So(err, ShouldBeNil)

				var o map[string]interface{}
				err = util.FromJSON(resp.Object, &o)
				So(err, ShouldBeNil)

				So(o["id"], ShouldNotBeEmpty)
				So(o["key"], ShouldEqual, object.MakeIdentityKey(newIdentity.Email))
				So(o["creator_id"], ShouldNotBeEmpty)
				So(o["owner_id"], ShouldNotBeEmpty)
				So(o["hash"], ShouldNotBeEmpty)
				So(o["prev_hash"], ShouldNotBeEmpty)
				So(o["partition_id"], ShouldNotBeEmpty)
				So(o["protected"], ShouldEqual, true)

				Convey("Should not error if existing identity is not confirmed", func() {
					_, err := rpcServer.CreateIdentity(context.Background(), newIdentity)
					So(err, ShouldBeNil)
				})

				Convey("Should return error if existing identity has been confirmed", func() {
					systemIdentity, err := rpcServer.getSystemIdentity()
					So(err, ShouldBeNil)

					email := "email_xyz@e.com"
					identity := object.MakeIdentityObject(systemIdentity.ID, systemIdentity.CreatorID, email, "some_pass", false)
					identity.Ref3 = "confirmed"
					err = cdb.Create(identity)
					So(err, ShouldBeNil)

					newIdentity := &proto_rpc.CreateIdentityMsg{Email: email, Password: "some_password"}
					_, err = rpcServer.CreateIdentity(context.Background(), newIdentity)
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, `{"Message":"Email is not available","Errors":{"errors":[{"status":"400","code":"used_email","field":"email","message":"Email is not available"}]},"StatusCode":400}`)
				})
			})
		})

		Convey(".GetIdentity", func() {
			Convey("Should return error if identity is not found", func() {
				identity, err := rpcServer.GetIdentity(context.Background(), &proto_rpc.GetIdentityMsg{
					ID: "some_unknown_id",
				})
				So(identity, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, `{"Message":"Identity not found","Errors":{"errors":[{"status":"404","message":"Identity not found"}]},"StatusCode":404}`)
			})

			Convey("Should successfully return an existing identity", func() {
				systemIdentity, err := rpcServer.getSystemIdentity()
				So(err, ShouldBeNil)

				email := util.RandString(5) + "@e.com"
				identity := object.MakeIdentityObject(systemIdentity.ID, systemIdentity.ID, email, "some_pass", false).Init()
				err = cdb.Create(identity)
				So(err, ShouldBeNil)

				resp, err := rpcServer.GetIdentity(context.Background(), &proto_rpc.GetIdentityMsg{ID: identity.ID})
				So(err, ShouldBeNil)

				var o map[string]interface{}
				util.FromJSON(resp.Object, &o)

				So(o["id"], ShouldEqual, identity.ID)
				So(o["key"], ShouldEqual, object.MakeIdentityKey(email))
				So(o["creator_id"], ShouldNotBeEmpty)
				So(o["owner_id"], ShouldNotBeEmpty)
				So(o["prev_hash"], ShouldNotBeEmpty)
				So(o["protected"], ShouldEqual, false)
			})
		})

		Convey(".validateObjects", func() {

			owner := object.MakeIdentityObject(util.RandString(10), util.RandString(10), "e@email.com", "password", true)
			err := cdb.Create(owner)
			So(err, ShouldBeNil)

			Convey("Should return error if required fields are not provided", func() {
				t := []map[string]interface{}{
					{},
				}
				errs, err := rpcServer.validateObjects(t, nil)
				So(err, ShouldBeNil)
				So(len(errs), ShouldEqual, 2)
				So(errs[0].Message, ShouldEqual, `object 0: owner_id is required`)
				So(errs[1].Message, ShouldEqual, `object 0: key is required`)
			})

			Convey("Should return error if a field has invalid type", func() {
				t := []map[string]interface{}{
					{"owner_id": 1000},
				}
				errs, err := rpcServer.validateObjects(t, nil)
				So(err, ShouldBeNil)
				So(len(errs), ShouldEqual, 2)
				So(errs[0].Message, ShouldEqual, `object 0: owner_id must be a string`)
				So(errs[1].Message, ShouldEqual, `object 0: key is required`)
			})

			Convey("Should return error if required fields are not provided and uses the mapping field in error message", func() {
				t := []map[string]interface{}{
					{},
				}
				mapping := map[string]string{
					"user_id":   "owner_id",
					"recordKey": "key",
				}
				errs, err := rpcServer.validateObjects(t, mapping)
				So(err, ShouldBeNil)
				So(len(errs), ShouldEqual, 2)
				So(errs[0].Message, ShouldEqual, `object 0: user_id is required`)
				So(errs[1].Message, ShouldEqual, `object 0: recordKey is required`)
			})

			Convey("Should return error if owner of objects does not exists", func() {
				t := []map[string]interface{}{
					{"owner_id": "abc", "key": "some_key"},
				}
				errs, err := rpcServer.validateObjects(t, nil)
				So(err, ShouldBeNil)
				So(len(errs), ShouldEqual, 1)
				So(errs[0].Message, ShouldEqual, `owner of object(s) does not exist`)
			})

			Convey("Should return error if objects do not share the same owner id", func() {
				t := []map[string]interface{}{
					{"owner_id": owner.ID, "key": "some_key"},
					{"owner_id": "different_id", "key": "some_key"},
				}
				errs, err := rpcServer.validateObjects(t, nil)
				So(err, ShouldBeNil)
				So(len(errs), ShouldEqual, 1)
				So(errs[0].Message, ShouldEqual, `object 1: all objects must share the same owner_id`)
			})

			Convey("Should return error if object key starts with '$' character", func() {
				t := []map[string]interface{}{
					{"owner_id": owner.ID, "key": "$some_key"},
				}
				errs, err := rpcServer.validateObjects(t, nil)
				So(err, ShouldBeNil)
				So(len(errs), ShouldEqual, 1)
				So(errs[0].Message, ShouldEqual, `object 0: key cannot start with a '$' character`)
			})

			Convey("Should successfully validate objects", func() {
				t := []map[string]interface{}{
					{"owner_id": owner.ID, "key": "some_key"},
				}
				errs, err := rpcServer.validateObjects(t, nil)
				So(err, ShouldBeNil)
				So(len(errs), ShouldEqual, 0)
			})
		})

		Convey(".CreateDBSession", func() {
			Convey("Should successfully create a session", func() {
				ctx := context.WithValue(context.Background(), CtxIdentity, "some_id")
				session := &proto_rpc.DBSession{}

				req, err := rpcServer.CreateDBSession(ctx, session)
				So(err, ShouldBeNil)
				So(req.ID, ShouldNotBeEmpty)
				So(govalidator.IsUUIDv4(req.ID), ShouldBeTrue)
				So(rpcServer.dbSession.HasSession(req.ID), ShouldBeTrue)

				item, err := consulReg.Get(req.ID)
				So(err, ShouldBeNil)
				So(item, ShouldNotBeNil)
				So(item.SID, ShouldEqual, req.ID)
				rpcServer.dbSession.End(req.ID)
			})
		})

		Convey(".GetDBSession", func() {
			Convey("Should return `session not found` error if session does not exist", func() {
				ctx := context.WithValue(context.Background(), CtxIdentity, "some_id")
				session := &proto_rpc.DBSession{
					ID: "unknown_id",
				}
				res, err := rpcServer.GetDBSession(ctx, session)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, `{"Message":"session not found","Errors":{"errors":[{"status":"404","message":"session not found"}]},"StatusCode":404}`)
				So(res, ShouldBeNil)
			})

			Convey("Should successfully get a session", func() {
				ctx := context.WithValue(context.Background(), CtxIdentity, "some_id")
				created, err := rpcServer.CreateDBSession(ctx, &proto_rpc.DBSession{})
				So(err, ShouldBeNil)
				So(created.ID, ShouldNotBeEmpty)

				existingSession, err := rpcServer.GetDBSession(ctx, &proto_rpc.DBSession{ID: created.ID})
				So(err, ShouldBeNil)
				So(existingSession.ID, ShouldEqual, created.ID)

				rpcServer.dbSession.End(existingSession.ID)
			})
		})

		Convey(".DeleteDBSession", func() {
			Convey("Should return error if session id is not provided", func() {
				ctx := context.WithValue(context.Background(), CtxIdentity, "some_id")
				_, err := rpcServer.DeleteDBSession(ctx, &proto_rpc.DBSession{})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, `{"Message":"session id is required","Errors":{"errors":[{"status":"400","message":"session id is required"}]},"StatusCode":400}`)
			})

			Convey("Should return error if session does not exist locally and on the session registry", func() {
				ctx := context.WithValue(context.Background(), CtxIdentity, "some_id")
				_, err := rpcServer.DeleteDBSession(ctx, &proto_rpc.DBSession{ID: "unknown"})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, `{"Message":"session not found","Errors":{"errors":[{"status":"404","message":"session not found"}]},"StatusCode":404}`)
			})
		})

		Convey(".orderByToString", func() {
			aCase := []*proto_rpc.OrderBy{
				{Field: "a_field", Order: 0},
				{Field: "b_field", Order: 1},
				{Field: "c_field", Order: 0},
			}
			actual := orderByToString(aCase)
			So(actual, ShouldEqual, "a_field desc, b_field asc, c_field desc")
		})

		Convey(".GetObjects", func() {
			Convey("With no session id", func() {
				Convey("Should return error if owner does not exist", func() {
					ctx := context.WithValue(context.Background(), CtxIdentity, "some_id")
					_, err := rpcServer.GetObjects(ctx, &proto_rpc.GetObjectMsg{
						Owner: "unknown",
					})
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, `{"Message":"owner not found","Errors":{"errors":[{"status":"404","message":"owner not found"}]},"StatusCode":404}`)
				})

				Convey("Should return error if creator does not exist", func() {
					ctx := context.WithValue(context.Background(), CtxIdentity, "some_id")
					_, err := rpcServer.GetObjects(ctx, &proto_rpc.GetObjectMsg{
						Creator: "unknown",
					})
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, `{"Message":"creator not found","Errors":{"errors":[{"status":"404","message":"creator not found"}]},"StatusCode":404}`)
				})

				Convey("Should return error if owner is not the same as the developer identity", func() {

					ownerID := util.RandString(10)
					owner := object.MakeIdentityObject(ownerID, ownerID, util.RandString(5)+"@email.com", "password", true)
					err := cdb.Create(owner)
					So(err, ShouldBeNil)

					ctx := context.WithValue(context.Background(), CtxIdentity, "developer_id")
					_, err = rpcServer.GetObjects(ctx, &proto_rpc.GetObjectMsg{
						Owner: owner.ID,
					})
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, `{"Message":"permission denied: you are not authorized to access objects for the owner","Errors":{"errors":[{"status":"401","message":"permission denied: you are not authorized to access objects for the owner"}]},"StatusCode":401}`)
				})

				Convey("Should successfully get an object", func() {

					ownerID := util.RandString(10)
					owner := object.MakeIdentityObject(ownerID, ownerID, util.RandString(5)+"@email.com", "password", true)
					err := cdb.Create(owner)
					So(err, ShouldBeNil)

					obj := &tables.Object{
						OwnerID:   owner.ID,
						CreatorID: owner.ID,
						Key:       "key_abc",
						Value:     "abcdef",
						Ref1:      "something_abc",
					}
					obj.Init().ComputeHash()
					err = cdb.Create(obj)
					So(err, ShouldBeNil)

					ctx := context.WithValue(context.Background(), CtxIdentity, owner.ID)
					resp, err := rpcServer.GetObjects(ctx, &proto_rpc.GetObjectMsg{
						Query: util.MustStringify(map[string]interface{}{"key": obj.Key}),
						Owner: owner.ID,
					})
					So(err, ShouldBeNil)

					var objs []map[string]interface{}
					err = util.FromJSON(resp.Objects, &objs)
					So(err, ShouldBeNil)

					So(len(objs), ShouldEqual, 1)
					So(objs[0]["id"], ShouldEqual, obj.ID)

					Convey("With mapping", func() {
						mapName := util.RandString(5)
						mapObj := object.MakeMappingObject(owner.ID, mapName, `{ "my_key": "key", "my_val": "value", "auth_id": "ref1" }`)
						err = cdb.Create(mapObj)
						So(err, ShouldBeNil)

						ctx := context.WithValue(context.Background(), CtxIdentity, owner.ID)
						resp, err := rpcServer.GetObjects(ctx, &proto_rpc.GetObjectMsg{
							Query:   util.MustStringify(map[string]interface{}{"key": obj.Key}),
							Owner:   owner.ID,
							Mapping: mapName,
						})
						So(err, ShouldBeNil)

						var objs []map[string]interface{}
						err = util.FromJSON(resp.Objects, &objs)
						So(err, ShouldBeNil)
						So(len(objs), ShouldEqual, 1)
						So(objs[0]["my_key"], ShouldEqual, obj.Key)
						So(objs[0]["my_val"], ShouldEqual, obj.Value)
						So(objs[0]["auth_id"], ShouldEqual, obj.Ref1)
					})
				})
			})
		})

		Convey(".CountObjects", func() {
			Convey("With no session id", func() {
				Convey("Should return error if owner does not exist", func() {
					ctx := context.WithValue(context.Background(), CtxIdentity, "some_id")
					_, err := rpcServer.CountObjects(ctx, &proto_rpc.GetObjectMsg{
						Owner: "unknown",
					})
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, `{"Message":"owner not found","Errors":{"errors":[{"status":"404","message":"owner not found"}]},"StatusCode":404}`)
				})

				Convey("Should return error if creator does not exist", func() {
					ctx := context.WithValue(context.Background(), CtxIdentity, "some_id")
					_, err := rpcServer.CountObjects(ctx, &proto_rpc.GetObjectMsg{
						Creator: "unknown",
					})
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, `{"Message":"creator not found","Errors":{"errors":[{"status":"404","message":"creator not found"}]},"StatusCode":404}`)
				})

				Convey("Should return error if owner is not the same as the developer identity", func() {

					ownerID := util.RandString(10)
					owner := object.MakeIdentityObject(ownerID, ownerID, util.RandString(5)+"@email.com", "password", true)
					err := cdb.Create(owner)
					So(err, ShouldBeNil)

					ctx := context.WithValue(context.Background(), CtxIdentity, "developer_id")
					_, err = rpcServer.CountObjects(ctx, &proto_rpc.GetObjectMsg{
						Owner: owner.ID,
					})
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, `{"Message":"permission denied: you are not authorized to access objects for the owner","Errors":{"errors":[{"status":"401","message":"permission denied: you are not authorized to access objects for the owner"}]},"StatusCode":401}`)
				})

				Convey("Should successfully get an object", func() {

					ownerID := util.RandString(10)
					owner := object.MakeIdentityObject(ownerID, ownerID, util.RandString(5)+"@email.com", "password", true)
					err := cdb.Create(owner)
					So(err, ShouldBeNil)

					obj := &tables.Object{OwnerID: owner.ID, CreatorID: owner.ID, Key: "a_key"}
					obj.Init().ComputeHash()
					err = cdb.Create(obj)
					So(err, ShouldBeNil)
					ctx := context.WithValue(context.Background(), CtxIdentity, owner.ID)
					resp, err := rpcServer.CountObjects(ctx, &proto_rpc.GetObjectMsg{
						Query: util.MustStringify(map[string]interface{}{"key": "a_key"}),
						Owner: owner.ID,
					})
					So(err, ShouldBeNil)
					So(resp.Count, ShouldEqual, 1)
				})
			})
		})

		Convey("Mappings", func() {
			Convey(".validateMapping", func() {
				Convey("Case 1: returns error if no mapping is provided", func() {
					mapping := map[string]string{}
					errs := validateMapping(mapping)
					So(errs, ShouldNotBeEmpty)
					So(errs[0].Field, ShouldEqual, "mapping")
					So(errs[0].Message, ShouldEqual, "requires at least one mapped field")
				})

				Convey("Case 2: returns error if custom field is mapped to an unknown column", func() {
					mapping := map[string]string{
						"custom_field": "unknown_column",
					}
					errs := validateMapping(mapping)
					So(errs, ShouldNotBeEmpty)
					So(errs[0].Field, ShouldEqual, "custom_field")
					So(errs[0].Message, ShouldEqual, "column name 'unknown_column' is unknown")
				})
			})

			Convey(".CreateMapping", func() {
				Convey("Should return error if mapping json is malformed", func() {
					ctx := context.WithValue(context.Background(), CtxIdentity, "some_id")
					_, err := rpcServer.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{
						Name:    "abc",
						Mapping: []byte("{"),
					})
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, `{"Message":"malformed mapping. Expected json object.","Errors":{"errors":[{"status":"404","message":"malformed mapping. Expected json object."}]},"StatusCode":404}`)
				})

				Convey("Should return error if mapping json failed validation", func() {
					ctx := context.WithValue(context.Background(), CtxIdentity, "some_id")
					_, err := rpcServer.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{
						Name:    "abc",
						Mapping: []byte(`{}`),
					})
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, `{"Message":"mapping error","Errors":{"errors":[{"field":"mapping","message":"requires at least one mapped field"}]},"StatusCode":400}`)
				})

				Convey("Should successfully create a mapping", func() {
					err := rpcServer.createSystemResources()
					So(err, ShouldBeNil)
					identity, err := rpcServer.getSystemIdentity()
					So(err, ShouldBeNil)

					ctx := context.WithValue(context.Background(), CtxIdentity, identity.ID)
					resp, err := rpcServer.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{
						Name:    "map1",
						Mapping: []byte(`{ "custom_name": "ref1" }`),
					})
					So(err, ShouldBeNil)
					So(resp.Name, ShouldEqual, "map1")
					So(resp.ID, ShouldNotBeEmpty)

					Convey(".GetMapping", func() {
						Convey("Should return error if mapping name is not provided", func() {
							ctx := context.WithValue(context.Background(), CtxIdentity, "some_id")
							_, err := rpcServer.GetMapping(ctx, &proto_rpc.GetMappingMsg{})
							So(err, ShouldNotBeNil)
							So(err.Error(), ShouldEqual, `{"Message":"name is required","Errors":{"errors":[{"status":"400","field":"id","message":"name is required"}]},"StatusCode":400}`)
						})

						Convey("Should return error if a mapping with matching name does not exist", func() {
							ctx := context.WithValue(context.Background(), CtxIdentity, identity.ID)
							_, err = rpcServer.GetMapping(ctx, &proto_rpc.GetMappingMsg{Name: "unknown_name"})
							So(err, ShouldNotBeNil)
							So(err.Error(), ShouldEqual, `{"Message":"mapping not found","Errors":{"errors":[{"status":"404","message":"mapping not found"}]},"StatusCode":404}`)
						})

						Convey("Should successfully return an existing mapping", func() {
							ctx := context.WithValue(context.Background(), CtxIdentity, identity.ID)
							resp, err := rpcServer.GetMapping(ctx, &proto_rpc.GetMappingMsg{Name: "map1"})
							So(err, ShouldBeNil)

							var mapping map[string]interface{}
							err = util.FromJSON(resp.Mapping, &mapping)
							So(err, ShouldBeNil)

							So(mapping["value"], ShouldResemble, `{ "custom_name": "ref1" }`)
						})
					})

					Convey(".GetMappings", func() {
						ctx := context.WithValue(context.Background(), CtxIdentity, identity.ID)
						resp, err := rpcServer.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{
							Name:    "map2",
							Mapping: []byte(`{ "custom_name": "ref2" }`),
						})
						So(err, ShouldBeNil)
						So(resp.Name, ShouldEqual, "map2")

						ctx = context.WithValue(context.Background(), CtxIdentity, identity.ID)
						resp, err = rpcServer.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{
							Name:    "map2",
							Mapping: []byte(`{ "custom_name": "ref2" }`),
						})
						So(err, ShouldBeNil)
						So(resp.Name, ShouldEqual, "map2")

						Convey("Should return all mappings belonging to the logged in identity", func() {
							ctx := context.WithValue(context.Background(), CtxIdentity, identity.ID)
							resp, err := rpcServer.GetAllMapping(ctx, &proto_rpc.GetAllMappingMsg{})
							So(err, ShouldBeNil)

							var mappings []map[string]interface{}
							err = util.FromJSON(resp.Mappings, &mappings)
							So(err, ShouldBeNil)

							So(len(mappings), ShouldEqual, 3)
							So(mappings[0]["key"], ShouldEqual, object.MakeMappingKey("map2"))
							So(mappings[0]["value"], ShouldResemble, `{ "custom_name": "ref2" }`)
							So(mappings[1]["key"], ShouldEqual, object.MakeMappingKey("map2"))
							So(mappings[1]["value"], ShouldResemble, `{ "custom_name": "ref2" }`)
							So(mappings[2]["key"], ShouldEqual, object.MakeMappingKey("map1"))
							So(mappings[2]["value"], ShouldResemble, `{ "custom_name": "ref1" }`)
						})

						Convey("Should return all mappings matching a specific name", func() {
							ctx := context.WithValue(context.Background(), CtxIdentity, identity.ID)
							resp, err := rpcServer.GetAllMapping(ctx, &proto_rpc.GetAllMappingMsg{
								Name: "map2",
							})
							So(err, ShouldBeNil)

							var mappings []map[string]interface{}
							err = util.FromJSON(resp.Mappings, &mappings)
							So(err, ShouldBeNil)

							So(len(mappings), ShouldEqual, 2)
							So(mappings[0]["key"], ShouldEqual, object.MakeMappingKey("map2"))
							So(mappings[0]["value"], ShouldResemble, `{ "custom_name": "ref2" }`)
							So(mappings[1]["key"], ShouldEqual, object.MakeMappingKey("map2"))
							So(mappings[1]["value"], ShouldResemble, `{ "custom_name": "ref2" }`)
						})

						Convey("Should return all mappings matching a specific name and limit by 1", func() {
							ctx := context.WithValue(context.Background(), CtxIdentity, identity.ID)
							resp, err := rpcServer.GetAllMapping(ctx, &proto_rpc.GetAllMappingMsg{
								Name:  "map2",
								Limit: 1,
							})
							So(err, ShouldBeNil)

							var mappings []map[string]interface{}
							err = util.FromJSON(resp.Mappings, &mappings)
							So(err, ShouldBeNil)

							So(len(mappings), ShouldEqual, 1)
							So(mappings[0]["key"], ShouldEqual, object.MakeMappingKey("map2"))
							So(mappings[0]["value"], ShouldResemble, `{ "custom_name": "ref2" }`)
						})

					})

					Reset(func() {
						clearTable(cdb.GetConn().(*gorm.DB), "objects")
					})
				})
			})
		})

		Convey(".CreateObject", func() {

			Convey("Should return parsing error if no object data is provided", func() {
				ctx := context.WithValue(context.Background(), CtxIdentity, "some_id")
				req := &proto_rpc.CreateObjectsMsg{}
				_, err := rpcServer.CreateObjects(ctx, req)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, `{"Message":"failed to parse objects","Errors":{"errors":[{"status":"400","message":"failed to parse objects"}]},"StatusCode":400}`)
			})

			Convey("Should return error if no object is provided", func() {
				ctx := context.WithValue(context.Background(), CtxIdentity, "some_id")
				req := &proto_rpc.CreateObjectsMsg{
					Objects: []byte("[]"),
				}
				_, err := rpcServer.CreateObjects(ctx, req)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, `{"Message":"validation errors","Errors":{"errors":[{"message":"no object provided. At least one object is required"}]},"StatusCode":400}`)
			})

			Convey("Should return error if unable to parse json encoded Objects", func() {
				ctx := context.WithValue(context.Background(), CtxIdentity, "some_id")
				req := &proto_rpc.CreateObjectsMsg{
					Objects: []byte("[{ "),
				}
				_, err := rpcServer.CreateObjects(ctx, req)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, `{"Message":"failed to parse objects","Errors":{"errors":[{"status":"400","message":"failed to parse objects"}]},"StatusCode":400}`)
			})

			Convey("Should return error objects in request exceeds MaxObjectPerPut", func() {
				MaxObjectPerPut = 2
				req := &proto_rpc.CreateObjectsMsg{
					Objects: util.MustStringify([]*proto_rpc.Object{
						{ID: "some_id"},
						{ID: "some_id2"},
						{ID: "some_id3"},
					}),
				}
				ctx := context.WithValue(context.Background(), CtxIdentity, "some_id")
				_, err := rpcServer.CreateObjects(ctx, req)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, `{"Message":"validation errors","Errors":{"errors":[{"message":"too many objects. Only a maximum of 2 can be created at once"}]},"StatusCode":400}`)
			})

			Convey("With valid owner identity", func() {
				ownerID := util.RandString(10)
				owner := object.MakeIdentityObject(ownerID, ownerID, "e@email.com", "password", true)
				err := cdb.Create(owner)
				So(err, ShouldBeNil)

				Convey("Should return error if developer identity is not permitted to PUT object", func() {
					ctx := context.WithValue(context.Background(), CtxIdentity, "some_id")
					req := &proto_rpc.CreateObjectsMsg{
						Objects: util.MustStringify([]*proto_rpc.Object{
							{ID: "some_id", OwnerID: owner.ID, Key: "some_key"},
						}),
					}
					_, err := rpcServer.CreateObjects(ctx, req)
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, `{"Message":"permission denied: you are not authorized to create objects for the owner","Errors":{"errors":[{"status":"401","message":"permission denied: you are not authorized to create objects for the owner"}]},"StatusCode":401}`)
				})

				Convey("Should return error if owner of object has no partition", func() {
					ctx := context.WithValue(context.Background(), CtxIdentity, owner.ID)
					req := &proto_rpc.CreateObjectsMsg{
						Objects: util.MustStringify([]*proto_rpc.Object{
							{ID: "some_id", OwnerID: owner.ID, Key: "some_key"},
						}),
					}
					_, err := rpcServer.CreateObjects(ctx, req)
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, `{"Message":"failed to put object(s): owner has no partition","Errors":{"errors":[{"status":"500","code":"put_error","field":"objects","message":"failed to put object(s): owner has no partition"}]},"StatusCode":500}`)
				})

				Convey("Should successfully create object belonging to the authenticated identity(developer)", func() {
					partitions, err := object.NewObject(cdb).CreatePartitions(1, owner.ID, owner.ID)
					So(err, ShouldBeNil)
					So(len(partitions), ShouldEqual, 1)

					ctx := context.WithValue(context.Background(), CtxIdentity, owner.ID)
					req := &proto_rpc.CreateObjectsMsg{
						Objects: util.MustStringify([]*proto_rpc.Object{
							{ID: "some_id", OwnerID: owner.ID, Key: "some_key"},
						}),
					}
					resp, err := rpcServer.CreateObjects(ctx, req)
					So(err, ShouldBeNil)

					var objs []map[string]interface{}
					err = util.FromJSON(resp.Objects, &objs)
					So(err, ShouldBeNil)

					So(len(objs), ShouldEqual, 1)
					So(objs[0]["owner_id"], ShouldEqual, owner.ID)
				})

				Convey("Should successfully create object belonging to the authenticated identity(developer) even when the owner_id of the object is not set", func() {
					partitions, err := object.NewObject(cdb).CreatePartitions(1, owner.ID, owner.ID)
					So(err, ShouldBeNil)
					So(len(partitions), ShouldEqual, 1)

					ctx := context.WithValue(context.Background(), CtxIdentity, owner.ID)
					req := &proto_rpc.CreateObjectsMsg{
						Objects: util.MustStringify([]*proto_rpc.Object{
							{ID: "some_id2", Key: "some_key"},
						}),
					}
					resp, err := rpcServer.CreateObjects(ctx, req)
					So(err, ShouldBeNil)

					var objs []map[string]interface{}
					err = util.FromJSON(resp.Objects, &objs)
					So(err, ShouldBeNil)

					So(len(objs), ShouldEqual, 1)
					So(objs[0]["owner_id"], ShouldEqual, owner.ID)
				})

				Convey("With mapping", func() {
					mapName := util.RandString(5)
					mapObj := object.MakeMappingObject(owner.ID, mapName, `{ "my_key": "key", "my_val": "value" }`)
					err = cdb.Create(mapObj)
					So(err, ShouldBeNil)

					Convey("Should successfully unmap object, create the object and reapply mapping to the returned object", func() {
						partitions, err := object.NewObject(cdb).CreatePartitions(1, owner.ID, owner.ID)
						So(err, ShouldBeNil)
						So(len(partitions), ShouldEqual, 1)

						ctx := context.WithValue(context.Background(), CtxIdentity, owner.ID)
						req := &proto_rpc.CreateObjectsMsg{
							Mapping: mapName,
							Objects: util.MustStringify([]map[string]interface{}{{
								"owner_id": owner.ID,
								"id":       "some_id_abc",
								"my_key":   "abcdef",
								"my_val":   "xyz",
							}}),
						}
						resp, err := rpcServer.CreateObjects(ctx, req)
						So(err, ShouldBeNil)

						var objs []map[string]interface{}
						err = util.FromJSON(resp.Objects, &objs)
						So(err, ShouldBeNil)

						So(len(objs), ShouldEqual, 1)
						So(objs[0]["owner_id"], ShouldEqual, owner.ID)
						So(objs[0]["id"], ShouldEqual, "some_id_abc")
						So(objs[0]["my_key"], ShouldEqual, "abcdef")
						So(objs[0]["my_val"], ShouldEqual, "xyz")
					})
				})

				Convey("With session id", func() {
					partitions, err := object.NewObject(cdb).CreatePartitions(1, owner.ID, owner.ID)
					So(err, ShouldBeNil)
					So(len(partitions), ShouldEqual, 1)

					Convey("Should return error if session id is not found in in-memory or session registry", func() {
						ctx := context.WithValue(context.Background(), CtxIdentity, owner.ID)
						ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", "unknown"))
						req := &proto_rpc.CreateObjectsMsg{
							Objects: util.MustStringify([]*proto_rpc.Object{
								{ID: "some_id", OwnerID: owner.ID, Key: "some_key"},
							}),
						}
						_, err := rpcServer.CreateObjects(ctx, req)
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldEqual, `{"Message":"session not found","Errors":{"errors":[{"status":"404","message":"session not found"}]},"StatusCode":404}`)
					})
				})
			})
		})

	})
}
