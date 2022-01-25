package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"sync"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/graph/generated"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/graph/input"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/config"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/db"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/models"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

func (r *cryptogotchiResolver) ID(ctx context.Context, obj *models.Cryptogotchi) (string, error) {
	return obj.Id.String(), nil
}

func (r *cryptogotchiResolver) OwnerID(ctx context.Context, obj *models.Cryptogotchi) (string, error) {
	return obj.OwnerId.String(), nil
}

func (r *eventResolver) ID(ctx context.Context, obj *models.Event) (string, error) {
	return obj.Id.String(), nil
}

func (r *eventResolver) Type(ctx context.Context, obj *models.Event) (string, error) {
	return string(obj.Type), nil
}

func (r *eventResolver) Payload(ctx context.Context, obj *models.Event) (map[string]interface{}, error) {
	return obj.Payload, nil
}

func (r *eventResolver) CryptogotchiID(ctx context.Context, obj *models.Event) (string, error) {
	return obj.CryptogotchiId.String(), nil
}

func (r *gameStatResolver) ID(ctx context.Context, obj *models.GameStat) (string, error) {
	return obj.Id.String(), nil
}

func (r *gameStatResolver) Type(ctx context.Context, obj *models.GameStat) (string, error) {
	return string(obj.Type), nil
}

func (r *mutationResolver) HandleNewEvent(ctx context.Context, event input.NewEvent) (*models.Cryptogotchi, error) {
	currentUser := ctx.Value(config.USER_CTX_KEY).(*models.User)
	cryptogotchi, err := r.cryptogotchiSvc.GetCryptogotchiById(event.CryptogotchiID)

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

	r.eventSvc.Save(&newEvent)
	newEvent.Apply(&cryptogotchi)

	return &cryptogotchi, nil
}

func (r *mutationResolver) ChangeCryptogotchiName(ctx context.Context, id string, newName string) (*models.Cryptogotchi, error) {
	cryptogotchi, err := r.cryptogotchiSvc.GetCryptogotchiById(id)
	if err != nil {
		orchardclient.Logger.Errorf("could not find cryptogotchi with id:%s", id)
		return nil, gqlerror.Errorf("could not find cryptogotchi with id %s", id)
	}

	// check if cryptogotchi is alive.
	cryptogotchi.Replay()
	if !cryptogotchi.IsAlive {
		return nil, gqlerror.Errorf("cryptogotchi is dead")
	}

	cryptogotchi.Name = &newName
	err = r.cryptogotchiSvc.Save(&cryptogotchi)
	return &cryptogotchi, err
}

func (r *queryResolver) Cryptogotchies(ctx context.Context) ([]*models.Cryptogotchi, error) {
	currentUser := ctx.Value(config.USER_CTX_KEY).(*models.User)
	cryptogotchies, err := r.cryptogotchiSvc.GetCryptogotchiesByUserId(currentUser.Id.String())
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
			res[index] = cryptogotchi.Replay()
		}(cryptogotchi, i)
	}

	wg.Wait()

	return res, nil
}

func (r *queryResolver) User(ctx context.Context) (*models.User, error) {
	res := make(chan *models.User)
	go func() {
		user := ctx.Value(config.USER_CTX_KEY).(*models.User)
		// replay all events concurrently
		wg := sync.WaitGroup{}
		for i, cryptogotchi := range user.Cryptogotchies {
			wg.Add(1)
			go func(c models.Cryptogotchi, index int) {
				user.Cryptogotchies[index] = *c.Replay()
			}(cryptogotchi, i)

		}

		wg.Wait()
		res <- user
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-res:
		return result, nil
	}
}

func (r *userResolver) ID(ctx context.Context, obj *models.User) (string, error) {
	return obj.Id.String(), nil
}

// Cryptogotchi returns generated.CryptogotchiResolver implementation.
func (r *Resolver) Cryptogotchi() generated.CryptogotchiResolver { return &cryptogotchiResolver{r} }

// Event returns generated.EventResolver implementation.
func (r *Resolver) Event() generated.EventResolver { return &eventResolver{r} }

// GameStat returns generated.GameStatResolver implementation.
func (r *Resolver) GameStat() generated.GameStatResolver { return &gameStatResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type cryptogotchiResolver struct{ *Resolver }
type eventResolver struct{ *Resolver }
type gameStatResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
