package service

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/generator"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/models"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/repositories"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/util"
)

type CryptogotchiSvc interface {
	repositories.CryptogotchiRepository
	GenerateCryptogotchiForUser(user *models.User) (models.Cryptogotchi, error)
}

type CryptogotchiService struct {
	repositories.CryptogotchiRepository
	generator generator.Generator
}

func NewCryptogotchiService(rep repositories.CryptogotchiRepository) CryptogotchiSvc {
	return &CryptogotchiService{
		CryptogotchiRepository: rep,
	}
}

func (svc *CryptogotchiService) GenerateCryptogotchiForUser(user *models.User) (models.Cryptogotchi, error) {
	foodValue := 50.
	foodDrainValue := 0.5
	now := time.Now()

	id, err := uuid.NewRandom()
	if err != nil {
		return models.Cryptogotchi{}, err
	}

	koi, _ := svc.generator.GetKoi(id.String())
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
