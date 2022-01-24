package resolver

import (
	"context"
	"sync"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/graph/input"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/config"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/db"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/models"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/repositories"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/service"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
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

func (c *CryptogotchiResolver) HandleNewEvent(ctx context.Context, event input.NewEvent) (*models.Cryptogotchi, error) {

	currentUser := ctx.Value(config.USER_CTX_KEY).(*models.User)
	cryptogotchi, err := c.cryptogotchiSvc.GetCryptogotchiById(event.CryptogotchiID)

	if err != nil {
		orchardclient.Logger.Warnf("cryptogotchi not found: %e", err)
		return nil, gqlerror.Errorf("could not find cryptogotchi with id %s", event.CryptogotchiID)
	}

	// check if the user is allowed to update the cryptogotchi
	if cryptogotchi.OwnerId != currentUser.Id {
		orchardclient.Logger.Warnf("user %s is not allowed to update cryptogotchi %s", currentUser.Id, cryptogotchi.Id)
		return nil, gqlerror.Errorf("user %s is not allowed to update cryptogotchi %s", currentUser.Id, event.CryptogotchiID)
	}

	isAlive, _ := cryptogotchi.ReplayEvents()
	if !isAlive {
		return nil, gqlerror.Errorf("cryptogotchi is dead")
	}

	// cryptogotchi is still alive.
	// apply the new event and save it inside the database.
	newEvent := models.NewEventFromInput(event)

	c.eventSvc.Save(&newEvent)
	newEvent.Apply(&cryptogotchi)

	return &cryptogotchi, nil
}
