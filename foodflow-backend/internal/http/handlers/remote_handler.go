package handlers

import (
	"net/http"

	"foodflow/internal/core"
	"foodflow/internal/lib"

	"github.com/gin-gonic/gin"
)

type RemoteHandler struct {
	remoteService core.RemoteService
}

func NewRemoteHandler(remoteService core.RemoteService) *RemoteHandler {
	return &RemoteHandler{
		remoteService: remoteService,
	}
}

// CreateRemoteOrg godoc
// @Summary Create remote organization
// @Description Create a new remote organization assignment
// @Tags Remote Organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.CreateRemoteOrgRequest true "Remote organization data"
// @Success 201 {object} core.RemoteOrgResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /remote-orgs [post]
func (h *RemoteHandler) CreateRemoteOrg(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		lib.RespondError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req core.CreateRemoteOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lib.RespondError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	remoteOrg, err := h.remoteService.CreateRemoteOrg(c.Request.Context(), userID, &req)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusCreated, remoteOrg)
}

// GetRemoteOrg godoc
// @Summary Get remote organization
// @Description Get remote organization details by ID
// @Tags Remote Organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Remote Organization ID"
// @Success 200 {object} core.RemoteOrgResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Failure 404 {object} lib.ErrorResponse
// @Router /remote-orgs/{id} [get]
func (h *RemoteHandler) GetRemoteOrg(c *gin.Context) {
	remoteOrgID := c.Param("id")
	if remoteOrgID == "" {
		lib.RespondError(c, http.StatusBadRequest, "Remote Organization ID is required")
		return
	}

	remoteOrg, err := h.remoteService.GetRemoteOrg(c.Request.Context(), remoteOrgID)
	if err != nil {
		lib.RespondError(c, http.StatusNotFound, "Remote Organization not found")
		return
	}

	lib.RespondSuccess(c, http.StatusOK, remoteOrg)
}

// UpdateRemoteOrg godoc
// @Summary Update remote organization
// @Description Update remote organization details
// @Tags Remote Organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Remote Organization ID"
// @Param request body core.UpdateRemoteOrgRequest true "Remote organization update data"
// @Success 200 {object} core.RemoteOrgResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Failure 404 {object} lib.ErrorResponse
// @Router /remote-orgs/{id} [put]
func (h *RemoteHandler) UpdateRemoteOrg(c *gin.Context) {
	remoteOrgID := c.Param("id")
	if remoteOrgID == "" {
		lib.RespondError(c, http.StatusBadRequest, "Remote Organization ID is required")
		return
	}

	var req core.UpdateRemoteOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lib.RespondError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	remoteOrg, err := h.remoteService.UpdateRemoteOrg(c.Request.Context(), remoteOrgID, &req)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, remoteOrg)
}

// ListRemoteOrgs godoc
// @Summary List remote organizations
// @Description Get a list of remote organizations with pagination
// @Tags Remote Organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} core.RemoteOrgListResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /remote-orgs [get]
func (h *RemoteHandler) ListRemoteOrgs(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	remoteOrgs, err := h.remoteService.ListRemoteOrgs(c.Request.Context(), page, limit)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, remoteOrgs)
}

// AssignRemoteOrg godoc
// @Summary Assign remote organization
// @Description Assign a remote organization to handle claims for a specific area
// @Tags Remote Organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.AssignRemoteOrgRequest true "Assignment data"
// @Success 200 {object} core.RemoteOrgAssignmentResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /remote-orgs/assign [post]
func (h *RemoteHandler) AssignRemoteOrg(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		lib.RespondError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req core.AssignRemoteOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lib.RespondError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	assignment, err := h.remoteService.AssignRemoteOrg(c.Request.Context(), userID, &req)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, assignment)
}
