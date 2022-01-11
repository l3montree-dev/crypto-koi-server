package repositories

import (
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/models"
	"gorm.io/gorm"
)

type EventRepository interface {
	Save(record *models.Event) error
}

type GormEventRepository struct {
	db *gorm.DB
}

func NewGormEventRepository(db *gorm.DB) EventRepository {
	return &GormEventRepository{db: db}
}

func (rep *GormEventRepository) Save(record *models.Event) error {
	return rep.db.Create(record).Error
}
