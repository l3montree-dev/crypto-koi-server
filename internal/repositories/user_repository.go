package repositories

import (
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/entities"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByDeviceId(deviceId string) (entities.User, error)
	GetByWalletAddress(address string) (entities.User, error)
	GetById(id string) (entities.User, error)
	Save(*entities.User) error
}

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

func (rep *GormUserRepository) GetByDeviceId(deviceId string) (entities.User, error) {
	var user entities.User
	err := rep.db.Preload("cryptogotchi").Where("device_id = ?", deviceId).First(&user).Error
	return user, err
}

func (rep *GormUserRepository) GetByWalletAddress(address string) (entities.User, error) {
	var user entities.User
	err := rep.db.Preload("cryptogotchi").Where("wallet_address = ?", address).First(&user).Error
	return user, err
}

func (rep *GormUserRepository) GetById(id string) (entities.User, error) {
	var user entities.User
	err := rep.db.Preload("cryptogotchi").Where("id = ?", id).First(&user).Error
	return user, err
}

func (rep *GormUserRepository) Save(user *entities.User) error {
	return rep.db.Create(user).Error
}
