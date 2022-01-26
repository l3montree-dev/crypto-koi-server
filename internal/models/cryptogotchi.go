package models

import (
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/config"
)

type Cryptogotchi struct {
	Base

	// aggregate of the Affection, boringness and hunger state variables
	IsAlive bool `json:"isAlive" gorm:"default:true"`

	Name *string `json:"name" gorm:"type:varchar(255);default:null"`

	OwnerId uuid.UUID `json:"owner" gorm:"type:char(36); not null"`

	// values between 100 and 0.
	// Affection float64 `json:"affection" gorm:"default:100"`
	// values between 100 and 0.
	// Fun float64 `json:"fun" gorm:"default:100"`
	// values between 100 and 0.
	Food float64 `json:"food" gorm:"default:100"`

	// drain per minute
	FoodDrain float64 `json:"foodDrain" gorm:"default:0.5"`
	// drain per minute
	// FunDrain float64 `json:"funDrain" gorm:"default:0.5"`
	// drain per minute
	// AffectionDrain float64 `json:"affectionDrain" gorm:"default:0.5"`

	// the id of the token - might be changed in the future.
	// stored inside the blockchain
	TokenId *string `json:"token_id" gorm:"type:varchar(255);default:null"`
	// mapping to the event struct.
	Events []Event `json:"events"`

	GameStats []GameStat `json:"game_stats" gorm:"foreignKey:cryptogotchi_id"`

	LastAggregated      time.Time `json:"-" gorm:"-"`
	ProcessedEventsTill time.Time `json:"-" gorm:"-"`
}

func (c *Cryptogotchi) ToOpenseaNFT() OpenseaNFT {
	return OpenseaNFT{}
}

// the metabolism value is used to calculate the hunger value to a given time.
// if the value is higher, the Food value will decrease faster
// if the value is lower, the Food value will decrease slower
func (c *Cryptogotchi) GetMetabolism() float64 {
	return c.FoodDrain
}

func (c *Cryptogotchi) GetMinutesLeft() float64 {
	c.Replay()
	return c.Food / c.GetMetabolism()
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
	timeDiff := nextTime.Sub(lastAggregatedTime)
	// calculate the time difference in seconds
	timeDiffMinutes := timeDiff.Minutes()

	// calculate the food value
	nextFood := c.Food - timeDiffMinutes*c.GetMetabolism()
	// calculate the affection value
	c.IsAlive = nextFood > 0

	c.LastAggregated = nextTime

	deathDate := time.Time{}

	if !c.IsAlive {
		// the cryptogotchi did die.
		// lets check the death date.
		minutes := c.Food / c.GetMetabolism()
		deathDate = lastAggregatedTime.Add(time.Duration(minutes) * time.Minute)
	}

	c.Food = nextFood
	return c.IsAlive, deathDate
}

func (c *Cryptogotchi) ReplayEvents() (bool, time.Time) {
	for _, event := range c.Events {
		if c.ProcessedEventsTill.Before(event.CreatedAt) {
			// mutates the cryptogotchi
			stillAlive, deathDate := event.Apply(c)
			c.ProcessedEventsTill = event.CreatedAt
			if !stillAlive {
				// the cryptogotchi did die already.
				return stillAlive, deathDate
			}
		}
	}
	return true, time.Time{}
}

func (c *Cryptogotchi) GetNextFeedingTime() time.Time {
	c.sortEvents()
	// get the last feeding event
	// iterate over the slice in the other direction
	for i := len(c.Events) - 1; i >= 0; i-- {
		event := c.Events[i]
		if event.Type == FeedEventType {
			// this is the last event.
			fmt.Println("Decided for event: ", event.Id, "Created at: ", event.CreatedAt, "therefore next time:", event.CreatedAt.Add(config.TIME_BETWEEN_FEEDINGS))
			return event.CreatedAt.Add(config.TIME_BETWEEN_FEEDINGS)
		}
	}
	// if there is no event - just return the current time
	return time.Now()
}

func (c *Cryptogotchi) sortEvents() {
	// order by createdAt
	sort.SliceStable(c.Events, func(i, j int) bool {
		return c.Events[i].CreatedAt.Before(c.Events[j].CreatedAt)
	})
}

func (c *Cryptogotchi) AddEventToHistory(event Event) {
	c.Events = append(c.Events, event)
	c.sortEvents()
}

func (c *Cryptogotchi) Replay() *Cryptogotchi {
	// replay the events
	stillAlive, _ := c.ReplayEvents()
	if !stillAlive {
		c.IsAlive = false
		return c
	}
	c.ProgressUntil(time.Now())
	return c
}

func NewCryptogotchi(user *User) Cryptogotchi {
	return Cryptogotchi{OwnerId: user.Id}
}
