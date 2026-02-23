package order

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/example/ppo/pkg/apperror"
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
		handleError(c, err)
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
		handleError(c, err)
		return
	}

	response.OK(c, gin.H{"message": "order cancelled and refunded"})
}

func handleError(c *gin.Context, err error) {
	var appErr *apperror.Error
	if !errors.As(err, &appErr) {
		response.Err(c, http.StatusInternalServerError, "INTERNAL", "unexpected error")
		return
	}

	switch appErr.Kind {
	case apperror.KindNotFound:
		response.Err(c, http.StatusNotFound, "NOT_FOUND", appErr.Message)
	case apperror.KindValidation:
		response.Err(c, http.StatusBadRequest, "VALIDATION", appErr.Message)
	case apperror.KindConflict:
		response.Err(c, http.StatusConflict, "CONFLICT", appErr.Message)
	case apperror.KindUpstream:
		response.Err(c, http.StatusBadGateway, "UPSTREAM_ERROR", appErr.Message)
	default:
		response.Err(c, http.StatusInternalServerError, "INTERNAL", appErr.Message)
	}
}
