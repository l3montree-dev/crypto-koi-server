// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package input

import (
	"time"
)

type NewEvent struct {
	CryptogotchiID string                 `json:"cryptogotchiId"`
	Type           string                 `json:"type"`
	Payload        map[string]interface{} `json:"payload"`
	CreatedAt      time.Time              `json:"createdAt"`
	UpdatedAt      time.Time              `json:"updatedAt"`
}