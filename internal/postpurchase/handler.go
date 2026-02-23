package postpurchase

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

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
	rg.GET("/users/:userId/installments", h.GetInstallments)
	rg.POST("/installments/pay", h.PayInstallment)
}

func (h *Handler) GetInstallments(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		response.Err(c, http.StatusBadRequest, "INVALID_REQUEST", "user_id is required")
		return
	}

	installments, err := h.svc.GetInstallments(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, installments)
}

func (h *Handler) PayInstallment(c *gin.Context) {
	var req PayInstallmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	resp, err := h.svc.PayInstallment(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
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
	case apperror.KindUpstream:
		response.Err(c, http.StatusBadGateway, "UPSTREAM_ERROR", appErr.Message)
	default:
		response.Err(c, http.StatusInternalServerError, "INTERNAL", appErr.Message)
	}
}
