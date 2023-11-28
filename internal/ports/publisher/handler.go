package publisher

import (
	"deez-nats/internal/service/publisher"
	"deez-nats/pkg/logging"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IHandler interface {
	PublishMessage(ctx *gin.Context)
}

type Handler struct {
	publisher publisher.IPublisher
	l         *logging.Logger
}

func NewHandler(publisher publisher.IPublisher, l logging.Logger) *Handler {
	return &Handler{publisher: publisher, l: &l}
}

func (h *Handler) PublishMessage(ctx *gin.Context) {
	var payload map[string]interface{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.publisher.PublishData(payload, "subject")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
