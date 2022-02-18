package service

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/config"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/generator"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/models"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/repositories"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/util"
)

type CryptogotchiSvc interface {
	repositories.CryptogotchiRepository
	GenerateCryptogotchiForUser(user *models.User) (models.Cryptogotchi, error)
	GenerateWithFixedTokenId(user *models.User, id uuid.UUID) (models.Cryptogotchi, error)
}

type CryptogotchiService struct {
	repositories.CryptogotchiRepository
	generator *generator.Generator
}

func NewCryptogotchiService(rep repositories.CryptogotchiRepository, g *generator.Generator) CryptogotchiSvc {
	return &CryptogotchiService{
		CryptogotchiRepository: rep,
		generator:              g,
	}
}

func (svc *CryptogotchiService) GenerateWithFixedTokenId(user *models.User, id uuid.UUID) (models.Cryptogotchi, error) {
	foodValue := config.DEFAULT_FOOD_VALUE
	foodDrainValue := config.DEFAULT_FOOD_DRAIN
	now := time.Now()

	tokenId, err := util.TokenIdToIntString(id.String())

	if err != nil {
		return models.Cryptogotchi{}, err
	}

	koi, _ := svc.generator.GetKoi(tokenId)
	name := strings.Title((koi.GetType()))

	newCrypt := models.Cryptogotchi{
		// TODO: generate a random name
		Base: models.Base{
			Id: id,
		},
		Name:               util.Str(name),
		OwnerId:            user.Id,
		Food:               foodValue,
		FoodDrain:          foodDrainValue,
		PredictedDeathDate: now.Add(time.Duration(foodValue/foodDrainValue) * time.Minute),
		SnapshotValid:      now,
	}
	err = svc.Save(&newCrypt)
	return newCrypt, err
}

func (svc *CryptogotchiService) GenerateCryptogotchiForUser(user *models.User) (models.Cryptogotchi, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return models.Cryptogotchi{}, err
	}
	return svc.GenerateWithFixedTokenId(user, id)
}
