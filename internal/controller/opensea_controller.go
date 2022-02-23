package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/generator"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/http_util"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/repositories"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/service"
)

type OpenseaController struct {
	imageBaseUrl    string
	eventSvc        service.EventSvc
	cryptogotchiSvc service.CryptogotchiSvc
	generator       *generator.Generator
}

func NewOpenseaController(imageBaseUrl string, generator *generator.Generator, eventRepository repositories.EventRepository, cryptogotchiSvc service.CryptogotchiSvc) OpenseaController {
	return OpenseaController{
		eventSvc:        service.NewEventService(eventRepository),
		cryptogotchiSvc: cryptogotchiSvc,
		imageBaseUrl:    imageBaseUrl,
		generator:       generator,
	}
}

func (c *OpenseaController) GetCryptogotchi(w http.ResponseWriter, req *http.Request) {
	tokenId := chi.URLParam(req, "tokenId")
	// fetch the correct cryptogotchi using the token.
	cryptogotchi, err := c.cryptogotchiSvc.GetCryptogotchiByUint256(tokenId)
	if err != nil {
		http_util.WriteHttpError(w, http.StatusNotFound, "could not get cryptogotchi: %e", err)
		return
	}

	// transform the cryptogotchi to an opensea-NFT compatible json.
	nft, err := cryptogotchi.ToOpenseaNFT(c.imageBaseUrl, c.generator)
	if err != nil {
		http_util.WriteHttpError(w, http.StatusInternalServerError, "could not transform cryptogotchi to opensea-NFT: %e", err)
		return
	}
	http_util.WriteJSON(w, nft)
}
