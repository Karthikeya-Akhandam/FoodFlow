package lib

import (
	"net/http"

	"foodflow/internal/core"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Responder struct{}

func NewResponder() *Responder {
	return &Responder{}
}

func (r *Responder) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func (r *Responder) Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

func (r *Responder) Accepted(c *gin.Context, data interface{}) {
	c.JSON(http.StatusAccepted, data)
}

func (r *Responder) NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func (r *Responder) Error(c *gin.Context, err error) {
	appErr, ok := err.(*core.AppError)
	if !ok {
		// Convert unknown errors to internal server error
		appErr = core.NewInternalError("Internal server error", err)
	}

	log.Error().
		Err(err).
		Str("code", appErr.Code).
		Int("status", appErr.StatusCode).
		Msg("Request error")

	errorResponse := core.ErrorResponse{
		Error: core.ErrorDetail{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		},
	}

	c.JSON(appErr.StatusCode, errorResponse)
}

func (r *Responder) ValidationError(c *gin.Context, message string, details interface{}) {
	r.Error(c, core.NewValidationError(message, details))
}

func (r *Responder) Unauthorized(c *gin.Context, message string) {
	r.Error(c, core.NewUnauthorizedError(message))
}

func (r *Responder) Forbidden(c *gin.Context, message string) {
	r.Error(c, core.NewForbiddenError(message))
}

func (r *Responder) NotFound(c *gin.Context, message string) {
	r.Error(c, core.NewNotFoundError(message))
}

func (r *Responder) Conflict(c *gin.Context, message string) {
	r.Error(c, core.NewConflictError(message))
}

func (r *Responder) RateLimited(c *gin.Context, message string) {
	r.Error(c, core.NewRateLimitedError(message))
}

// Global helper functions for convenience
func RespondSuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

func RespondError(c *gin.Context, statusCode int, message string) {
	errorResponse := core.ErrorResponse{
		Error: core.ErrorDetail{
			Code:    "ERROR",
			Message: message,
		},
	}
	c.JSON(statusCode, errorResponse)
}
