package handlers

import (
	"strconv"

	"foodflow/internal/core"
	"foodflow/internal/http/middleware"
	"foodflow/internal/lib"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type OffersHandler struct {
	offerService    core.OfferService
	matchingService core.MatchingService
	validator       *validator.Validate
	responder       *lib.Responder
}

func NewOffersHandler(offerService core.OfferService, matchingService core.MatchingService) *OffersHandler {
	return &OffersHandler{
		offerService:    offerService,
		matchingService: matchingService,
		validator:       validator.New(),
		responder:       lib.NewResponder(),
	}
}

// CreateOffer handles donation offer creation
// @Summary Create a donation offer
// @Description Create a new donation offer (Collaborator only)
// @Tags offers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.CreateOfferRequest true "Create offer request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 403 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /v1/offers [post]
func (h *OffersHandler) CreateOffer(c *gin.Context) {
	var req core.CreateOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.responder.ValidationError(c, "Invalid request body", err.Error())
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		h.responder.ValidationError(c, "Validation failed", err.Error())
		return
	}

	userID, _ := middleware.GetUserID(c)

	offer, err := h.offerService.CreateOffer(c.Request.Context(), userID, &req)
	if err != nil {
		h.responder.Error(c, err)
		return
	}

	h.responder.Created(c, gin.H{
		"offer_id": offer.OfferID,
		"status":   offer.Status,
	})
}

// GetOffers handles listing offers
// @Summary List offers
// @Description Get list of offers (Collaborator sees own offers, Admin sees all)
// @Tags offers
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status" Enums(OPEN,PENDING_CONFIRM,CLAIMED,EXPIRED,CANCELLED)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} core.PaginatedResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /v1/offers [get]
func (h *OffersHandler) GetOffers(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	userRole, _ := middleware.GetUserRole(c)

	// Parse query parameters
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > core.MaxPageLimit {
		limit = core.DefaultPageLimit
	}

	offset := (page - 1) * limit

	offers, total, err := h.offerService.GetOffers(c.Request.Context(), userID, userRole, status, limit, offset)
	if err != nil {
		h.responder.Error(c, err)
		return
	}

	// Convert to response format
	items := make([]interface{}, len(offers))
	for i, offer := range offers {
		items[i] = gin.H{
			"offer_id":           offer.OfferID,
			"title":              offer.Title,
			"status":             offer.Status,
			"estimated_servings": offer.EstimatedServings,
			"purpose":            offer.Purpose,
			"ready_from":         offer.ReadyFrom,
			"expires_at":         offer.ExpiresAt,
		}
	}

	h.responder.Success(c, core.PaginatedResponse{
		Items: items,
		Page:  page,
		Limit: limit,
		Total: total,
	})
}

// GetOffer handles getting a single offer
// @Summary Get offer details
// @Description Get detailed information about a specific offer
// @Tags offers
// @Produce json
// @Security BearerAuth
// @Param id path string true "Offer ID"
// @Success 200 {object} core.OfferResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 403 {object} core.ErrorResponse
// @Failure 404 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /v1/offers/{id} [get]
func (h *OffersHandler) GetOffer(c *gin.Context) {
	offerID := c.Param("id")

	offer, err := h.offerService.GetOfferByID(c.Request.Context(), offerID)
	if err != nil {
		h.responder.Error(c, err)
		return
	}

	h.responder.Success(c, offer)
}

// UpdateOffer handles offer updates
// @Summary Update an offer
// @Description Update an existing donation offer (Collaborator only)
// @Tags offers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Offer ID"
// @Param request body core.UpdateOfferRequest true "Update offer request"
// @Success 200 {object} core.OfferResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 403 {object} core.ErrorResponse
// @Failure 404 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /v1/offers/{id} [put]
func (h *OffersHandler) UpdateOffer(c *gin.Context) {
	offerID := c.Param("id")
	
	var req core.UpdateOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.responder.ValidationError(c, "Invalid request body", err.Error())
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		h.responder.ValidationError(c, "Validation failed", err.Error())
		return
	}

	offer, err := h.offerService.UpdateOffer(c.Request.Context(), offerID, &req)
	if err != nil {
		h.responder.Error(c, err)
		return
	}

	h.responder.Success(c, offer)
}

// CancelOffer handles offer cancellation
// @Summary Cancel an offer
// @Description Cancel a donation offer (Owner or Admin only)
// @Tags offers
// @Produce json
// @Security BearerAuth
// @Param id path string true "Offer ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 403 {object} core.ErrorResponse
// @Failure 404 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /v1/offers/{id}/cancel [post]
func (h *OffersHandler) CancelOffer(c *gin.Context) {
	offerID := c.Param("id")
	userID, _ := middleware.GetUserID(c)

	offer, err := h.offerService.CancelOffer(c.Request.Context(), offerID, userID)
	if err != nil {
		h.responder.Error(c, err)
		return
	}

	h.responder.Success(c, gin.H{
		"offer_id": offer.OfferID,
		"status":   offer.Status,
	})
}

// GetNearbyOffers handles getting nearby offers for organizations
// @Summary Get nearby offers
// @Description Get offers near the organization's location with priority scores
// @Tags offers
// @Produce json
// @Security BearerAuth
// @Param scope query string true "Search scope" Enums(pincode,city,state)
// @Param purpose query string false "Filter by purpose" Enums(CHILDREN,ELDERLY,WOMEN,GENERAL,EMERGENCY)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} core.PaginatedResponse
// @Failure 400 {object} core.ErrorResponse
// @Failure 401 {object} core.ErrorResponse
// @Failure 403 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /v1/offers/nearby [get]
func (h *OffersHandler) GetNearbyOffers(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	purpose := c.Query("purpose")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > core.MaxPageLimit {
		limit = core.DefaultPageLimit
	}

	offset := (page - 1) * limit

	offers, total, err := h.offerService.GetNearbyOffersForUser(c.Request.Context(), userID, purpose, limit, offset)
	if err != nil {
		h.responder.Error(c, err)
		return
	}

	// Convert to response format
	items := make([]interface{}, len(offers))
	for i, offer := range offers {
		items[i] = offer
	}

	h.responder.Success(c, core.PaginatedResponse{
		Items: items,
		Page:  page,
		Limit: limit,
		Total: total,
	})
}
