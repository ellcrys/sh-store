package session

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/ellcrys/elldb/servers/common"
	"github.com/ellcrys/elldb/servers/db"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	. "github.com/smartystreets/goconvey/convey"
)

var err error
var testDB *sql.DB
var dbName string
var conStr = "postgresql://root@localhost:26257?sslmode=disable"
var conStrWithDB string

func setupDBAgent() {
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

func TestAgent(t *testing.T) {

	setupDBAgent()

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

	Convey("TestAgent", t, func() {
		Convey(".put", func() {
			Convey("Should return error if Start() has not been called", func() {
				a := NewAgent(dbCon, make(chan *Op))
				a.curOp = &Op{PutObjects: []*db.Object{}}
				err := a.put()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "agent has not started. Did you call Start()?")
			})
		})

		Convey(".get", func() {
			Convey("Should return error if Start() has not been called", func() {
				a := NewAgent(dbCon, make(chan *Op))
				err := a.get()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "agent has not started. Did you call Start()?")
			})

			Convey("Should return error if `curOp` of agent is nil", func() {
				a := NewAgent(dbCon, make(chan *Op))
				go a.Start(nil)
				time.Sleep(10 * time.Millisecond)
				err := a.get()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "operation object is required")
			})

			Convey("Should return error if query is not valid json", func() {
				a := NewAgent(dbCon, make(chan *Op))
				go a.Start(nil)
				time.Sleep(10 * time.Millisecond)
				a.curOp = &Op{QueryWithJSQ: "{"}
				err := a.get()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "query must be a json string")
			})

			Convey("Should return error if json query was not parsed", func() {
				a := NewAgent(dbCon, make(chan *Op))
				go a.Start(nil)
				time.Sleep(10 * time.Millisecond)
				a.curOp = &Op{QueryWithJSQ: `{ "$not": [] }`}
				err := a.get()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "parser:unknown top level operator: $not")
			})
		})

		Convey(".count", func() {
			Convey("Should return error if Start() has not been called", func() {
				a := NewAgent(dbCon, make(chan *Op))
				err := a.count()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "agent has not started. Did you call Start()?")
			})

			Convey("Should return error if `curOp` of agent is nil", func() {
				a := NewAgent(dbCon, make(chan *Op))
				go a.Start(nil)
				time.Sleep(10 * time.Millisecond)
				err := a.count()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "operation object is required")
			})

			Convey("Should return error if query is not valid json", func() {
				a := NewAgent(dbCon, make(chan *Op))
				a.curOp = &Op{QueryWithJSQ: "{"}
				go a.Start(nil)
				time.Sleep(10 * time.Millisecond)
				err := a.count()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "query must be a json string")
			})

			Convey("Should return error if json query was not parsed", func() {
				a := NewAgent(dbCon, make(chan *Op))
				a.curOp = &Op{QueryWithJSQ: `{ "$not": [] }`}
				go a.Start(nil)
				time.Sleep(10 * time.Millisecond)
				err := a.count()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "parser:unknown top level operator: $not")
			})
		})

		Convey(".update", func() {
			Convey("Should return error if Start() has not been called", func() {
				a := NewAgent(dbCon, make(chan *Op))
				err := a.update()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "agent has not started. Did you call Start()?")
			})

			Convey("Should return error if `curOp` of agent is nil", func() {
				a := NewAgent(dbCon, make(chan *Op))
				go a.Start(nil)
				time.Sleep(10 * time.Millisecond)
				err := a.update()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "operation object is required")
			})

			Convey("Should return error if query is not valid json", func() {
				a := NewAgent(dbCon, make(chan *Op))
				a.curOp = &Op{QueryWithJSQ: "{"}
				go a.Start(nil)
				time.Sleep(10 * time.Millisecond)
				err := a.update()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "query must be a json string")
			})

			Convey("Should return error if json query was not parsed", func() {
				a := NewAgent(dbCon, make(chan *Op))
				a.curOp = &Op{QueryWithJSQ: `{ "$not": [] }`}
				go a.Start(nil)
				time.Sleep(10 * time.Millisecond)
				err := a.update()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "parser:unknown top level operator: $not")
			})
		})
	})
}
