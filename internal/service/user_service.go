package service

import "gitlab.com/l3montree/cryptogotchi/clodhopper/internal/repositories"

type UserSvc interface {
	repositories.UserRepository
}

type UserService struct {
	repositories.UserRepository
}

func NewUserService(rep repositories.UserRepository) UserSvc {
	return &UserService{UserRepository: rep}
}
