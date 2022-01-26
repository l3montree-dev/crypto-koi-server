package config

import (
	"time"
)

type CTX_KEYS string

const (
	USER_CTX_KEY CTX_KEYS = "user"
)

const TIME_BETWEEN_FEEDINGS = 10 * time.Minute
const DEFAULT_FEED_VALUE = 10
