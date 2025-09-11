package handlers

import (
	"net/http"

	"foodflow/internal/core"
	"foodflow/internal/lib"

	"github.com/gin-gonic/gin"
)

type CreditsHandler struct {
	creditService core.CreditService
}

func NewCreditsHandler(creditService core.CreditService) *CreditsHandler {
	return &CreditsHandler{
		creditService: creditService,
	}
}

// GetCredits godoc
// @Summary Get user credits
// @Description Get the current user's credit balance and history
// @Tags Credits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} core.CreditsResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /credits [get]
func (h *CreditsHandler) GetCredits(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		lib.RespondError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	credits, err := h.creditService.GetUserCredits(c.Request.Context(), userID)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, credits)
}

// GetCreditHistory godoc
// @Summary Get credit history
// @Description Get the user's credit transaction history
// @Tags Credits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} core.CreditHistoryResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /credits/history [get]
func (h *CreditsHandler) GetCreditHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		lib.RespondError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	history, err := h.creditService.GetCreditHistory(c.Request.Context(), userID, page, limit)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, history)
}

// SpendCredits godoc
// @Summary Spend credits
// @Description Spend credits for a specific purpose
// @Tags Credits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.SpendCreditsRequest true "Credit spending data"
// @Success 200 {object} core.CreditTransactionResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /credits/spend [post]
func (h *CreditsHandler) SpendCredits(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		lib.RespondError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req core.SpendCreditsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lib.RespondError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	transaction, err := h.creditService.SpendCredits(c.Request.Context(), userID, &req)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, transaction)
}
