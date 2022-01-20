package repositories

import (
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/models"
	"gorm.io/gorm"
)

type GameStatRepository interface {
	Create(gameStat *models.GameStat) error
	FindAllByUserId(userId string) ([]models.GameStat, error)
}

type GormGameStatRepository struct {
	db *gorm.DB
}

func NewGormGameStatRepository(db *gorm.DB) GameStatRepository {
	return &GormGameStatRepository{db: db}
}

func (rep *GormGameStatRepository) Create(gameStat *models.GameStat) error {
	return rep.db.Create(gameStat).Error
}

func (rep *GormGameStatRepository) FindAllByUserId(userId string) ([]models.GameStat, error) {
	var gameStats []models.GameStat
	err := rep.db.Where("user_id = ?", userId).Find(&gameStats).Error
	return gameStats, err
}
