package repositories

import (
	"time"

	"github.com/ethereum/go-ethereum/common/math"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/models"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/util"
	"gorm.io/gorm"
)

type CryptogotchiRepository interface {
	GetCryptogotchiByUint256(tokenId string) (models.Cryptogotchi, error)
	GetCryptogotchiesByUserId(userId string) ([]models.Cryptogotchi, error)
	GetCryptogotchiById(id string) (models.Cryptogotchi, error)
	Save(*models.Cryptogotchi) error
	GetLeaderboard(offset, limit int) ([]models.Cryptogotchi, error)
}

type GormCryptogotchiRepository struct {
	db *gorm.DB
}

func NewGormCryptogotchiRepository(db *gorm.DB) CryptogotchiRepository {
	return &GormCryptogotchiRepository{db: db}
}

func (rep *GormCryptogotchiRepository) GetCryptogotchiByUint256(tokenId string) (models.Cryptogotchi, error) {
	bigInt := math.MustParseBig256(tokenId)

	id, err := util.Uint256ToUuid(bigInt)
	if err != nil {
		return models.Cryptogotchi{}, err
	}

	return rep.GetCryptogotchiById(id.String())
}

func (rep *GormCryptogotchiRepository) Save(m *models.Cryptogotchi) error {
	return rep.db.Save(m).Error
}

func (rep *GormCryptogotchiRepository) GetCryptogotchiesByUserId(userId string) ([]models.Cryptogotchi, error) {
	var cryptogotchies []models.Cryptogotchi
	err := rep.db.Where("user_id = ?", userId).Find(&cryptogotchies).Error
	return cryptogotchies, err
}

func (rep *GormCryptogotchiRepository) GetCryptogotchiById(id string) (models.Cryptogotchi, error) {
	var cryptogotchi models.Cryptogotchi
	err := rep.db.Where("id = ?", id).First(&cryptogotchi).Error
	return cryptogotchi, err
}

func (rep *GormCryptogotchiRepository) GetLeaderboard(offset, limit int) ([]models.Cryptogotchi, error) {
	var cryptogotchies []models.Cryptogotchi
	err := rep.db.Where("predicted_death_date > ?", time.Now()).Order("created_at ASC").Offset(offset).Limit(limit).Find(&cryptogotchies).Error
	return cryptogotchies, err
}
