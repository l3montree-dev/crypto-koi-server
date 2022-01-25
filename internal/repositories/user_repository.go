package repositories

import (
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByDeviceId(deviceId string) (models.User, error)
	GetByWalletAddress(address string) (models.User, error)
	GetById(id string) (models.User, error)
	Save(*models.User) error
	GetByRefreshToken(refreshToken string) (models.User, error)
}

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

func (rep *GormUserRepository) GetByDeviceId(deviceId string) (models.User, error) {
	var user models.User
	err := rep.db.Preload("Cryptogotchies").Where("device_id = ?", deviceId).First(&user).Error
	return user, err
}

func (rep *GormUserRepository) GetByWalletAddress(address string) (models.User, error) {
	var user models.User
	err := rep.db.Preload("Cryptogotchies").Where("wallet_address = ?", address).First(&user).Error
	return user, err
}

func (rep *GormUserRepository) GetById(id string) (models.User, error) {
	var user models.User
	err := rep.db.Preload("Cryptogotchies").Where("id = ?", id).First(&user).Error
	return user, err
}

func (rep *GormUserRepository) Save(user *models.User) error {
	return rep.db.Save(user).Error
}

func (rep *GormUserRepository) GetByRefreshToken(refreshToken string) (models.User, error) {
	var user models.User
	err := rep.db.Where("refresh_token = ?", refreshToken).First(&user).Error
	return user, err
}
