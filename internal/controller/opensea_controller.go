package controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/repositories"
)

type OpenseaController struct {
	recordRepository       repositories.RecordRepository
	cryptogotchiRepository repositories.CryptogotchiRepository
}

func NewOpenseaController(recordRepository repositories.RecordRepository, cryptogotchiRepository repositories.CryptogotchiRepository) OpenseaController {
	return OpenseaController{
		recordRepository:       recordRepository,
		cryptogotchiRepository: cryptogotchiRepository,
	}
}

func (c *OpenseaController) GetCryptogotchi(ctx *fiber.Ctx) error {
	tokenId := ctx.Params("tokenId")
	// fetch the correct cryptogotchi using the token.
	cryptogotchi, err := c.cryptogotchiRepository.GetCryptogotchiByTokenId(tokenId)
	if err != nil {
		return err
	}
	// transform the cryptogotchi to an opensea-NFT compatible json.
	nft := cryptogotchi.ToOpenseaNFT()
	ctx.JSON(nft)
	return nil
}
