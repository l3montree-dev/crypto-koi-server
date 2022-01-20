package service

import (
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/dto"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/models"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/repositories"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

type AuthSvc interface {
	repositories.UserRepository
	CreateTokenForUser(user *models.User) (dto.TokenResponse, error)
	GetSigningKey() []byte
}

type AuthService struct {
	repositories.UserRepository
	privateKey []byte
}

func NewAuthService(rep repositories.UserRepository) AuthSvc {
	privateKeyPath := os.Getenv("PRIVATE_KEY_PATH")

	privateKey, err := os.ReadFile(privateKeyPath)
	orchardclient.FailOnError(err, "Failed to read private key")

	return &AuthService{
		UserRepository: rep,
		privateKey:     privateKey,
	}
}

func (svc *AuthService) GetSigningKey() []byte {
	return svc.privateKey
}

func (svc *AuthService) generateRefreshToken() string {
	return uuid.NewString()
}

func (svc *AuthService) CreateTokenForUser(user *models.User) (dto.TokenResponse, error) {
	claims := jwt.MapClaims{
		"id": user.Id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	t, err := token.SignedString(svc.privateKey)

	if err != nil {
		return dto.TokenResponse{}, err
	}

	// create a refresh token for the user.
	user.RefreshToken = svc.generateRefreshToken()

	return dto.TokenResponse{
		AccessToken:  t,
		RefreshToken: user.RefreshToken,
	}, nil
}
