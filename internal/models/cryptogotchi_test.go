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
		Name:           util.Str("Tabito"),
		Food:           100,
		FoodDrain:      1,
		Fun:            100,
		FunDrain:       1,
		Affection:      100,
		AffectionDrain: 1,
		Events:         []models.Event{},
		GameStats:      []models.GameStat{},
		Base: models.Base{
			CreatedAt: time.Now().Add(time.Minute * -10),
		},
	}

	cryptogotchi.ProgressUntil(time.Now())
	assert.LessOrEqual(t, float64(89), cryptogotchi.Food)
	assert.LessOrEqual(t, float64(89), cryptogotchi.Affection)
	assert.LessOrEqual(t, float64(89), cryptogotchi.Fun)
}

func TestDeathDate(t *testing.T) {
	cryptogotchi := models.Cryptogotchi{
		Name:           util.Str("Tabito"),
		Food:           10,
		FoodDrain:      1,
		Fun:            10,
		FunDrain:       1,
		Affection:      10,
		AffectionDrain: 1,
		Events:         []models.Event{},
		GameStats:      []models.GameStat{},
		Base: models.Base{
			CreatedAt: time.Now().Add(time.Minute * -100),
		},
	}

	isAlive, deathDate := cryptogotchi.ProgressUntil(time.Now())
	assert.False(t, isAlive)
	assert.Equal(t, time.Now().Add(time.Minute*-90).Unix(), deathDate.Unix())
}
