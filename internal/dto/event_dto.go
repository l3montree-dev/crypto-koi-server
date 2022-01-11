package dto

import (
	"encoding/json"

	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/entities"
	"gorm.io/datatypes"
)

type EventDTO struct {
	Type    entities.EventType     `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

func (e *EventDTO) ToEvent() (entities.Event, error) {
	payloadStr, err := json.Marshal(e.Payload)
	if err != nil {
		return entities.Event{}, err
	}
	return entities.Event{
		Type:    e.Type,
		Payload: datatypes.JSON(payloadStr),
	}, nil
}
