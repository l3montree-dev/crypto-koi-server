package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	Id        uuid.UUID `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	// check if no id is set - if there is an id set, do not overwrite it.
	// cryptogotchis are getting created by providing an id value.
	if b.Id.String() == "00000000-0000-0000-0000-000000000000" {
		b.Id = uuid
	}

	return nil
}

func (b *Base) BeforeUpdate(tx *gorm.DB) error {
	b.UpdatedAt = time.Now()
	return nil
}
