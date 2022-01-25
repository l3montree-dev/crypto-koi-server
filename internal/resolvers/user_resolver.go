package resolver

import (
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/service"
)

type UserResolver struct {
	userSvc service.UserSvc
}

func NewUserResolver(userSvc service.UserService) UserResolver {
	return UserResolver{
		userSvc: userSvc,
	}
}
