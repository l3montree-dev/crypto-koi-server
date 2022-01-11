package models

import "time"

type EventType string

const (
	FeedEventType   EventType = "feed"
	PlayEventType   EventType = "play"
	CuddleEventType EventType = "cuddle"
)

type Event struct {
	Base
	Type    EventType              `json:"type" gorm:"type:varchar(255)"`
	Payload map[string]interface{} `json:"payload"`
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
