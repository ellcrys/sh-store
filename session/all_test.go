package session

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ellcrys/util"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/ncodes/patchain/cockroach"
	"github.com/ncodes/patchain/cockroach/tables"
	"github.com/ncodes/patchain/object"
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

func TestAgent(t *testing.T) {

	if err := createDb(t); err != nil {
		t.Fatalf("failed to create test database. %s", err)
	}
	defer dropDB(t)

	cdb := cockroach.NewDB()
	cdb.ConnectionString = conStrWithDB
	cdb.NoLogging()
	if err := cdb.Connect(0, 5); err != nil {
		t.Fatalf("failed to connect to database. %s", err)
	}

	if err := cdb.CreateTables(); err != nil {
		t.Fatalf("failed to create tables. %s", err)
	}

	sessionReg, err := NewConsulRegistry()
	if err != nil {
		t.Fatalf("failed to connect to session registry")
	}

	Convey("TestAgent", t, func() {
		Convey(".put", func() {
			Convey("Should return error if Start() has not been called", func() {
				a := NewAgent(cdb, make(chan *Op))
				a.curOp = &Op{PutData: tables.Object{}}
				err := a.put()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "agent has not started. Did you call Start()?")
			})
		})

		Convey(".get", func() {
			Convey("Should return error if Start() has not been called", func() {
				a := NewAgent(cdb, make(chan *Op))
				err := a.get()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "agent has not started. Did you call Start()?")
			})

			Convey("Should return error if `curOp` of agent is nil", func() {
				a := NewAgent(cdb, make(chan *Op))
				go a.Start(nil)
				time.Sleep(10 * time.Millisecond)
				err := a.get()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "operation object is required")
			})

			Convey("Should return error if query is not valid json", func() {
				a := NewAgent(cdb, make(chan *Op))
				go a.Start(nil)
				time.Sleep(10 * time.Millisecond)
				a.curOp = &Op{QueryWithJSQ: "{"}
				err := a.get()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "query must be a json string")
			})

			Convey("Should return error if json query was not parsed", func() {
				a := NewAgent(cdb, make(chan *Op))
				go a.Start(nil)
				time.Sleep(10 * time.Millisecond)
				a.curOp = &Op{QueryWithJSQ: `{ "$not": [] }`}
				err := a.get()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unknown top level operator: $not")
			})
		})

		Convey(".count", func() {
			Convey("Should return error if Start() has not been called", func() {
				a := NewAgent(cdb, make(chan *Op))
				err := a.count()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "agent has not started. Did you call Start()?")
			})

			Convey("Should return error if `curOp` of agent is nil", func() {
				a := NewAgent(cdb, make(chan *Op))
				go a.Start(nil)
				time.Sleep(10 * time.Millisecond)
				err := a.count()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "operation object is required")
			})

			Convey("Should return error if query is not valid json", func() {
				a := NewAgent(cdb, make(chan *Op))
				a.curOp = &Op{QueryWithJSQ: "{"}
				go a.Start(nil)
				time.Sleep(10 * time.Millisecond)
				err := a.count()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "query must be a json string")
			})

			Convey("Should return error if json query was not parsed", func() {
				a := NewAgent(cdb, make(chan *Op))
				a.curOp = &Op{QueryWithJSQ: `{ "$not": [] }`}
				go a.Start(nil)
				time.Sleep(10 * time.Millisecond)
				err := a.count()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unknown top level operator: $not")
			})
		})
	})

	Convey("TestSession", t, func() {
		MaxSessionIdleTime = 1 * time.Second

		Convey(".SetDB", func() {
			session := NewSession(sessionReg)
			So(session.db, ShouldBeNil)
			session.SetDB(cdb)
			So(session.db, ShouldNotBeNil)
			So(session.db, ShouldResemble, cdb)
		})

		Convey(".CreateSession", func() {
			session := NewSession(sessionReg)
			session.SetDB(cdb)
			_, err := session.CreateSession("sid", "")
			So(err, ShouldBeNil)
		})

		Convey(".NumSession", func() {
			session := NewSession(sessionReg)
			session.SetDB(cdb)
			Convey("Should contain 1 session", func() {
				sid := "abc"
				_, err := session.CreateSession(sid, "")
				So(err, ShouldBeNil)
				So(len(session.agents), ShouldEqual, 1)

				Convey("Should contain 2 sessions", func() {
					sid, err := session.CreateSession("xyz", "")
					So(err, ShouldBeNil)
					So(sid, ShouldNotBeEmpty)
					So(len(session.agents), ShouldEqual, 2)
				})
			})
		})

		Convey(".HasSession", func() {
			session := NewSession(sessionReg)
			session.SetDB(cdb)

			Convey("Should include a session with id `abc`", func() {
				sid := "abc"
				sid, err := session.CreateSession(sid, "")
				So(err, ShouldBeNil)
				So(session.HasSession(sid), ShouldEqual, true)
				So(session.HasSession("xyz"), ShouldEqual, false)
			})
		})

		Convey(".GetAgent", func() {
			sid := "abc"
			session := NewSession(sessionReg)
			session.SetDB(cdb)
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

		Convey(".removeAgent", func() {
			session := NewSession(sessionReg)
			session.SetDB(cdb)

			Convey("Should remove agent session reference", func() {
				sid := "abcdef"
				sid, err := session.CreateSession(sid, "")
				So(err, ShouldBeNil)
				So(sid, ShouldNotBeEmpty)
				session.removeAgent(sid)
				So(len(session.agents), ShouldEqual, 0)
			})

			Convey("Should remove session reference when session exits", func() {
				sid := "abcdefgh"
				sid, err := session.CreateSession(sid, "")
				So(err, ShouldBeNil)
				So(session.HasSession(sid), ShouldEqual, true)
				time.Sleep(3 * time.Second)
				So(session.HasSession(sid), ShouldEqual, false)
			})
		})

		Convey(".End", func() {
			Convey("Should remove agent from memory and registry", func() {
				sid := "abc"
				session := NewSession(sessionReg)
				session.SetDB(cdb)
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
			session.SetDB(cdb)

			Convey("Should return error if session does not exist", func() {
				err := session.SendOp("unknown", nil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "session not found")
			})

			Convey("Should return error if agent is busy", func() {
				sid := "abc"
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

				ownerID := util.RandString(10)
				owner := object.MakeIdentityObject(ownerID, ownerID, "e@email.com", "password", true)
				err := cdb.Create(owner)
				So(err, ShouldBeNil)

				partitions, err := object.NewObject(cdb).CreatePartitions(1, owner.ID, owner.ID)
				So(err, ShouldBeNil)
				So(len(partitions), ShouldEqual, 1)

				Convey("OpPutObjects", func() {
					MaxSessionIdleTime = 1 * time.Second
					Convey("Should successfully put an object", func() {
						sid := "abc"
						sid, err := session.CreateSession(sid, "")
						So(err, ShouldBeNil)
						err = session.SendOp(sid, &Op{
							OpType: OpPutObjects,
							Done:   make(chan struct{}),
							PutData: &tables.Object{
								ID:        "my_obj_1",
								OwnerID:   owner.ID,
								CreatorID: owner.ID,
							},
						})
						So(err, ShouldBeNil)
						session.CommitEnd(sid)

						Convey("OpGetObjects", func() {
							Convey("Should successfully get an object using JSQ", func() {
								sid := "abc"
								sid, err := session.CreateSession(sid, "")
								So(err, ShouldBeNil)
								var out []tables.Object
								err = session.SendOp(sid, &Op{
									OpType:       OpGetObjects,
									Done:         make(chan struct{}),
									QueryWithJSQ: `{ "id": "my_obj_1" }`,
									Out:          &out,
								})
								So(err, ShouldBeNil)
								So(len(out), ShouldEqual, 1)
								session.CommitEnd(sid)
							})

							Convey("Should successfully get an object using tables.Object", func() {
								sid := "abc"
								sid, err := session.CreateSession(sid, "")
								So(err, ShouldBeNil)
								var out []tables.Object
								err = session.SendOp(sid, &Op{
									OpType: OpGetObjects,
									Done:   make(chan struct{}),
									QueryWithObject: &tables.Object{
										ID: "my_obj_1",
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
								sid := "abc"
								sid, err := session.CreateSession(sid, "")
								So(err, ShouldBeNil)
								var count int64
								err = session.SendOp(sid, &Op{
									OpType:       OpCountObjects,
									Done:         make(chan struct{}),
									QueryWithJSQ: `{ "id": "my_obj_1" }`,
									Out:          &count,
								})
								So(err, ShouldBeNil)
								So(count, ShouldEqual, 1)
								session.CommitEnd(sid)
							})

							Convey("Should successfully count an object by querying with tables.Object", func() {
								sid := "abc"
								sid, err := session.CreateSession(sid, "")
								So(err, ShouldBeNil)
								var count int64
								err = session.SendOp(sid, &Op{
									OpType: OpCountObjects,
									Done:   make(chan struct{}),
									QueryWithObject: &tables.Object{
										ID: "my_obj_1",
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
					clearTable(cdb.GetConn().(*gorm.DB), "objects")
				})
			})

			Reset(func() {
				clearTable(cdb.GetConn().(*gorm.DB), "objects")
			})
		})
	})
}
