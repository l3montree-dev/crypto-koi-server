package controller

import (
	"net/http"

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
}

func NewAuthController(userRepository repositories.UserRepository, cryptogotchiSvc service.CryptogotchiSvc, authSvc service.AuthSvc) AuthController {
	return AuthController{
		authSvc:         authSvc,
		cryptogotchiSvc: cryptogotchiSvc,
	}
}

func (c *AuthController) Refresh(w http.ResponseWriter, req *http.Request) {
	var refreshRequest http_dto.RefreshRequest
	err := http_util.ParseBody(req, &refreshRequest)

	if err != nil {
		orchardclient.Logger.Errorf("could not parse body: %e", err)
		http_util.WriteHttpError(w, http.StatusBadRequest, "could not parse body: %e", err)
		return
	}

	user, err := c.authSvc.GetByRefreshToken(refreshRequest.RefreshToken)

	if db.IsNotFound(err) {
		orchardclient.Logger.Warn("refresh token is not valid")
		http_util.WriteHttpError(w, http.StatusForbidden, "refresh token not valid: %e", err)
		return
	}

	if err != nil {
		orchardclient.Logger.Errorf("could not get user: %e", err)
		http_util.WriteHttpError(w, http.StatusInternalServerError, "could not get user: %e", err)
		return
	}

	// create new refresh token for user.
	res, err := c.authSvc.CreateTokenForUser(&user)
	if err != nil {
		orchardclient.Logger.Errorf("could not generate tokens: %e", err)
		http_util.WriteHttpError(w, http.StatusInternalServerError, "could not generate tokens: %e", err)
		return
	}

	http_util.WriteJSON(w, res)
}

func (c *AuthController) Login(w http.ResponseWriter, req *http.Request) {
	// check if the user would
	var loginRequest http_dto.LoginRequest
	err := http_util.ParseBody(req, &loginRequest)
	if err != nil {
		orchardclient.Logger.Warnf("could not parse body: %e", err)
		http_util.WriteHttpError(w, http.StatusBadRequest, "could not parse body: %e", err)
		return
	}

	var user models.User
	switch loginRequest.Type {
	case http_dto.LoginTypeDeviceId:
		user, err = c.authSvc.GetByDeviceId(loginRequest.DeviceId)
	case http_dto.LoginTypeWalletAddress:
		user, err = c.authSvc.GetByWalletAddress(loginRequest.WalletAddress)
	}

	if db.IsNotFound(err) {
		// first time the user logs in.
		// create the user
		user = models.User{}
		switch loginRequest.Type {
		case http_dto.LoginTypeDeviceId:
			user.DeviceId = loginRequest.DeviceId
		case http_dto.LoginTypeWalletAddress:
			user.WalletAddress = &loginRequest.WalletAddress
		}
		err := c.authSvc.Save(&user)

		if err != nil {
			orchardclient.Logger.Errorf("could not save user: %e", err)
			http_util.WriteHttpError(w, http.StatusInternalServerError, "could not save user: %e", err)
			return
		}

		// generate a new cryptogotchi for the user.
		_, err = c.cryptogotchiSvc.GenerateCryptogotchiForUser(&user)

		if err != nil {
			orchardclient.Logger.Errorf("could not generate cryptogotchi: %e", err)
			// delete the created user to avoid having a user without a cryptogotchi.
			c.authSvc.Delete(&user)
			http_util.WriteHttpError(w, http.StatusInternalServerError, "could not generate cryptogotchi: %e", err)
			return
		}
	} else if err != nil {
		orchardclient.Logger.Errorf("could not get user: %e", err)
		http_util.WriteHttpError(w, http.StatusInternalServerError, "could not get user: %e", err)
		return
	}

	// the user is logged in.
	// return a token for the user.
	res, err := c.authSvc.CreateTokenForUser(&user)
	if err != nil {
		orchardclient.Logger.Errorf("could not generate tokens: %e", err)
		http_util.WriteHttpError(w, http.StatusInternalServerError, "could not generate tokens: %e", err)
		return
	}

	http_util.WriteJSON(w, res)
}
