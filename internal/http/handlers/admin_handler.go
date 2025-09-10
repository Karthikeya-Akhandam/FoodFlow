package handlers

import (
	"net/http"

	"foodflow/internal/core"
	"foodflow/internal/lib"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminService core.AdminService
}

func NewAdminHandler(adminService core.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// GetSystemStats godoc
// @Summary Get system statistics
// @Description Get overall system statistics and metrics
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} core.SystemStatsResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /admin/stats [get]
func (h *AdminHandler) GetSystemStats(c *gin.Context) {
	stats, err := h.adminService.GetSystemStats(c.Request.Context())
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, stats)
}

// GetUserStats godoc
// @Summary Get user statistics
// @Description Get user-related statistics and metrics
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} core.UserStatsResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /admin/users/stats [get]
func (h *AdminHandler) GetUserStats(c *gin.Context) {
	stats, err := h.adminService.GetUserStats(c.Request.Context())
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, stats)
}

// GetDonationStats godoc
// @Summary Get donation statistics
// @Description Get donation-related statistics and metrics
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} core.DonationStatsResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /admin/donations/stats [get]
func (h *AdminHandler) GetDonationStats(c *gin.Context) {
	stats, err := h.adminService.GetDonationStats(c.Request.Context())
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, stats)
}

// ListUsers godoc
// @Summary List all users
// @Description Get a list of all users with pagination and filters
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param role query string false "Filter by role"
// @Param status query string false "Filter by status"
// @Success 200 {object} core.UserListResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /admin/users [get]
func (h *AdminHandler) ListUsers(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	role := c.Query("role")
	status := c.Query("status")

	users, err := h.adminService.ListUsers(c.Request.Context(), page, limit, role, status)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, users)
}

// UpdateUserStatus godoc
// @Summary Update user status
// @Description Update the status of a user (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body core.UpdateUserStatusRequest true "Status update data"
// @Success 200 {object} core.UserResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Failure 404 {object} lib.ErrorResponse
// @Router /admin/users/{id}/status [put]
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		lib.RespondError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	var req core.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lib.RespondError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.adminService.UpdateUserStatus(c.Request.Context(), userID, &req)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, user)
}

// GetAuditLogs godoc
// @Summary Get audit logs
// @Description Get system audit logs with pagination and filters
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param action query string false "Filter by action"
// @Param user_id query string false "Filter by user ID"
// @Success 200 {object} core.AuditLogResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /admin/audit-logs [get]
func (h *AdminHandler) GetAuditLogs(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	action := c.Query("action")
	userID := c.Query("user_id")

	logs, err := h.adminService.GetAuditLogs(c.Request.Context(), page, limit, action, userID)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, logs)
}
