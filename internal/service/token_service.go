package service

import (
	"crypto/ecdsa"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

type TokenSvc interface {
	CreateSignedToken(claims jwt.Claims) (string, error)
	GetSigningKey() *ecdsa.PrivateKey
	ParseToken(token string) (jwt.Claims, error)
}

type TokenService struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func NewTokenService() TokenSvc {
	privateKeyPath := os.Getenv("PRIVATE_KEY_PATH")
	publicKeyPath := os.Getenv("PUBLIC_KEY_PATH")

	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	orchardclient.FailOnError(err, "Failed to read public key")

	publicKey, err := jwt.ParseECPublicKeyFromPEM(publicKeyBytes)
	orchardclient.FailOnError(err, "Failed to parse public key")

	pem, err := os.ReadFile(privateKeyPath)
	orchardclient.FailOnError(err, "Failed to read private key")
	privateKey, err := jwt.ParseECPrivateKeyFromPEM(pem)
	orchardclient.FailOnError(err, "Failed to parse private key")
	return &TokenService{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (svc *TokenService) CreateSignedToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	return token.SignedString(svc.privateKey)
}

func (svc *TokenService) GetSigningKey() *ecdsa.PrivateKey {
	return svc.privateKey
}

func (svc *TokenService) ParseToken(token string) (jwt.Claims, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return svc.publicKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodES256.Alg()}))

	return t.Claims, err
}
