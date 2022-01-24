package resolver

import (
	"context"
	"sync"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/config"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/db"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/dto"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/models"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/repositories"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/service"
)

type CryptogotchiResolver struct {
	eventSvc        service.EventSvc
	cryptogotchiSvc service.CryptogotchiSvc
}

func NewCryptogotchiResolver(eventRepository repositories.EventRepository, cryptogotchiRepository repositories.CryptogotchiRepository) CryptogotchiResolver {
	return CryptogotchiResolver{
		eventSvc:        service.NewEventService(eventRepository),
		cryptogotchiSvc: service.NewCryptogotchiService(cryptogotchiRepository),
	}
}

func (c *CryptogotchiResolver) Cryptogotchis(ctx context.Context) ([]*models.Cryptogotchi, error) {
	currentUser := ctx.Value(config.USER_CTX_KEY).(*models.User)
	cryptogotchies, err := c.cryptogotchiSvc.GetCryptogotchiesByUserId(currentUser.Id.String())
	if db.IsNotFound(err) {
		// the user does not have a cryptogotchi yet.
		cryptogotchies = []models.Cryptogotchi{models.NewCryptogotchi(currentUser)}
	}

	res := make([]*models.Cryptogotchi, len(cryptogotchies))

	// replay all events concurrently
	wg := sync.WaitGroup{}
	for i, cryptogotchi := range cryptogotchies {
		wg.Add(1)
		go func(cryptogotchi models.Cryptogotchi, index int) {
			defer wg.Done()
			res[i] = &cryptogotchi
			res[i].ReplayEvents()
		}(cryptogotchi, i)
	}

	wg.Wait()

	return res, nil
}

func (c *CryptogotchiResolver) HandleNewEvent(ctx *fiber.Ctx) error {
	var body dto.EventDTO

	err := ctx.BodyParser(body)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	currentUser := ctx.Locals("currentUser").(*models.User)
	cryptogotchi, err := c.cryptogotchiSvc.GetCryptogotchiByUserId(currentUser.Id.String())

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
	c.eventSvc.Save(&newEvent)
	newEvent.Apply(&cryptogotchi)

	return ctx.JSON(cryptogotchi)
}
