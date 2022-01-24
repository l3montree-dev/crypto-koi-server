package repositories

import (
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/models"
	"gorm.io/gorm"
)

type CryptogotchiRepository interface {
	GetCryptogotchiByTokenId(tokenId string) (models.Cryptogotchi, error)
	GetCryptogotchiesByUserId(userId string) ([]models.Cryptogotchi, error)
}

type GormCryptogotchiRepository struct {
	db *gorm.DB
}

func NewGormCryptogotchiRepository(db *gorm.DB) CryptogotchiRepository {
	return &GormCryptogotchiRepository{db: db}
}

func (rep *GormCryptogotchiRepository) GetCryptogotchiByTokenId(tokenId string) (models.Cryptogotchi, error) {
	var cryptogotchi models.Cryptogotchi
	err := rep.db.Preload("Events").Where("token_id = ?", tokenId).First(&cryptogotchi).Error
	return cryptogotchi, err
}

func (rep *GormCryptogotchiRepository) GetCryptogotchiesByUserId(userId string) ([]models.Cryptogotchi, error) {
	var cryptogotchies []models.Cryptogotchi
	err := rep.db.Preload("Events").Where("user_id = ?", userId).Find(&cryptogotchies).Error
	return cryptogotchies, err
}
