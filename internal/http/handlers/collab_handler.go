package handlers

import (
	"net/http"

	"foodflow/internal/core"
	"foodflow/internal/lib"

	"github.com/gin-gonic/gin"
)

type CollabHandler struct {
	collabService core.CollabService
}

func NewCollabHandler(collabService core.CollabService) *CollabHandler {
	return &CollabHandler{
		collabService: collabService,
	}
}

// CreateCollaborator godoc
// @Summary Create collaborator
// @Description Create a new collaborator profile
// @Tags Collaborator
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.CreateCollabRequest true "Collaborator data"
// @Success 201 {object} core.CollabResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /collaborators [post]
func (h *CollabHandler) CreateCollaborator(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		lib.RespondError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req core.CreateCollabRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lib.RespondError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	collab, err := h.collabService.CreateCollaborator(c.Request.Context(), userID, &req)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusCreated, collab)
}

// GetCollaborator godoc
// @Summary Get collaborator
// @Description Get collaborator details by ID
// @Tags Collaborator
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Collaborator ID"
// @Success 200 {object} core.CollabResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Failure 404 {object} lib.ErrorResponse
// @Router /collaborators/{id} [get]
func (h *CollabHandler) GetCollaborator(c *gin.Context) {
	collabID := c.Param("id")
	if collabID == "" {
		lib.RespondError(c, http.StatusBadRequest, "Collaborator ID is required")
		return
	}

	collab, err := h.collabService.GetCollaborator(c.Request.Context(), collabID)
	if err != nil {
		lib.RespondError(c, http.StatusNotFound, "Collaborator not found")
		return
	}

	lib.RespondSuccess(c, http.StatusOK, collab)
}

// UpdateCollaborator godoc
// @Summary Update collaborator
// @Description Update collaborator details
// @Tags Collaborator
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Collaborator ID"
// @Param request body core.UpdateCollabRequest true "Collaborator update data"
// @Success 200 {object} core.CollabResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Failure 404 {object} lib.ErrorResponse
// @Router /collaborators/{id} [put]
func (h *CollabHandler) UpdateCollaborator(c *gin.Context) {
	collabID := c.Param("id")
	if collabID == "" {
		lib.RespondError(c, http.StatusBadRequest, "Collaborator ID is required")
		return
	}

	var req core.UpdateCollabRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lib.RespondError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	collab, err := h.collabService.UpdateCollaborator(c.Request.Context(), collabID, &req)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, collab)
}

// ListCollaborators godoc
// @Summary List collaborators
// @Description Get a list of collaborators with pagination
// @Tags Collaborator
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} core.CollabListResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /collaborators [get]
func (h *CollabHandler) ListCollaborators(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	collabs, err := h.collabService.ListCollaborators(c.Request.Context(), page, limit)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, collabs)
}
