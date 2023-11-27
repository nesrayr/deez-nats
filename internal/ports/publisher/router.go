package publisher

import (
	"deez-nats/internal/adapters/publisher"
	"deez-nats/pkg/logging"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(publisher publisher.IPublisher, logger logging.Logger) *gin.Engine {
	router := gin.Default()

	handler := NewHandler(publisher, logger)

	router.POST("/publish", handler.PublishMessage)

	return router
}
