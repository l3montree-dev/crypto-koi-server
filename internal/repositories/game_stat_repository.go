package repositories

import (
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/models"
	"gorm.io/gorm"
)

type GameStatRepository interface {
	Repository[models.GameStat]
	FindAllByUserId(userId string) ([]models.GameStat, error)
}

type GormGameStatRepository struct {
	db *gorm.DB
}

func NewGormGameStatRepository(db *gorm.DB) GameStatRepository {
	return &GormGameStatRepository{db: db}
}

func (rep *GormGameStatRepository) Save(gameStat *models.GameStat) error {
	return rep.db.Save(gameStat).Error
}

func (rep *GormGameStatRepository) FindAllByUserId(userId string) ([]models.GameStat, error) {
	var gameStats []models.GameStat
	err := rep.db.Where("user_id = ?", userId).Find(&gameStats).Error
	return gameStats, err
}

func (rep *GormGameStatRepository) GetById(id string) (models.GameStat, error) {
	var gameStat models.GameStat
	err := rep.db.Where("id = ?", id).Find(&gameStat).Error
	return gameStat, err
}
