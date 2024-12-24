package main

import (
	"log"
	"net/http"

	"github.com/heissonwillen/event-go/internal/config"
	"github.com/heissonwillen/event-go/internal/database"
	"github.com/heissonwillen/event-go/internal/models"
	"github.com/heissonwillen/event-go/internal/routes"
)

var version = "dev"

func main() {
	log.Printf("Starting event-go: %s", version)

	config := config.LoadConfig()

	db := database.InitDB(config)
	db.AutoMigrate(&models.Event{})

	router := routes.SetupRouter(config, db)

	server := &http.Server{
		Addr:    config.ListenAddr,
		Handler: router,
	}

	log.Printf("Listening on %s", config.ListenAddr)
	log.Fatal(server.ListenAndServe())
}
