package models

import (
	"time"

	"github.com/google/uuid"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/graph/input"
	"gorm.io/datatypes"
)

type EventType string

const (
	FeedEventType   EventType = "feed"
	PlayEventType   EventType = "play"
	CuddleEventType EventType = "cuddle"
)

type Event struct {
	Base
	Type           EventType         `json:"type" gorm:"type:varchar(255)"`
	Payload        datatypes.JSONMap `json:"payload"`
	CryptogotchiId uuid.UUID         `json:"cryptogotchiId" gorm:"type:char(36)"`
}

func (e Event) Apply(c *Cryptogotchi) (bool, time.Time) {
	isAlive, deathDate := c.ProgressUntil(e.CreatedAt)
	if !isAlive {
		return isAlive, deathDate
	}

	switch e.Type {
	case FeedEventType:
		c.Food += 10
	case CuddleEventType:
		c.Affection += 10
	case PlayEventType:
		c.Fun += 10
	}
	return true, time.Time{}
}

func NewEventFromInput(newEvent input.NewEvent) Event {
	return Event{
		Type:           EventType(newEvent.Type),
		Payload:        newEvent.Payload,
		CryptogotchiId: uuid.MustParse(newEvent.CryptogotchiID),
	}
}