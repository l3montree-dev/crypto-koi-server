package service

import (
	"crypto/ecdsa"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/http_dto"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/models"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/repositories"
)

type AuthSvc interface {
	repositories.UserRepository
	CreateTokenForUser(user *models.User) (http_dto.TokenResponse, error)
	GetSigningKey() *ecdsa.PrivateKey
}

type AuthService struct {
	repositories.UserRepository
	tokenSvc TokenSvc
}

func NewAuthService(rep repositories.UserRepository) AuthSvc {
	return &AuthService{
		UserRepository: rep,
		tokenSvc:       NewTokenService(),
	}
}

func (svc *AuthService) GetSigningKey() *ecdsa.PrivateKey {
	return svc.tokenSvc.GetSigningKey()
}

func (svc *AuthService) generateRefreshToken() string {
	return uuid.NewString()
}

func (svc *AuthService) CreateTokenForUser(user *models.User) (http_dto.TokenResponse, error) {
	claims := jwt.RegisteredClaims{
		Subject:  user.Id.String(),
		Issuer:   "clodhopper",
		Audience: []string{"cattleshow-app"},
	}

	t, err := svc.tokenSvc.CreateSignedToken(claims)

	if err != nil {
		return http_dto.TokenResponse{}, err
	}

	// create a refresh token for the user.
	user.RefreshToken = svc.generateRefreshToken()
	// store the token inside the database
	svc.Save(user)

	return http_dto.TokenResponse{
		AccessToken:  t,
		RefreshToken: user.RefreshToken,
	}, nil
}
