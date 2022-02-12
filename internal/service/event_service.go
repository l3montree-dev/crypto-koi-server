package service

import "gitlab.com/l3montree/crypto-koi/crypto-koi-api/internal/repositories"

type EventSvc interface {
	repositories.EventRepository
}

type EventService struct {
	repositories.EventRepository
}

func NewEventService(rep repositories.EventRepository) EventSvc {
	return &EventService{
		EventRepository: rep,
	}
}
