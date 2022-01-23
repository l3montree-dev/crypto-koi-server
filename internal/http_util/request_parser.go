package http_util

import (
	"encoding/json"
	"net/http"
)

func ParseBody(req *http.Request, v interface{}) error {
	if err := json.NewDecoder(req.Body).Decode(v); err != nil {
		return err
	}
	return nil
}
