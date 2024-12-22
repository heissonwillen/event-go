package models

import "time"

type Send struct {
	// TODO: make 'Id' an UUID
	Id        string `gorm:"primary_key"`
	Message   string
	CreatedAt time.Time
}
