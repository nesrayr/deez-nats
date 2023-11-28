package subscriber

import (
	"deez-nats/internal/repo"
	"deez-nats/pkg/logging"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(repo repo.IRepository, logger logging.Logger) *gin.Engine {
	router := gin.Default()

	handler := NewHandler(repo, logger)

	router.GET("/order/:id", handler.GetOrderByID)

	return router
}
