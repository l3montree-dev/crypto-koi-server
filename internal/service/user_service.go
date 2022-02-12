package service

import "gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/repositories"

type UserSvc interface {
	repositories.UserRepository
}

type UserService struct {
	repositories.UserRepository
}

func NewUserService(rep repositories.UserRepository) UserSvc {
	return &UserService{UserRepository: rep}
}
