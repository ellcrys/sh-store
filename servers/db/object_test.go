package db

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/ellcrys/elldb/servers/common"
	"github.com/ellcrys/util"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	. "github.com/smartystreets/goconvey/convey"
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

func TestDB(t *testing.T) {

	db, err := gorm.Open("postgres", conStrWithDB)
	if err != nil {
		t.Fatal("failed to connect to test database")
	}

	err = db.CreateTable(&Bucket{}, &Object{}, &Account{}).Error
	if err != nil {
		t.Fatalf("failed to create database tables. %s", err)
	}

	defer func() {
		common.DropDB(testDB, dbName)
	}()

	Convey("object.go", t, func() {
		Convey(".GetValidObjectFields", func() {
			Convey("Should return expected fields", func() {
				fields := GetValidObjectFields()
				So(fields, ShouldResemble, []string{"id", "owner_id", "creator_id", "bucket", "key", "value", "created_at", "ref1", "ref2", "ref3", "ref4", "ref5", "ref6", "ref7", "ref8", "ref9", "ref10"})
			})
		})

		Convey(".CreateBucket", func() {
			bucket := NewBucket()
			bucket.Name = util.RandString(5)
			bucket.Account = "account1"
			bucket.Immutable = false
			err := CreateBucket(db, bucket)
			So(err, ShouldBeNil)
			So(bucket.SN, ShouldNotBeEmpty)
		})

		Convey(".CreateObjects", func() {
			Convey("Should return error if bucket does not exist", func() {
				obj := NewObject()
				obj.OwnerID = "abc"
				obj.CreatorID = "abc"
				obj.Bucket = "unknown"
				err := CreateObjects(db, "unknown", []*Object{obj})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "bucket not found")
			})

			Convey("Should successfully create objects", func() {
				bucket := NewBucket()
				bucket.Name = util.RandString(5)
				err := CreateBucket(db, bucket)
				So(err, ShouldBeNil)

				objs := []*Object{NewObject(), NewObject()}
				err = CreateObjects(db, bucket.Name, objs)
				So(err, ShouldBeNil)
			})
		})

		Convey(".QueryObjects", func() {
			bucket := NewBucket()
			bucket.Name = util.RandString(5)
			err := CreateBucket(db, bucket)
			So(err, ShouldBeNil)

			objs := []*Object{NewObject(), NewObject()}
			objs[0].Key = "obj1"
			objs[1].Key = "obj2"
			err = CreateObjects(db, bucket.Name, objs)

			Convey("Should query objects using Object struct", func() {
				var objs []*Object
				err := QueryObjects(db, &Object{Key: "obj1"}, &objs)
				So(err, ShouldBeNil)
				So(len(objs), ShouldEqual, 1)
				So(objs[0].Key, ShouldEqual, "obj1")
			})

			Convey("Should query objects using JSQ expression", func() {
				var objs []*Object
				err := QueryObjects(db, &Object{
					QueryParams: QueryParams{
						Expr: Expr{
							Expr: "key = ?",
							Args: []interface{}{"obj2"},
						},
					},
				}, &objs)
				So(err, ShouldBeNil)
				So(len(objs), ShouldEqual, 1)
				So(objs[0].Key, ShouldEqual, "obj2")
			})

			Convey("Should query object and apply limit", func() {
				var objs []*Object
				err := QueryObjects(db, &Object{
					QueryParams: QueryParams{
						Limit: 1,
					},
				}, &objs)
				So(err, ShouldBeNil)
				So(len(objs), ShouldEqual, 1)
			})

			Convey("Should query object and apply ordering", func() {
				var objs []*Object
				err := QueryObjects(db, &Object{
					QueryParams: QueryParams{
						OrderBy: "created_at desc",
					},
				}, &objs)
				So(err, ShouldBeNil)
				So(len(objs), ShouldEqual, 2)
				So(objs[0].Key, ShouldEqual, "obj2")
				So(objs[1].Key, ShouldEqual, "obj1")
			})

			Convey("Should return all objects if no query is included", func() {
				var objs []*Object
				err := QueryObjects(db, &Object{}, &objs)
				So(err, ShouldBeNil)
				So(len(objs), ShouldEqual, 2)
			})
		})

		Convey(".CountObjects", func() {
			bucket := NewBucket()
			bucket.Name = util.RandString(5)
			err := CreateBucket(db, bucket)
			So(err, ShouldBeNil)

			objs := []*Object{NewObject(), NewObject()}
			objs[0].Key = "obj1"
			objs[1].Key = "obj2"
			err = CreateObjects(db, bucket.Name, objs)

			Convey("Should query objects and return expected count", func() {
				var count int
				err := CountObjects(db, &Object{Key: "obj1"}, &count)
				So(err, ShouldBeNil)
				So(count, ShouldEqual, 1)

				err = CountObjects(db, &Object{}, &count)
				So(err, ShouldBeNil)
				So(count, ShouldEqual, 2)
			})

			Convey("Should query objects and return expected count using JSQ expression", func() {
				var count int
				err := CountObjects(db, &Object{
					QueryParams: QueryParams{
						Expr: Expr{
							Expr: "key = ?",
							Args: []interface{}{"obj1"},
						},
					},
				}, &count)
				So(err, ShouldBeNil)
				So(count, ShouldEqual, 1)
			})
		})

		Reset(func() {
			common.ClearTable(db.DB(), []string{"buckets", "objects"}...)
		})
	})
}
