package models

import "gitlab.com/l3montree/cryptogotchi/clodhopper/internal/cqrs"

type Record struct {
	Base
	Type string `json:"type" gorm:"type:varchar(255)"`
	// contains the CQRS Event
	Payload string `json:"payload" gorm:"type:text"`
}

func FromEvent(event cqrs.Event) Record {
	record := Record{
		Type:    event.GetType(),
		Payload: event.ToJSON(),
	}
	return record
}

func (r Record) Parse() cqrs.Event {
	switch r.Type {

	}
	// could not parse the event payload
	return nil
}
