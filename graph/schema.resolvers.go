package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"sync"
	"time"

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

func (r *cryptogotchiResolver) MinutesTillDeath(ctx context.Context, obj *models.Cryptogotchi) (float64, error) {
	return obj.GetMinutesLeft(), nil
}

func (r *cryptogotchiResolver) MaxLifetimeMinutes(ctx context.Context, obj *models.Cryptogotchi) (float64, error) {
	return 100 / obj.FoodDrain, nil
}

func (r *cryptogotchiResolver) DeathDate(ctx context.Context, obj *models.Cryptogotchi) (*time.Time, error) {
	isAlive, deathDate := obj.ReplayEvents()
	if !isAlive {
		return &deathDate, nil
	}

	isAlive, deathDate = obj.ProgressUntil(time.Now())
	if !isAlive {
		return &deathDate, nil
	}
	return nil, nil
}

func (r *cryptogotchiResolver) OwnerID(ctx context.Context, obj *models.Cryptogotchi) (string, error) {
	return obj.OwnerId.String(), nil
}

func (r *cryptogotchiResolver) NextFeeding(ctx context.Context, obj *models.Cryptogotchi) (*time.Time, error) {
	res := obj.GetNextFeedingTime()
	return &res, nil
}

func (r *eventResolver) ID(ctx context.Context, obj *models.Event) (string, error) {
	return obj.Id.String(), nil
}

func (r *eventResolver) Type(ctx context.Context, obj *models.Event) (string, error) {
	return string(obj.Type), nil
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

func (r *gameStatResolver) CryptogotchiID(ctx context.Context, obj *models.GameStat) (string, error) {
	return obj.CryptogotchiId.String(), nil
}

func (r *mutationResolver) Feed(ctx context.Context, cryptogotchiID string) (*models.Cryptogotchi, error) {
	// check if we are allowed to feed
	cryptogotchi, err := r.checkCryptogotchiInteractable(ctx, cryptogotchiID)

	if err != nil {
		return nil, err
	}
	// the user is allowed to feed it.
	// get the last time it was fed.
	/*nextFeedingTime := cryptogotchi.GetNextFeedingTime()
	if !nextFeedingTime.Before(time.Now()) {
		return &cryptogotchi, gqlerror.Errorf("it is not time to feed yet")
	}*/

	// finally feed it.
	feedEvent := models.NewFeedEvent()
	feedEvent.CryptogotchiId = cryptogotchi.Id

	r.eventSvc.Save(&feedEvent)

	feedEvent.Apply(&cryptogotchi)

	// otherwise other resolver functions might not have the complete event history.
	cryptogotchi.AddEventToHistory(feedEvent)
	return &cryptogotchi, nil
}

func (r *mutationResolver) StartGame(ctx context.Context, cryptogotchiID string, gameType string) (*input.GameStartResponse, error) {
	// start a new game
	cryptogotchi, err := r.checkCryptogotchiInteractable(ctx, cryptogotchiID)
	if err != nil {
		return nil, err
	}

	// check if valid game type.
	parsedGameType, err := models.IsGameType(gameType)
	if err != nil {
		return nil, err
	}

	_, token, err := r.gameSvc.StartGame(&cryptogotchi, models.GameType(parsedGameType))
	if err != nil {
		return nil, err
	}

	return &input.GameStartResponse{
		Token: token,
	}, nil
}

func (r *mutationResolver) FinishGame(ctx context.Context, token string, score float64) (*models.Cryptogotchi, error) {
	game, err := r.gameSvc.GetGameByToken(token)
	if err != nil {
		return nil, err
	}

	// check if the cryptogotchi is interactable
	cryptogotchi, err := r.checkCryptogotchiInteractable(ctx, game.CryptogotchiId.String())
	if err != nil {
		return nil, err
	}

	// finally finish the game
	event, err := r.gameSvc.FinishGame(token, score)
	if err != nil {
		return nil, err
	}

	event.Apply(&cryptogotchi)
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

func (r *queryResolver) Events(ctx context.Context, cryptogotchiID string, offset int, limit int) ([]*models.Event, error) {
	cryptogotchi, err := r.cryptogotchiSvc.GetCryptogotchiByIdWithoutPreload(cryptogotchiID)
	if err != nil {
		return nil, err
	}

	// check if the user is the owner of this cryptogotchi.
	// if not - return an error.
	if cryptogotchi.OwnerId.String() != ctx.Value(config.USER_CTX_KEY).(*models.User).Id.String() {
		return nil, gqlerror.Errorf("you are not the owner of this cryptogotchi")
	}

	events, err := r.eventSvc.GetPaginated(cryptogotchiID, offset, limit)
	if err != nil {
		return nil, err
	}
	eventPointer := make([]*models.Event, len(events))
	for i, event := range events {
		tmp := event
		eventPointer[i] = &tmp
	}
	return eventPointer, nil
}

func (r *queryResolver) Cryptogotchi(ctx context.Context, cryptogotchiID string) (*models.Cryptogotchi, error) {
	currentUser := ctx.Value(config.USER_CTX_KEY).(*models.User)
	cryptogotchi, err := r.cryptogotchiSvc.GetCryptogotchiById(cryptogotchiID)
	if db.IsNotFound(err) {
		return nil, gqlerror.Errorf("could not find cryptogotchi with id %s", cryptogotchiID)
	}
	cryptogotchi.Replay()
	// check if the cryptogotchi belongs to the current user
	if cryptogotchi.OwnerId != currentUser.Id {
		// remove the events from the history
		// privacy policy :-)
		cryptogotchi.Events = nil
	}
	return &cryptogotchi, err
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
				defer wg.Done()
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
