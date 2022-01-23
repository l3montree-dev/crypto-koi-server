package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/http_util"
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

func (c *OpenseaController) GetCryptogotchi(w http.ResponseWriter, req *http.Request) {
	tokenId := chi.URLParam(req, "tokenId")
	// fetch the correct cryptogotchi using the token.
	cryptogotchi, err := c.cryptogotchiSvc.GetCryptogotchiByTokenId(tokenId)
	if err != nil {
		http_util.WriteHttpError(w, http.StatusNotFound, "could not get cryptogotchi: %e", err)
		return
	}
	// mutates the cryptogotchi struct
	cryptogotchi.ReplayEvents()
	// transform the cryptogotchi to an opensea-NFT compatible json.
	nft := cryptogotchi.ToOpenseaNFT()
	http_util.WriteJSON(w, nft)
}
