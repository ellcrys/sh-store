package db

import (
	"fmt"
	"testing"

	"github.com/ellcrys/elldb/servers/common"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	. "github.com/smartystreets/goconvey/convey"
)

func SetupDB() {

	conStr = "postgresql://root@localhost:26257?sslmode=disable"

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

func TestSQL(t *testing.T) {

	SetupDB()

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

	Convey("sql.go", t, func() {
		Convey(".seed", func() {
			Convey("Should successfully create system account and bucket", func() {
				err := seed(db)
				So(err, ShouldBeNil)

				var count int
				err = db.Model(Account{}).Where("email = ?", SystemEmail).Count(&count).Error
				So(err, ShouldBeNil)
				So(count, ShouldEqual, 1)

				err = db.Model(Bucket{}).Where("name = ?", SystemBucket).Count(&count).Error
				So(err, ShouldBeNil)
				So(count, ShouldEqual, 1)
			})
		})
	})
}
