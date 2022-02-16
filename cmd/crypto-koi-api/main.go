package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/db"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/server"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

type SentryErrorLoggingHook struct {
}

func (hook *SentryErrorLoggingHook) Fire(entry *logrus.Entry) error {
	sentry.CaptureMessage(entry.Message)
	return nil
}

func (hook *SentryErrorLoggingHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	}
}

func mustReadFile(filepath string) []byte {
	bytes, err := ioutil.ReadFile(filepath)
	orchardclient.FailOnError(err, fmt.Sprintf("could not read file: %s", filepath))
	return bytes
}

func main() {
	err := godotenv.Load()
	if err != nil {
		orchardclient.Logger.Fatal("Error loading .env file")
	}

	isDev := os.Getenv("DEV") == "true"
	if !isDev {
		err = sentry.Init(sentry.ClientOptions{
			Dsn: "https://e56b8f4eedcf451e9b1cec93799f4443@sentry.l3montree.com/50",
		})
		if err != nil {
			orchardclient.Logger.Fatalf("sentry.Init: %s", err)
		}
	}

	orchardclient.Logger.AddHook(&SentryErrorLoggingHook{})

	db, err := db.NewMySQL(db.MySQLConfig{
		User:     os.Getenv("DB_USER"),
		Password: strings.TrimSpace(string(mustReadFile(os.Getenv("DB_PASSWORD_FILE_PATH")))),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
		Host:     os.Getenv("DB_HOST"),
	})
	orchardclient.FailOnError(err, "could not connect to database")

	baseImagePath := os.Getenv("BASE_IMAGE_PATH")
	if baseImagePath == "" {
		orchardclient.Logger.Fatal("BASE_IMAGE_PATH is not set")
	}
	server := server.NewGraphqlServer(db, baseImagePath)
	server.Start()
}