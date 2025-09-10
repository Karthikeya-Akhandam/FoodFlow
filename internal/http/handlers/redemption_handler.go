package handlers

import (
	"net/http"

	"foodflow/internal/core"
	"foodflow/internal/lib"

	"github.com/gin-gonic/gin"
)

type RedemptionHandler struct {
	redemptionService core.RedemptionService
}

func NewRedemptionHandler(redemptionService core.RedemptionService) *RedemptionHandler {
	return &RedemptionHandler{
		redemptionService: redemptionService,
	}
}

// CreateRedemption godoc
// @Summary Create redemption
// @Description Create a new redemption for an organization
// @Tags Redemption
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.CreateRedemptionRequest true "Redemption data"
// @Success 201 {object} core.RedemptionResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /redemptions [post]
func (h *RedemptionHandler) CreateRedemption(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		lib.RespondError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req core.CreateRedemptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lib.RespondError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	redemption, err := h.redemptionService.CreateRedemption(c.Request.Context(), userID, &req)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusCreated, redemption)
}

// GetRedemption godoc
// @Summary Get redemption
// @Description Get redemption details by ID
// @Tags Redemption
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Redemption ID"
// @Success 200 {object} core.RedemptionResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Failure 404 {object} lib.ErrorResponse
// @Router /redemptions/{id} [get]
func (h *RedemptionHandler) GetRedemption(c *gin.Context) {
	redemptionID := c.Param("id")
	if redemptionID == "" {
		lib.RespondError(c, http.StatusBadRequest, "Redemption ID is required")
		return
	}

	redemption, err := h.redemptionService.GetRedemption(c.Request.Context(), redemptionID)
	if err != nil {
		lib.RespondError(c, http.StatusNotFound, "Redemption not found")
		return
	}

	lib.RespondSuccess(c, http.StatusOK, redemption)
}

// UpdateRedemptionStatus godoc
// @Summary Update redemption status
// @Description Update the status of a redemption
// @Tags Redemption
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Redemption ID"
// @Param request body core.UpdateRedemptionStatusRequest true "Status update data"
// @Success 200 {object} core.RedemptionResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Failure 404 {object} lib.ErrorResponse
// @Router /redemptions/{id}/status [put]
func (h *RedemptionHandler) UpdateRedemptionStatus(c *gin.Context) {
	redemptionID := c.Param("id")
	if redemptionID == "" {
		lib.RespondError(c, http.StatusBadRequest, "Redemption ID is required")
		return
	}

	var req core.UpdateRedemptionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lib.RespondError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	redemption, err := h.redemptionService.UpdateRedemptionStatus(c.Request.Context(), redemptionID, &req)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, redemption)
}

// ListRedemptions godoc
// @Summary List redemptions
// @Description Get a list of redemptions with pagination and filters
// @Tags Redemption
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status"
// @Param org_id query string false "Filter by organization ID"
// @Success 200 {object} core.RedemptionListResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /redemptions [get]
func (h *RedemptionHandler) ListRedemptions(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	status := c.Query("status")
	orgID := c.Query("org_id")

	redemptions, err := h.redemptionService.ListRedemptions(c.Request.Context(), page, limit, status, orgID)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, redemptions)
}
