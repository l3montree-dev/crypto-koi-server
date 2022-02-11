package models

import (
	"time"

	"github.com/google/uuid"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/config"
)

type Cryptogotchi struct {
	Base
	Name               *string    `json:"name" gorm:"type:varchar(255);default:null"`
	OwnerId            uuid.UUID  `json:"owner" gorm:"type:char(36); not null"`
	PredictedDeathDate time.Time  `json:"-" gorm:"not null"`
	LastFed            *time.Time `json:"-" gorm:"default:null"`
	// values between 100 and 0.
	Food float64 `json:"food" gorm:"default:100"`
	// drain per minute
	FoodDrain float64 `json:"foodDrain" gorm:"default:0.5"`
	// the id of the token - might be changed in the future.
	// mapping to the event struct.
	Events    []Event    `json:"events"`
	GameStats []GameStat `json:"game_stats" gorm:"foreignKey:cryptogotchi_id"`
	// the timestamp of the current snapshot stored inside the database.
	// in most cases this equals the LastFeed value. - nevertheless to build the struct a bit more
	// future proof, we store the timestamp of the snapshot in the database as a separate column.
	// currently this affects only the food value.
	SnapshotValid time.Time `json:"-" gorm:"not null"`
}

func (c *Cryptogotchi) ToOpenseaNFT() OpenseaNFT {
	return OpenseaNFT{}
}

// make sure to only call this function after the food value has been updated.
func (c *Cryptogotchi) PredictNewDeathDate() time.Time {
	return time.Now().Add(time.Duration(c.Food/c.FoodDrain) * time.Minute)
}

func (c *Cryptogotchi) IsAlive() bool {
	return c.PredictedDeathDate.After(time.Now())
}

// the cryptogotchi has a few time dependentant state variables.
// this function updates the state variables according to the provided time
// returns if the cryptogotchi is still alive
func (c *Cryptogotchi) ProgressUntil(nextTime time.Time) (bool, time.Time) {
	if c.PredictedDeathDate.Before(nextTime) {
		c.Food = 0
		return false, c.PredictedDeathDate
	}
	// calculate the time difference
	timeDiffMinutes := nextTime.Sub(c.SnapshotValid).Minutes()

	// calculate the food value
	// just calculating the current state does not change the predicted death date.
	c.Food -= timeDiffMinutes * c.FoodDrain
	// update the snapshot validity since we mutated the cryptogotchi
	c.SnapshotValid = nextTime
	return true, time.Time{}
}

func (c *Cryptogotchi) GetNextFeedingTime() time.Time {

	if c.LastFed == nil {
		return time.Now()
	}

	return (*c.LastFed).Add(config.TIME_BETWEEN_FEEDINGS)
}

func NewCryptogotchi(user *User) Cryptogotchi {
	return Cryptogotchi{OwnerId: user.Id}
}
