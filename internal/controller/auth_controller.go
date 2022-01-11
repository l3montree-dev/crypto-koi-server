package controller

import (
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/db"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/dto"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/entities"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/repositories"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

type AuthController struct {
	userRepository repositories.UserRepository
	privateKey     []byte
}

func NewAuthController(userRepository repositories.UserRepository) AuthController {

	privateKeyPath := os.Getenv("PRIVATE_KEY_PATH")

	privateKey, err := os.ReadFile(privateKeyPath)
	orchardclient.FailOnError(err, "Failed to read private key")

	return AuthController{
		userRepository: userRepository,
		privateKey:     privateKey,
	}
}

func (c *AuthController) GetSigningKey() []byte {
	return c.privateKey
}

func (c *AuthController) AuthMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: c.privateKey,
	})
}

func (c *AuthController) CurrentUserMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token := ctx.Locals("user").(*jwt.Token)
		if token == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}
		claims := token.Claims.(jwt.MapClaims)
		userId := claims["id"]
		user, err := c.userRepository.GetById(userId.(string))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to get user")
		}
		// set the current user in the locals
		ctx.Locals("currentUser", user)
		return ctx.Next()
	}
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	// check if the user would
	var loginRequest dto.LoginRequest
	err := ctx.BodyParser(loginRequest)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	var user entities.User
	switch loginRequest.Type {
	case "deviceId":
		user, err = c.userRepository.GetByDeviceId(loginRequest.DeviceId)
	case "walletAddress":
		user, err = c.userRepository.GetByWalletAddress(loginRequest.WalletAddress)
	}

	if db.IsNotFound(err) {
		// first time the user logs in.
		// create the user
		user = entities.User{}
		switch loginRequest.Type {
		case "deviceId":
			user.DeviceId = loginRequest.DeviceId
		case "walletAddress":
			user.WalletAddress = loginRequest.WalletAddress
		}
		err := c.userRepository.Save(&user)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
		}
	} else if err != nil {
		return fiber.NewError(fiber.StatusForbidden, "Invalid credentials")
	}

	// the user is logged in.
	// return a token for the user.
	claims := jwt.MapClaims{
		"id": user.Id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	t, err := token.SignedString(c.privateKey)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to sign token")
	}

	return ctx.JSON(fiber.Map{
		"accessToken": t,
	})
}
