package routes

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heissonwillen/event-go/internal/config"
	"gorm.io/gorm"
)

type Event struct {
	Data          chan string
	Type          chan string
	NewClients    chan chan string
	ClosedClients chan chan string
	TotalClients  map[chan string]bool
}

type PostEventRequestBody struct {
	Data string `json:"data" binding:"required"`
	Type string `json:"type" binding:"required"`
}

type ClientChan chan string

// TODO: store data on DB
func PostEvent(config config.Config, db *gorm.DB, stream *Event) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestBody PostEventRequestBody
		if err := ctx.ShouldBindJSON(&requestBody); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Parsed request")
		stream.Data <- requestBody.Data
		// stream.Type <- requestBody.Type
		log.Printf("Got stream.Data")
		stream.Type <- "temperature"
		log.Printf("Sending JSON back")

		ctx.JSON(http.StatusOK, gin.H{
			"data": requestBody.Data,
			"type": requestBody.Type,
		})
	}
}

// TODO: return record from DB if no stream ever happened
func GetEvents(config config.Config, db *gorm.DB) gin.HandlerFunc {
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
			if data, ok := <-clientChan; ok {
				ctx.SSEvent("message", data)
				return true
			}
			return false
		})
	}
}

func NewServer() (event *Event) {
	event = &Event{
		Data:          make(chan string),
		Type:          make(chan string),
		NewClients:    make(chan chan string),
		ClosedClients: make(chan chan string),
		TotalClients:  make(map[chan string]bool),
	}

	go event.listen()

	return
}

func (stream *Event) listen() {
	for {
		select {
		case client := <-stream.NewClients:
			stream.TotalClients[client] = true
			log.Printf("Client added. %d registered clients", len(stream.TotalClients))

		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client)
			close(client)
			log.Printf("Removed client. %d registered clients", len(stream.TotalClients))

		case eventData := <-stream.Data:
			for clientDataChan := range stream.TotalClients {
				clientDataChan <- eventData
			}

		case eventType := <-stream.Type:
			for clientTypeChan := range stream.TotalClients {
				clientTypeChan <- eventType
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
