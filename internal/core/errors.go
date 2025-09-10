package core

import (
	"fmt"
	"net/http"
)

// AppError represents an application error
type AppError struct {
	Code       string
	Message    string
	Details    interface{}
	StatusCode int
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Error constructors
func NewValidationError(message string, details interface{}) *AppError {
	return &AppError{
		Code:       ErrorCodeValidation,
		Message:    message,
		Details:    details,
		StatusCode: http.StatusBadRequest,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:       ErrorCodeUnauthorized,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:       ErrorCodeForbidden,
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:       ErrorCodeNotFound,
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Code:       ErrorCodeConflict,
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Code:       ErrorCodeInternalError,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

func NewRateLimitedError(message string) *AppError {
	return &AppError{
		Code:       ErrorCodeRateLimited,
		Message:    message,
		StatusCode: http.StatusTooManyRequests,
	}
}

func NewIdempotencyKeyError(message string) *AppError {
	return &AppError{
		Code:       ErrorCodeIdempotencyKey,
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

func NewNotImplementedError(message string) *AppError {
	return &AppError{
		Code:       "NOT_IMPLEMENTED",
		Message:    message,
		StatusCode: http.StatusNotImplemented,
	}
}

// Common error messages
const (
	ErrMsgUserNotFound           = "user not found"
	ErrMsgInvalidCredentials     = "invalid email or password"
	ErrMsgUserAlreadyExists      = "user with this email already exists"
	ErrMsgInvalidRole            = "invalid user role"
	ErrMsgProfileNotFound        = "profile not found"
	ErrMsgOrganizationNotFound   = "organization not found"
	ErrMsgCollaboratorNotFound   = "collaborator not found"
	ErrMsgOfferNotFound          = "donation offer not found"
	ErrMsgClaimNotFound          = "claim not found"
	ErrMsgRedemptionNotFound     = "redemption not found"
	ErrMsgCreditNotFound         = "credit not found"
	ErrMsgTokenNotFound          = "token not found"
	ErrMsgRemoteAssignmentNotFound = "remote assignment not found"
	ErrMsgInsufficientCredits    = "insufficient credits"
	ErrMsgOfferExpired           = "donation offer has expired"
	ErrMsgOfferAlreadyClaimed    = "donation offer is already claimed"
	ErrMsgClaimAlreadyExists     = "claim already exists for this offer"
	ErrMsgInvalidServings        = "invalid number of servings"
	ErrMsgInvalidOfferStatus     = "invalid offer status for this operation"
	ErrMsgInvalidClaimStatus     = "invalid claim status for this operation"
	ErrMsgRemoteAssignmentExists = "remote assignment already exists"
	ErrMsgInvalidIdempotencyKey  = "invalid or expired idempotency key"
	ErrMsgRateLimitExceeded      = "rate limit exceeded"
)
