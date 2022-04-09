package repositories

import (
	"time"

	"github.com/ethereum/go-ethereum/common/math"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/graph/input"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/models"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/util"
	"gorm.io/gorm"
)

type CryptogotchiRepository interface {
	Repository[models.Cryptogotchi]
	GetCryptogotchiByUint256(tokenId string) (models.Cryptogotchi, error)
	GetCryptogotchiesByUserId(userId string) ([]models.Cryptogotchi, error)
	GetLeaderboard() ([]models.Cryptogotchi, error)
	GetCachedLeaderboard(offset, limit int) ([]models.Cryptogotchi, error)
	GetCryptogotchies(query *input.SearchQuery, offset, limit int) ([]models.Cryptogotchi, error)
	Create(m *models.Cryptogotchi) error
	GetCryptogotchiesWithPredictedDeathDateBetween(start, end time.Time) ([]models.Cryptogotchi, error)
}

type GormCryptogotchiRepository struct {
	db *gorm.DB
}

func NewGormCryptogotchiRepository(db *gorm.DB) CryptogotchiRepository {
	return &GormCryptogotchiRepository{db: db}
}

func onlyActive(db *gorm.DB) *gorm.DB {
	return db.Where("active = ?", true)
}

func (rep *GormCryptogotchiRepository) GetCryptogotchies(query *input.SearchQuery, offset, limit int) ([]models.Cryptogotchi, error) {
	var cryptogotchies []models.Cryptogotchi
	q := rep.db.Scopes(onlyActive).Order("`rank` asc").Offset(offset).Limit(limit)
	if query != nil {
		q.Where("name LIKE ?", "%"+query.Name+"%")
	}
	err := q.Find(&cryptogotchies).Error
	return cryptogotchies, err
}

func (rep *GormCryptogotchiRepository) GetCryptogotchiByUint256(tokenId string) (models.Cryptogotchi, error) {
	bigInt := math.MustParseBig256(tokenId)

	id, err := util.Uint256ToUuid(bigInt)
	if err != nil {
		return models.Cryptogotchi{}, err
	}

	return rep.GetById(id.String())
}

func (rep *GormCryptogotchiRepository) Save(m *models.Cryptogotchi) error {
	return rep.db.Save(m).Error
}

func (rep *GormCryptogotchiRepository) Create(m *models.Cryptogotchi) error {
	return rep.db.Create(m).Error
}

func (rep *GormCryptogotchiRepository) GetCryptogotchiesByUserId(userId string) ([]models.Cryptogotchi, error) {
	var cryptogotchies []models.Cryptogotchi
	err := rep.db.Scopes(onlyActive).Where("user_id = ?", userId).Find(&cryptogotchies).Error
	return cryptogotchies, err
}

func (rep *GormCryptogotchiRepository) GetById(id string) (models.Cryptogotchi, error) {
	var cryptogotchi models.Cryptogotchi
	err := rep.db.Where("id = ?", id).First(&cryptogotchi).Error
	return cryptogotchi, err
}

func (rep *GormCryptogotchiRepository) GetCryptogotchiesWithPredictedDeathDateBetween(start, end time.Time) ([]models.Cryptogotchi, error) {
	var cryptogotchies []models.Cryptogotchi
	err := rep.db.Where("predicted_death_date >= ? AND predicted_death_date < ?", start, end).Find(&cryptogotchies).Error
	return cryptogotchies, err
}

func (rep *GormCryptogotchiRepository) GetLeaderboard() ([]models.Cryptogotchi, error) {
	var cryptogotchies []models.Cryptogotchi
	err := rep.db.Where("predicted_death_date > ?", time.Now()).Order("created_at ASC").Find(&cryptogotchies).Error
	return cryptogotchies, err
}

func (rep *GormCryptogotchiRepository) GetCachedLeaderboard(offset, limit int) ([]models.Cryptogotchi, error) {
	var cryptogotchies []models.Cryptogotchi
	err := rep.db.Where("predicted_death_date > ? AND `rank` > -1", time.Now()).Order("`rank` asc").Offset(offset).Limit(limit).Find(&cryptogotchies).Error
	return cryptogotchies, err
}
