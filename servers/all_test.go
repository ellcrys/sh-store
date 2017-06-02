package servers

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"google.golang.org/grpc/metadata"

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
				o, err := rpcServer.CreateIdentity(context.Background(), newIdentity)
				So(err, ShouldBeNil)
				So(o.Data.Attributes.Key, ShouldEqual, object.MakeIdentityKey(newIdentity.Email))
				So(o.Data.ID, ShouldNotBeEmpty)
				So(o.Data.Type, ShouldNotBeEmpty)
				So(o.Data.Attributes.CreatorID, ShouldNotBeEmpty)
				So(o.Data.Attributes.OwnerID, ShouldNotBeEmpty)
				So(o.Data.Attributes.Hash, ShouldNotBeEmpty)
				So(o.Data.Attributes.PrevHash, ShouldNotBeEmpty)
				So(o.Data.Attributes.PartitionID, ShouldNotBeEmpty)
				So(o.Data.Attributes.Protected, ShouldEqual, true)

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
				identity := object.MakeIdentityObject(systemIdentity.ID, systemIdentity.ID, email, "some_pass", false)
				identity.Ref3 = "confirmed"
				err = cdb.Create(identity)
				So(err, ShouldBeNil)

				resp, err := rpcServer.GetIdentity(context.Background(), &proto_rpc.GetIdentityMsg{ID: identity.ID})
				So(err, ShouldBeNil)
				So(resp.Data.Attributes.ID, ShouldEqual, identity.ID)
			})
		})

		Convey(".validateObjects", func() {

			owner := object.MakeIdentityObject(util.RandString(10), util.RandString(10), "e@email.com", "password", true)
			err := cdb.Create(owner)
			So(err, ShouldBeNil)

			Convey("Should test and fail all validations", func() {
				t := []*proto_rpc.Object{
					{},
				}
				errs := rpcServer.validateObjects(t)
				So(len(errs), ShouldEqual, 1)
				So(errs[0].Message, ShouldEqual, `object 0: owner id is required`)

				t = []*proto_rpc.Object{
					{OwnerID: owner.ID},
				}
				errs = rpcServer.validateObjects(t)
				So(len(errs), ShouldEqual, 1)
				So(errs[0].Message, ShouldEqual, `object 0: key is required`)

				t = []*proto_rpc.Object{
					{OwnerID: "abc", Key: "some_key"},
				}
				errs = rpcServer.validateObjects(t)
				So(len(errs), ShouldEqual, 1)
				So(errs[0].Message, ShouldEqual, `object 0: owner id does not exist`)

				t = []*proto_rpc.Object{
					{OwnerID: owner.ID, Key: "some_key"},
				}
				errs = rpcServer.validateObjects(t)
				So(len(errs), ShouldEqual, 0)
			})
		})

		Convey(".CreateObject", func() {

			Convey("Should return error if no object is provided", func() {
				ctx := context.WithValue(context.Background(), CtxIdentity, "some_id")
				req := &proto_rpc.CreateObjectsMsg{}
				_, err := rpcServer.CreateObjects(ctx, req)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, `{"Message":"no object provided. At least one object is required","Errors":{"errors":[{"status":"400","message":"no object provided. At least one object is required"}]},"StatusCode":400}`)
			})

			Convey("Should return error objects in request exceeds MaxObjectPerPut", func() {
				MaxObjectPerPut = 2
				req := &proto_rpc.CreateObjectsMsg{
					Objects: []*proto_rpc.Object{
						{ID: "some_id"},
						{ID: "some_id2"},
						{ID: "some_id3"},
					},
				}
				ctx := context.WithValue(context.Background(), CtxIdentity, "some_id")
				_, err := rpcServer.CreateObjects(ctx, req)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, `{"Message":"too many objects. Only a maximum of 2 can be created at once","Errors":{"errors":[{"status":"400","message":"too many objects. Only a maximum of 2 can be created at once"}]},"StatusCode":400}`)
			})

			Convey("With valid owner identity", func() {
				ownerID := util.RandString(10)
				owner := object.MakeIdentityObject(ownerID, ownerID, "e@email.com", "password", true)
				err := cdb.Create(owner)
				So(err, ShouldBeNil)

				Convey("Should return error if developer identity is not permitted to PUT object", func() {
					ctx := context.WithValue(context.Background(), CtxIdentity, "some_id")
					req := &proto_rpc.CreateObjectsMsg{
						Objects: []*proto_rpc.Object{
							{ID: "some_id", OwnerID: owner.ID, Key: "some_key"},
						},
					}
					_, err := rpcServer.CreateObjects(ctx, req)
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, `{"Message":"permission denied: you are not authorized to create objects for the owner","Errors":{"errors":[{"status":"401","message":"permission denied: you are not authorized to create objects for the owner"}]},"StatusCode":401}`)
				})

				Convey("Should return error if owner of object has no partition", func() {
					ctx := context.WithValue(context.Background(), CtxIdentity, owner.ID)
					req := &proto_rpc.CreateObjectsMsg{
						Objects: []*proto_rpc.Object{
							{ID: "some_id", OwnerID: owner.ID, Key: "some_key"},
						},
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
						Objects: []*proto_rpc.Object{
							{ID: "some_id", OwnerID: owner.ID, Key: "some_key"},
						},
					}
					resp, err := rpcServer.CreateObjects(ctx, req)
					So(err, ShouldBeNil)
					So(len(resp.Data), ShouldEqual, 1)
				})

				Convey("With session id", func() {
					partitions, err := object.NewObject(cdb).CreatePartitions(1, owner.ID, owner.ID)
					So(err, ShouldBeNil)
					So(len(partitions), ShouldEqual, 1)

					Convey("Should return error if session id is not found in in-memory or session registry", func() {
						ctx := context.WithValue(context.Background(), CtxIdentity, owner.ID)
						ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", "unknown"))
						req := &proto_rpc.CreateObjectsMsg{
							Objects: []*proto_rpc.Object{
								{ID: "some_id", OwnerID: owner.ID, Key: "some_key"},
							},
						}
						_, err := rpcServer.CreateObjects(ctx, req)
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldEqual, `{"Message":"session not found","Errors":{"errors":[{"status":"404","message":"session not found"}]},"StatusCode":404}`)
					})
				})
			})
		})

		Convey(".makeDBSessionID", func() {
			Convey("Should return session id in the format: [identityID].[sessionID]", func() {
				So(makeDBSessionID("some_id", "some_id"), ShouldEqual, "some_id.some_id")
			})
		})

		Convey(".CreateDBSession", func() {
			Convey("Should successfully create a session", func() {
				ctx := context.WithValue(context.Background(), CtxIdentity, "some_id")
				session := &proto_rpc.DBSession{
					ID: util.RandString(5),
				}

				req, err := rpcServer.CreateDBSession(ctx, session)
				So(err, ShouldBeNil)
				So(req.ID, ShouldEqual, session.ID)
				fullSid := makeDBSessionID("some_id", session.ID)
				So(rpcServer.dbSession.HasSession(fullSid), ShouldEqual, true)

				item, err := consulReg.Get(fullSid)
				So(err, ShouldBeNil)
				So(item, ShouldNotBeNil)
				So(item.SID, ShouldEqual, fullSid)
				rpcServer.dbSession.End(fullSid)
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
				session := &proto_rpc.DBSession{
					ID: "some_id",
				}
				res, err := rpcServer.CreateDBSession(ctx, session)
				So(err, ShouldBeNil)
				So(res.ID, ShouldEqual, "some_id")

				res, err = rpcServer.GetDBSession(ctx, session)
				So(err, ShouldBeNil)
				So(res.ID, ShouldEqual, session.ID)

				rpcServer.dbSession.End(makeDBSessionID("some_id", "some_id"))
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

					obj := &tables.Object{OwnerID: owner.ID, CreatorID: owner.ID, Key: "a_key"}
					obj.Init().ComputeHash()
					err = cdb.Create(obj)
					So(err, ShouldBeNil)
					ctx := context.WithValue(context.Background(), CtxIdentity, owner.ID)
					resp, err := rpcServer.GetObjects(ctx, &proto_rpc.GetObjectMsg{
						Query: util.MustStringify(map[string]interface{}{"key": "a_key"}),
						Owner: owner.ID,
					})
					So(err, ShouldBeNil)
					So(len(resp.Data), ShouldEqual, 1)
					So(resp.Data[0].ID, ShouldEqual, obj.ID)
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
					mapping := map[string]interface{}{}
					errs := validateMapping(mapping)
					So(errs, ShouldNotBeEmpty)
					So(errs[0].Field, ShouldEqual, "mapping")
					So(errs[0].Message, ShouldEqual, "requires at least one mapped field")
				})

				Convey("Case 2: returns error if mapping value type is not string or map", func() {
					mapping := map[string]interface{}{
						"custom_field": 2,
					}
					errs := validateMapping(mapping)
					So(errs, ShouldNotBeEmpty)
					So(errs[0].Field, ShouldEqual, "custom_field")
					So(errs[0].Message, ShouldEqual, "invalid value type. Expected string or json type")
				})

				Convey("Case 2: returns error if custom field is mapped to an unknown column", func() {
					mapping := map[string]interface{}{
						"custom_field": "unknown_column",
					}
					errs := validateMapping(mapping)
					So(errs, ShouldNotBeEmpty)
					So(errs[0].Field, ShouldEqual, "custom_field")
					So(errs[0].Message, ShouldEqual, "column name 'unknown_column' is unknown")
				})

				Convey("Case 3: returns error if custom field with map value has no `field` property", func() {
					mapping := map[string]interface{}{
						"custom_field": map[string]interface{}{},
					}
					errs := validateMapping(mapping)
					So(errs, ShouldNotBeEmpty)
					So(errs[0].Field, ShouldEqual, "custom_field")
					So(errs[0].Message, ShouldEqual, "'field' property is required")
				})

				Convey("Case 4: returns error if `custom_field.field` has invalid value type", func() {
					mapping := map[string]interface{}{
						"custom_field": map[string]interface{}{
							"field": 1,
						},
					}
					errs := validateMapping(mapping)
					So(errs, ShouldNotBeEmpty)
					So(errs[0].Field, ShouldEqual, "custom_field.field")
					So(errs[0].Message, ShouldEqual, "invalid value type. Expected string type")
				})

				Convey("Case 5: returns error if `custom_field.protected` property has invalid value type", func() {
					mapping := map[string]interface{}{
						"custom_field": map[string]interface{}{
							"field":     "id",
							"protected": 1,
						},
					}
					errs := validateMapping(mapping)
					So(errs, ShouldNotBeEmpty)
					So(errs[0].Field, ShouldEqual, "custom_field.protected")
					So(errs[0].Message, ShouldEqual, "invalid value type. Expected boolean type")
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
							So(resp.Mapping, ShouldResemble, []byte(`{ "custom_name": "ref1" }`))
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
							So(len(resp.Mappings), ShouldEqual, 3)
							So(resp.Mappings[0].Name, ShouldEqual, object.MakeMappingKey("map2"))
							So(resp.Mappings[0].Mapping, ShouldResemble, []byte(`{ "custom_name": "ref2" }`))
							So(resp.Mappings[1].Name, ShouldEqual, object.MakeMappingKey("map2"))
							So(resp.Mappings[1].Mapping, ShouldResemble, []byte(`{ "custom_name": "ref2" }`))
							So(resp.Mappings[2].Name, ShouldEqual, object.MakeMappingKey("map1"))
							So(resp.Mappings[2].Mapping, ShouldResemble, []byte(`{ "custom_name": "ref1" }`))
						})

						Convey("Should return all mappings matching a specific name", func() {
							ctx := context.WithValue(context.Background(), CtxIdentity, identity.ID)
							resp, err := rpcServer.GetAllMapping(ctx, &proto_rpc.GetAllMappingMsg{
								Name: "map2",
							})
							So(err, ShouldBeNil)
							So(len(resp.Mappings), ShouldEqual, 2)
							So(resp.Mappings[0].Name, ShouldEqual, object.MakeMappingKey("map2"))
							So(resp.Mappings[0].Mapping, ShouldResemble, []byte(`{ "custom_name": "ref2" }`))
							So(resp.Mappings[1].Name, ShouldEqual, object.MakeMappingKey("map2"))
							So(resp.Mappings[1].Mapping, ShouldResemble, []byte(`{ "custom_name": "ref2" }`))
						})

						Convey("Should return all mappings matching a specific name and limit by 1", func() {
							ctx := context.WithValue(context.Background(), CtxIdentity, identity.ID)
							resp, err := rpcServer.GetAllMapping(ctx, &proto_rpc.GetAllMappingMsg{
								Name:  "map2",
								Limit: 1,
							})
							So(err, ShouldBeNil)
							So(len(resp.Mappings), ShouldEqual, 1)
							So(resp.Mappings[0].Name, ShouldEqual, object.MakeMappingKey("map2"))
							So(resp.Mappings[0].Mapping, ShouldResemble, []byte(`{ "custom_name": "ref2" }`))
						})

					})

					Reset(func() {
						clearTable(cdb.GetConn().(*gorm.DB), "objects")
					})
				})
			})
		})
	})
}
