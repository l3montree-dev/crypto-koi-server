package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/http_util"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/models"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/repositories"
	"gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/service"
)

type OpenseaController struct {
	imageBaseUrl    string
	eventSvc        service.EventSvc
	cryptogotchiSvc service.CryptogotchiSvc
}

func NewOpenseaController(imageBaseUrl string, eventRepository repositories.EventRepository, cryptogotchiSvc service.CryptogotchiSvc) OpenseaController {
	return OpenseaController{
		eventSvc:        service.NewEventService(eventRepository),
		cryptogotchiSvc: cryptogotchiSvc,
		imageBaseUrl:    imageBaseUrl,
	}
}

func (c *OpenseaController) GetCryptogotchi(w http.ResponseWriter, req *http.Request) {
	tokenId := chi.URLParam(req, "tokenId")
	// fetch the correct cryptogotchi using the token.
	cryptogotchi, err := c.cryptogotchiSvc.GetCryptogotchiByUint256(tokenId)
	if err != nil {
		http_util.WriteHttpError(w, http.StatusNotFound, fmt.Sprintf("could not get cryptogotchi: %e", err))
		return
	}

	// transform the cryptogotchi to an opensea-NFT compatible json.
	nft, err := cryptogotchi.ToOpenseaNFT(c.imageBaseUrl)
	if err != nil {
		http_util.WriteHttpError(w, http.StatusInternalServerError, fmt.Sprintf("could not transform cryptogotchi to opensea-NFT: %e", err))
		return
	}
	http_util.WriteJSON(w, nft)
}

// returns the metadata for a cryptogotchi which does not exist.
// based on the provided tokenId.
func (c *OpenseaController) GetFakeCryptogotchi(w http.ResponseWriter, req *http.Request) {
	tokenId := chi.URLParam(req, "tokenId")

	nft, err := models.ToOpenseaNFT(c.imageBaseUrl, tokenId, true, "Fake", time.Now())
	if err != nil {
		http_util.WriteHttpError(w, http.StatusInternalServerError, fmt.Sprintf("could not transform cryptogotchi to opensea-NFT: %e", err))
		return
	}
	http_util.WriteJSON(w, nft)
}
