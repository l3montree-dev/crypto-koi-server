package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/config"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/models"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/util"
)

func TestCryptogotchiProgressUntil(t *testing.T) {
	cryptogotchi := models.Cryptogotchi{
		Name:      util.Str("Tabito"),
		Food:      100,
		FoodDrain: 1,
		Events:    []models.Event{},
		GameStats: []models.GameStat{},
		// set the snapshot validity to now - 10 minutes
		SnapshotValid: time.Now().Add(time.Minute * -10),
		Base: models.Base{
			CreatedAt: time.Now().Add(time.Minute * -10),
		},
		PredictedDeathDate: time.Now().Add(100 * time.Minute),
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
		// died 10 minutes ago
		PredictedDeathDate: time.Now().Add(-10 * time.Minute),
	}

	isAlive, deathDate := cryptogotchi.ProgressUntil(time.Now())

	assert.False(t, isAlive)
	assert.Equal(t, time.Now().Add(time.Minute*-10).Unix(), deathDate.Unix())
}

func TestNextFeedingTime(t *testing.T) {
	cryptogotchi := models.Cryptogotchi{
		Name:      util.Str("Tabito"),
		Food:      10,
		FoodDrain: 1,

		Events:             []models.Event{},
		GameStats:          []models.GameStat{},
		PredictedDeathDate: time.Now().Add(100 * config.TIME_BETWEEN_FEEDINGS),
		Base: models.Base{
			CreatedAt: time.Now().Add(time.Minute * -10),
		},
	}

	nextFeeding := cryptogotchi.GetNextFeedingTime()
	// if no events - than it should be now
	assert.Equal(t, time.Now().Unix(), nextFeeding.Unix())

	// the cryptogotchi was feed 5 minutes ago
	// the next feeding time would be - in 5 minutes
	event := models.Event{
		Type: models.FeedEventType,
	}
	event.Apply(&cryptogotchi)

	nextFeeding = cryptogotchi.GetNextFeedingTime()
	// feeding time should be now + 10 minutes (or whatever the config value is set to)
	assert.Equal(t, time.Now().Add(config.TIME_BETWEEN_FEEDINGS).Unix(), nextFeeding.Unix())
}
