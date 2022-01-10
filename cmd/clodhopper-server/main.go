package main

import (
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/db"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/server"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

func main() {
	db, err := db.NewMySQL(db.MySQLConfig{})
	orchardclient.FailOnError(err, "could not connect to database")

	server := server.NewGameserver(db)
	server.Start()
}
