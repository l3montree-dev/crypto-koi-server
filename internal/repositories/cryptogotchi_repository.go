package repositories

import (
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/entities"
	"gorm.io/gorm"
)

type CryptogotchiRepository interface {
	GetCryptogotchiByTokenId(tokenId string) (entities.Cryptogotchi, error)
	GetCryptogotchiByUserId(userId string) (entities.Cryptogotchi, error)
}

type GormCryptogotchiRepository struct {
	db *gorm.DB
}

func NewGormCryptogotchiRepository(db *gorm.DB) CryptogotchiRepository {
	return &GormCryptogotchiRepository{db: db}
}

func (rep *GormCryptogotchiRepository) GetCryptogotchiByTokenId(tokenId string) (entities.Cryptogotchi, error) {
	var cryptogotchi entities.Cryptogotchi
	err := rep.db.Preload("Records").Where("token_id = ?", tokenId).First(&cryptogotchi).Error
	return cryptogotchi, err
}

func (rep *GormCryptogotchiRepository) GetCryptogotchiByUserId(userId string) (entities.Cryptogotchi, error) {
	var cryptogotchi entities.Cryptogotchi
	err := rep.db.Preload("Records").Where("user_id = ?", userId).First(&cryptogotchi).Error
	return cryptogotchi, err
}
