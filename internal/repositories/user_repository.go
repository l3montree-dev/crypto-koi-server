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
}

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

func (rep *GormUserRepository) GetByDeviceId(deviceId string) (models.User, error) {
	var user models.User
	err := rep.db.Preload("cryptogotchi").Where("device_id = ?", deviceId).First(&user).Error
	return user, err
}

func (rep *GormUserRepository) GetByWalletAddress(address string) (models.User, error) {
	var user models.User
	err := rep.db.Preload("cryptogotchi").Where("wallet_address = ?", address).First(&user).Error
	return user, err
}

func (rep *GormUserRepository) GetById(id string) (models.User, error) {
	var user models.User
	err := rep.db.Preload("cryptogotchi").Where("id = ?", id).First(&user).Error
	return user, err
}

func (rep *GormUserRepository) Save(user *models.User) error {
	return rep.db.Create(user).Error
}
