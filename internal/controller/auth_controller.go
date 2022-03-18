package controller

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/db"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/http_dto"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/http_util"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/models"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/repositories"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/service"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

type AuthController struct {
	authSvc         service.AuthSvc
	cryptogotchiSvc service.CryptogotchiSvc
	logger          *logrus.Entry
}

func NewAuthController(userRepository repositories.UserRepository, cryptogotchiSvc service.CryptogotchiSvc, authSvc service.AuthSvc) AuthController {
	return AuthController{
		authSvc:         authSvc,
		cryptogotchiSvc: cryptogotchiSvc,
		logger:          orchardclient.Logger.WithField("component", "AuthController"),
	}
}

func (c *AuthController) Refresh(w http.ResponseWriter, req *http.Request) {
	var refreshRequest http_dto.RefreshRequest
	err := http_util.ParseBody(req, &refreshRequest)

	if err != nil {
		c.logger.Errorf("could not parse body: %e", err)
		http_util.WriteHttpError(w, http.StatusBadRequest, fmt.Sprintf("could not parse body: %e", err))
		return
	}

	user, err := c.authSvc.GetByRefreshToken(refreshRequest.RefreshToken)

	if db.IsNotFound(err) {
		c.logger.Warn("refresh token is not valid")
		http_util.WriteHttpError(w, http.StatusForbidden, fmt.Sprintf("refresh token not valid: %e", err))
		return
	}

	if err != nil {
		c.logger.Errorf("could not get user: %e", err)
		http_util.WriteHttpError(w, http.StatusInternalServerError, fmt.Sprintf("could not get user: %e", err))
		return
	}

	// create new refresh token for user.
	res, err := c.authSvc.CreateTokenForUser(&user)
	if err != nil {
		c.logger.Errorf("could not generate tokens: %e", err)
		http_util.WriteHttpError(w, http.StatusInternalServerError, fmt.Sprintf("could not generate tokens: %e", err))
		return
	}

	http_util.WriteJSON(w, res)
}

func (c *AuthController) DestroyAccount(w http.ResponseWriter, req *http.Request) {
	user := http_util.GetUserFromContext(req)
	if user == nil {
		c.logger.Warn("user is not authenticated")
		http_util.WriteHttpError(w, http.StatusForbidden, "user is not authenticated")
		return
	}

	// delete the user account.
	c.authSvc.Delete(user)
	http_util.WriteJSON(w, http.StatusOK)
}

func (c *AuthController) Login(w http.ResponseWriter, req *http.Request) {
	// check if the user would
	var loginRequest http_dto.LoginRequest
	err := http_util.ParseBody(req, &loginRequest)
	if err != nil {
		c.logger.Warnf("could not parse body: %e", err)
		http_util.WriteHttpError(w, http.StatusBadRequest, fmt.Sprintf("could not parse body: %e", err))
		return
	}

	if loginRequest.WalletAddress == nil && loginRequest.DeviceId == nil {
		c.logger.Warnf("called with empty wallet address and empty device token")
		http_util.WriteHttpError(w, http.StatusBadRequest, "wallet address and device token is empty")
		return
	}

	var user models.User
	if loginRequest.WalletAddress != nil {
		user, err = c.authSvc.GetByWalletAddress(*loginRequest.WalletAddress)
	} else {
		user, err = c.authSvc.GetByDeviceId(*loginRequest.DeviceId)
	}

	if db.IsNotFound(err) {
		// first time the user logs in.
		// create the user
		user = models.User{}
		if loginRequest.WalletAddress != nil {
			user.WalletAddress = loginRequest.WalletAddress
		} else {
			user.DeviceId = loginRequest.DeviceId
		}

		err := c.authSvc.Save(&user)

		if err != nil {
			c.logger.Errorf("could not save user: %e", err)
			http_util.WriteHttpError(w, http.StatusInternalServerError, fmt.Sprintf("could not save user: %e", err))
			return
		}

		// generate a new cryptogotchi for the user.
		_, err = c.cryptogotchiSvc.GenerateCryptogotchiForUser(&user, true)

		if err != nil {
			c.logger.Errorf("could not generate cryptogotchi: %e", err)
			// delete the created user to avoid having a user without a cryptogotchi.
			c.authSvc.Delete(&user)
			http_util.WriteHttpError(w, http.StatusInternalServerError, fmt.Sprintf("could not generate cryptogotchi: %e", err))
			return
		}
	} else if err != nil {
		c.logger.Errorf("could not get user: %e", err)
		http_util.WriteHttpError(w, http.StatusInternalServerError, fmt.Sprintf("could not get user: %e", err))
		return
	}

	// the user is logged in.
	// return a token for the user.
	res, err := c.authSvc.CreateTokenForUser(&user)
	if err != nil {
		c.logger.Errorf("could not generate tokens: %e", err)
		http_util.WriteHttpError(w, http.StatusInternalServerError, fmt.Sprintf("could not generate tokens: %e", err))
		return
	}

	http_util.WriteJSON(w, res)
}
