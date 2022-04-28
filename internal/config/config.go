package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"gitlab.com/l3montree/microservices/libs/orchardclient"
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
const DEFAULT_FOOD_DRAIN = 100. / (36 * 60) /* 1.5 days */

// the amount of food each cryptogotchi has when created. Value between 0 and 100
const DEFAULT_FOOD_VALUE = 75.

type Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type NotificationDefinition struct {
	Notifications    []Notification `json:"notifications"`
	HoursBeforeDeath int            `json:"hoursBeforeDeath"`
}

type PreloadedNotifications = map[string]NotificationDefinition

var preloadedNotifications PreloadedNotifications

func loadNotifications() PreloadedNotifications {
	// load the notifications from the json file.
	path := os.Getenv("NOTIFICATION_JSON_FILE_PATH")
	if path == "" {
		orchardclient.Logger.Fatal("NOTIFICATION_JSON_FILE_PATH is not set")
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		orchardclient.Logger.Fatal(err)
	}

	var notifications PreloadedNotifications

	err = json.Unmarshal(b, &notifications)
	if err != nil {
		orchardclient.Logger.Fatal(err)
	}
	return notifications
}

func GetNotifications() PreloadedNotifications {
	if preloadedNotifications == nil {
		preloadedNotifications = loadNotifications()
	}
	return preloadedNotifications
}
