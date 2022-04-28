package repositories

import (
	"strings"

	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/graph/input"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Repository[models.User]
	GetByDeviceId(deviceId string) (models.User, error)
	GetByWalletAddress(address string) (models.User, error)
	GetByRefreshToken(refreshToken string) (models.User, error)
	GetUsers(query *input.SearchQuery, offset, limit int) ([]models.User, error)
	GetByEmail(email string) (models.User, error)
	Delete(*models.User) error
}

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

func (rep *GormUserRepository) GetUsers(query *input.SearchQuery, offset, limit int) ([]models.User, error) {
	var users []models.User
	q := rep.db.Order("`name` asc").Offset(offset).Limit(limit)
	if query != nil {
		q.Where("name LIKE ?", "%"+query.Name+"%")
	}
	err := q.Find(&users).Error
	return users, err

}

func (rep *GormUserRepository) Delete(u *models.User) error {
	return rep.db.Delete(u).Error
}

func (rep *GormUserRepository) GetByEmail(email string) (models.User, error) {
	var user models.User
	return user, rep.db.Where("email = ?", email).First(&user).Error
}

func (rep *GormUserRepository) GetByDeviceId(deviceId string) (models.User, error) {
	var user models.User
	err := rep.db.Preload("Cryptogotchies").Preload("Cryptogotchies.Events", orderEventsASC).Preload("Cryptogotchies.GameStats").Where("device_id = ?", deviceId).First(&user).Error
	return user, err
}

func (rep *GormUserRepository) GetByWalletAddress(address string) (models.User, error) {
	var user models.User
	err := rep.db.Preload("Cryptogotchies").Preload("Cryptogotchies.Events", orderEventsASC).Preload("Cryptogotchies.GameStats").Where("wallet_address = ?", strings.ToLower(address)).First(&user).Error
	return user, err
}

func (rep *GormUserRepository) GetById(id string) (models.User, error) {
	var user models.User
	err := rep.db.Preload("Cryptogotchies").Preload("Cryptogotchies.Events", orderEventsASC).Preload("Cryptogotchies.GameStats").Where("id = ?", id).First(&user).Error
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
