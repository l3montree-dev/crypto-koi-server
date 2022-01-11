package controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/db"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/dto"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/models"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/repositories"
)

type CryptogotchiController struct {
	eventRepository        repositories.EventRepository
	cryptogotchiRepository repositories.CryptogotchiRepository
}

func NewCryptogotchiController(eventRepository repositories.EventRepository, cryptogotchiRepository repositories.CryptogotchiRepository) CryptogotchiController {
	return CryptogotchiController{
		eventRepository:        eventRepository,
		cryptogotchiRepository: cryptogotchiRepository,
	}
}

func (c *CryptogotchiController) GetCryptogotchi(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals("currentUser").(*models.User)
	cryptogotchi, err := c.cryptogotchiRepository.GetCryptogotchiByUserId(currentUser.Id.String())
	if db.IsNotFound(err) {
		// the user does not have a cryptogotchi yet.
		cryptogotchi = models.NewCryptogotchi(currentUser)
	}

	cryptogotchi.ReplayEvents()
	return ctx.JSON(cryptogotchi)
}

func (c *CryptogotchiController) HandleNewEvent(ctx *fiber.Ctx) error {
	var body dto.EventDTO

	err := ctx.BodyParser(body)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	currentUser := ctx.Locals("currentUser").(*models.User)
	cryptogotchi, err := c.cryptogotchiRepository.GetCryptogotchiByUserId(currentUser.Id.String())

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "no cryptogotchi found but user tries to add a new event")
	}

	isAlive, deathDate := cryptogotchi.ReplayEvents()
	if !isAlive {
		return fiber.NewError(fiber.StatusBadRequest, "could not add new event. The cryptogotchi died already at: "+deathDate.String())
	}

	// cryptogotchi is still alive.
	// apply the new event and save it inside the database.
	newEvent, err := body.ToEvent()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "could not convert event to entity")
	}
	c.eventRepository.Save(&newEvent)
	newEvent.Apply(&cryptogotchi)

	return ctx.JSON(cryptogotchi)
}
