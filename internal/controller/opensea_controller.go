package controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/repositories"
)

type OpenseaController struct {
	eventRepository        repositories.EventRepository
	cryptogotchiRepository repositories.CryptogotchiRepository
}

func NewOpenseaController(eventRepository repositories.EventRepository, cryptogotchiRepository repositories.CryptogotchiRepository) OpenseaController {
	return OpenseaController{
		eventRepository:        eventRepository,
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
	// mutates the cryptogotchi struct
	cryptogotchi.ReplayEvents()
	// transform the cryptogotchi to an opensea-NFT compatible json.
	nft := cryptogotchi.ToOpenseaNFT()
	return ctx.JSON(nft)
}
