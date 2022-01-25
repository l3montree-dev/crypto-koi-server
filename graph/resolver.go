package graph

import (
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/service"
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
}

func NewResolver(
	userSvc service.UserSvc,
	eventSvc service.EventSvc,
	cryptogotchiSvc service.CryptogotchiSvc,
	gameSvc service.GameSvc,
	authSvc service.AuthSvc,
) Resolver {
	return Resolver{
		eventSvc:        eventSvc,
		userSvc:         userSvc,
		cryptogotchiSvc: cryptogotchiSvc,
		gameSvc:         gameSvc,
		authSvc:         authSvc,
	}
}
