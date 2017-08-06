package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

// ApplyCallbacks applies model callbacks
func ApplyCallbacks(db *gorm.DB) {

	db.Callback().Create().Before("gorm:create").Register("update_dates", func(s *gorm.Scope) {
		if !s.HasError() {
			now := time.Now().UTC()
			if createdAtField, ok := s.FieldByName("CreatedAt"); ok {
				if createdAtField.IsBlank {
					createdAtField.Set(now.Format(time.RFC3339))
				}
			}
			if updatedAtField, ok := s.FieldByName("UpdatedAt"); ok {
				if updatedAtField.IsBlank {
					updatedAtField.Set(now.Format(time.RFC3339))
				}
			}
		}
	})

	db.Callback().Create().Before("gorm:update").Register("update_updated_at", func(s *gorm.Scope) {
		if !s.HasError() {
			if updatedAtField, ok := s.FieldByName("UpdatedAt"); ok {
				if updatedAtField.IsBlank {
					updatedAtField.Set(time.Now().UTC().Format(time.RFC3339))
				}
			}
		}
	})

}
