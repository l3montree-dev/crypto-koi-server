package models

type GameType string

const (
	SNAKE  GameType = "snake"
	HOCKEY GameType = "hockey"
)

type GameStat struct {
	Base
	UserId string   `json:"user_id" gorm:"type:varchar(255)"`
	Type   GameType `json:"type" gorm:"type:varchar(255)"`
	Score  int      `json:"score" gorm:"type:int"`
}
