package session

import (
	"fmt"
	"testing"
	"time"

	"github.com/ellcrys/elldb/core/servers/common"
	"github.com/ellcrys/elldb/core/servers/db"
	"github.com/ellcrys/util"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAgent(t *testing.T) {

	dbName, err := common.CreateRandomDB()
	if err != nil {
		panic(fmt.Errorf("failed to create test database: %s", err))
	}

	conStrWithDB := "postgresql://postgres@localhost:5432/" + dbName + "?sslmode=disable"
	dbCon, err := gorm.Open("postgres", conStrWithDB)
	if err != nil {
		t.Fatal("failed to connect to test database")
	}

	db.ApplyCallbacks(dbCon)

	defer common.DropDB(dbName)
	defer dbCon.Close()

	err = dbCon.CreateTable(&db.Bucket{}, &db.Object{}, &db.Account{}).Error
	if err != nil {
		t.Fatalf("failed to create database tables. %s", err)
	}

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
				defer a.Stop()
				err := a.get()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "operation object is required")
			})

			Convey("Should return error if query is not valid json", func() {
				a := NewAgent(dbCon, make(chan *Op))
				go a.Start(nil)
				time.Sleep(10 * time.Millisecond)
				defer a.Stop()
				a.curOp = &Op{QueryWithJSQ: "{"}
				err := a.get()
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, &QueryError{Message: "query must be a json string"})
			})

			Convey("Should return error if json query was not parsed", func() {
				a := NewAgent(dbCon, make(chan *Op))
				go a.Start(nil)
				time.Sleep(10 * time.Millisecond)
				defer a.Stop()
				a.curOp = &Op{QueryWithJSQ: `{ "$not": [] }`}
				err := a.get()
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, &QueryError{Message: "unknown top level operator: $not"})
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
				defer a.Stop()
				err := a.count()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "operation object is required")
			})

			Convey("Should return error if query is not valid json", func() {
				a := NewAgent(dbCon, make(chan *Op))
				a.curOp = &Op{QueryWithJSQ: "{"}
				go a.Start(nil)
				defer a.Stop()
				time.Sleep(10 * time.Millisecond)
				err := a.count()
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, &QueryError{"query must be a json string"})
			})

			Convey("Should return error if json query was not parsed", func() {
				a := NewAgent(dbCon, make(chan *Op))
				a.curOp = &Op{QueryWithJSQ: `{ "$not": [] }`}
				go a.Start(nil)
				defer a.Stop()
				time.Sleep(10 * time.Millisecond)
				err := a.count()
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, &QueryError{"unknown top level operator: $not"})
			})
		})

		Convey(".update", func() {
			Convey("Should return error if Start() has not been called", func() {
				a := NewAgent(dbCon, make(chan *Op))
				err := a.update()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldResemble, "agent has not started. Did you call Start()?")
			})

			Convey("Should return error if `curOp` of agent is nil", func() {
				a := NewAgent(dbCon, make(chan *Op))
				go a.Start(nil)
				time.Sleep(10 * time.Millisecond)
				defer a.Stop()
				err := a.update()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "operation object is required")
			})

			Convey("Should return error if query is not valid json", func() {
				a := NewAgent(dbCon, make(chan *Op))
				a.curOp = &Op{QueryWithJSQ: "{"}
				go a.Start(nil)
				defer a.Stop()
				time.Sleep(10 * time.Millisecond)
				err := a.update()
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, &QueryError{"query must be a json string"})
			})

			Convey("Should return error if json query was not parsed", func() {
				a := NewAgent(dbCon, make(chan *Op))
				a.curOp = &Op{QueryWithJSQ: `{ "$not": [] }`}
				go a.Start(nil)
				defer a.Stop()
				time.Sleep(10 * time.Millisecond)
				err := a.update()
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, &QueryError{"unknown top level operator: $not"})
			})
		})

		Convey(".rollback", func() {
			Convey("Should successfully rollback a transaction", func() {
				a := NewAgent(dbCon, make(chan *Op))
				go a.Start(nil)
				time.Sleep(10 * time.Millisecond)
				defer a.Stop()

				obj := db.NewObject()
				obj.Key = util.RandString(5)
				err := a.tx.Create(obj).Error
				So(err, ShouldBeNil)
				a.rollback()

				Convey("Object must not exist", func() {
					var count int64
					err = a.tx.Model(db.Object{}).Where("key = ?", obj.Key).Count(&count).Error
					So(err, ShouldBeNil)
					So(count, ShouldEqual, 0)
				})

				Convey("Agent must be immediately reusable", func() {
					var count int64
					err = a.tx.Model(db.Object{}).Where("key = ?", obj.Key).Count(&count).Error
					So(err, ShouldBeNil)
					So(count, ShouldEqual, 0)
				})
			})
		})
	})
}
