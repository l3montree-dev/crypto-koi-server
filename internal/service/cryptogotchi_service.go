package service

import "gitlab.com/l3montree/cryptogotchi/clodhopper/internal/repositories"

type CryptogotchiSvc interface {
	repositories.CryptogotchiRepository
}

type CryptogotchiService struct {
	repositories.CryptogotchiRepository
}

func NewCryptogotchiService(rep repositories.CryptogotchiRepository) CryptogotchiSvc {
	return &CryptogotchiService{
		CryptogotchiRepository: rep,
	}
}
