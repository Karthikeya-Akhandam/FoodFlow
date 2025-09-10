package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"foodflow/internal/core"
	"foodflow/internal/core/repos"
)

type TokenService struct {
	tokenRepo  repos.TokenRepo
	collabRepo repos.CollabRepo
}

func NewTokenService(tokenRepo repos.TokenRepo, collabRepo repos.CollabRepo) core.TokenService {
	return &TokenService{
		tokenRepo:  tokenRepo,
		collabRepo: collabRepo,
	}
}

func (s *TokenService) GetUserTokens(ctx context.Context, userID string) (*core.TokensResponse, error) {
	// Get collaborator for this user
	collab, err := s.collabRepo.GetCollaboratorByUserID(ctx, userID)
	if err != nil {
		return nil, core.NewNotFoundError("Collaborator not found for user")
	}

	// Get all tokens for the collaborator
	allTokens, err := s.tokenRepo.GetTokensByCollaborator(ctx, collab.ID.String())
	if err != nil {
		return nil, core.NewInternalError("Failed to retrieve tokens", err)
	}

	// Calculate total tokens
	totalTokens := 0
	monthlyTokens := make([]core.TokenResponse, len(allTokens))
	
	for i, token := range allTokens {
		totalTokens += token.TokensEarned
		
		monthlyTokens[i] = core.TokenResponse{
			Month:        token.Month,
			TokensEarned: token.TokensEarned,
		}
	}

	return &core.TokensResponse{
		TotalTokens:   totalTokens,
		MonthlyTokens: monthlyTokens,
	}, nil
}

func (s *TokenService) GetTokenHistory(ctx context.Context, userID, page, limit string) (*core.TokenHistoryResponse, error) {
	// Parse pagination parameters
	pageNum := 1
	limitNum := 20
	if page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageNum = p
		}
	}
	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			limitNum = l
		}
	}

	// Get collaborator for this user
	collab, err := s.collabRepo.GetCollaboratorByUserID(ctx, userID)
	if err != nil {
		return nil, core.NewNotFoundError("Collaborator not found for user")
	}

	// Get token history for the collaborator
	tokenHistory, err := s.tokenRepo.GetTokenHistory(ctx, collab.ID.String(), limitNum, (pageNum-1)*limitNum)
	if err != nil {
		return nil, core.NewInternalError("Failed to retrieve token history", err)
	}

	// Convert to response format
	transactions := make([]core.TokenTransactionResponse, len(tokenHistory))
	for i, tx := range tokenHistory {
		transactions[i] = core.TokenTransactionResponse{
			TransactionID: tx.ID,
			Amount:        tx.TokensEarned,
			Type:          "EARNED",
			Description:   fmt.Sprintf("Tokens earned for %s", tx.Month),
			CreatedAt:     tx.CreatedAt,
		}
	}

	return &core.TokenHistoryResponse{
		Items: transactions,
		Page:  pageNum,
		Limit: limitNum,
		Total: len(transactions), // In real app, would be separate count query
	}, nil
}

func (s *TokenService) RedeemTokens(ctx context.Context, userID string, req *core.RedeemTokensRequest) (*core.TokenRedemptionResponse, error) {
	// Get collaborator for this user
	collab, err := s.collabRepo.GetCollaboratorByUserID(ctx, userID)
	if err != nil {
		return nil, core.NewNotFoundError("Collaborator not found for user")
	}

	// Get current month's tokens
	currentMonth := time.Now().Format("200601")
	token, err := s.tokenRepo.GetTokenByCollabAndMonth(ctx, collab.ID.String(), currentMonth)
	if err != nil {
		return nil, core.NewNotFoundError("No tokens available for current month")
	}

	// Check if sufficient tokens available
	if token.TokensEarned < req.Amount {
		return nil, core.NewValidationError("Insufficient tokens", map[string]interface{}{
			"requested": req.Amount,
			"available": token.TokensEarned,
		})
	}

	// Calculate credits earned (tokens to credits conversion)
	creditsEarned := req.Amount / core.DefaultTokenMultiplier

	// Redeem the tokens
	transaction, err := s.tokenRepo.RedeemTokens(ctx, collab.ID.String(), currentMonth, req.Amount, req.Reason)
	if err != nil {
		return nil, core.NewInternalError("Failed to redeem tokens", err)
	}

	return &core.TokenRedemptionResponse{
		RedemptionID:  transaction.ID,
		TokensSpent:   req.Amount,
		CreditsEarned: creditsEarned,
		Status:        "CONFIRMED",
	}, nil
}

func (s *TokenService) AwardTokens(ctx context.Context, collabID string, creditsSpent int) error {
	// Calculate tokens to award
	tokensToAward := creditsSpent * core.DefaultTokenMultiplier
	currentMonth := time.Now().Format("200601")

	// Update or create tokens for collaborator
	token, err := s.tokenRepo.GetTokenByCollabAndMonth(ctx, collabID, currentMonth)
	if err != nil {
		// Create new token record
		_, err = s.tokenRepo.CreateToken(ctx, collabID, currentMonth, tokensToAward)
		if err != nil {
			return err
		}
	} else {
		// Update existing token record
		err = s.tokenRepo.UpdateTokenEarned(ctx, token.ID.String(), token.TokensEarned+tokensToAward)
		if err != nil {
			return err
		}
	}

	return nil
}
