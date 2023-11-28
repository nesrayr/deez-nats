package subscriber

import (
	"deez-nats/internal/repo"
	"deez-nats/pkg/logging"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IHandler interface {
	GetOrderByID(ctx *gin.Context)
}

type Handler struct {
	repo repo.IRepository
	l    *logging.Logger
}

func NewHandler(repo repo.IRepository, l logging.Logger) *Handler {
	return &Handler{repo: repo, l: &l}
}

func (h *Handler) GetOrderByID(ctx *gin.Context) {
	orderID := ctx.Param("id")

	order, err := h.repo.GetOrderById(orderID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
	}

	ctx.JSON(http.StatusOK, gin.H{"order": order})
}
