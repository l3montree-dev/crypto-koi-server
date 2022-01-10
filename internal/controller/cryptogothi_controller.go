package controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/repositories"
)

type CryptogotchiController struct {
	recordRepository       repositories.RecordRepository
	cryptogotchiRepository repositories.CryptogotchiRepository
}

func NewCryptogotchiController(recordRepository repositories.RecordRepository, cryptogotchiRepository repositories.CryptogotchiRepository) CryptogotchiController {
	return CryptogotchiController{
		recordRepository:       recordRepository,
		cryptogotchiRepository: cryptogotchiRepository,
	}
}

// opensea.io integration
func (c *CryptogotchiController) GetCryptogotchi(ctx *fiber.Ctx) error {
	return nil
}
