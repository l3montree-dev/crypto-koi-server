package service

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/models"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/repositories"
)

type GameSvc interface {
	repositories.GameStatRepository

	// starting a game will generate a new GameStat instance
	// the GameStat instance is populated with a generated token.
	// the token needs to get resend.
	StartGame(cryptogotchi *models.Cryptogotchi, gameType models.GameType) (models.GameStat, string, error)
	GetGameByToken(token string) (models.GameStat, error)
	StopGame(token string, score float64) error
}

type GameService struct {
	repositories.GameStatRepository
	tokenSvc TokenSvc
}

func NewGameService(rep repositories.GameStatRepository, tokenSvc TokenSvc) GameSvc {
	return &GameService{
		GameStatRepository: rep,
		tokenSvc:           tokenSvc,
	}
}

func (svc *GameService) StartGame(cryptogotchi *models.Cryptogotchi, gameType models.GameType) (models.GameStat, string, error) {
	// generate a new token.
	gameStat := models.GameStat{
		CryptogotchiId: cryptogotchi.Id.String(),
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

func (svc *GameService) StopGame(token string, score float64) error {
	game, err := svc.GetGameByToken(token)
	if err != nil {
		return err
	}

	game.Score = score
	game.GameFinished = time.Now()
	err = svc.Save(&game)
	return err
}
