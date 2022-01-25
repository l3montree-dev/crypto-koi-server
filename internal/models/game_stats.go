package models

import "time"

type GameType string

const (
	SNAKE  GameType = "snake"
	HOCKEY GameType = "hockey"
)

type GameStat struct {
	Base
	CryptogotchiId string     `json:"cryptogotchiId" gorm:"type:varchar(255)"`
	Type           GameType   `json:"type" gorm:"type:varchar(255)"`
	Score          *float64   `json:"score" gorm:"default:null"`
	GameFinished   *time.Time `json:"gameFinished" gorm:"type:datetime;default:null"`
}
