package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/heissonwillen/event-go/internal/config"
)

func SetupRouter(config config.Config, db *gorm.DB) *gin.Engine {
	router := gin.Default()
	stream := NewServer()

	authorized := router.Group("/authorized", gin.BasicAuth(gin.Accounts{
		config.BasicAuthUser: config.BasicAuthPassword,
	}))

	router.GET("/events", EventStreamHeadersMiddleware(), stream.serveHTTP(), GetEvents(config, db))
	authorized.POST("/events", PostEvent(config, db, stream))

	return router
}
