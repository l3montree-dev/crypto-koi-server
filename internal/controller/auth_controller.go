package controller

import (
	"net/http"

	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/db"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/http_dto"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/http_util"
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

func (c *AuthController) Refresh(w http.ResponseWriter, req *http.Request) {
	var refreshRequest http_dto.RefreshRequest
	err := http_util.ParseBody(req, refreshRequest)

	if err != nil {
		http_util.WriteHttpError(w, http.StatusBadRequest, "could not parse body: %e", err)
		return
	}

	user, err := c.authSvc.GetByRefreshToken(refreshRequest.RefreshToken)

	if db.IsNotFound(err) {
		http_util.WriteHttpError(w, http.StatusForbidden, "refresh token not valid: %e", err)
		return
	}

	if err != nil {
		http_util.WriteHttpError(w, http.StatusInternalServerError, "could not get user: %e", err)
		return
	}

	// create new refresh token for user.
	res, err := c.authSvc.CreateTokenForUser(&user)
	if err != nil {
		http_util.WriteHttpError(w, http.StatusInternalServerError, "could not generate tokens: %e", err)
		return
	}

	http_util.WriteJSON(w, res)
}

func (c *AuthController) Login(w http.ResponseWriter, req *http.Request) {
	// check if the user would
	var loginRequest http_dto.LoginRequest
	err := http_util.ParseBody(req, &loginRequest)
	http_util.WriteHttpError(w, http.StatusBadRequest, "could not parse body: %e", err)

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
			user.WalletAddress = loginRequest.WalletAddress
		}
		err := c.authSvc.Save(&user)
		if err != nil {
			http_util.WriteHttpError(w, http.StatusInternalServerError, "could not save user: %e", err)
			return
		}
	} else if err != nil {
		http_util.WriteHttpError(w, http.StatusInternalServerError, "could not get user: %e", err)
		return
	}

	// the user is logged in.
	// return a token for the user.
	res, err := c.authSvc.CreateTokenForUser(&user)
	if err != nil {
		http_util.WriteHttpError(w, http.StatusInternalServerError, "could not generate tokens: %e", err)
		return
	}

	http_util.WriteJSON(w, res)
}
