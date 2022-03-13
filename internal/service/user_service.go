package service

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/repositories"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

type UserSvc interface {
	repositories.UserRepository
}

type UserService struct {
	repositories.UserRepository
	notificationSvc NotificationSvc
	logger          *logrus.Entry
}

func NewUserService(rep repositories.UserRepository) UserSvc {
	logger := orchardclient.Logger.WithField("component", "UserService")
	return &UserService{UserRepository: rep, logger: logger}
}
