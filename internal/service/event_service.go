package service

import "gitlab.com/l3montree/cryptogotchi/clodhopper/internal/repositories"

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
