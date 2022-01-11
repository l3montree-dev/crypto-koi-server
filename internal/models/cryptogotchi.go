package models

import (
	"time"

	"github.com/google/uuid"
)

type Cryptogotchi struct {
	// aggregate of the Affection, boringness and hunger state variables
	IsAlive bool `json:"isAlive"`

	Base
	Name string `json:"name" gorm:"type:varchar(255) not null"`

	OwnerId uuid.UUID `json:"owner" gorm:"type:char(36) not null"`

	// values between 100 and 0.
	Affection float64 `json:"affection"`
	// values between 100 and 0.
	Fun float64 `json:"fun"`
	// values between 100 and 0.
	Food float64 `json:"food"`

	// the id of the token - might be changed in the future.
	// stored inside the blockchain
	TokenId string `json:"token_id" gorm:"type:varchar(255)"`
	// mapping to the event struct.
	Events []Event `json:"events"`

	LastAggregated time.Time
}

func (c *Cryptogotchi) ToOpenseaNFT() OpenseaNFT {
	return OpenseaNFT{}
}

// the metabolism value is used to calculate the hunger value to a given time.
// if the value is higher, the hunger value will decrease faster
// if the value is lower, the hunger value will decrease slower
func (c *Cryptogotchi) GetMetabolism() float64 {
	return 0.1
}

// the loner value is used to calculate the affection value to a given time.
func (c *Cryptogotchi) GetLonerValue() float64 {
	return 0.1
}

// the need of love value is used to calculate the affection value to a given time.
func (c *Cryptogotchi) GetNeedOfLove() float64 {
	return 0.1
}

// the cryptogotchi has a few time dependentant state variables.
// this function updates the state variables according to the provided time
// returns if the cryptogotchi is still alive
func (c *Cryptogotchi) ProgressUntil(nextTime time.Time) (bool, time.Time) {
	lastAggregatedTime := c.LastAggregated
	if lastAggregatedTime.IsZero() {
		lastAggregatedTime = c.CreatedAt
	}
	// calculate the time difference
	timeDiff := nextTime.Sub(c.LastAggregated)
	// calculate the time difference in seconds
	timeDiffSeconds := timeDiff.Seconds()

	// calculate the food value
	nextFood := c.Food - timeDiffSeconds*c.GetMetabolism()
	// calculate the affection value
	nextAffection := c.Affection - timeDiffSeconds*c.GetLonerValue()
	// calculate the fun value
	nextFun := c.Fun - timeDiffSeconds*c.GetNeedOfLove()

	c.IsAlive = nextFood > 0 && nextAffection > 0 && nextFun > 0

	c.LastAggregated = nextTime

	deathDate := time.Time{}

	if !c.IsAlive {
		// the cryptogotchi did die.
		// lets check the death date.
		seconds := c.Food / c.GetMetabolism()
		deathDate = lastAggregatedTime.Add(time.Duration(seconds) * time.Second)
	}

	c.Food = nextFood
	c.Affection = nextAffection
	c.Fun = nextFun

	return c.IsAlive, deathDate
}

func (c *Cryptogotchi) ReplayEvents() (bool, time.Time) {
	for _, event := range c.Events {
		// mutates the cryptogotchi
		stillAlive, deathDate := event.Apply(c)
		if !stillAlive {
			// the cryptogotchi did die already.
			return stillAlive, deathDate
		}
	}
	return c.IsAlive, time.Time{}
}

func NewCryptogotchi(user *User) Cryptogotchi {
	return Cryptogotchi{OwnerId: user.Id}
}
