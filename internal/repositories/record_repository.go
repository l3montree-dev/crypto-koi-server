package repositories

import (
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/models"
	"gorm.io/gorm"
)

type EventRepository interface {
	Save(record *models.Event) error
	GetPaginated(cryptogotchiId string, offset int, amount int) ([]models.Event, error)
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

func (rep *GormEventRepository) GetPaginated(cryptogotchiId string, offset int, amount int) ([]models.Event, error) {
	var events []models.Event
	err := rep.db.Where("cryptogotchi_id = ?", cryptogotchiId).Order("created_at desc").Offset(offset).Limit(amount).Find(&events).Error
	return events, err
}
