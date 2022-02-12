package service

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/models"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/repositories"
)

type GameSvc interface {
	repositories.GameStatRepository

	// starting a game will generate a new GameStat instance
	// the GameStat instance is populated with a generated token.
	// the token needs to get resend.
	StartGame(cryptogotchi *models.Cryptogotchi, gameType models.GameType) (models.GameStat, string, error)
	GetGameByToken(token string) (models.GameStat, error)
	FinishGame(token string, score float64) (models.Event, error)
}

type GameService struct {
	repositories.GameStatRepository
	tokenSvc TokenSvc
	eventSvc EventSvc
}

func NewGameService(rep repositories.GameStatRepository, eventSvc EventSvc, tokenSvc TokenSvc) GameSvc {
	return &GameService{
		GameStatRepository: rep,
		tokenSvc:           tokenSvc,
		eventSvc:           eventSvc,
	}
}

func (svc *GameService) StartGame(cryptogotchi *models.Cryptogotchi, gameType models.GameType) (models.GameStat, string, error) {
	// generate a new token.
	gameStat := models.GameStat{
		CryptogotchiId: cryptogotchi.Id,
		Type:           gameType,
	}
	err := svc.Save(&gameStat)

	if err != nil {
		return models.GameStat{}, "", err
	}

	// generate the token based on the game stat.
	claims := jwt.MapClaims{
		"gameStatId": gameStat.Id.String(),
		"exp":        time.Now().Add(time.Hour * 1).Unix(),
	}

	// the token is used to send the game score afterwards.
	// when a game is started, a token gets generated and send to the client
	// when the client would like to complete the game, the token needs to be resend.
	// it can be checked against the token generation time to ensure that the token is valid.
	// the token should be signed.
	token, err := svc.tokenSvc.CreateSignedToken(claims)
	if err != nil {
		return models.GameStat{}, "", err
	}

	return gameStat, token, nil
}

func (svc *GameService) GetGameByToken(token string) (models.GameStat, error) {
	claims, err := svc.tokenSvc.ParseToken(token)
	if err != nil {
		return models.GameStat{}, err
	}

	mapClaims := claims.(jwt.MapClaims)
	// check exp of token.
	if err = mapClaims.Valid(); err != nil {
		return models.GameStat{}, err
	}

	// get the game stat by the token.
	gameStatId := mapClaims["gameStatId"]
	gameStat, err := svc.GetById(gameStatId.(string))

	return gameStat, err
}

func (svc *GameService) FinishGame(token string, score float64) (models.Event, error) {
	game, err := svc.GetGameByToken(token)
	if err != nil {
		return models.Event{}, err
	}

	game.Score = &score
	now := time.Now()
	game.GameFinished = &now
	// create an event from the game stat
	event, err := game.ToEvent()
	if err != nil {
		return models.Event{}, err
	}

	err = svc.Save(&game)
	if err != nil {
		return models.Event{}, err
	}

	err = svc.eventSvc.Save(&event)
	return event, err
}
