package db

import (
	"fmt"
	"strings"

	"github.com/ellcrys/util"
	"github.com/fatih/structs"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/gorm"
)

var blacklistedFields = []string{}

// GetValidObjectFields returns a list of valid object fields that can
// be included in a query. The field name is retrieved from the json tag
// of Object struct and converted to snake case. Fields included in
// blacklistedFields will be not be returned.
func GetValidObjectFields() (fields []string) {
	var fieldNames = structs.New(Object{}).Fields()
	for _, f := range fieldNames {
		field := strcase.ToSnake(f.Tag("json"))
		field = strings.Split(field, ",")[0]
		if !util.InStringSlice(blacklistedFields, field) {
			fields = append(fields, field)
		}
	}
	return
}

// CreateBucket creates a new bucket
func CreateBucket(db *gorm.DB, bucket *Bucket) error {
	return db.Create(bucket).Error
}

// CreateObjects creates objects.
func CreateObjects(db *gorm.DB, bucket string, objs []*Object) error {

	// get bucket
	var b Bucket
	if err := db.Where("name = ?", bucket).First(&b).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("bucket not found")
		}
		return err
	}

	// assign id (if not assigned), bucket to objects and create the object
	for i, o := range objs {
		if len(o.ID) == 0 {
			o.ID = util.UUID4()
		}
		o.Bucket = b.Name
		if err := db.Create(o).Error; err != nil {
			return fmt.Errorf("object %d: %s", i, err)
		}
		o.SN = 0 // set to zero, so it doesn't get returned
	}

	return nil
}

// QueryObjects fetches objects
func QueryObjects(db *gorm.DB, q Query, out interface{}) error {

	qp := q.GetQueryParams()

	// no custom expression, use full object
	if qp == nil || len(qp.Expr.Expr) == 0 {
		db = db.Where(q)
	} else { // use expression and arguments
		db = db.Where(qp.Expr.Expr, qp.Expr.Args...)
	}

	// apply ordering
	if len(qp.OrderBy) > 0 {
		db = db.Order(qp.OrderBy)
	} else if qp.Limit > 0 {
		db = db.Limit(qp.Limit)
	}

	return db.Find(out).Error
}

// CountObjects counts the number of objects that match a query
func CountObjects(db *gorm.DB, q Query, out interface{}) error {

	db = db.Model(Object{})
	qp := q.GetQueryParams()

	// no custom expression, use full object
	if qp == nil || len(qp.Expr.Expr) == 0 {
		db = db.Where(q)
	} else { // use expression and arguments
		db = db.Where(qp.Expr.Expr, qp.Expr.Args...)
	}

	return db.Count(out).Error
}

// UpdateObjects updates objects match the query
func UpdateObjects(db *gorm.DB, q Query, update interface{}) (int64, error) {

	qp := q.GetQueryParams()
	db = db.Model(Object{})

	// no custom expression, use full object
	if qp == nil || len(qp.Expr.Expr) == 0 {
		db = db.Where(q)
	} else { // use expression and arguments
		db = db.Where(qp.Expr.Expr, qp.Expr.Args...)
	}

	r := db.UpdateColumns(update)
	return r.RowsAffected, r.Error
}

// DeleteObjects deletes objects match the query
func DeleteObjects(db *gorm.DB, q Query) (int64, error) {

	qp := q.GetQueryParams()
	db = db.Model(Object{})

	// no custom expression, use full object
	if qp == nil || len(qp.Expr.Expr) == 0 {
		db = db.Where(q)
	} else { // use expression and arguments
		db = db.Where(qp.Expr.Expr, qp.Expr.Args...)
	}

	r := db.Delete(Object{})
	return r.RowsAffected, r.Error
}
