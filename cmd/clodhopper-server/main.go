package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/db"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/server"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := db.NewMySQL(db.MySQLConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
		Host:     os.Getenv("DB_HOST"),
	})
	orchardclient.FailOnError(err, "could not connect to database")

	server := server.NewGameserver(db)
	server.Start()
}
