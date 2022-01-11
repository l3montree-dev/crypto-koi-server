package dto

import "gitlab.com/l3montree/cryptogotchi/clodhopper/internal/models"

type EventDTO struct {
	Type    models.EventType       `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

func (e *EventDTO) ToEvent() models.Event {
	return models.Event{
		Type:    e.Type,
		Payload: e.Payload,
	}
}
