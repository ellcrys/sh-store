package servers

import (
	"context"
	"testing"
	"time"

	"github.com/ellcrys/elldb/servers/db"
	"github.com/ellcrys/elldb/servers/oauth"
	"github.com/ellcrys/elldb/servers/proto_rpc"
	"github.com/ellcrys/elldb/session"
	"github.com/ellcrys/util"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/metadata"
)

func TestSession(t *testing.T) {
	setup(t, func(rpc, rpc2 *RPC) {
		Convey("Session", t, func() {
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

			Convey("session.go", func() {

				Convey(".CreateSession", func() {

					Convey("Should return error if session id is not a UUID4", func() {
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
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
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						id := util.UUID4()
						resp, err := rpc.CreateSession(ctx, &proto_rpc.Session{ID: id})
						So(err, ShouldBeNil)
						So(resp.ID, ShouldEqual, id)
					})

					Convey("Should successfully create a session with an implicitly assigned ID when ID is not provide", func() {
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
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

					Convey("Should return permission error if authenticated account/caller is not the owner of the session", func() {
						sessionID := rpc.dbSession.CreateUnregisteredSession(util.UUID4(), "some_id")
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
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
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
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

					Convey("Using a remote session; Should return permission error if authenticated account/caller is not the owner of the session", func() {
						sid := util.UUID4()
						err := rpc.sessionReg.Add(session.RegItem{SID: sid,
							Address: "localhost",
							Port:    9000,
							Meta: map[string]interface{}{
								"account": "some_account",
							},
						})
						So(err, ShouldBeNil)
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
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

					Convey("With a remote session; Should successfully get a session", func() {
						sid := util.UUID4()
						err := rpc.sessionReg.Add(session.RegItem{SID: sid,
							Address: "localhost",
							Port:    9000,
							Meta: map[string]interface{}{
								"account": account["id"].(string),
							},
						})
						So(err, ShouldBeNil)
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						s, err := rpc.GetSession(ctx, &proto_rpc.Session{ID: sid})
						So(err, ShouldBeNil)
						So(s.ID, ShouldEqual, sid)
					})

					Convey("With a unregistered session; Should successfully get a session", func() {
						sid := rpc.dbSession.CreateUnregisteredSession(util.UUID4(), account["id"].(string))
						So(err, ShouldBeNil)
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						s, err := rpc.GetSession(ctx, &proto_rpc.Session{ID: sid})
						So(err, ShouldBeNil)
						So(s.ID, ShouldEqual, sid)
					})
				})

				Convey(".DeleteSession", func() {
					Convey("Should return error if session id is not provided", func() {
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						_, err := rpc.DeleteSession(ctx, &proto_rpc.Session{})
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

					Convey("Should return permission error if authenticated account/caller is not the owner of the session", func() {
						sessionID := rpc.dbSession.CreateUnregisteredSession(util.UUID4(), "some_id")
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						_, err := rpc.DeleteSession(ctx, &proto_rpc.Session{ID: sessionID})
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

					Convey("With an unregistered session; Should successfully delete session", func() {
						sessionID := rpc.dbSession.CreateUnregisteredSession(util.UUID4(), account["id"].(string))
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						resp, err := rpc.DeleteSession(ctx, &proto_rpc.Session{ID: sessionID})
						So(err, ShouldBeNil)
						So(resp.ID, ShouldEqual, sessionID)
						So(rpc.dbSession.HasSession(sessionID), ShouldBeFalse)
					})

					Convey("With an remote session; Should successfully delete session", func() {
						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"id":   account["client_id"],
							"type": oauth.TokenTypeApp,
							"iat":  time.Now().Unix(),
						})

						sessionID := util.UUID4()
						rpc2.dbSession.CreateSession(sessionID, account["id"].(string))

						ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "bearer "+token))
						ctx = context.WithValue(ctx, CtxAccount, account["id"])
						resp, err := rpc.DeleteSession(ctx, &proto_rpc.Session{ID: sessionID})
						So(err, ShouldBeNil)
						So(resp.ID, ShouldEqual, sessionID)
						So(rpc.dbSession.HasSession(sessionID), ShouldBeFalse)
					})

					Convey("With remote session; Should return error if session exists in registry but not found in the node hosting the session", func() {
						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"id":   account["client_id"],
							"type": oauth.TokenTypeApp,
							"iat":  time.Now().Unix(),
						})

						sessionID := util.UUID4()
						rpc2.dbSession.CreateSession(sessionID, account["id"].(string))
						rpc2.dbSession.RemoveAgent(sessionID)

						ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "bearer "+token))
						ctx = context.WithValue(ctx, CtxAccount, account["id"])
						_, err := rpc.DeleteSession(ctx, &proto_rpc.Session{ID: sessionID})
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

					Convey("Should return error if session does not exist locally or in session registry", func() {
						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"id":   account["client_id"],
							"type": oauth.TokenTypeApp,
							"iat":  time.Now().Unix(),
						})

						sessionID := util.UUID4()

						ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "bearer "+token))
						ctx = context.WithValue(ctx, CtxAccount, account["id"])
						_, err := rpc.DeleteSession(ctx, &proto_rpc.Session{ID: sessionID})
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

					Convey("With an remote session; Should return error if authenticated account/caller is not the owner of the session", func() {
						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"id":   account["client_id"],
							"type": oauth.TokenTypeApp,
							"iat":  time.Now().Unix(),
						})

						sessionID := util.UUID4()
						rpc2.dbSession.CreateSession(sessionID, "some_id")

						ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "bearer "+token))
						ctx = context.WithValue(ctx, CtxAccount, account["id"])
						_, err := rpc.DeleteSession(ctx, &proto_rpc.Session{ID: sessionID})
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
				})

				Convey(".CommitSession", func() {
					Convey("Should return error if session id is not provide", func() {
						_, err := rpc.CommitSession(context.Background(), &proto_rpc.Session{})
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

					Convey("With local session; Should return error if authenticated account is not the owner of the session", func() {
						sessionID := rpc.dbSession.CreateUnregisteredSession(util.UUID4(), "some_id")
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						_, err := rpc.CommitSession(ctx, &proto_rpc.Session{ID: sessionID})
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

					Convey("With local session; Should successfully commit a session", func() {
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						sessionID := rpc.dbSession.CreateUnregisteredSession(util.UUID4(), account["id"].(string))

						obj1 := db.NewObject()
						obj1.Key = util.RandString(10)
						obj1.Value = "myvalue"
						ctx2 := metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID))
						_, err := rpc.CreateObjects(ctx2, &proto_rpc.CreateObjectsMsg{
							Bucket:  b.Name,
							Objects: util.MustStringify([]*db.Object{obj1}),
						})
						So(err, ShouldBeNil)

						_, err = rpc.CommitSession(ctx, &proto_rpc.Session{ID: sessionID})
						So(err, ShouldBeNil)

						var count int
						err = rpc.db.Model(db.Object{}).Where("key = ?", obj1.Key).Count(&count).Error
						So(err, ShouldBeNil)
						So(count, ShouldEqual, 1)
					})

					Convey("Should return error if session does not exist locally or in session registry", func() {
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						_, err = rpc.CommitSession(ctx, &proto_rpc.Session{ID: "unknown"})
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

					Convey("With remote session; Should return error if authenticated account is not the owner of the session", func() {
						sessionID, _ := rpc2.dbSession.CreateSession(util.UUID4(), "some_id")

						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						_, err := rpc.CommitSession(ctx, &proto_rpc.Session{ID: sessionID})
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

					Convey("With remote session; Should successfully commit a session", func() {
						sessionID, _ := rpc2.dbSession.CreateSession(util.UUID4(), account["id"].(string))

						obj1 := db.NewObject()
						obj1.Key = util.RandString(10)
						obj1.Value = "myvalue"

						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"id":   account["client_id"],
							"type": oauth.TokenTypeApp,
							"iat":  time.Now().Unix(),
						})

						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						ctx2 := metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID, "authorization", "bearer "+token))
						_, err := rpc.CreateObjects(ctx2, &proto_rpc.CreateObjectsMsg{
							Bucket:  b.Name,
							Objects: util.MustStringify([]*db.Object{obj1}),
						})
						So(err, ShouldBeNil)

						_, err = rpc.CommitSession(ctx2, &proto_rpc.Session{ID: sessionID})
						So(err, ShouldBeNil)

						var count int
						err = rpc.db.Model(db.Object{}).Where("key = ?", obj1.Key).Count(&count).Error
						So(err, ShouldBeNil)
						So(count, ShouldEqual, 1)
					})

					Convey("With remote session; Should return error if session was not found in remode node", func() {
						sessionID, _ := rpc2.dbSession.CreateSession(util.UUID4(), account["id"].(string))

						// remove the session from the remote node
						rpc2.dbSession.RemoveAgent(sessionID)

						obj1 := db.NewObject()
						obj1.Key = util.RandString(10)
						obj1.Value = "myvalue"

						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"id":   account["client_id"],
							"type": oauth.TokenTypeApp,
							"iat":  time.Now().Unix(),
						})

						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						ctx2 := metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID, "authorization", "bearer "+token))

						_, err := rpc.CommitSession(ctx2, &proto_rpc.Session{ID: sessionID})
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
				})

				Convey(".RollbackSession", func() {
					Convey("Should return error if session id is not provided", func() {
						_, err := rpc.RollbackSession(context.Background(), &proto_rpc.Session{})
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

					Convey("With local session; Should return error if authenticated account is not the owner of the session", func() {
						sessionID := rpc.dbSession.CreateUnregisteredSession(util.UUID4(), "some_id")
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						_, err := rpc.RollbackSession(ctx, &proto_rpc.Session{ID: sessionID})
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

					Convey("With local session; Should successfully rollback a session", func() {
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						sessionID := rpc.dbSession.CreateUnregisteredSession(util.UUID4(), account["id"].(string))

						obj1 := db.NewObject()
						obj1.Key = util.RandString(10)
						obj1.Value = "myvalue"
						ctx2 := metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID))
						_, err := rpc.CreateObjects(ctx2, &proto_rpc.CreateObjectsMsg{
							Bucket:  b.Name,
							Objects: util.MustStringify([]*db.Object{obj1}),
						})
						So(err, ShouldBeNil)

						_, err = rpc.RollbackSession(ctx, &proto_rpc.Session{ID: sessionID})
						So(err, ShouldBeNil)

						var count int
						err = rpc.db.Model(db.Object{}).Where("key = ?", obj1.Key).Count(&count).Error
						So(err, ShouldBeNil)
						So(count, ShouldEqual, 0)
					})

					Convey("Should return error if session does not exist locally or in session registry", func() {
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						_, err = rpc.RollbackSession(ctx, &proto_rpc.Session{ID: "unknown"})
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

					Convey("With remote session; Should return error if authenticated account is not the owner of the session", func() {
						sessionID, _ := rpc2.dbSession.CreateSession(util.UUID4(), "some_id")

						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						_, err := rpc.RollbackSession(ctx, &proto_rpc.Session{ID: sessionID})
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

					Convey("With remote session; Should successfully rollback a session", func() {
						sessionID, _ := rpc2.dbSession.CreateSession(util.UUID4(), account["id"].(string))

						obj1 := db.NewObject()
						obj1.Key = util.RandString(10)
						obj1.Value = "myvalue"

						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"id":   account["client_id"],
							"type": oauth.TokenTypeApp,
							"iat":  time.Now().Unix(),
						})

						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						ctx2 := metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID, "authorization", "bearer "+token))
						_, err := rpc.CreateObjects(ctx2, &proto_rpc.CreateObjectsMsg{
							Bucket:  b.Name,
							Objects: util.MustStringify([]*db.Object{obj1}),
						})
						So(err, ShouldBeNil)

						_, err = rpc.RollbackSession(ctx2, &proto_rpc.Session{ID: sessionID})
						So(err, ShouldBeNil)

						var count int
						err = rpc.db.Model(db.Object{}).Where("key = ?", obj1.Key).Count(&count).Error
						So(err, ShouldBeNil)
						So(count, ShouldEqual, 0)
					})

					Convey("With remote session; Should return error if session was not found in remode node", func() {
						sessionID, _ := rpc2.dbSession.CreateSession(util.UUID4(), account["id"].(string))

						// remove the session from the remote node
						rpc2.dbSession.RemoveAgent(sessionID)

						obj1 := db.NewObject()
						obj1.Key = util.RandString(10)
						obj1.Value = "myvalue"

						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"id":   account["client_id"],
							"type": oauth.TokenTypeApp,
							"iat":  time.Now().Unix(),
						})

						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						ctx2 := metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID, "authorization", "bearer "+token))

						_, err := rpc.RollbackSession(ctx2, &proto_rpc.Session{ID: sessionID})
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
				})
			})
		})
	})
}
