package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/controller"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/repositories"
	"gorm.io/gorm"
)

type Gameserver struct {
	db *gorm.DB
}

func NewGameserver(db *gorm.DB) Server {
	return &Gameserver{db: db}
}

func (s *Gameserver) Start() {
	cryptogotchiRepository := repositories.NewGormCryptogotchiRepository(s.db)
	recordRepository := repositories.NewGormRecordRepository(s.db)

	cryptogotchiController := controller.NewCryptogotchiController(recordRepository, cryptogotchiRepository)
	openseaController := controller.NewOpenseaController(recordRepository, cryptogotchiRepository)

	app := fiber.New()

	// register all middlewares
	// allow cross origin request
	app.Use(cors.New())

	// register the controller

	// opensea.io integration.
	// gets called by their API and wallet applications.
	app.Get("integrations/opensea/:tokenId", openseaController.GetCryptogotchi)
}
