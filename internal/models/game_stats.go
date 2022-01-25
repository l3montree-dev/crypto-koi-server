package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type GameType string

const (
	SNAKE  GameType = "snake"
	HOCKEY GameType = "hockey"
)

func IsGameType(stringToCheck string) (GameType, error) {
	switch GameType(stringToCheck) {
	case SNAKE:
		return SNAKE, nil
	case HOCKEY:
		return HOCKEY, nil
	default:
		return "", fmt.Errorf("unknown game type: %s", stringToCheck)
	}
}

type GameStat struct {
	Base
	CryptogotchiId uuid.UUID  `json:"cryptogotchiId" gorm:"type:varchar(255)"`
	Type           GameType   `json:"type" gorm:"type:varchar(255)"`
	Score          *float64   `json:"score" gorm:"default:null"`
	GameFinished   *time.Time `json:"gameFinished" gorm:"type:datetime;default:null"`
}

// To event returns game won events
func (gameStat *GameStat) ToEvent() (Event, error) {
	if gameStat.Score == nil {
		return Event{}, fmt.Errorf("score is nil")
	}
	return Event{
		Type:           GameWonEventType,
		CryptogotchiId: gameStat.CryptogotchiId,
		Payload:        *gameStat.Score,
	}, nil
}
