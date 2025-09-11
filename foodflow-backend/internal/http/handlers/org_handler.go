package handlers

import (
	"net/http"

	"foodflow/internal/core"
	"foodflow/internal/lib"

	"github.com/gin-gonic/gin"
)

type OrgHandler struct {
	orgService core.OrgService
}

func NewOrgHandler(orgService core.OrgService) *OrgHandler {
	return &OrgHandler{
		orgService: orgService,
	}
}

// CreateOrganization godoc
// @Summary Create organization
// @Description Create a new organization profile
// @Tags Organization
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.CreateOrgRequest true "Organization data"
// @Success 201 {object} core.OrgResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /organizations [post]
func (h *OrgHandler) CreateOrganization(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		lib.RespondError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req core.CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lib.RespondError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	org, err := h.orgService.CreateOrganization(c.Request.Context(), userID, &req)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusCreated, org)
}

// GetOrganization godoc
// @Summary Get organization
// @Description Get organization details by ID
// @Tags Organization
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organization ID"
// @Success 200 {object} core.OrgResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Failure 404 {object} lib.ErrorResponse
// @Router /organizations/{id} [get]
func (h *OrgHandler) GetOrganization(c *gin.Context) {
	orgID := c.Param("id")
	if orgID == "" {
		lib.RespondError(c, http.StatusBadRequest, "Organization ID is required")
		return
	}

	org, err := h.orgService.GetOrganization(c.Request.Context(), orgID)
	if err != nil {
		lib.RespondError(c, http.StatusNotFound, "Organization not found")
		return
	}

	lib.RespondSuccess(c, http.StatusOK, org)
}

// UpdateOrganization godoc
// @Summary Update organization
// @Description Update organization details
// @Tags Organization
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organization ID"
// @Param request body core.UpdateOrgRequest true "Organization update data"
// @Success 200 {object} core.OrgResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Failure 404 {object} lib.ErrorResponse
// @Router /organizations/{id} [put]
func (h *OrgHandler) UpdateOrganization(c *gin.Context) {
	orgID := c.Param("id")
	if orgID == "" {
		lib.RespondError(c, http.StatusBadRequest, "Organization ID is required")
		return
	}

	var req core.UpdateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lib.RespondError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	org, err := h.orgService.UpdateOrganization(c.Request.Context(), orgID, &req)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, org)
}

// ListOrganizations godoc
// @Summary List organizations
// @Description Get a list of organizations with pagination
// @Tags Organization
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} core.OrgListResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /organizations [get]
func (h *OrgHandler) ListOrganizations(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	orgs, err := h.orgService.ListOrganizations(c.Request.Context(), page, limit)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, orgs)
}
