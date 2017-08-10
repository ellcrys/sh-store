package servers

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/ellcrys/elldb/core/servers/db"
	"github.com/ellcrys/elldb/core/servers/oauth"
	"github.com/ellcrys/elldb/core/servers/proto_rpc"
	"github.com/ellcrys/elldb/core/session"
	"github.com/ellcrys/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSession(t *testing.T) {
	setup(t, func(rpc, rpc2 *RPC) {
		Convey("Session", t, func() {
			var account = db.Account{ID: util.UUID4(), FirstName: "john", LastName: "Doe", Email: util.RandString(5) + "@example.com", Password: "something"}
			err := rpc.db.Create(&account).Error
			So(err, ShouldBeNil)

			var contract = db.Contract{ID: util.UUID4(), Creator: account.ID, Name: util.RandString(5), ClientID: util.RandString(5), ClientSecret: util.RandString(5)}
			err = rpc.db.Create(&contract).Error
			So(err, ShouldBeNil)

			ctx := context.WithValue(context.Background(), CtxContract, &contract)
			b := &proto_rpc.CreateBucketMsg{Name: util.RandString(5)}
			bucket, err := rpc.CreateBucket(ctx, b)
			So(err, ShouldBeNil)
			So(bucket.Name, ShouldEqual, b.Name)
			So(bucket.ID, ShouldHaveLength, 36)

			Convey("session.go", func() {

				Convey(".CreateSession", func() {

					Convey("Should return error if session id is not a UUID4", func() {
						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						_, err := rpc.CreateSession(ctx, &proto_rpc.Session{ID: "invalid_id"})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"detail": "id is invalid. Expected UUIDv4 value",
							"status": "400",
						})
					})

					Convey("Should successfully create a session with an explicitly provided ID", func() {
						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						id := util.UUID4()
						resp, err := rpc.CreateSession(ctx, &proto_rpc.Session{ID: id})
						So(err, ShouldBeNil)
						So(resp.ID, ShouldEqual, id)
						rpc.dbSession.GetAgent(id).Agent.Stop() // stop agent. if we don't do this, the agent transaction will remain alive
					})

					Convey("Should successfully create a session with an implicitly assigned ID when ID is not provide", func() {
						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						resp, err := rpc.CreateSession(ctx, &proto_rpc.Session{})
						So(err, ShouldBeNil)
						So(resp.ID, ShouldNotBeEmpty)
						rpc.dbSession.GetAgent(resp.ID).Agent.Stop() // stop agent. if we don't do this, the agent transaction will remain alive
					})
				})

				Convey(".GetSession", func() {
					Convey("Should return error if session id is not provided", func() {
						_, err := rpc.GetSession(context.Background(), &proto_rpc.Session{})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"status": "400",
							"detail": "session id is required",
						})
					})

					Convey("Should return permission error if authenticated contract is not the owner of the session", func() {
						sessionID := rpc.dbSession.CreateUnregisteredSession(util.UUID4(), "some_id")
						defer rpc.dbSession.GetAgent(sessionID).Agent.Stop() // stop agent. if we don't do this, the agent transaction will remain alive
						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						_, err := rpc.GetSession(ctx, &proto_rpc.Session{ID: sessionID})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"detail": "permission denied: you don't have permission to perform this operation",
							"status": "401",
						})
					})

					Convey("Should return error is session does not exists", func() {
						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						_, err := rpc.GetSession(ctx, &proto_rpc.Session{ID: "unknown"})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"detail": "session not found",
							"status": "404",
						})
					})

					Convey("Using a remote session; Should return permission error if authenticated contract is not the owner of the session", func() {
						sid := util.UUID4()
						err := rpc.sessionReg.Add(session.RegItem{SID: sid,
							Address: "localhost",
							Port:    9000,
							Meta: map[string]interface{}{
								"owner_id": "some_account",
							},
						})
						So(err, ShouldBeNil)
						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						_, err = rpc.GetSession(ctx, &proto_rpc.Session{ID: sid})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"status": "401",
							"detail": "permission denied: you don't have permission to perform this operation",
						})
					})

					Convey("With a remote session; Should successfully get a session", func() {
						sid := util.UUID4()
						err := rpc.sessionReg.Add(session.RegItem{SID: sid,
							Address: "localhost",
							Port:    9000,
							Meta: map[string]interface{}{
								"owner_id": contract.ID,
							},
						})
						So(err, ShouldBeNil)
						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						s, err := rpc.GetSession(ctx, &proto_rpc.Session{ID: sid})
						So(err, ShouldBeNil)
						So(s.ID, ShouldEqual, sid)
					})

					Convey("With a unregistered session; Should successfully get a session", func() {
						sid := rpc.dbSession.CreateUnregisteredSession(util.UUID4(), contract.ID)
						defer rpc.dbSession.GetAgent(sid).Agent.Stop() // stop agent. if we don't do this, the agent transaction will remain alive
						So(err, ShouldBeNil)
						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						s, err := rpc.GetSession(ctx, &proto_rpc.Session{ID: sid})
						So(err, ShouldBeNil)
						So(s.ID, ShouldEqual, sid)
					})
				})

				Convey(".DeleteSession", func() {
					Convey("Should return error if session id is not provided", func() {
						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						_, err := rpc.DeleteSession(ctx, &proto_rpc.Session{})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"detail": "session id is required",
							"status": "400",
						})
					})

					Convey("Should return permission error if authenticated contract is not the owner of the session", func() {
						sessionID := rpc.dbSession.CreateUnregisteredSession(util.UUID4(), "some_id")
						defer rpc.dbSession.GetAgent(sessionID).Agent.Stop() // stop agent. if we don't do this, the agent transaction will remain alive
						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						_, err := rpc.DeleteSession(ctx, &proto_rpc.Session{ID: sessionID})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"detail": "permission denied: you don't have permission to perform this operation",
							"status": "401",
						})
					})

					Convey("With an unregistered session; Should successfully delete session", func() {
						sessionID := rpc.dbSession.CreateUnregisteredSession(util.UUID4(), contract.ID)
						defer rpc.dbSession.GetAgent(sessionID).Agent.Stop() // stop agent. if we don't do this, the agent transaction will remain alive
						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						resp, err := rpc.DeleteSession(ctx, &proto_rpc.Session{ID: sessionID})
						So(err, ShouldBeNil)
						So(resp.ID, ShouldEqual, sessionID)
						So(rpc.dbSession.HasSession(sessionID), ShouldBeFalse)
					})

					Convey("With an remote session; Should successfully delete session", func() {
						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"client_id": contract.ClientID,
							"type":      oauth.TokenTypeApp,
							"iat":       time.Now().Unix(),
						})

						sessionID := util.UUID4()
						rpc2.dbSession.CreateSession(sessionID, contract.ID)
						defer rpc2.dbSession.GetAgent(sessionID).Agent.Stop() // stop agent. if we don't do this, the agent transaction will remain alive

						ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "bearer "+token))
						ctx = context.WithValue(ctx, CtxContract, &contract)
						resp, err := rpc.DeleteSession(ctx, &proto_rpc.Session{ID: sessionID})
						So(err, ShouldBeNil)
						So(resp.ID, ShouldEqual, sessionID)
						So(rpc.dbSession.HasSession(sessionID), ShouldBeFalse)
					})

					Convey("With remote session; Should return error if session exists in registry but not found in the node hosting the session", func() {
						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"client_id": contract.ClientID,
							"type":      oauth.TokenTypeApp,
							"iat":       time.Now().Unix(),
						})

						sessionID := util.UUID4()
						rpc2.dbSession.CreateSession(sessionID, contract.ID)

						defer rpc2.dbSession.GetAgent(sessionID).Agent.Stop() // stop agent. if we don't do this, the agent transaction will remain alive
						rpc2.dbSession.RemoveAgent(sessionID)                 // remove agent from session record

						ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "bearer "+token))
						ctx = context.WithValue(ctx, CtxContract, &contract)
						_, err := rpc.DeleteSession(ctx, &proto_rpc.Session{ID: sessionID})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"detail": "session not found",
							"status": "404",
						})
					})

					Convey("Should return error if session does not exist locally or in session registry", func() {
						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"client_id": contract.ClientID,
							"type":      oauth.TokenTypeApp,
							"iat":       time.Now().Unix(),
						})

						sessionID := util.UUID4()

						ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "bearer "+token))
						ctx = context.WithValue(ctx, CtxContract, &contract)
						_, err := rpc.DeleteSession(ctx, &proto_rpc.Session{ID: sessionID})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"detail": "session not found",
							"status": "404",
						})
					})

					Convey("With an remote session; Should return error if authenticated contract is not the owner of the session", func() {
						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"client_id": contract.ClientID,
							"type":      oauth.TokenTypeApp,
							"iat":       time.Now().Unix(),
						})

						sessionID := util.UUID4()
						rpc2.dbSession.CreateSession(sessionID, "some_id")
						defer rpc2.dbSession.GetAgent(sessionID).Agent.Stop() // stop agent. if we don't do this, the agent transaction will remain alive

						ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "bearer "+token))
						ctx = context.WithValue(ctx, CtxContract, &contract)
						_, err := rpc.DeleteSession(ctx, &proto_rpc.Session{ID: sessionID})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"detail": "permission denied: you don't have permission to perform this operation",
							"status": "401",
						})
					})
				})

				Convey(".CommitSession", func() {
					Convey("Should return error if session id is not provide", func() {
						_, err := rpc.CommitSession(context.Background(), &proto_rpc.Session{})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"detail": "session id is required",
							"status": "400",
						})
					})

					Convey("With local session; Should return error if authenticated account is not the owner of the session", func() {
						sessionID := rpc.dbSession.CreateUnregisteredSession(util.UUID4(), "some_id")
						defer rpc.dbSession.GetAgent(sessionID).Agent.Stop() // stop agent. if we don't do this, the agent transaction will remain alive

						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						_, err := rpc.CommitSession(ctx, &proto_rpc.Session{ID: sessionID})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"detail": "permission denied: you don't have permission to perform this operation",
							"status": "401",
						})
					})

					Convey("With local session; Should successfully commit a session", func() {
						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						sessionID := rpc.dbSession.CreateUnregisteredSession(util.UUID4(), contract.ID)
						defer rpc.dbSession.GetAgent(sessionID).Agent.Stop() // stop agent. if we don't do this, the agent transaction will remain alive

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
						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						_, err = rpc.CommitSession(ctx, &proto_rpc.Session{ID: "unknown"})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"detail": "session not found",
							"status": "404",
						})
					})

					Convey("With remote session; Should return error if authenticated account is not the owner of the session", func() {
						sessionID, _ := rpc2.dbSession.CreateSession(util.UUID4(), "some_id")
						defer rpc2.dbSession.GetAgent(sessionID).Agent.Stop() // stop agent. if we don't do this, the agent transaction will remain alive

						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						_, err := rpc.CommitSession(ctx, &proto_rpc.Session{ID: sessionID})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"detail": "permission denied: you don't have permission to perform this operation",
							"status": "401",
						})
					})

					Convey("With remote session; Should successfully commit a session", func() {
						sessionID, _ := rpc2.dbSession.CreateSession(util.UUID4(), contract.ID)
						defer rpc2.dbSession.GetAgent(sessionID).Agent.Stop() // stop agent. if we don't do this, the agent transaction will remain alive

						obj1 := db.NewObject()
						obj1.Key = util.RandString(10)
						obj1.Value = "myvalue"

						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"client_id": contract.ClientID,
							"type":      oauth.TokenTypeApp,
							"iat":       time.Now().Unix(),
						})

						ctx := context.WithValue(context.Background(), CtxContract, &contract)
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
						sessionID, _ := rpc2.dbSession.CreateSession(util.UUID4(), contract.ID)
						defer rpc2.dbSession.GetAgent(sessionID).Agent.Stop() // stop agent. if we don't do this, the agent transaction will remain alive

						// remove the session from the remote node
						rpc2.dbSession.RemoveAgent(sessionID)

						obj1 := db.NewObject()
						obj1.Key = util.RandString(10)
						obj1.Value = "myvalue"

						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"client_id": contract.ClientID,
							"type":      oauth.TokenTypeApp,
							"iat":       time.Now().Unix(),
						})

						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						ctx2 := metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID, "authorization", "bearer "+token))

						_, err := rpc.CommitSession(ctx2, &proto_rpc.Session{ID: sessionID})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"status": "404",
							"detail": "session not found",
						})
					})
				})

				Convey(".RollbackSession", func() {
					Convey("Should return error if session id is not provided", func() {
						_, err := rpc.RollbackSession(context.Background(), &proto_rpc.Session{})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"detail": "session id is required",
							"status": "400",
						})
					})

					Convey("With local session; Should return error if authenticated account is not the owner of the session", func() {
						sessionID := rpc.dbSession.CreateUnregisteredSession(util.UUID4(), "some_id")
						defer rpc.dbSession.GetAgent(sessionID).Agent.Stop() // stop agent. if we don't do this, the agent transaction will remain alive

						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						_, err := rpc.RollbackSession(ctx, &proto_rpc.Session{ID: sessionID})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"detail": "permission denied: you don't have permission to perform this operation",
							"status": "401",
						})
					})

					Convey("With local session; Should successfully rollback a session", func() {
						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						sessionID := rpc.dbSession.CreateUnregisteredSession(util.UUID4(), contract.ID)
						defer rpc.dbSession.GetAgent(sessionID).Agent.Stop() // stop agent. if we don't do this, the agent transaction will remain alive

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
						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						_, err = rpc.RollbackSession(ctx, &proto_rpc.Session{ID: "unknown"})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"detail": "session not found",
							"status": "404",
						})
					})

					Convey("With remote session; Should return error if authenticated account is not the owner of the session", func() {
						sessionID, _ := rpc2.dbSession.CreateSession(util.UUID4(), "some_id")
						defer rpc2.dbSession.GetAgent(sessionID).Agent.Stop() // stop agent. if we don't do this, the agent transaction will remain alive

						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						_, err := rpc.RollbackSession(ctx, &proto_rpc.Session{ID: sessionID})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"detail": "permission denied: you don't have permission to perform this operation",
							"status": "401",
						})
					})

					Convey("With remote session; Should successfully rollback a session", func() {
						sessionID, _ := rpc2.dbSession.CreateSession(util.UUID4(), contract.ID)
						defer rpc2.dbSession.GetAgent(sessionID).Agent.Stop() // stop agent. if we don't do this, the agent transaction will remain alive

						obj1 := db.NewObject()
						obj1.Key = util.RandString(10)
						obj1.Value = "myvalue"

						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"client_id": contract.ClientID,
							"type":      oauth.TokenTypeApp,
							"iat":       time.Now().Unix(),
						})

						ctx := context.WithValue(context.Background(), CtxContract, &contract)
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
						sessionID, _ := rpc2.dbSession.CreateSession(util.UUID4(), contract.ID)
						defer rpc2.dbSession.GetAgent(sessionID).Agent.Stop() // stop agent. if we don't do this, the agent transaction will remain alive

						// remove the session from the remote node
						rpc2.dbSession.RemoveAgent(sessionID)

						obj1 := db.NewObject()
						obj1.Key = util.RandString(10)
						obj1.Value = "myvalue"

						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"client_id": contract.ClientID,
							"type":      oauth.TokenTypeApp,
							"iat":       time.Now().Unix(),
						})

						ctx := context.WithValue(context.Background(), CtxContract, &contract)
						ctx2 := metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID, "authorization", "bearer "+token))

						_, err := rpc.RollbackSession(ctx2, &proto_rpc.Session{ID: sessionID})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"detail": "session not found",
							"status": "404",
						})
					})
				})
			})
		})
	})
}
