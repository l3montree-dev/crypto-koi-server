package service

import (
	"time"

	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/models"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/repositories"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/util"
)

type CryptogotchiSvc interface {
	repositories.CryptogotchiRepository
	GenerateCryptogotchiForUser(user *models.User) (models.Cryptogotchi, error)
}

type CryptogotchiService struct {
	repositories.CryptogotchiRepository
}

func NewCryptogotchiService(rep repositories.CryptogotchiRepository) CryptogotchiSvc {
	return &CryptogotchiService{
		CryptogotchiRepository: rep,
	}
}

func (svc *CryptogotchiService) GenerateCryptogotchiForUser(user *models.User) (models.Cryptogotchi, error) {
	foodValue := 100.
	foodDrainValue := 0.5
	newCrypt := models.Cryptogotchi{
		// TODO: generate a random name
		Name:               util.Str("Tabito"),
		OwnerId:            user.Id,
		Food:               foodValue,
		FoodDrain:          foodDrainValue,
		PredictedDeathDate: time.Now().Add(time.Duration(foodValue/foodDrainValue) * time.Minute),
	}
	err := svc.Save(&newCrypt)
	return newCrypt, err
}
