package db

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres driver
)

// SystemEmail is the default system account email
var SystemEmail = "system@ellcrys.co"

// SystemBucket is the default system bucket name
var SystemBucket = "@system"

// Connect connects to a database and returns the DB object
func Connect(conStr string) (db *gorm.DB, err error) {
	db, err = gorm.Open("postgres", conStr)
	return
}

// Initialize creates default tables if they haven't been created.
func Initialize(db *gorm.DB) error {
	if err := db.AutoMigrate(&Account{}, &Bucket{}, &Object{}, &Mapping{}).Error; err != nil {
		return err
	}
	if err := seed(db); err != nil {
		return fmt.Errorf("failed to create system account. %s", err)
	}
	return nil
}

// Seed creates system database records
func seed(db *gorm.DB) error {

	var system Account

	// create system account if not existing
	err := db.Where(Account{Email: SystemEmail}).Attrs(Account{
		FirstName: "",
		LastName:  "",
		Email:     SystemEmail,
		Developer: true,
		Confirmed: true,
		CreatedAt: time.Now().UnixNano(),
	}).FirstOrCreate(&system).Error
	if err != nil {
		return fmt.Errorf("failed to create system account. %s", err)
	}

	// create system bucket
	err = db.Where(Bucket{Name: SystemBucket}).Attrs(Bucket{
		Account:  system.ID,
		Name:      SystemBucket,
		CreatedAt: time.Now().UnixNano(),
	}).FirstOrCreate(&Bucket{}).Error
	if err != nil {
		return fmt.Errorf("failed to create system bucket. %s", err)
	}

	return nil
}
