package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Event struct {
	Id        string    `gorm:"type:text;primary_key"`
	Data      string    `gorm:"not null"`
	Type      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// BeforeCreate generates a UUID before inserting the record into the database
func (e *Event) BeforeCreate(tx *gorm.DB) (err error) {
	if e.Id == "" {
		e.Id = uuid.NewString()
	}
	return
}
