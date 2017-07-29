package session

import (
	"fmt"
	"testing"
	"time"

	"github.com/ellcrys/elldb/servers/common"
	"github.com/ellcrys/elldb/servers/db"
	"github.com/ellcrys/util"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	. "github.com/smartystreets/goconvey/convey"
)

func setupDB() {
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

func TestSession(t *testing.T) {

	setupDB()

	dbCon, err := gorm.Open("postgres", conStrWithDB)
	if err != nil {
		t.Fatal("failed to connect to test database")
	}

	err = dbCon.CreateTable(&db.Bucket{}, &db.Object{}, &db.Account{}).Error
	if err != nil {
		t.Fatalf("failed to create database tables. %s", err)
	}

	defer func() {
		common.DropDB(testDB, dbName)
	}()

	sessionReg, err := NewConsulRegistry()
	if err != nil {
		t.Fatalf("failed to connect to session registry")
	}

	Convey("TestSession", t, func() {
		MaxSessionIdleTime = 1 * time.Second

		Convey(".SetDB", func() {
			session := NewSession(sessionReg)
			So(session.db, ShouldBeNil)
			session.SetDB(dbCon)
			So(session.db, ShouldNotBeNil)
			So(session.db, ShouldResemble, dbCon)
		})

		Convey(".CreateSession", func() {

			Convey("Should return error if session id is not a UUID v4 value", func() {
				session := NewSession(sessionReg)
				session.SetDB(dbCon)
				_, err := session.CreateSession("sid", "")
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "id is invalid. Expected a UUIDv4 value")
			})

			Convey("Should successfully create a session", func() {
				session := NewSession(sessionReg)
				session.SetDB(dbCon)
				sessionID := util.UUID4()
				sid, err := session.CreateSession(sessionID, "some_id")
				So(err, ShouldBeNil)
				So(sid, ShouldEqual, sessionID)
			})

			Convey("Should implicitly create session ID if not provided", func() {
				session := NewSession(sessionReg)
				session.SetDB(dbCon)
				sid, err := session.CreateSession("", "some_id")
				So(err, ShouldBeNil)
				So(sid, ShouldHaveLength, 36)
			})
		})

		Convey(".NumSession", func() {
			session := NewSession(sessionReg)
			session.SetDB(dbCon)
			Convey("Should contain 1 session", func() {
				sid := util.UUID4()
				_, err := session.CreateSession(sid, "")
				So(err, ShouldBeNil)
				So(session.NumSessions(), ShouldEqual, 1)

				Convey("Should contain 2 sessions", func() {
					sid, err := session.CreateSession(util.UUID4(), "")
					So(err, ShouldBeNil)
					So(sid, ShouldNotBeEmpty)
					So(session.NumSessions(), ShouldEqual, 2)
				})
			})
		})

		Convey(".HasSession", func() {
			session := NewSession(sessionReg)
			session.SetDB(dbCon)

			Convey("Should return true if session exists and false if not", func() {
				sid, err := session.CreateSession(util.UUID4(), util.UUID4())
				So(err, ShouldBeNil)
				So(session.HasSession(sid), ShouldEqual, true)
				So(session.HasSession("xyz"), ShouldEqual, false)
			})
		})

		Convey(".GetAgent", func() {
			sid := util.UUID4()
			session := NewSession(sessionReg)
			session.SetDB(dbCon)
			_, err := session.CreateSession(sid, "")
			So(err, ShouldBeNil)

			Convey("Should return an agent for the id `abc`", func() {
				agent := session.GetAgent(sid)
				So(agent, ShouldNotBeEmpty)

				Convey("Should return nil if agent `xyz` does not exists", func() {
					So(session.GetAgent("xyz"), ShouldBeNil)
				})
			})
		})

		Convey(".GetAgents", func() {
			sid := util.UUID4()
			session := NewSession(sessionReg)
			session.SetDB(dbCon)
			_, err := session.CreateSession(sid, "")
			So(err, ShouldBeNil)
			sid = util.UUID4()
			_, err = session.CreateSession(sid, "")
			So(err, ShouldBeNil)

			Convey("Should return all agents", func() {
				agents := session.GetAgents()
				So(agents, ShouldHaveLength, 2)
			})
		})

		Convey(".Stop", func() {
			sid := util.UUID4()
			session := NewSession(sessionReg)
			session.SetDB(dbCon)
			_, err := session.CreateSession(sid, "")
			So(err, ShouldBeNil)

			Convey("Should remove session agent and delete session from registry", func() {
				session.Stop()
				So(session.NumSessions(), ShouldEqual, 0)
				_, err := session.sessionReg.Get(sid)
				So(err.Error(), ShouldEqual, "not found")
			})
		})

		Convey(".RollbackEnd", func() {
			sid := util.UUID4()
			session := NewSession(sessionReg)
			session.SetDB(dbCon)
			_, err := session.CreateSession(sid, "")
			So(err, ShouldBeNil)

			Convey("Should remove session agent and delete session from registry", func() {
				session.RollbackEnd(sid)
				So(session.NumSessions(), ShouldEqual, 0)
				_, err := session.sessionReg.Get(sid)
				So(err.Error(), ShouldEqual, "not found")
			})
		})

		Convey(".CommitEnd", func() {
			sid := util.UUID4()
			session := NewSession(sessionReg)
			session.SetDB(dbCon)
			_, err := session.CreateSession(sid, "")
			So(err, ShouldBeNil)

			Convey("Should remove session agent and delete session from registry", func() {
				session.CommitEnd(sid)
				So(session.NumSessions(), ShouldEqual, 0)
				_, err := session.sessionReg.Get(sid)
				So(err.Error(), ShouldEqual, "not found")
			})
		})

		Convey(".Rollback", func() {
			sid := util.UUID4()
			session := NewSession(sessionReg)
			session.SetDB(dbCon)
			_, err := session.CreateSession(sid, "")
			So(err, ShouldBeNil)

			Convey("Should not remove session agent and not delete session from registry", func() {
				session.Rollback(sid)
				So(session.NumSessions(), ShouldEqual, 1)
				regItem, err := session.sessionReg.Get(sid)
				So(err, ShouldBeNil)
				So(regItem, ShouldNotBeNil)
				So(regItem.SID, ShouldEqual, sid)
			})
		})

		Convey(".Commit", func() {
			sid := util.UUID4()
			session := NewSession(sessionReg)
			session.SetDB(dbCon)
			_, err := session.CreateSession(sid, "")
			So(err, ShouldBeNil)

			Convey("Should not remove session agent and not delete session from registry", func() {
				session.Commit(sid)
				So(session.NumSessions(), ShouldEqual, 1)
				regItem, err := session.sessionReg.Get(sid)
				So(err, ShouldBeNil)
				So(regItem, ShouldNotBeNil)
				So(regItem.SID, ShouldEqual, sid)
			})
		})

		Convey(".createUnregisteredSession", func() {
			sid := util.UUID4()
			session := NewSession(sessionReg)
			session.SetDB(dbCon)
			session.CreateUnregisteredSession(sid, "account_id")

			Convey("Should not be registered on the session registry", func() {
				_, err := session.sessionReg.Get(sid)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "not found")
			})
		})

		Convey(".RemoveAgent", func() {
			session := NewSession(sessionReg)
			session.SetDB(dbCon)

			Convey("Should remove agent session reference", func() {
				sid := util.UUID4()
				sid, err := session.CreateSession(sid, "")
				So(err, ShouldBeNil)
				So(sid, ShouldNotBeEmpty)
				session.RemoveAgent(sid)
				So(len(session.agents), ShouldEqual, 0)
			})

			Convey("Should remove session reference when session exits", func() {
				sid := util.UUID4()
				sid, err := session.CreateSession(sid, "")
				So(err, ShouldBeNil)
				So(session.HasSession(sid), ShouldEqual, true)
				time.Sleep(3 * time.Second)
				So(session.HasSession(sid), ShouldEqual, false)
			})
		})

		Convey(".End", func() {
			Convey("Should remove agent from memory and registry", func() {
				sid := util.UUID4()
				session := NewSession(sessionReg)
				session.SetDB(dbCon)
				sid, err := session.CreateSession(sid, "")
				So(err, ShouldBeNil)
				So(session.HasSession(sid), ShouldEqual, true)
				session.End(sid)
				So(session.HasSession(sid), ShouldEqual, false)
				regItem, err := session.sessionReg.Get(sid)
				So(regItem, ShouldBeNil)
				So(err, ShouldEqual, ErrNotFound)
			})
		})

		Convey(".SendOp", func() {
			session := NewSession(sessionReg)
			session.SetDB(dbCon)

			Convey("Should return error if session does not exist", func() {
				err := session.SendOp("unknown", nil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "session not found")
			})

			Convey("Should return error if agent is busy", func() {
				sid := util.UUID4()
				sid, err := session.CreateSession(sid, "")
				So(err, ShouldBeNil)
				agentInfo := session.GetAgent(sid)
				So(agentInfo, ShouldNotBeNil)
				agentInfo.Agent.busy = true

				err = session.SendOp(sid, nil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "session is busy")
			})

			Convey("Perform operations", func() {
				bucket := db.NewBucket()
				bucket.Name = "some_name"
				err := dbCon.Create(bucket).Error
				So(err, ShouldBeNil)

				owner := db.NewObject()
				owner.OwnerID = util.RandString(10)
				owner.Bucket = bucket.Name
				err = dbCon.Create(owner).Error
				So(err, ShouldBeNil)

				Convey("OpPutObjects", func() {
					MaxSessionIdleTime = 1 * time.Second
					Convey("Should successfully put an object", func() {
						sid := util.UUID4()
						sid, err := session.CreateSession(sid, "")
						So(err, ShouldBeNil)
						err = session.SendOp(sid, &Op{
							OpType: OpPutObjects,
							Bucket: bucket.Name,
							Done:   make(chan struct{}),
							PutObjects: []*db.Object{
								{
									Key:       "my_obj_1",
									OwnerID:   owner.ID,
									CreatorID: owner.ID,
								},
							},
						})
						So(err, ShouldBeNil)
						session.CommitEnd(sid)

						Convey("OpGetObjects", func() {
							Convey("Should successfully get an object using JSQ", func() {
								sid := util.UUID4()
								sid, err := session.CreateSession(sid, "")
								So(err, ShouldBeNil)
								var out []db.Object
								err = session.SendOp(sid, &Op{
									OpType:       OpGetObjects,
									Done:         make(chan struct{}),
									QueryWithJSQ: `{ "key": "my_obj_1" }`,
									Out:          &out,
								})
								So(err, ShouldBeNil)
								So(len(out), ShouldEqual, 1)
								session.CommitEnd(sid)
							})

							Convey("Should successfully get an object using db.Object", func() {
								sid := util.UUID4()
								sid, err := session.CreateSession(sid, "")
								So(err, ShouldBeNil)
								var out []db.Object
								err = session.SendOp(sid, &Op{
									OpType: OpGetObjects,
									Done:   make(chan struct{}),
									QueryWithObject: &db.Object{
										Key: "my_obj_1",
									},
									Out: &out,
								})
								So(err, ShouldBeNil)
								So(len(out), ShouldEqual, 1)
								session.CommitEnd(sid)
							})
						})

						Convey("OpCountObjects", func() {
							Convey("Should successfully count an object by querying with JSQ", func() {
								sid := util.UUID4()
								sid, err := session.CreateSession(sid, "")
								So(err, ShouldBeNil)
								var count int64
								err = session.SendOp(sid, &Op{
									OpType:       OpCountObjects,
									Done:         make(chan struct{}),
									QueryWithJSQ: `{ "key": "my_obj_1" }`,
									Out:          &count,
								})
								So(err, ShouldBeNil)
								So(count, ShouldEqual, 1)
								session.CommitEnd(sid)
							})

							Convey("Should successfully count an object by querying with db.Object", func() {
								sid := util.UUID4()
								sid, err := session.CreateSession(sid, "")
								So(err, ShouldBeNil)
								var count int64
								err = session.SendOp(sid, &Op{
									OpType: OpCountObjects,
									Done:   make(chan struct{}),
									QueryWithObject: &db.Object{
										Key: "my_obj_1",
									},
									Out: &count,
								})
								So(err, ShouldBeNil)
								So(count, ShouldEqual, 1)
								session.CommitEnd(sid)
							})
						})
					})
				})

				Reset(func() {
					common.ClearTable(dbCon.DB(), "objects")
				})
			})

			Reset(func() {
				common.ClearTable(dbCon.DB(), "objects")
			})
		})
	})
}
