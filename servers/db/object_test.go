package db

import (
	"fmt"
	"testing"

	"github.com/ellcrys/elldb/servers/common"
	"github.com/ellcrys/util"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	. "github.com/smartystreets/goconvey/convey"
)

var err error
var conStr = "postgresql://postgres@localhost:5432?sslmode=disable"

func TestDB(t *testing.T) {

	dbName, err := common.CreateRandomDB()
	if err != nil {
		panic(fmt.Errorf("failed to create test database: %s", err))
	}

	conStrWithDB := "postgresql://postgres@localhost:5432/" + dbName + "?sslmode=disable"

	db, err := gorm.Open("postgres", conStrWithDB)
	if err != nil {
		t.Fatal("failed to connect to test database")
	}

	ApplyCallbacks(db)
	defer common.DropDB(dbName)
	defer db.Close()

	err = db.CreateTable(&Bucket{}, &Object{}, &Account{}).Error
	if err != nil {
		t.Fatalf("failed to create database tables. %s", err)
	}

	Convey("object.go", t, func() {
		Convey(".GetValidObjectFields", func() {
			Convey("Should return expected fields", func() {
				fields := GetValidObjectFields()
				So(fields, ShouldResemble, []string{"id", "owner_id", "creator_id", "bucket", "key", "value", "created_at", "updated_at", "ref1", "ref2", "ref3", "ref4", "ref5", "ref6", "ref7", "ref8", "ref9", "ref10"})
			})
		})

		Convey(".CreateBucket", func() {
			bucket := NewBucket()
			bucket.Name = util.RandString(5)
			bucket.Creator = "account1"
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

		Convey(".UpdateObjects", func() {
			bucket := NewBucket()
			bucket.Name = util.RandString(5)
			err := CreateBucket(db, bucket)
			So(err, ShouldBeNil)

			objs := []*Object{NewObject(), NewObject()}
			objs[0].Key = "obj1"
			objs[0].Value = "abc"
			objs[1].Key = "obj2"
			err = CreateObjects(db, bucket.Name, objs)

			Convey("Should successfully update object", func() {
				affected, err := UpdateObjects(db, &Object{Key: objs[0].Key}, Object{Value: "xyz"})
				So(affected, ShouldEqual, 1)
				So(err, ShouldBeNil)
				var obj Object
				err = db.Where("key = ?", objs[0].Key).First(&obj).Error
				So(err, ShouldBeNil)
				So(obj.Value, ShouldEqual, "xyz")
			})

			Convey("Should successfully update object using jsq expression", func() {
				affected, err := UpdateObjects(db, &Object{
					QueryParams: QueryParams{
						Expr: Expr{
							Expr: "key = ?",
							Args: []interface{}{objs[0].Key},
						},
					},
				}, Object{Value: "xyz"})
				So(affected, ShouldEqual, 1)
				So(err, ShouldBeNil)
				var obj Object
				err = db.Where("key = ?", objs[0].Key).First(&obj).Error
				So(err, ShouldBeNil)
				So(obj.Value, ShouldEqual, "xyz")
			})
		})

		Convey(".DeleteObjects", func() {
			bucket := NewBucket()
			bucket.Name = util.RandString(5)
			err := CreateBucket(db, bucket)
			So(err, ShouldBeNil)

			objs := []*Object{NewObject(), NewObject()}
			objs[0].Key = "obj1"
			objs[0].Value = "abc"
			objs[1].Key = "obj2"
			err = CreateObjects(db, bucket.Name, objs)

			Convey("Should successfully delete object", func() {
				affected, err := DeleteObjects(db, &Object{Key: objs[0].Key})
				So(affected, ShouldEqual, 1)
				So(err, ShouldBeNil)
				var obj Object
				err = db.Where("key = ?", objs[0].Key).First(&obj).Error
				So(err, ShouldEqual, gorm.ErrRecordNotFound)
			})

			Convey("Should successfully delete object using jsq expression", func() {
				affected, err := DeleteObjects(db, &Object{
					QueryParams: QueryParams{
						Expr: Expr{
							Expr: "key = ?",
							Args: []interface{}{objs[0].Key},
						},
					},
				})
				So(affected, ShouldEqual, 1)
				So(err, ShouldBeNil)
				var obj Object
				err = db.Where("key = ?", objs[0].Key).First(&obj).Error
				So(err, ShouldEqual, gorm.ErrRecordNotFound)
			})
		})

		Reset(func() {
			common.ClearTable(db.DB(), []string{"buckets", "objects"}...)
		})
	})
}
