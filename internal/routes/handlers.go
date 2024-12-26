package routes

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heissonwillen/event-go/internal/config"
	"github.com/heissonwillen/event-go/internal/models"
	"gorm.io/gorm"
)

type Event struct {
	Messages      chan EventMessage
	NewClients    chan chan EventMessage
	ClosedClients chan chan EventMessage
	TotalClients  map[chan EventMessage]bool
}

type EventMessage struct {
	Data string
	Type string
}

type PostEventRequestBody struct {
	Data string `json:"data" binding:"required"`
	Type string `json:"type" binding:"required"`
}

type ClientChan chan EventMessage

func PostEvent(config config.Config, db *gorm.DB, stream *Event) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestBody PostEventRequestBody
		if err := ctx.ShouldBindJSON(&requestBody); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		eventMessage := EventMessage(requestBody)

		stream.Messages <- eventMessage

		eventDB := models.Event{
			Data: requestBody.Data,
			Type: requestBody.Type,
		}
		db.Create(&eventDB)

		ctx.JSON(http.StatusOK, gin.H{
			"data": requestBody.Data,
			"type": requestBody.Type,
		})
	}
}

// GetEvents streams events to connected clients
func GetEvents(config config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		v, ok := ctx.Get("clientChan")
		if !ok {
			return
		}
		clientChan, ok := v.(ClientChan)
		if !ok {
			return
		}
		ctx.Stream(func(w io.Writer) bool {
			if event, ok := <-clientChan; ok {
				ctx.SSEvent(event.Type, event.Data)
				return true
			}
			return false
		})
	}
}

// NewServer initializes the Event struct and starts listening for events
func NewServer(db *gorm.DB) (event *Event) {
	event = &Event{
		Messages:      make(chan EventMessage),
		NewClients:    make(chan chan EventMessage),
		ClosedClients: make(chan chan EventMessage),
		TotalClients:  make(map[chan EventMessage]bool),
	}

	go event.listen(db)

	return
}

// listen handles broadcasting events to clients and managing client connections
func (stream *Event) listen(db *gorm.DB) {
	for {
		select {
		case client := <-stream.NewClients:
			stream.TotalClients[client] = true
			log.Printf("Client added. %d registered clients", len(stream.TotalClients))

			// Propagate the latest event of each type to all clients once there's a new connection
			var latestEvents []models.Event

			// Query to get the latest event of each type
			if err := db.Raw(`
				SELECT *
				FROM (
					SELECT *, ROW_NUMBER() OVER (PARTITION BY type ORDER BY created_at DESC) AS row_num
					FROM events
				) subquery
				WHERE row_num = 1
			`).Scan(&latestEvents).Error; err == nil {
				for _, latestEvent := range latestEvents {
					message := EventMessage{
						Data: latestEvent.Data,
						Type: latestEvent.Type,
					}
					for clientChan := range stream.TotalClients {
						clientChan <- message
					}
				}
			}

		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client)
			close(client)
			log.Printf("Removed client. %d registered clients", len(stream.TotalClients))

		case message := <-stream.Messages:
			for clientChan := range stream.TotalClients {
				clientChan <- message
			}
		}
	}
}

func (stream *Event) serveHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientChan := make(ClientChan)
		stream.NewClients <- clientChan

		defer func() {
			go func() {
				for range clientChan {
				}
			}()
			stream.ClosedClients <- clientChan
		}()

		c.Set("clientChan", clientChan)
		c.Next()
	}
}
