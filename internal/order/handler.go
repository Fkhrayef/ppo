package order

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/example/ppo/pkg/response"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	orders := rg.Group("/orders")
	orders.POST("", h.CreateOrder)
	orders.POST("/:id/cancel", h.CancelOrder)
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	o, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.Created(c, ToResponse(o))
}

func (h *Handler) CancelOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Err(c, http.StatusBadRequest, "INVALID_ID", "invalid order id")
		return
	}

	if err := h.svc.Cancel(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}

	response.OK(c, gin.H{"message": "order cancelled and refunded"})
}
