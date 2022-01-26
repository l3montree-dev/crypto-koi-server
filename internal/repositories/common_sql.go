package repositories

import "gorm.io/gorm"

func orderEventsASC(db *gorm.DB) *gorm.DB {
	return db.Order("events.created_at ASC")
}
