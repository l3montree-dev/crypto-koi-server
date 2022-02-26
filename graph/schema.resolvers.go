package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/graph/generated"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/graph/input"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/config"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/cryptokoi"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/db"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/models"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/util"
)

func (r *cryptogotchiResolver) ID(ctx context.Context, obj *models.Cryptogotchi) (string, error) {
	return obj.Id.String(), nil
}

func (r *cryptogotchiResolver) MinutesTillDeath(ctx context.Context, obj *models.Cryptogotchi) (float64, error) {
	return time.Until(obj.PredictedDeathDate).Minutes(), nil
}

func (r *cryptogotchiResolver) MaxLifetimeMinutes(ctx context.Context, obj *models.Cryptogotchi) (float64, error) {
	return 100 / obj.FoodDrain, nil
}

func (r *cryptogotchiResolver) DeathDate(ctx context.Context, obj *models.Cryptogotchi) (*time.Time, error) {
	if time.Now().After(obj.PredictedDeathDate) {
		return &obj.PredictedDeathDate, nil
	}
	return nil, nil
}

func (r *cryptogotchiResolver) NextFeeding(ctx context.Context, obj *models.Cryptogotchi) (*time.Time, error) {
	if obj.LastFed == nil {
		now := time.Now()
		return &now, nil
	}
	next := obj.LastFed.Add(config.TIME_BETWEEN_FEEDINGS)
	return &next, nil
}

func (r *cryptogotchiResolver) Color(ctx context.Context, obj *models.Cryptogotchi) (string, error) {
	tokenId, err := util.UuidToUint256(obj.Id.String())
	if err != nil {
		return "", err
	}
	koi := cryptokoi.NewKoi(tokenId.String())
	attributes := koi.GetAttributes()

	return util.ConvertColor2Hex(attributes.PrimaryColor), nil
}

func (r *cryptogotchiResolver) OwnerAddress(ctx context.Context, obj *models.Cryptogotchi) (string, error) {
	owner, err := r.userSvc.GetById(obj.OwnerId.String())
	return owner.WalletAddress, err
}

func (r *cryptogotchiResolver) Attributes(ctx context.Context, obj *models.Cryptogotchi) (*input.CryptogotchiAttributes, error) {
	uInt, err := util.UuidToUint256(obj.Id.String())
	if err != nil {
		return nil, err
	}

	koi := cryptokoi.NewKoi(uInt.String())

	attributes := koi.GetAttributes()

	return &input.CryptogotchiAttributes{
		PrimaryColor:    util.ConvertColor2Hex(attributes.PrimaryColor),
		BodyColor:       util.ConvertColor2Hex(attributes.BodyColor),
		FinColor:        util.ConvertColor2Hex(attributes.FinColor),
		PatternQuantity: len(attributes.HeadImages) + len(attributes.BodyImages) + len(attributes.FinImages),
		Food:            obj.Food,
		Species:         attributes.KoiType,
		Birthday:        int(obj.CreatedAt.Unix()),
	}, nil
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
	nextFeedingTime := cryptogotchi.GetNextFeedingTime()
	if !nextFeedingTime.Before(time.Now()) {
		return &cryptogotchi, gqlerror.Errorf("it is not time to feed yet")
	}

	// finally feed it.
	feedEvent := models.NewFeedEvent()
	feedEvent.CryptogotchiId = cryptogotchi.Id

	err = r.eventSvc.Save(&feedEvent)

	if err != nil {
		return nil, err
	}

	feedEvent.Apply(&cryptogotchi)

	err = r.cryptogotchiSvc.Save(&cryptogotchi)

	return &cryptogotchi, err
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
		return nil, gqlerror.Errorf("could not find cryptogotchi with id %s", id)
	}

	cryptogotchi.Name = &newName
	err = r.cryptogotchiSvc.Save(&cryptogotchi)
	return &cryptogotchi, err
}

func (r *mutationResolver) ChangeUserName(ctx context.Context, newName string) (*models.User, error) {
	currentUser := ctx.Value(config.USER_CTX_KEY).(*models.User)
	currentUser.Name = &newName
	err := r.userSvc.Save(currentUser)
	return currentUser, err
}

func (r *mutationResolver) GetNftSignature(ctx context.Context, id string, address string) (*input.NftData, error) {
	cryptogotchi, err := r.checkCryptogotchiInteractable(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not find cryptogotchi with id %s", id)
	}

	user := ctx.Value(config.USER_CTX_KEY).(*models.User)
	if user.Id.String() != cryptogotchi.OwnerId.String() {
		// the current user is not the owner.
		// only the owner is allowed to transform the cryptogotchi into an nft.
		return nil, fmt.Errorf("only the owner is allowed to transform the cryptogotchi into an nft")
	}

	signature, tokenId, err := r.cryptokoiApi.GetNftSignatureForCryptogotchi(cryptogotchi.Id.String(), user.WalletAddress)

	if err != nil {
		return nil, err
	}

	return &input.NftData{
		Signature: signature,
		TokenID:   tokenId,
		Address:   user.WalletAddress,
		ChainID:   1,
	}, nil
}

func (r *queryResolver) Leaderboard(ctx context.Context, offset int, limit int) ([]*models.Cryptogotchi, error) {
	cryptogotchis, err := r.cryptogotchiSvc.GetLeaderboard(offset, limit)
	if err != nil {
		return nil, err
	}
	res := make([]*models.Cryptogotchi, len(cryptogotchis))
	for i, c := range cryptogotchis {
		tmp := c
		res[i] = &tmp
	}
	return res, nil
}

func (r *queryResolver) Events(ctx context.Context, cryptogotchiID string, offset int, limit int) ([]*models.Event, error) {
	cryptogotchi, err := r.cryptogotchiSvc.GetCryptogotchiById(cryptogotchiID)
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

	// check if the cryptogotchi belongs to the current user
	if cryptogotchi.OwnerId != currentUser.Id {
		// remove the events from the history
		// privacy policy :-)
		cryptogotchi.Events = nil
	}
	return &cryptogotchi, err
}

func (r *queryResolver) User(ctx context.Context) (*models.User, error) {
	user := ctx.Value(config.USER_CTX_KEY).(*models.User)
	return user, nil
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

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *cryptogotchiResolver) OwnerID(ctx context.Context, obj *models.Cryptogotchi) (string, error) {
	return obj.OwnerId.String(), nil
}
