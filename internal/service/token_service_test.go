package service_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/service"
)

func TestTokenService(t *testing.T) {
	privKeyPath, _ := filepath.Abs(filepath.Join("../../testdata/key.pem"))
	pubKeyPath, _ := filepath.Abs(filepath.Join("../../testdata/public.pem"))
	os.Setenv("PRIVATE_KEY_PATH", privKeyPath)
	os.Setenv("PUBLIC_KEY_PATH", pubKeyPath)
	tokenSvc := service.NewTokenService()

	claims := jwt.MapClaims{
		"sub": "1234567890",
	}
	signedToken, err := tokenSvc.CreateSignedToken(claims)
	if err != nil {
		t.Fatalf("%v", err)
	}
	parsed, err := tokenSvc.ParseToken(signedToken)
	if err != nil {
		t.Fatalf("%v", err)
	}
	assert.Equal(t, "1234567890", parsed.(jwt.MapClaims)["sub"])

}
