package http_util

import (
	"encoding/json"
	"net/http"

	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/config"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/models"
)

func ParseBody(req *http.Request, v interface{}) error {
	if err := json.NewDecoder(req.Body).Decode(v); err != nil {
		return err
	}
	return nil
}

func GetUserFromContext(r *http.Request) *models.User {
	user := r.Context().Value(config.USER_CTX_KEY)
	if user == nil {
		return nil
	}

	return user.(*models.User)
}
