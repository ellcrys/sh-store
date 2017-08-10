package db

import (
	"fmt"
	"testing"

	"github.com/ellcrys/elldb/core/servers/common"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSQL(t *testing.T) {

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

	Convey("sql.go", t, func() {
		Convey(".seed", func() {
			Convey("Should successfully create system account", func() {
				err := seed(db)
				So(err, ShouldBeNil)

				var count int
				err = db.Model(Account{}).Where("email = ?", SystemEmail).Count(&count).Error
				So(err, ShouldBeNil)
				So(count, ShouldEqual, 1)
			})
		})
	})
}
