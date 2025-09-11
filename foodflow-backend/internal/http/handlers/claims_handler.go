package handlers

import (
	"net/http"

	"foodflow/internal/core"
	"foodflow/internal/lib"

	"github.com/gin-gonic/gin"
)

type ClaimsHandler struct {
	claimService core.ClaimService
}

func NewClaimsHandler(claimService core.ClaimService) *ClaimsHandler {
	return &ClaimsHandler{
		claimService: claimService,
	}
}

// CreateClaim godoc
// @Summary Create donation claim
// @Description Create a new donation claim for an organization
// @Tags Claims
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.CreateClaimRequest true "Claim data"
// @Success 201 {object} core.ClaimResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /claims [post]
func (h *ClaimsHandler) CreateClaim(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		lib.RespondError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req core.CreateClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lib.RespondError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	offerID := c.Param("offer_id")
	if offerID == "" {
		lib.RespondError(c, http.StatusBadRequest, "Offer ID is required")
		return
	}
	
	claim, err := h.claimService.CreateClaim(c.Request.Context(), userID, offerID, &req)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusCreated, claim)
}

// GetClaim godoc
// @Summary Get donation claim
// @Description Get claim details by ID
// @Tags Claims
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Claim ID"
// @Success 200 {object} core.ClaimResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Failure 404 {object} lib.ErrorResponse
// @Router /claims/{id} [get]
func (h *ClaimsHandler) GetClaim(c *gin.Context) {
	claimID := c.Param("id")
	if claimID == "" {
		lib.RespondError(c, http.StatusBadRequest, "Claim ID is required")
		return
	}

	claim, err := h.claimService.GetClaim(c.Request.Context(), claimID)
	if err != nil {
		lib.RespondError(c, http.StatusNotFound, "Claim not found")
		return
	}

	lib.RespondSuccess(c, http.StatusOK, claim)
}

// UpdateClaimStatus godoc
// @Summary Update claim status
// @Description Update the status of a donation claim
// @Tags Claims
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Claim ID"
// @Param request body core.UpdateClaimStatusRequest true "Status update data"
// @Success 200 {object} core.ClaimResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Failure 404 {object} lib.ErrorResponse
// @Router /claims/{id}/status [put]
func (h *ClaimsHandler) UpdateClaimStatus(c *gin.Context) {
	claimID := c.Param("id")
	if claimID == "" {
		lib.RespondError(c, http.StatusBadRequest, "Claim ID is required")
		return
	}

	var req core.UpdateClaimStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lib.RespondError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	claim, err := h.claimService.UpdateClaimStatus(c.Request.Context(), claimID, &req)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, claim)
}

// ListClaims godoc
// @Summary List claims
// @Description Get a list of claims with pagination and filters
// @Tags Claims
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status"
// @Param org_id query string false "Filter by organization ID"
// @Success 200 {object} core.ClaimListResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /claims [get]
func (h *ClaimsHandler) ListClaims(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	status := c.Query("status")
	orgID := c.Query("org_id")

	claims, err := h.claimService.ListClaims(c.Request.Context(), page, limit, status, orgID)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, claims)
}
