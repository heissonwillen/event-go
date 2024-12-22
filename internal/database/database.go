package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/heissonwillen/event-go/internal/config"
)

func InitDB(config config.Config) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(config.SQLiteDBPath), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
