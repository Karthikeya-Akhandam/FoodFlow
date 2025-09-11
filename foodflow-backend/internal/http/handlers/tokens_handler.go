package handlers

import (
	"net/http"

	"foodflow/internal/core"
	"foodflow/internal/lib"

	"github.com/gin-gonic/gin"
)

type TokensHandler struct {
	tokenService core.TokenService
}

func NewTokensHandler(tokenService core.TokenService) *TokensHandler {
	return &TokensHandler{
		tokenService: tokenService,
	}
}

// GetTokens godoc
// @Summary Get user tokens
// @Description Get the current user's token balance and history
// @Tags Tokens
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} core.TokensResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /tokens [get]
func (h *TokensHandler) GetTokens(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		lib.RespondError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	tokens, err := h.tokenService.GetUserTokens(c.Request.Context(), userID)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, tokens)
}

// GetTokenHistory godoc
// @Summary Get token history
// @Description Get the user's token transaction history
// @Tags Tokens
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} core.TokenHistoryResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /tokens/history [get]
func (h *TokensHandler) GetTokenHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		lib.RespondError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	history, err := h.tokenService.GetTokenHistory(c.Request.Context(), userID, page, limit)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, history)
}

// RedeemTokens godoc
// @Summary Redeem tokens
// @Description Redeem tokens for rewards
// @Tags Tokens
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.RedeemTokensRequest true "Token redemption data"
// @Success 200 {object} core.TokenRedemptionResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 401 {object} lib.ErrorResponse
// @Router /tokens/redeem [post]
func (h *TokensHandler) RedeemTokens(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		lib.RespondError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req core.RedeemTokensRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lib.RespondError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	redemption, err := h.tokenService.RedeemTokens(c.Request.Context(), userID, &req)
	if err != nil {
		lib.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	lib.RespondSuccess(c, http.StatusOK, redemption)
}
