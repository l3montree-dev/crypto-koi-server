package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/config"
)

type EventType string

const (
	FeedEventType    EventType = "feed"
	GameWonEventType EventType = "game-won"
	// PlayEventType   EventType = "play"
	// CuddleEventType EventType = "cuddle"
)

func IsEventType(stringToCheck string) (EventType, error) {
	switch EventType(stringToCheck) {
	case FeedEventType:
		return FeedEventType, nil
	case GameWonEventType:
		return GameWonEventType, nil
	default:
		return "", fmt.Errorf("unknown event type: %s", stringToCheck)
	}
}

type Event struct {
	Base
	Type           EventType `json:"type" gorm:"type:varchar(255)"`
	CryptogotchiId uuid.UUID `json:"cryptogotchiId" gorm:"type:char(36)"`
	// the value to increment.
	// a regular feed event will contain the value 10
	Payload float64
}

func (e Event) Apply(c *Cryptogotchi) (bool, time.Time) {
	isAlive, deathDate := c.ProgressUntil(e.CreatedAt)
	if !isAlive {
		return isAlive, deathDate
	}

	c.Food += e.Payload
	// limit the food to 100
	if c.Food > 100 {
		c.Food = 100
	} else if c.Food < 0 {
		c.Food = 0
	}

	c.PredictedDeathDate = c.PredictNewDeathDate()
	now := time.Now()
	c.LastFed = &now
	return true, time.Time{}
}

func NewFeedEvent() Event {
	return Event{
		Type:    FeedEventType,
		Payload: config.DEFAULT_FEED_VALUE,
	}
}
