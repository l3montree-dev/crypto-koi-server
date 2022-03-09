package graph

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/config"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/cryptokoi"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/generator"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/models"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/service"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	eventSvc        service.EventSvc
	cryptogotchiSvc service.CryptogotchiSvc
	userSvc         service.UserSvc
	gameSvc         service.GameSvc
	authSvc         service.AuthSvc
	cryptokoiApi    cryptokoi.CryptoKoiApi
	generator       generator.Generator
	logger          *logrus.Entry
	chainId         int
}

func NewResolver(
	chainId int,
	userSvc service.UserSvc,
	eventSvc service.EventSvc,
	cryptogotchiSvc service.CryptogotchiSvc,
	gameSvc service.GameSvc,
	authSvc service.AuthSvc,
	cryptokoiApi cryptokoi.CryptoKoiApi,
	generator generator.Generator,
) Resolver {
	return Resolver{
		chainId:         chainId,
		eventSvc:        eventSvc,
		userSvc:         userSvc,
		cryptogotchiSvc: cryptogotchiSvc,
		gameSvc:         gameSvc,
		authSvc:         authSvc,
		cryptokoiApi:    cryptokoiApi,
		generator:       generator,
		logger:          orchardclient.Logger.WithField("package", "graph"),
	}
}

func (r *Resolver) checkCryptogotchiInteractable(ctx context.Context, cryptogotchiId string) (models.Cryptogotchi, error) {
	// check if we are allowed to interact
	cryptogotchi, err := r.cryptogotchiSvc.GetCryptogotchiById(cryptogotchiId)
	if err != nil {
		return cryptogotchi, err
	}
	currentUser := ctx.Value(config.USER_CTX_KEY).(*models.User)
	if cryptogotchi.OwnerId != currentUser.Id {
		return cryptogotchi, gqlerror.Errorf("you are not the owner of this cryptogotchi")
	}

	if !cryptogotchi.IsAlive() {
		return cryptogotchi, (gqlerror.Errorf("this cryptogotchi is already dead"))
	}
	return cryptogotchi, nil
}
