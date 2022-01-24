package dto

import (
	"encoding/json"

	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/models"
	"gorm.io/datatypes"
)

type EventDTO struct {
	Type           models.EventType       `json:"type"`
	Payload        map[string]interface{} `json:"payload"`
	CryptogotchiId string                 `json:"cryptogotchiId"`
}

func (e *EventDTO) ToEvent() (models.Event, error) {
	payloadStr, err := json.Marshal(e.Payload)
	if err != nil {
		return models.Event{}, err
	}
	return models.Event{
		Type:    e.Type,
		Payload: datatypes.JSON(payloadStr),
	}, nil
}
