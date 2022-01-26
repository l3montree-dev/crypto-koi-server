package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/models"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/util"
)

func TestCryptogotchiProgressUntil(t *testing.T) {
	cryptogotchi := models.Cryptogotchi{
		Name:      util.Str("Tabito"),
		Food:      100,
		FoodDrain: 1,
		Events:    []models.Event{},
		GameStats: []models.GameStat{},
		Base: models.Base{
			CreatedAt: time.Now().Add(time.Minute * -10),
		},
	}

	cryptogotchi.ProgressUntil(time.Now())
	assert.LessOrEqual(t, float64(89), cryptogotchi.Food)

}

func TestDeathDate(t *testing.T) {
	cryptogotchi := models.Cryptogotchi{
		Name:      util.Str("Tabito"),
		Food:      10,
		FoodDrain: 1,

		Events:    []models.Event{},
		GameStats: []models.GameStat{},
		Base: models.Base{
			CreatedAt: time.Now().Add(time.Minute * -100),
		},
	}

	isAlive, deathDate := cryptogotchi.ProgressUntil(time.Now())
	assert.False(t, isAlive)
	assert.Equal(t, time.Now().Add(time.Minute*-90).Unix(), deathDate.Unix())
}

func generateFeedEvent(createdAt time.Time) models.Event {
	return models.Event{
		Type: models.FeedEventType,
		Base: models.Base{
			CreatedAt: createdAt,
		},
	}
}
func TestNextFeedingTime(t *testing.T) {
	cryptogotchi := models.Cryptogotchi{
		Name:      util.Str("Tabito"),
		Food:      10,
		FoodDrain: 1,

		Events:    []models.Event{},
		GameStats: []models.GameStat{},
		Base: models.Base{
			CreatedAt: time.Now().Add(time.Minute * -10),
		},
	}

	nextFeeding := cryptogotchi.GetNextFeedingTime()
	// if no events - than it should be now
	assert.Equal(t, time.Now().Unix(), nextFeeding.Unix())

	// the cryptogotchi was feed 5 minutes ago
	// the next feeding time would be - in 5 minutes
	cryptogotchi.Events = append(cryptogotchi.Events, generateFeedEvent(time.Now().Add(time.Minute*-5)))
	nextFeeding = cryptogotchi.GetNextFeedingTime()
	assert.Equal(t, time.Now().Add(time.Minute*5).Unix(), nextFeeding.Unix())

	// the cryptogotchi was feed 20 minutes ago and 5 minutes ago - the next feeding time would still be in 5 minutes
	cryptogotchi.Events = append(cryptogotchi.Events, generateFeedEvent(time.Now().Add(time.Minute*-20)))
	nextFeeding = cryptogotchi.GetNextFeedingTime()
	assert.Equal(t, time.Now().Add(time.Minute*5).Unix(), nextFeeding.Unix())
}
