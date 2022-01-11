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
	eventRepository := repositories.NewGormEventRepository(s.db)
	userRepository := repositories.NewGormUserRepository(s.db)

	authController := controller.NewAuthController(userRepository)
	cryptogotchiController := controller.NewCryptogotchiController(eventRepository, cryptogotchiRepository)
	openseaController := controller.NewOpenseaController(eventRepository, cryptogotchiRepository)

	app := fiber.New()

	// register all middlewares
	// allow cross origin request
	app.Use(cors.New())

	// register the controller
	app.Post("/auth/login", authController.Login)
	// opensea.io integration.
	// gets called by their API and wallet applications.
	app.Get("integrations/opensea/:tokenId", openseaController.GetCryptogotchi)
	// add the authentication middleware
	app.Use(authController.AuthMiddleware())
	app.Use(authController.CurrentUserMiddleware())
	app.Get("/cryptogotchi", cryptogotchiController.GetCryptogotchi)
	app.Post("/cryptogotchi", cryptogotchiController.HandleNewEvent)
}
