package controller

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/db"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/dto"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/models"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/repositories"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/service"
)

type AuthController struct {
	authSvc service.AuthSvc
}

func NewAuthController(userRepository repositories.UserRepository) AuthController {

	return AuthController{
		authSvc: service.NewAuthService(userRepository),
	}
}

func (c *AuthController) AuthMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: c.authSvc.GetSigningKey(),
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
		user, err := c.authSvc.GetById(userId.(string))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to get user")
		}
		// set the current user in the locals
		ctx.Locals("currentUser", user)
		return ctx.Next()
	}
}

func (c *AuthController) Refresh(ctx *fiber.Ctx) error {
	var refreshRequest dto.RefreshRequest
	err := ctx.BodyParser(refreshRequest)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse request body")
	}

	user, err := c.authSvc.GetByRefreshToken(refreshRequest.RefreshToken)

	if db.IsNotFound(err) {
		return fiber.NewError(fiber.StatusUnauthorized, "Refresh token not found")
	}

	if err != nil {
		// TODO: Log it.
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get user (refreshing)")
	}

	// create new refresh token for user.
	res, err := c.authSvc.CreateTokenForUser(&user)
	if err != nil {
		// TODO: Log it
		return fiber.NewError(fiber.StatusInternalServerError, "Could not generate tokens")
	}

	return ctx.JSON(res)
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	// check if the user would
	var loginRequest dto.LoginRequest
	err := ctx.BodyParser(loginRequest)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	var user models.User
	switch loginRequest.Type {
	case dto.LoginTypeDeviceId:
		user, err = c.authSvc.GetByDeviceId(loginRequest.DeviceId)
	case dto.LoginTypeWalletAddress:
		user, err = c.authSvc.GetByWalletAddress(loginRequest.WalletAddress)
	}

	if db.IsNotFound(err) {
		// first time the user logs in.
		// create the user
		user = models.User{}
		switch loginRequest.Type {
		case dto.LoginTypeDeviceId:
			user.DeviceId = loginRequest.DeviceId
		case dto.LoginTypeWalletAddress:
			user.WalletAddress = loginRequest.WalletAddress
		}
		err := c.authSvc.Save(&user)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
		}
	} else if err != nil {
		return fiber.NewError(fiber.StatusForbidden, "Invalid credentials")
	}

	// the user is logged in.
	// return a token for the user.
	res, err := c.authSvc.CreateTokenForUser(&user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to sign token")
	}

	return ctx.JSON(res)
}
