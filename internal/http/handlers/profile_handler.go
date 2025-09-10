package handlers

import (
	"net/http"

	"foodflow/internal/core"
	"foodflow/internal/lib"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profileService core.ProfileService
}

func NewProfileHandler(profileService core.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get the current user's profile information
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} core.ProfileResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Failure 404 {object} lib.ErrorResponse
// @Router /profile [get]
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		lib.RespondError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	profile, err := h.profileService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		lib.RespondError(c, http.StatusNotFound, "Profile not found")
		return
	}

	lib.RespondSuccess(c, http.StatusOK, profile)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the current user's profile information
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} core.ProfileResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Failure 404 {object} lib.ErrorResponse
// @Router /profile [put]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		lib.RespondError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req core.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lib.RespondError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	profile, err := h.profileService.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, profile)
}
