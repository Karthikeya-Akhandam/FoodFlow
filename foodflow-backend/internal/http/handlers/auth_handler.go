package handlers

import (
	"foodflow/internal/core"
	"foodflow/internal/http/middleware"
	"foodflow/internal/lib"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService core.AuthService
	validator   *validator.Validate
	responder   *lib.Responder
}

func NewAuthHandler(authService core.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator.New(),
		responder:   lib.NewResponder(),
	}
}

// Signup handles user registration
// @Summary Sign up a new user
// @Description Register a new user with email, password and role
// @Tags auth
// @Accept json
// @Produce json
// @Param request body core.SignupRequest true "Signup request"
// @Success 201 {object} core.UserResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 409 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /v1/auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
	var req core.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.responder.ValidationError(c, "Invalid request body", err.Error())
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		h.responder.ValidationError(c, "Validation failed", err.Error())
		return
	}

	user, err := h.authService.Signup(c.Request.Context(), &req)
	if err != nil {
		h.responder.Error(c, err)
		return
	}

	h.responder.Created(c, user)
}

// Login handles user authentication
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body core.LoginRequest true "Login request"
// @Success 200 {object} core.AuthResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req core.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.responder.ValidationError(c, "Invalid request body", err.Error())
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		h.responder.ValidationError(c, "Validation failed", err.Error())
		return
	}

	auth, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		h.responder.Error(c, err)
		return
	}

	h.responder.Success(c, auth)
}

// GetMe returns current user information
// @Summary Get current user
// @Description Get information about the currently authenticated user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} core.UserResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 404 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /v1/auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		h.responder.Unauthorized(c, "User ID not found in context")
		return
	}

	user, err := h.authService.GetMe(c.Request.Context(), userID)
	if err != nil {
		h.responder.Error(c, err)
		return
	}

	h.responder.Success(c, user)
}
