package repositories

import (
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/models"
	"gorm.io/gorm"
)

type CryptogotchiRepository interface {
	GetCryptogotchiByTokenId(tokenId string) (models.Cryptogotchi, error)
	GetCryptogotchiesByUserId(userId string) ([]models.Cryptogotchi, error)
	GetCryptogotchiById(id string) (models.Cryptogotchi, error)
	GetCryptogotchiByIdWithoutPreload(id string) (models.Cryptogotchi, error)
	Save(*models.Cryptogotchi) error
}

type GormCryptogotchiRepository struct {
	db *gorm.DB
}

func NewGormCryptogotchiRepository(db *gorm.DB) CryptogotchiRepository {
	return &GormCryptogotchiRepository{db: db}
}

func (rep *GormCryptogotchiRepository) GetCryptogotchiByTokenId(tokenId string) (models.Cryptogotchi, error) {
	var cryptogotchi models.Cryptogotchi
	err := rep.db.Preload("Events", orderEventsASC).Where("token_id = ?", tokenId).First(&cryptogotchi).Error
	return cryptogotchi, err
}

func (rep *GormCryptogotchiRepository) Save(m *models.Cryptogotchi) error {
	return rep.db.Save(m).Error
}

func (rep *GormCryptogotchiRepository) GetCryptogotchiesByUserId(userId string) ([]models.Cryptogotchi, error) {
	var cryptogotchies []models.Cryptogotchi
	err := rep.db.Preload("Events", orderEventsASC).Where("user_id = ?", userId).Find(&cryptogotchies).Error
	return cryptogotchies, err
}

func (rep *GormCryptogotchiRepository) GetCryptogotchiById(id string) (models.Cryptogotchi, error) {
	var cryptogotchi models.Cryptogotchi
	err := rep.db.Preload("Events", orderEventsASC).Where("id = ?", id).First(&cryptogotchi).Error
	return cryptogotchi, err
}

func (rep *GormCryptogotchiRepository) GetCryptogotchiByIdWithoutPreload(id string) (models.Cryptogotchi, error) {
	var cryptogotchi models.Cryptogotchi
	err := rep.db.Where("id = ?", id).First(&cryptogotchi).Error
	return cryptogotchi, err
}
