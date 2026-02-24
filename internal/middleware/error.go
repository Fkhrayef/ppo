package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/example/ppo/pkg/apperror"
	"github.com/example/ppo/pkg/response"
)

// ErrorHandler is a Gin middleware that translates apperror.Kind values into
// the correct HTTP status code. Handlers only need to call c.Error(err) and
// return — this middleware takes care of writing the JSON error response.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

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
}
