package repositories

import "gorm.io/gorm"

type RecordRepository interface {
}

type GormRecordRepository struct {
	db *gorm.DB
}

func NewGormRecordRepository(db *gorm.DB) RecordRepository {
	return &GormRecordRepository{db: db}
}
