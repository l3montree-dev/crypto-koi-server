package controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/repositories"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/service"
)

type OpenseaController struct {
	eventSvc        service.EventSvc
	cryptogotchiSvc service.CryptogotchiSvc
}

func NewOpenseaController(eventRepository repositories.EventRepository, cryptogotchiRepository repositories.CryptogotchiRepository) OpenseaController {
	return OpenseaController{
		eventSvc:        service.NewEventService(eventRepository),
		cryptogotchiSvc: service.NewCryptogotchiService(cryptogotchiRepository),
	}
}

func (c *OpenseaController) GetCryptogotchi(ctx *fiber.Ctx) error {
	tokenId := ctx.Params("tokenId")
	// fetch the correct cryptogotchi using the token.
	cryptogotchi, err := c.cryptogotchiSvc.GetCryptogotchiByTokenId(tokenId)
	if err != nil {
		return err
	}
	// mutates the cryptogotchi struct
	cryptogotchi.ReplayEvents()
	// transform the cryptogotchi to an opensea-NFT compatible json.
	nft := cryptogotchi.ToOpenseaNFT()
	return ctx.JSON(nft)
}
