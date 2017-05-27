package db

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ellcrys/util"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/ncodes/patchain/cockroach"
	"github.com/ncodes/patchain/cockroach/tables"
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

func TestAgent(t *testing.T) {

	if err := createDb(t); err != nil {
		t.Fatalf("failed to create test database. %s", err)
	}
	defer dropDB(t)

	cdb := cockroach.NewDB()
	cdb.ConnectionString = conStrWithDB
	cdb.NoLogging()
	if err := cdb.Connect(10, 5); err != nil {
		t.Fatalf("failed to connect to database. %s", err)
	}

	if err := cdb.CreateTables(); err != nil {
		t.Fatalf("failed to create tables. %s", err)
	}

	Convey("TestAgent", t, func() {
		Convey(".put", func() {
			Convey("Should return error if Begin() has not been called", func() {
				a := NewAgent(cdb, make(chan *Op))
				err := a.put(tables.Object{})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "agent has not started. Did you can Begin()?")
			})
		})

		Convey(".get", func() {
			Convey("Should return error if Begin() has not been called", func() {
				a := NewAgent(cdb, make(chan *Op))
				err := a.get(nil, nil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "agent has not started. Did you can Begin()?")
			})

			Convey("Should return error if `op` argument is nil", func() {
				a := NewAgent(cdb, make(chan *Op))
				go a.Begin(nil)
				time.Sleep(10 * time.Millisecond)
				err := a.get(nil, nil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "operation object is required")
			})

			Convey("Should return error if query is not valid json", func() {
				a := NewAgent(cdb, make(chan *Op))
				go a.Begin(nil)
				time.Sleep(10 * time.Millisecond)
				err := a.get(&Op{Data: "{"}, nil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "query must be a json string")
			})

			Convey("Should return error if json query was not parsed", func() {
				a := NewAgent(cdb, make(chan *Op))
				go a.Begin(nil)
				time.Sleep(10 * time.Millisecond)
				err := a.get(&Op{Data: `{ "$not": [] }`}, nil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "failed to parse query: unknown top level operator: $not")
			})
		})
	})

	Convey("TestSession", t, func() {
		session := NewSession(conStrWithDB)
		Convey(".Connect", func() {
			Convey("Should connect successfully", func() {
				err := session.Connect(cdb)
				So(err, ShouldBeNil)
			})
		})

		Convey(".SetDB", func() {
			session := NewSession(conStrWithDB)
			So(session.db, ShouldBeNil)
			session.SetDB(cdb)
			So(session.db, ShouldNotBeNil)
			So(session.db, ShouldResemble, cdb)
		})

		Convey(".CreateSession", func() {

			Convey("Should panic if session object is not connected to a database", func() {
				session := NewSession(conStrWithDB)
				sid := "abc"
				So(func() {
					session.CreateSession(sid)
				}, ShouldPanicWith, "not connected to database")
			})

			Convey("Should successfully create a session agent", func() {
				session := NewSession(conStrWithDB)
				session.SetDB(cdb)
				sid := "abc"
				MaxSessionIdleTime = 5 * time.Second
				sid = session.CreateSession(sid)
				So(sid, ShouldNotBeEmpty)
				So(sid, ShouldEqual, sid)

				Convey(".NumSession", func() {
					Convey("Should contain 1 session", func() {
						So(len(session.agents), ShouldEqual, 1)
					})
				})

				Convey(".HasSession", func() {
					Convey("Should include a session with id `abc`", func() {
						So(session.HasSession(sid), ShouldEqual, true)
					})
					Convey("Should not include session with id `xyz`", func() {
						So(session.HasSession("xyz"), ShouldEqual, false)
					})
				})

				Convey(".GetAgent", func() {
					Convey("Should return an agent for the id `abc`", func() {
						agent := session.GetAgent(sid)
						So(agent, ShouldNotBeEmpty)
					})
					Convey("Should return nil if agent `xyz` does not exists", func() {
						So(session.GetAgent("xyz"), ShouldBeNil)
					})
				})

				Convey(".removeAgentRef", func() {
					Convey("Should remove agent session reference", func() {
						session.removeAgentRef(sid)
						So(len(session.agents), ShouldEqual, 0)
					})

					Convey("Should remove session reference when session exits", func() {
						MaxSessionIdleTime = 1 * time.Second
						sid := session.CreateSession("abcd")
						So(session.HasSession(sid), ShouldEqual, true)
						time.Sleep(2 * time.Second)
						So(session.HasSession(sid), ShouldEqual, false)
					})
				})
			})
		})
	})
}
