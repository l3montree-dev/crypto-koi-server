package config

import (
	"time"
)

type CTX_KEYS string

const (
	USER_CTX_KEY CTX_KEYS = "user"
)

// the time between feedings
const TIME_BETWEEN_FEEDINGS = 1 * time.Hour

// the amount of food the cryptogotchi eats per feeding
const DEFAULT_FEED_VALUE = 50

// the amount of food each cryptogotchi loses per minute
const DEFAULT_FOOD_DRAIN = 100. / (48 * 60) /* 2 days */

// the amount of food each cryptogotchi has when created. Value between 0 and 100
const DEFAULT_FOOD_VALUE = 75.
