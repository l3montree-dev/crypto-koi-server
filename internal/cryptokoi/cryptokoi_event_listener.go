package cryptokoi

import (
	"time"

	"github.com/ethereum/go-ethereum/event"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

type CryptoKoiEventListener struct {
	binding *CryptoKoiBinding
}

type CryptoKoiEvent struct {
	TokenId string
	From    string
	To      string
}

func NewCryptoKoiEventListener(binding *CryptoKoiBinding) *CryptoKoiEventListener {
	return &CryptoKoiEventListener{
		binding: binding,
	}
}

func (c *CryptoKoiEventListener) init() (event.Subscription, chan *CryptoKoiBindingTransfer, error) {
	transfers := make(chan *CryptoKoiBindingTransfer)
	sub, err := c.binding.WatchTransfer(nil, transfers, nil, nil, nil)
	return sub, transfers, err
}

func (c *CryptoKoiEventListener) connect(eventChan chan<- CryptoKoiEvent) {
	sub, ch, err := c.init()
	if err != nil {
		orchardclient.Logger.Error(err)
		// try to reconnect.
		time.Sleep(time.Second * 5)
		orchardclient.Logger.Info("reconnecting...")
		c.connect(eventChan)
		return
	}

	orchardclient.Logger.Info("websocket connection established")
	for {
		select {
		case transfer := <-ch:
			orchardclient.Logger.Info("Transfer: ", transfer.TokenId.String(), " ", transfer.From.String(), " ", transfer.To.String())
			eventChan <- CryptoKoiEvent{
				TokenId: transfer.TokenId.String(),
				From:    transfer.From.String(),
				To:      transfer.To.String(),
			}

		case err := <-sub.Err():
			orchardclient.Logger.Error(err)
			// try to reconnect.
			time.Sleep(time.Second * 5)
			orchardclient.Logger.Info("reconnecting...")
			c.connect(eventChan)
			return
		}
	}
}

// has basic reconnection logic
func (c *CryptoKoiEventListener) StartListener() chan<- CryptoKoiEvent {
	eventChan := make(chan CryptoKoiEvent)
	go c.connect(eventChan)
	return eventChan
}
